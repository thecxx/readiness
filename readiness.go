package readiness

import (
	"sync/atomic"
	"time"

	cache "github.com/patrickmn/go-cache"
)

type PullHandlerFunc func(string) (interface{}, error)

type object struct {
	pull    PullHandlerFunc
	value   interface{}
	expire  time.Duration
	syncAt  time.Time
	syncing int32
}

type Readiness struct {
	// Local cache fields.
	storage *cache.Cache
	// Some callback.
	onPullFailed func(string, error)
}

var (
	// A default object.
	r = New()
)

// Get a new Readiness.
func New(opts ...Option) *Readiness {
	r := &Readiness{
		storage:      cache.New(0, 0),
		onPullFailed: nil,
	}
	// Apply option
	for _, f := range opts {
		f(r)
	}
	return r
}

// Get value with key.
func (r *Readiness) Get(key string) interface{} {
	return r.GetDefault(key, nil)
}

// Get value with key.
func (r *Readiness) GetDefault(key string, def interface{}) interface{} {
	// Check the cache
	value, ok := r.storage.Get(key)
	if !ok {
		return def
	}
	obj := value.(*object)
	// If there is no expiration time or no expiration
	if obj.value != nil && (obj.expire <= 0 || obj.syncAt.Add(obj.expire).After(time.Now())) {
		return obj.value
	}
	// From remote
	if value := r.sync(key, obj); value != nil {
		return value
	}
	// From cache
	if obj.value != nil {
		return obj.value
	}
	return def
}

// Get value with key.
func Get(key string) interface{} {
	return r.Get(key)
}

// Get value with key.
func GetDefault(key string, def interface{}) interface{} {
	return r.GetDefault(key, def)
}

// Sync a key-value.
func (r *Readiness) sync(key string, obj *object) interface{} {
	// Sync
	if !atomic.CompareAndSwapInt32(&obj.syncing, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&obj.syncing, 0)
	// Getting
	value, err := obj.pull(key)
	if err != nil {
		if r.onPullFailed != nil {
			r.onPullFailed(key, err)
		}
		return nil
	}
	// Save
	if value != nil {
		obj.value = value
		obj.syncAt = time.Now()
	}
	return value
}

// Register a new data source.
func (r *Readiness) Register(key string, pull PullHandlerFunc, expire time.Duration) {
	r.storage.SetDefault(key, &object{pull: pull, value: nil, expire: expire, syncing: 0})
}

// Register a new data source.
func Register(key string, pull PullHandlerFunc, expire time.Duration) {
	r.Register(key, pull, expire)
}
