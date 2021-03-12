package readiness

import (
	"testing"
	"time"
)

func TestReadiness_Get(t *testing.T) {
	d := ""
	r := New()
	r.Register("unit_test",
		func(key string) (interface{}, error) {
			d = time.Now().String()
			return d, nil
		},
		2*time.Second,
	)
	for i := 0; i < 10; i++ {
		value := r.Get("unit_test")
		if value != nil && value.(string) == d {
			t.Logf("Value is %s\n", value.(string))
		} else {
			t.Errorf("Getting failed\n")
		}
		time.Sleep(1 * time.Second)
	}

}
