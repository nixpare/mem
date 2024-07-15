package mem

import (
	"fmt"
	"unsafe"
)

type Slice[T any] struct {
	array unsafe.Pointer
	len int
	cap int
	subslice bool
	offset uintptr
}

func NewSlice[T any](len, cap int) Slice[T] {
	var x T
	size := unsafe.Sizeof(x)

	if size == 0 {
		panic("makeslice: unsupported 0-byte elements")
	}

	mem, overflow := mulUintptr(size, uintptr(cap))
	if overflow || len < 0 || len > cap {
		_, overflow := mulUintptr(size, uintptr(len))
		if overflow || len < 0 {
			panic("makeslice: len out of range")
		}
		panic("makeslice: cap out of range")
	}

	return Slice[T]{
		array: Malloc(int(mem)),
		len: len,
		cap: cap,
		offset: 0,
		subslice: false,
	}
}

func (v Slice[T]) Len() int {
	return v.len
}

func (v Slice[T]) Cap() int {
	return v.cap
}

func (v Slice[T]) Get(idx int) T {
	if idx >= v.len {
		panic(fmt.Errorf("index out of range: idx = %d with len = %d", idx, v.len))
	}

	var tmp T
	addr := uintptr(v.array) + v.offset + unsafe.Sizeof(tmp) * uintptr(idx)

	return *(*T)(unsafe.Pointer(addr))
}

func (v Slice[T]) Set(idx int, value T) {
	if idx >= v.len {
		panic(fmt.Errorf("index out of range: idx = %d with len = %d", idx, v.len))
	}

	addr := uintptr(v.array) + v.offset + unsafe.Sizeof(value) * uintptr(idx)
	*(*T)(unsafe.Pointer(addr)) = value
}

func (v Slice[T]) Slice() []T {
	addr := uintptr(v.array) + v.offset
	s := unsafe.Slice((*T)(unsafe.Pointer(addr)), v.cap)
	return s[:v.len]
}

func (v Slice[T]) Subslice(start, end int) Slice[T] {
	if start < 0 || end < 0 || end < start || start >= v.len || end > v.len {
		panic(fmt.Errorf("subslice out of range: trying [ %d : %d ) on [ %d : %d )", start, end, 0, v.len))
	}

	var tmp T

	return Slice[T]{
		array: v.array,
		len: end - start,
		cap: v.cap - start,
		offset: uintptr(start) * unsafe.Sizeof(tmp),
		subslice: true,
	}
}

func (v *Slice[T]) Append(elems ...T) {
	oldLen := v.len
	newLen := oldLen + len(elems)
	if newLen > v.cap {
		v.growslice(newLen)
	}

	copy(v.Slice()[oldLen:], elems)
}

func (v *Slice[T]) Free() {
	Free(v.array)
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

	oldCap := v.cap
	newCap := nextslicecap(newLen, oldCap)
	
	capmem := newCap * int(size)

	if !v.subslice {
		v.array = Realloc(v.array, capmem)
		v.len = newLen
		v.cap = newCap

		return
	}
	
	old := *v
	*v = Slice[T]{
		array: Malloc(capmem),
		len: newLen,
		cap: newCap,
		offset: 0,
		subslice: false,
	}

	copy(v.Slice(), old.Slice())
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
