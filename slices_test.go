package mem

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlices(t *testing.T) {
	v := NewSlice[int](0, 2)
	defer v.Free()
	assert.Equal(t, 0, v.Len())
	assert.Equal(t, 2, v.Cap())

	v.Append(10, 20, 30)
	assert.Equal(t, 3, v.Len())
	assert.Equal(t, 4, v.Cap())
	assert.Equal(t, 10, v.Get(0))
	assert.Equal(t, 20, v.Get(1))
	assert.Equal(t, 30, v.Get(2))
	assert.Equal(t, true, slices.Equal([]int{10, 20, 30}, v.Slice()))
}

func TestLongSliceUsage(t *testing.T) {
	n := 1024 * 1024 * 40
	v := NewSlice[int](n, n)
	defer v.Free()

	for i := range v.Slice() {
		v.Set(i, i+1)
	}

	for range 100 {
		for i, x := range v.Slice() {
			v.Set(i, (x+1) % v.Len())
		}
	}
}