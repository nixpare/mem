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
	assert.Equal(t, 4, a.Len())
	assert.Equal(t, 6, a.Cap())
	assert.Equal(t, []int{0, 0, 0, 0}, a.ToGoSlice())

	a.Set(0, 1)
	a.Set(1, 2)
	a.Set(2, 3)
	a.Set(3, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, a.ToGoSlice())

	b := a.Subslice(2, 6)
	defer b.Free()
	assert.Equal(t, 4, b.Len())
	assert.Equal(t, 4, b.Cap())
	assert.Equal(t, []int{3, 4, 0, 0}, b.ToGoSlice())

	b = b.Subslice(0, 2)
	defer b.Free()
	assert.Equal(t, []int{3, 4}, b.ToGoSlice())
}