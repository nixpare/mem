package mem

import (
	"unsafe"
)

type Slice[T any] []T

func NewSlice[T any](len, cap int, allocStrategy func(n int, sizeof uintptr) unsafe.Pointer) Slice[T] {
	var x T
	size := unsafe.Sizeof(x)

	if size == 0 {
		panic("makeslice: unsupported 0-byte elements")
	}

	_, overflow := mulUintptr(size, uintptr(cap))
	if overflow || len < 0 || len > cap {
		_, overflow := mulUintptr(size, uintptr(len))
		if overflow || len < 0 {
			panic("makeslice: len out of range")
		}
		panic("makeslice: cap out of range")
	}

	array := allocStrategy(cap, size)
	return unsafe.Slice((*T)(array), cap)[:len]
}

func (v *Slice[T]) Append(elems ...T) {
	oldLen := len(*v)
	newLen := oldLen + len(elems)
	if newLen > cap(*v) {
		v.growslice(newLen)
	} else {
		*v = unsafe.Slice((*T)(v.Pointer()), cap(*v))[:newLen]
	}

	copy((*v)[oldLen:], elems)
}

func (v Slice[T]) Pointer() unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(v))
}

func (v *Slice[T]) growslice(newLen int) {
	var x T
	size := unsafe.Sizeof(x)
	if size == 0 {
		panic("growslice: unsupported 0-byte elements")
	}

	if newLen < 0 {
		panic("growslice: len out of range")
	}

	oldCap := cap(*v)
	newCap := nextslicecap(newLen, oldCap)
	
	capmem := uintptr(newCap) * size
	array := Realloc(v.Pointer(), capmem)
	*v = unsafe.Slice((*T)(array), newCap)[:newLen]
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
