package mem

import "unsafe"

type AllocStrategy func(sizeof, alignof uintptr) unsafe.Pointer
type AllocNStrategy func(n int, sizeof, alignof uintptr) unsafe.Pointer
type FreeStrategy func(p unsafe.Pointer)

func New[T any](alloc AllocStrategy) *T {
	var tmp T
	p := alloc(unsafe.Sizeof(tmp), unsafe.Alignof(tmp))
	if p == nil {
		panic("new: allocation failed")
	}
	return (*T)(unsafe.Pointer(p))
}

func ObjectPointer[T any](obj *T) unsafe.Pointer {
	return unsafe.Pointer(obj)
}

func FreeObject[T any](obj *T, free FreeStrategy) {
	free(ObjectPointer(obj))
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
