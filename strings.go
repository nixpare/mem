package mem

import (
	"unsafe"
)

type String struct {
	str unsafe.Pointer
	len int
}

func (s String) Len() int {
	return int(s.len)
}

func (s String) Bytes() []byte {
	return unsafe.Slice((*byte)(s.str), s.len)
}

func (s String) String() string {
	return unsafe.String((*byte)(s.str), s.len)
}

func (s String) Clone() String {
	return s.Concat()
}

func (s String) Free() {
	Free(s.str)
}

func (s *String) Set(str string) {
	strLen := len(str)
	if s.len != strLen {
		s.resize(strLen)
	}

	copy(s.Bytes(), str)
}

func (s *String) Append(a ...String) {
	oldL := s.len
	newL := oldL
	for _, x := range a {
		n := x.len
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
		copy(b, unsafe.Slice((*byte)(x.str), x.len))
		b = b[x.len:]
	}
}

func (s String) Concat(a ...String) String {
	l := s.len
	for _, x := range a {
		n := x.len
		if n == 0 {
			continue
		}
		if l+n < l {
			panic("string concatenation too long")
		}
		l += n
	}
	
	str := newString(l, false)
	b := str.Bytes()
	
	copy(b, unsafe.Slice((*byte)(s.str), s.len))
	b = b[s.len:]
	for _, x := range a {
		copy(b, unsafe.Slice((*byte)(x.str), x.len))
		b = b[x.len:]
	}

	return str
}

func (s *String) resize(size int) {
	s.str, s.len = Realloc(s.str, size), size
}

func NewString(size int) String {
	return newString(size, true)
}

func StringFromGO(s string) String {
	return String{
		str: unsafe.Pointer(unsafe.StringData(s)),
		len: len(s),
	}.Clone()
}

func newString(size int, zeored bool) String {
	var p unsafe.Pointer
	if !zeored {
		p = Malloc(size)
	} else {
		p = Calloc(size, 1)
	}
	
	return String{ str: p, len: size }
}