package mem

import (
	"runtime"
	"runtime/debug"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlices(t *testing.T) {
	testAlloc, free, _, allocN := testAllocFree(t, true)
	defer testAlloc()

	v := NewSlice[int](0, 2, allocN)
	defer FreeSlice(&v, free)
	assert.Equal(t, 0, len(v))
	assert.Equal(t, 2, cap(v))

	v.Append(free, allocN, 10, 20, 30)
	assert.Equal(t, 3, len(v))
	assert.Equal(t, 4, cap(v))
	assert.Equal(t, 10, v[0])
	assert.Equal(t, 20, v[1])
	assert.Equal(t, 30, v[2])
	assert.Equal(t, true, slices.Equal([]int{10, 20, 30}, v), "expected %v, got %v", []int{10, 20, 30}, v)

	v2 := v[1:]
	assert.Equal(t, 2, len(v2))
	assert.Equal(t, 3, cap(v2))
	assert.NotEqual(t, v.Pointer(), v2.Pointer())
	assert.Equal(t, true, slices.Equal([]int{20, 30}, v2), "expected %v, got %v", []int{20, 30}, v2)

	v2[0] = 25
	assert.Equal(t, true, slices.Equal([]int{10, 25, 30}, v), "expected %v, got %v", []int{10, 25, 30}, v)
	assert.Equal(t, true, slices.Equal([]int{25, 30}, v2), "expected %v, got %v", []int{25, 30}, v2)

	v3 := NewSlice[int](0, 0, allocN)
	defer FreeSlice(&v3, free)
	assert.Equal(t, 0, len(v3))
	assert.Equal(t, 0, cap(v3))

	v3.Append(nil, allocN, 1)
	assert.Equal(t, 1, len(v3))
	assert.Equal(t, 1, cap(v3))
	assert.Equal(t, 1, v3[0])
}

func TestLongSliceUsage(t *testing.T) {
	testAlloc, free, _, allocN := testAllocFree(t, true)
	defer testAlloc()

	oldGC := debug.SetGCPercent(1)
	defer debug.SetGCPercent(oldGC)

	oldLimit := debug.SetMemoryLimit(1024)
	defer debug.SetMemoryLimit(oldLimit)

	n := 1024 * 1024 * 4
	v := NewSlice[int](n, n, allocN)
	defer FreeSlice(&v, free)

	for i := range v {
		v[i] = i
	}

	for j := range 10_000 {
		for i, x := range v {
			if i != x {
				assert.Equal(t, i, x)
			}
		}

		if j % 1000 == 0 {
			runtime.GC()
		}
	}
}