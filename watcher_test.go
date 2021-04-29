package readiness

import (
	"fmt"
	"testing"
	"time"
)

func TestWatcher_Get(t *testing.T) {

	w := NewWatcher(func() (i interface{}, err error) {
		return time.Now().String(), nil
	}, time.Second*2)

	for i := 0; i < 10; i++ {

		v := w.GetDefault(nil)

		fmt.Printf("%#v\n", v)

		time.Sleep(time.Second)
	}

}
