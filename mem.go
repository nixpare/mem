package mem

import "unsafe"

func New[T any](allocStrategy func(sizeof uintptr) unsafe.Pointer) *T {
	var tmp T
	ptr := allocStrategy(unsafe.Sizeof(tmp))
	return (*T)(unsafe.Pointer(ptr))
}

func ObjPointer[T any](obj *T) unsafe.Pointer {
	return unsafe.Pointer(obj)
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

func Malloc(n uintptr) unsafe.Pointer {
	return stdlibMalloc(n)
}

func MallocZero(n uintptr) unsafe.Pointer {
	return Calloc(1, n)
}

func MallocN(n int, sizeof uintptr) unsafe.Pointer {
	return Malloc(uintptr(n) * sizeof)
}

func Calloc(n int, sizeof uintptr) unsafe.Pointer {
	return stdlibCalloc(n, sizeof)
}

func Realloc(p unsafe.Pointer, newSize uintptr) unsafe.Pointer {
	return stdlibRealloc(p, newSize)
}

func Free(p unsafe.Pointer) {
	stdlibFree(p)
}
