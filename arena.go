package mem

import (
	"sync"
	"unsafe"
)

type region struct {
	start unsafe.Pointer
	end   uintptr
	next  uintptr
}

func newRegion(size uintptr, alloc AllocStrategy) region {
	p := alloc(size, 1)
	if p == nil {
		panic("arena: create new memory block failed")
	}

	ptr := uintptr(p)
	return region{
		start: p,
		end:   ptr + size,
		next:  ptr,
	}
}

func (r *region) allocate(n int, sizeof uintptr, alignof uintptr) uintptr {
	aligned := (r.next + (alignof - 1)) & ^(alignof - 1)

	newNext := aligned + uintptr(n)*sizeof
	if newNext > r.end {
		return 0
	}

	r.next = newNext
	return aligned
}

func (r region) free(free FreeStrategy) {
	free(r.start)
}

type Arena struct {
	regions   Slice[region]
	allocSize uintptr
	alloc     AllocStrategy
	free	  FreeStrategy
	m         sync.Mutex
}

func NewArena(allocSize uintptr, free FreeStrategy,alloc AllocStrategy) *Arena {
	a := New[Arena](alloc)
	*a = Arena{
		allocSize: allocSize,
		alloc:     alloc,
		free:      free,
	}

	return a
}

func (a *Arena) Alloc(sizeof, alignof uintptr) unsafe.Pointer {
	return a.AllocN(1, sizeof, alignof)
}

func (a *Arena) AllocN(n int, sizeof, alignof uintptr) unsafe.Pointer {
	a.m.Lock()
	defer a.m.Unlock()

	for i := range a.regions {
		ptr := a.regions[i].allocate(n, sizeof, alignof)
		if ptr != 0 {
			return unsafe.Pointer(ptr)
		}
	}

	a.allocRegion(n, sizeof)

	ptr := a.regions[len(a.regions)-1].allocate(n, sizeof, alignof)
	if ptr == 0 {
		panic("arena: allocation of new memory block failed")
	}

	return unsafe.Pointer(ptr)
}

func (a *Arena) Regions() int {
	return len(a.regions)
}

func (a *Arena) allocRegion(n int, sizeof uintptr) {
	memSize, reqSize := a.allocSize, uintptr(n)*sizeof
	if memSize < reqSize {
		memSize = reqSize
	}

	a.regions.Append(a.free, a.allocNAdapter, newRegion(memSize, a.alloc))
}

func (a *Arena) allocNAdapter(n int, sizeof, alignof uintptr) unsafe.Pointer {
	return a.alloc(uintptr(n)*sizeof, alignof)
}

func (a *Arena) Reset() {
	a.m.Lock()
	defer a.m.Unlock()

	for i := range a.regions {
		a.regions[i].next = uintptr(a.regions[i].start)
	}
}

func (a *Arena) Free() {
	a.m.Lock()

	for _, r := range a.regions {
		r.free(a.free)
	}
	FreeSlice(&a.regions, a.free)

	a.m.Unlock()
	FreeObject(a, a.free)
}
