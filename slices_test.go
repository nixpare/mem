package mem

import (
	"runtime"
	"runtime/debug"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlices(t *testing.T) {
	v := NewSlice[int](0, 2, MallocN)
	defer Free(v.Pointer())
	assert.Equal(t, 0, len(v))
	assert.Equal(t, 2, cap(v))

	v.Append(10, 20, 30)
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
}

func TestLongSliceUsage(t *testing.T) {
	oldGC := debug.SetGCPercent(1)
	defer debug.SetGCPercent(oldGC)

	oldLimit := debug.SetMemoryLimit(1024)
	defer debug.SetMemoryLimit(oldLimit)

	n := 1024 * 1024 * 4
	v := NewSlice[int](n, n, MallocN)
	defer Free(v.Pointer())

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