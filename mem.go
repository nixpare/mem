package mem

import "unsafe"

// #include <stdlib.h>
import "C"

var (
	Malloc func(n int) unsafe.Pointer = stdlibMalloc
	Calloc func(n int, sizeof uintptr) unsafe.Pointer = stdlibCalloc
	Realloc func(p unsafe.Pointer, n int) unsafe.Pointer = stdlibRealloc
	Free func(p unsafe.Pointer) = stdlibFree
)

func New[T any]() *T {
	var obj T
	ptr := Calloc(1, unsafe.Sizeof(obj))
	return (*T)(ptr)
}

func Dealloc[T any](obj *T) {
	Free(unsafe.Pointer(obj))
}

// GetRef allows, through the use of unsafe Go, to reference a variable
// allocated on the stack without the possibility that the compiler decides
// to move that variable automatically on the heap.
//
//go:inline
func GetRef[T any](p *T) *T {
	addr := uintptr(unsafe.Pointer(p))
	return (*T)(unsafe.Pointer(addr))
}

func stdlibMalloc(n int) unsafe.Pointer {
	return C.malloc(C.ulong(n))
}

func stdlibCalloc(n int, sizeof uintptr) unsafe.Pointer {
	return C.calloc(C.ulong(n), C.ulong(sizeof))
}

func stdlibRealloc(p unsafe.Pointer, n int) unsafe.Pointer {
	return C.realloc(p, C.ulong(n))
}

func stdlibFree(p unsafe.Pointer) {
	C.free(p)
}
