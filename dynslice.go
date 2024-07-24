package mem

import (
	"fmt"
	"unsafe"
)

type DynSlice[T any] struct {
	chunk  int
	len    int
	offset int
	data   Slice[Slice[T]]
	ownedData Slice[unsafe.Pointer]
	alloc  AllocNStrategy
	free   FreeStrategy
}

func NewDynSlice[T any](chunk int, len int, cap int, free FreeStrategy, alloc AllocNStrategy) *DynSlice[T] {
	a := New[DynSlice[T]](func(sizeof, alignof uintptr) unsafe.Pointer {
		return alloc(1, sizeof, alignof)
	})
	*a = DynSlice[T]{
		chunk: chunk,
		ownedData: NewSlice[unsafe.Pointer](1, 1, alloc),
		alloc: alloc,
		free:  free,
	}

	a.setDimentions(len, cap)
	return a
}

func (a *DynSlice[T]) getIndexes(i int) (int, int) {
	i += a.offset
	chunkIdx, idx := i/a.chunk, i%a.chunk
	return chunkIdx, idx
}

func (a *DynSlice[T]) Get(i int) T {
	if i >= a.len {
		panic(fmt.Sprintf("index out of range [%d] with length %d", i, a.len))
	}

	chunkIdx, idx := a.getIndexes(i)
	return a.data[chunkIdx][idx]
}

func (a *DynSlice[T]) Set(i int, value T) {
	if i >= a.len {
		panic(fmt.Sprintf("index out of range [%d] with length %d", i, a.len))
	}

	chunkIdx, idx := a.getIndexes(i)
	a.data[chunkIdx][idx] = value
}

func (a *DynSlice[T]) Append(values ...T) {
	oldLen := a.len
	a.EnsureLen(oldLen + len(values))

	for i, value := range values {
		a.Set(oldLen+i, value)
	}
}

func (a *DynSlice[T]) Len() int {
	return a.len
}

func (a *DynSlice[T]) Cap() int {
	return len(a.data)*a.chunk - a.offset
}

func (a *DynSlice[T]) SetLen(len int) {
	a.setDimentions(len, max(a.Cap(), len))
}

func (a *DynSlice[T]) SetCap(cap int) {
	a.setDimentions(min(a.len, cap), cap)
}

func (a *DynSlice[T]) setDimentions(newLen int, newCap int) {
	a.len = newLen
	oldCap := a.Cap()

	if newCap > oldCap {
		if mod := newCap % a.chunk; mod != 0 {
			newCap += a.chunk - mod
		}

		newChunkSize := newCap - oldCap
		newChunk := NewSlice[T](newChunkSize, newChunkSize, a.alloc)
		a.ownedData.Append(a.free, a.alloc, newChunk.Pointer())

		dataPtr := a.data.Pointer()
		for i := 0; i < newChunkSize; i += a.chunk {
			if a.data.Append(a.free, a.alloc, newChunk[i:i+a.chunk]) {
				dataPtr = a.data.Pointer()
			}
		}
		a.ownedData[0] = dataPtr
	} else if newCap < oldCap {
		for ; oldCap-a.chunk >= newCap; oldCap -= a.chunk {
			lastChunkPtr := a.data[len(a.data)-1].Pointer()
			
			if a.ownedData[len(a.ownedData)-1] == lastChunkPtr {
				a.free(lastChunkPtr)
				a.ownedData = a.ownedData[:len(a.ownedData)-1]
			}

			a.data = a.data[:len(a.data)-1]
		}
	}
}

func (a *DynSlice[T]) EnsureLen(len int) {
	if len > a.len {
		a.SetLen(len)
	}
}

func (a *DynSlice[T]) EnsureCap(cap int) {
	if cap > a.Cap() {
		a.SetCap(cap)
	}
}

func (a *DynSlice[T]) Subslice(start int, end int) *DynSlice[T] {
	cap := a.Cap()

	if start > end {
		panic(fmt.Sprintf("start index %d is greater than end index %d", start, end))
	}

	if start < 0 || start >= cap {
		panic(fmt.Sprintf("start index out of range [%d] with capacity %d", start, cap))
	}

	if end < 0 || end > cap {
		panic(fmt.Sprintf("end index out of range [%d] with capacity %d", end, cap))
	}

	slice := New[DynSlice[T]](func(sizeof, alignof uintptr) unsafe.Pointer {
		return a.alloc(1, sizeof, alignof)
	})
	*slice = DynSlice[T]{
		chunk:  a.chunk,
		len:    end - start,
		offset: a.offset + start,
		data:   a.data,
		ownedData: NewSlice[unsafe.Pointer](1, 1, a.alloc),
		alloc:  a.alloc,
		free:   a.free,
	}

	return slice
}

func NewDynSliceFromGoSlice[S ~[]E, E any](v S, free FreeStrategy, alloc AllocNStrategy) *DynSlice[E] {
	chunk := cap(v) / 2
	if chunk < 32 {
		chunk = cap(v)
	}

	a := NewDynSlice[E](chunk, len(v), cap(v), free, alloc)
	for i, x := range v {
		a.Set(i, x)
	}
	return a
}

func (a *DynSlice[T]) Free() {
	for _, x := range a.ownedData[1:] {
		a.free(x)
	}

	if a.data.Pointer() == a.ownedData[0] {
		a.free(a.data.Pointer())
	}
	a.free(a.ownedData.Pointer())

	a.free(ObjectPointer(a))
}
