package mem

import (
	"unsafe"
)

type Slice[T any] []T

func NewSlice[T any](len, cap int, allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer) Slice[T] {
	if len == 0 && cap == 0 {
		return nil
	}
	
	var tmp T
	size := unsafe.Sizeof(tmp)

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

	p := allocStrategy(cap, size, unsafe.Alignof(tmp))
	if uintptr(p) == 0 {
		panic("new: allocation failed")
	}

	return unsafe.Slice((*T)(p), cap)[:len]
}

func (v *Slice[T]) Append(
	freeStrategy func(unsafe.Pointer),
	allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer,
	elems ...T,
) {
	oldLen := len(*v)
	newLen := oldLen + len(elems)
	if newLen > cap(*v) {
		v.growslice(newLen, freeStrategy, allocStrategy)
	} else {
		*v = unsafe.Slice((*T)(v.Pointer()), cap(*v))[:newLen]
	}

	copy((*v)[oldLen:], elems)
}

func (v Slice[T]) Pointer() unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(v))
}

func (v *Slice[T]) growslice(
	newLen int,
	freeStrategy func(unsafe.Pointer),
	allocStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer,
) {
	var tmp T
	size := unsafe.Sizeof(tmp)
	if size == 0 {
		panic("growslice: unsupported 0-byte elements")
	}

	if newLen < 0 {
		panic("growslice: len out of range")
	}

	newCap := nextslicecap(newLen, cap(*v))
	
	p := allocStrategy(newCap, size, unsafe.Alignof(tmp))
	if uintptr(p) == 0 {
		panic("new: allocation failed")
	}
	
	newV := unsafe.Slice((*T)(p), newCap)[:newLen]
	copy(newV, *v)

	if freeStrategy != nil {
		freeStrategy(v.Pointer())
	}
	*v = newV
}

func FreeSlice[T any](v *Slice[T], freeStrategy func(p unsafe.Pointer)) {
	freeStrategy(v.Pointer())
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
