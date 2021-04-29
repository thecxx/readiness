package readiness

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	ErrorInvalidFetchFunction = errors.New("invalid fetch function")
)

type SyncFunc func() (interface{}, error)

type buffer struct {
	// Last sync time
	st int64
	// Last sync error
	se error
	// Last valid value
	vv interface{}
}

type Watcher struct {
	slider int32
	// 2 buffers
	buffers [2]buffer

	vv atomic.Value

	// Sync function
	fun func() (interface{}, error)
	// Atomic
	lock uint32
	// Expiration
	expire int64
}

func NewWatcher(fun SyncFunc, expire time.Duration) *Watcher {
	return &Watcher{
		buffers: [2]buffer{
			{0, nil, nil},
			{0, nil, nil},
		},
		slider: 0,
		fun:    fun,
		lock:   0,
		expire: expire.Nanoseconds(),
	}
}

func (w *Watcher) Get() (interface{}, error) {
	return w.get()
}

func (w *Watcher) GetDefault(def interface{}) interface{} {
	if v, err := w.get(); err == nil {
		return v
	}
	return def
}

func (w *Watcher) get() (interface{}, error) {

	v, t, e := w.getBuffer()

	if time.Now().UnixNano() < t+w.expire {
		return v, e
	}
	// Lock
	if !atomic.CompareAndSwapUint32(&w.lock, 0, 1) {
		return v, e
	}
	// Unlock
	atomic.StoreUint32(&w.lock, 0)

	// Fetch data
	value, err := w.fetch()
	if err != nil {
		value = v
	}
	// Update buffer
	w.setBuffer(value, err)

	return value, err
}

func (w *Watcher) fetch() (interface{}, error) {
	if w.fun != nil {
		return w.fun()
	}
	return nil, ErrorInvalidFetchFunction
}

func (w *Watcher) getBuffer() (interface{}, int64, error) {
	buf := w.buffers[atomic.LoadInt32(&w.slider)]
	return buf.vv, buf.st, buf.se
}

func (w *Watcher) setBuffer(value interface{}, err error) {
	// Other buffer
	other := atomic.LoadInt32(&w.slider) ^ 1
	// Update buffer
	w.buffers[other].se = err
	w.buffers[other].st = time.Now().UnixNano()
	w.buffers[other].vv = value
	// Slide to other buffer
	atomic.StoreInt32(&w.slider, other)
}
