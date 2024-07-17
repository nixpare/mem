package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	a1, a2, a3 := StringFromGO("ciao", Malloc), StringFromGO("bella", Malloc), StringFromGO("addio", Malloc)
	defer func() { Free(a1.Pointer()); Free(a2.Pointer()); Free(a3.Pointer()) }()

	s := a1.Clone()
	defer Free(s.Pointer())
	s.Append(a2, a3)
	
	assert.Equal(t, 14, len(s))
	assert.Equal(t, "ciaobellaaddio", string(s))

	s2 := s.Clone()
	defer Free(s2.Pointer())
	assert.NotEqual(t, s.Pointer(), s2.Pointer())

	s2.Set("saluti")
	assert.Equal(t, "ciaobellaaddio", string(s))
	assert.Equal(t, "saluti", string(s2))

	s3 := a1.Clone()
	defer Free(s3.Pointer())
	oldPtr := s3.Pointer()
	s3.Append(a2, a3)
	assert.Equal(t, oldPtr, s3.Pointer())
	assert.Equal(t, "ciaobellaaddio", string(s3))
}