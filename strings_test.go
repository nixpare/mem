package mem

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func testAllocString(n int, sizeof, alignof uintptr) unsafe.Pointer {
	return MallocN(n, sizeof)
}

func TestString(t *testing.T) {
	a1, a2, a3 := StringFromGO("ciao", testAllocString), StringFromGO("bella", testAllocString), StringFromGO("addio", testAllocString)
	defer func() { FreeString(&a1, Free); FreeString(&a2, Free); FreeString(&a3, Free) }()

	s := a1.Clone(testAllocString)
	defer FreeString(&s, Free)
	assert.Equal(t, 4, len(s))
	assert.Equal(t, "ciao", string(s))

	s.Append(Free, testAllocString, a2, a3)
	assert.Equal(t, 14, len(s))
	assert.Equal(t, "ciaobellaaddio", string(s))

	s2 := s.Clone(testAllocString)
	defer FreeString(&s2, Free)
	assert.NotEqual(t, s.Pointer(), s2.Pointer())

	s2.Set("saluti", Free, testAllocString)
	assert.Equal(t, "ciaobellaaddio", string(s))
	assert.Equal(t, "saluti", string(s2))

	s3 := a1.Clone(testAllocString)
	defer FreeString(&s3, Free)
	oldPtr := s3.Pointer()
	s3.Append(Free, testAllocString, a2, a3)
	assert.NotEqual(t, oldPtr, s3.Pointer())
	assert.Equal(t, "ciaobellaaddio", string(s3))
}
