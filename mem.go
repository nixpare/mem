package mem

import "unsafe"

var (
	Malloc func(n uintptr) unsafe.Pointer = stdlibMalloc
	Calloc func(n int, sizeof uintptr) unsafe.Pointer = stdlibCalloc
	Realloc func(p unsafe.Pointer, n uintptr) unsafe.Pointer = stdlibRealloc
	Free func(p unsafe.Pointer) = stdlibFree
)

type HeapObj[T any] struct {
	ptr unsafe.Pointer
}

func New[T any]() HeapObj[T] {
	var tmp T
	ptr := Calloc(1, unsafe.Sizeof(tmp))
	return HeapObj[T]{ ptr: ptr }
}

func (o HeapObj[T]) Value() *T {
	return (*T)(unsafe.Pointer(o.ptr))
}

func (o HeapObj[T]) Clone() HeapObj[T] {
	value := *o.Value()
	ptr := Malloc(unsafe.Sizeof(value))

	return HeapObj[T]{ ptr: ptr }
}

func (o HeapObj[T]) Free() {
	Free(unsafe.Pointer(o.ptr))
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
