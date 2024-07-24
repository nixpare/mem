package mem

import (
	"unsafe"
)

type String string

func NewString(size int, alloc AllocNStrategy) String {
	p := alloc(size, 1, 1)
	if p == nil {
		panic("string: allocation failed")
	}

	return String(unsafe.String((*byte)(p), size))
}

func (s String) Bytes() []byte {
	return unsafe.Slice((*byte)(s.Pointer()), len(s))
}

func (s String) Clone(alloc AllocNStrategy) String {
	str := NewString(len(s), alloc)
	copy(str.Bytes(), s)
	return str
}

func (s *String) Set(str string, free FreeStrategy, alloc AllocNStrategy) {
	if len(*s) != len(str) {
		s.resize(len(str), free, alloc)
	}

	copy(s.Bytes(), str)
}

func (s *String) Append(free FreeStrategy, alloc AllocNStrategy, a ...String) {
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
	
	s.resize(newL, free, alloc)
	b := s.Bytes()[oldL:]
	
	for _, x := range a {
		copy(b, x)
		b = b[len(x):]
	}
}

func (s *String) resize(size int, free FreeStrategy, alloc AllocNStrategy) {
	newStr := NewString(size, alloc)
	copy(newStr.Bytes(), *s)
	
	if free != nil {
		free(s.Pointer())
	}
	*s = newStr
}

func (s String) Pointer() unsafe.Pointer {
	return unsafe.Pointer(unsafe.StringData(string(s)))
}

func FreeString(s *String, free FreeStrategy) {
	free(s.Pointer())
}

func StringFromGO(s string, alloc AllocNStrategy) String {
	str := NewString(len(s), alloc)
	str.Set(s, nil, alloc)
	return str
}
