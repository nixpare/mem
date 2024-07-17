package mem

import (
	"unsafe"
)

type String string

func NewString(size int, allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer) String {
	p := allocStrategy(size, 1, 1)
	if uintptr(p) == 0 {
		panic("newstring: allocation failed")
	}

	return String(unsafe.String((*byte)(p), size))
}

func (s String) Bytes() []byte {
	return unsafe.Slice((*byte)(s.Pointer()), len(s))
}

func (s String) Clone(allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer) String {
	p := allocStrategy(len(s), 1, 1)
	if uintptr(p) == 0 {
		panic("clonestring: allocation failed")
	}

	str := NewString(len(s), allocStrategy)
	copy(str.Bytes(), s)
	return str
}

func (s *String) Set(
	str string,
	freeStrategy func(p unsafe.Pointer),
	allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer,
) {
	if len(*s) != len(str) {
		s.resize(len(str), freeStrategy, allocStrategy)
	}

	copy(s.Bytes(), str)
}

func (s *String) Append(
	freeStrategy func(p unsafe.Pointer),
	allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer,
	a ...String,
) {
	oldL := len(*s)
	newL := oldL
	for _, x := range a {
		n := len(x)
		if n == 0 {
			continue
		}
		if newL+n < newL {
			panic("string concatenation too long")
		}
		newL += n
	}

	if newL == oldL {
		return
	}
	
	s.resize(newL, freeStrategy, allocStrategy)
	b := s.Bytes()[oldL:]
	
	for _, x := range a {
		copy(b, x)
		b = b[len(x):]
	}
}

func (s *String) resize(
	size int,
	freeStrategy func(p unsafe.Pointer),
	allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer,
) {
	newStr := NewString(size, allocStrategy)
	copy(newStr.Bytes(), *s)
	
	if freeStrategy != nil {
		freeStrategy(s.Pointer())
	}
	*s = newStr
}

func (s String) Pointer() unsafe.Pointer {
	return unsafe.Pointer(unsafe.StringData(string(s)))
}

func FreeString(s *String, freeStrategy func(p unsafe.Pointer)) {
	freeStrategy(s.Pointer())
}

func StringFromGO(s string, allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer) String {
	str := NewString(len(s), allocStrategy)
	str.Set(s, nil, allocStrategy)
	return str
}
