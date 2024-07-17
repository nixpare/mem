package mem

import "unsafe"
//#include <stdlib.h>
import "C"

func stdlibMalloc(n uintptr) unsafe.Pointer {
	return C.malloc(C.ulonglong(n))
}

func stdlibCalloc(n int, sizeof uintptr) unsafe.Pointer {
	return C.calloc(C.ulonglong(n), C.ulonglong(sizeof))
}

func stdlibRealloc(p unsafe.Pointer, newSize uintptr) unsafe.Pointer {
	return C.realloc(p, C.ulonglong(newSize))
}

func stdlibFree(p unsafe.Pointer) {
	C.free(p)
}
