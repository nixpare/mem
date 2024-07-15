package mem

import (
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	x := New[int]()
	defer FreeObj(x)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		y := x

		go func() {
			time.Sleep(time.Second)
			
			assert.Equal(t, 10, *y)
			assert.Equal(t, *x, *y)
			assert.Equal(t, unsafe.Pointer(x), unsafe.Pointer(y))

			wg.Done()
		}()

		*y = 10
		wg.Done()
	}()

	wg.Wait()
}
