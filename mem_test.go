package mem

import (
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testAlloc, free, alloc, _ := testAllocFree(t, true)
	defer testAlloc()

	x := New[int](alloc)
	defer FreeObject(x, free)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		y := x

		go func() {
			time.Sleep(time.Second)
			
			assert.Equal(t, 10, *y)
			assert.Equal(t, *x, *y)
			assert.Equal(t, ObjectPointer(x), ObjectPointer(x))

			wg.Done()
		}()

		*y = 10
		wg.Done()
	}()

	wg.Wait()
}

func testAllocFree(t *testing.T, debug bool) (func(), FreeStrategy, AllocStrategy, AllocNStrategy) {
	m := make(map[unsafe.Pointer]bool)

	testFunc := func() {
		t.Helper()

		for ptr, freed := range m {
			if !freed {
				t.Errorf("memory leak detected: %v", ptr)
			}
		}
	}

	free := func(ptr unsafe.Pointer) {
		t.Helper()

		if debug {
			t.Logf("freeing %v", ptr)
		}

		if ptr == nil {
			return
		}

		freed, ok := m[ptr]
		if !ok {
			t.Errorf("freeing unknown pointer: %v", ptr)
			return
		}

		if freed {
			t.Errorf("double free detected: %v", ptr)
			return
		}

		m[ptr] = true
	}

	alloc := func(sizeof, alignof uintptr) unsafe.Pointer {
		t.Helper()

		ptr := Calloc(1, sizeof, alignof)
		m[ptr] = false

		if debug {
			t.Logf("allocating %v (%d : %d)", ptr, sizeof, alignof)
		}
		return ptr
	}

	allocN := func(n int, sizeof, alignof uintptr) unsafe.Pointer {
		t.Helper()

		ptr := Calloc(n, sizeof, alignof)
		m[ptr] = false

		if debug {
			t.Logf("allocating %v (%d x %d : %d)", ptr, n, sizeof, alignof)
		}
		return ptr
	}

	return testFunc, free, alloc, allocN
}
