package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynSlice(t *testing.T) {
	testAlloc, free, _, allocN := testAllocFree(t, true)
	defer testAlloc()

	a := NewDynSlice[int](3, 4, 4, free, allocN)
	defer a.Free()

	assert.Equal(t, 3, a.chunk)
	assert.Equal(t, []int{0, 0, 0, 0}, a.ToGoSlice())
	assert.Equal(t, 6, a.Cap())

	a.Set(0, 1)
	a.Set(1, 2)
	a.Set(2, 3)
	a.Set(3, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, a.ToGoSlice())
	assert.Equal(t, 6, a.Cap())

	b := a.Subslice(2, 6)
	defer b.Free()
	assert.Equal(t, []int{3, 4, 0, 0}, b.ToGoSlice())
	assert.Equal(t, 4, b.Cap())

	b = b.Subslice(0, 2)
	defer b.Free()
	assert.Equal(t, []int{3, 4}, b.ToGoSlice())

	a.Append(5, 6, 7, 8, 9, 10, 11, 12, 13)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, a.ToGoSlice())
	assert.Equal(t, 15, a.Cap())
}