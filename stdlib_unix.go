//go:build !windows
package mem

import "unsafe"
//#include <stdlib.h>
import "C"

func stdlibMalloc(n uintptr) unsafe.Pointer {
	return C.malloc(C.ulong(n))
}

func stdlibCalloc(n int, sizeof uintptr) unsafe.Pointer {
	return C.calloc(C.ulong(n), C.ulong(sizeof))
}

func stdlibRealloc(p unsafe.Pointer, oldSize, newSize uintptr) unsafe.Pointer {
	return C.realloc(p, C.ulong(newSize))
}

func stdlibFree(p unsafe.Pointer) {
	C.free(p)
}