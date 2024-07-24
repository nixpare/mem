package mem

import (
	"unsafe"
)

type Slice[T any] []T

func NewSlice[T any](len, cap int, alloc AllocNStrategy) Slice[T] {
	if len == 0 && cap == 0 {
		return nil
	}
	
	var tmp T
	size := unsafe.Sizeof(tmp)

	if size == 0 {
		panic("slice: unsupported 0-byte elements")
	}

	_, overflow := mulUintptr(size, uintptr(cap))
	if overflow || len < 0 || len > cap {
		_, overflow := mulUintptr(size, uintptr(len))
		if overflow || len < 0 {
			panic("makeslice: len out of range")
		}
		panic("makeslice: cap out of range")
	}

	p := alloc(cap, size, unsafe.Alignof(tmp))
	if p == nil {
		panic("slice: allocation failed")
	}

	return unsafe.Slice((*T)(p), cap)[:len]
}

func (v *Slice[T]) Append(free FreeStrategy, alloc AllocNStrategy, elems ...T) bool {
	oldLen := len(*v)
	newLen := oldLen + len(elems)

	var realloc bool
	if newLen > cap(*v) {
		realloc = true
		v.growslice(newLen, free, alloc)
	} else {
		*v = unsafe.Slice((*T)(v.Pointer()), cap(*v))[:newLen]
	}

	copy((*v)[oldLen:], elems)
	return realloc
}

func (v Slice[T]) Pointer() unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(v))
}

func (v *Slice[T]) growslice(newLen int, free FreeStrategy, alloc AllocNStrategy) {
	var tmp T
	size := unsafe.Sizeof(tmp)
	if size == 0 {
		panic("growslice: unsupported 0-byte elements")
	}

	if newLen < 0 {
		panic("growslice: len out of range")
	}

	newCap := nextslicecap(newLen, cap(*v))
	
	p := alloc(newCap, size, unsafe.Alignof(tmp))
	if p == nil {
		panic("slice: allocation failed")
	}
	
	newV := unsafe.Slice((*T)(p), newCap)[:newLen]
	copy(newV, *v)

	if free != nil {
		free(v.Pointer())
	}
	*v = newV
}

func ToSliceMatrix[T any](m Slice[Slice[T]]) [][]T {
	return unsafe.Slice(
		(*[]T)(m.Pointer()),
		cap(m),
	)[:len(m)]
}

func FromSliceMatrix[T any](m [][]T) Slice[Slice[T]] {
	return unsafe.Slice(
		(*Slice[T])(unsafe.SliceData(m)),
		cap(m),
	)[:len(m)]
}

func FreeSlice[T any](v *Slice[T], free FreeStrategy) {
	free(v.Pointer())
}

func nextslicecap(newLen, oldCap int) int {
	newcap := oldCap
	doublecap := newcap + newcap
	if newLen > doublecap {
		return newLen
	}

	const threshold = 256
	if oldCap < threshold {
		return doublecap
	}
	for {
		// Transition from growing 2x for small slices
		// to growing 1.25x for large slices. This formula
		// gives a smooth-ish transition between the two.
		newcap += (newcap + 3*threshold) >> 2

		// We need to check `newcap >= newLen` and whether `newcap` overflowed.
		// newLen is guaranteed to be larger than zero, hence
		// when newcap overflows then `uint(newcap) > uint(newLen)`.
		// This allows to check for both with the same comparison.
		if uint(newcap) >= uint(newLen) {
			break
		}
	}

	// Set newcap to the requested cap when
	// the newcap calculation overflowed.
	if newcap <= 0 {
		return newLen
	}
	return newcap
}
