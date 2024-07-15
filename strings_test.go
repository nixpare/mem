package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	a1, a2, a3 := StringFromGO("ciao"), StringFromGO("bella"), StringFromGO("addio")
	defer func() { a1.Free(); a2.Free(); a3.Free() }()

	s := a1.Concat(a2, a3)
	defer s.Free()
	
	assert.Equal(t, uintptr(14), s.len)
	assert.Equal(t, "ciaobellaaddio", s.String())

	s2 := s.Clone()
	defer s2.Free()
	assert.NotEqual(t, s.str, s2.str)

	s2.Set("saluti")
	assert.Equal(t, "ciaobellaaddio", s.String())
	assert.Equal(t, "saluti", s2.String())

	s3 := a1.Clone()
	defer s3.Free()
	oldPtr := s3.str
	s3.Append(a2, a3)
	assert.Equal(t, oldPtr, s3.str)
	assert.Equal(t, "ciaobellaaddio", s3.String())
}