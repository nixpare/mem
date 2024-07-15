package mem

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	x := New[int]()
	defer x.Free()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		y := x

		go func() {
			time.Sleep(time.Second)
			
			assert.Equal(t, 10, *y.Value())
			assert.Equal(t, *x.Value(), *y.Value())
			assert.Equal(t, x.ptr, y.ptr)

			wg.Done()
		}()

		*y.Value() = 10
		wg.Done()
	}()

	wg.Wait()
}
