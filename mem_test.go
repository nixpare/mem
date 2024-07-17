package mem

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	x := New[int](MallocZero)
	defer Free(ObjPointer(x))

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		y := x

		go func() {
			time.Sleep(time.Second)
			
			assert.Equal(t, 10, *y)
			assert.Equal(t, *x, *y)
			assert.Equal(t, ObjPointer(x), ObjPointer(x))

			wg.Done()
		}()

		*y = 10
		wg.Done()
	}()

	wg.Wait()
}
