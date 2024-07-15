package mem

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlices(t *testing.T) {
	v := NewSlice[int](0, 2)
	defer v.Free()
	assert.Equal(t, 0, len(v))
	assert.Equal(t, 2, cap(v))

	v.Append(10, 20, 30)
	assert.Equal(t, 3, len(v))
	assert.Equal(t, 4, cap(v))
	assert.Equal(t, true, slices.Equal([]int{10, 20, 30}, v), "Got v(%d, %d) -> %v", len(v), cap(v), v)
}