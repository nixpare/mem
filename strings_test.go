package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	a1, a2, a3 := StringFromGO("ciao"), StringFromGO("bella"), StringFromGO("addio")
	defer func() { a1.Free(); a2.Free(); a3.Free() }()

	s := a1.Clone()
	defer s.Free()
	s.Append(a2, a3)
	
	assert.Equal(t, 14, len(s))
	assert.Equal(t, "ciaobellaaddio", string(s))

	s2 := s.Clone()
	defer s2.Free()
	assert.NotEqual(t, s.pointer(), s2.pointer())

	s2.Set("saluti")
	assert.Equal(t, "ciaobellaaddio", string(s))
	assert.Equal(t, "saluti", string(s2))

	s3 := a1.Clone()
	defer s3.Free()
	oldPtr := s3.pointer()
	s3.Append(a2, a3)
	assert.Equal(t, oldPtr, s3.pointer())
	assert.Equal(t, "ciaobellaaddio", string(s3))
}