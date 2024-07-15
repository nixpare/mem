//go:build !windows
package mem

import "unsafe"
//#include <stdlib.h>
import "C"

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