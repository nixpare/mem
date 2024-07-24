package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	testAlloc, free, _, allocN := testAllocFree(t, true)
	defer testAlloc()

	a1, a2, a3 := StringFromGO("ciao", allocN), StringFromGO("bella", allocN), StringFromGO("addio", allocN)
	defer func() { FreeString(&a1, free); FreeString(&a2, free); FreeString(&a3, free) }()

	s := a1.Clone(allocN)
	defer FreeString(&s, free)
	assert.Equal(t, 4, len(s))
	assert.Equal(t, "ciao", string(s))

	s.Append(free, allocN, a2, a3)
	assert.Equal(t, 14, len(s))
	assert.Equal(t, "ciaobellaaddio", string(s))

	s2 := s.Clone(allocN)
	defer FreeString(&s2, free)
	assert.NotEqual(t, s.Pointer(), s2.Pointer())

	s2.Set("saluti", free, allocN)
	assert.Equal(t, "ciaobellaaddio", string(s))
	assert.Equal(t, "saluti", string(s2))

	s3 := a1.Clone(allocN)
	defer FreeString(&s3, free)
	oldPtr := s3.Pointer()
	s3.Append(free, allocN, a2, a3)
	assert.NotEqual(t, oldPtr, s3.Pointer())
	assert.Equal(t, "ciaobellaaddio", string(s3))
}
