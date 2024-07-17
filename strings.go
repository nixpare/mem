package mem

import (
	"unsafe"
)

type String string

func NewString(size int, allocStrategy func(size uintptr) unsafe.Pointer) String {
	strLen := uintptr(size)
	p := allocStrategy(strLen)
	return String(unsafe.String((*byte)(p), int(size)))
}

func (s String) Bytes() []byte {
	return unsafe.Slice(unsafe.StringData(string(s)), len(s))
}

func (s String) Clone() String {
	ptr := Malloc(uintptr(len(s)))
	str := String(unsafe.String((*byte)(ptr), len(s)))
	copy(str.Bytes(), s.Bytes())
	return str
}

func (s *String) Set(str string) {
	if len(*s) != len(str) {
		s.resize(len(str))
	}

	copy(s.Bytes(), str)
}

func (s *String) Append(a ...String) {
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
	
	s.resize(newL)
	b := s.Bytes()[oldL:]
	
	for _, x := range a {
		copy(b, x.Bytes())
		b = b[len(x):]
	}
}

func (s *String) resize(size int) {
	*s = String(unsafe.String(
		(*byte)(Realloc(s.Pointer(), uintptr(size))),
		size,
	))
}

func (s String) Pointer() unsafe.Pointer {
	return unsafe.Pointer(unsafe.StringData(string(s)))
}

func StringFromGO(s string, allocStrategy func(size uintptr) unsafe.Pointer) String {
	str := NewString(len(s), allocStrategy)
	str.Set(s)
	return str
}
