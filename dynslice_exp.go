//go:build goexperiment.rangefunc

package mem

import (
	"fmt"
	"iter"
	"strings"
)

func (a *DynSlice[T]) Iter() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {

		firstChunk, firstIdx := a.getIndexes(0)
		lastChunk, lastIdx := a.getIndexes(a.len - 1)
		var i int

		if firstChunk == lastChunk {
			for idx := firstIdx; idx < a.chunk && i < a.len; idx++ {
				if !yield(i, a.data[firstChunk][idx]) {
					return
				}
				i++
			}

			return
		}

		for idx := firstIdx; idx < a.chunk; idx++ {
			if !yield(i, a.data[firstChunk][idx]) {
				return
			}
			i++
		}

		for chunkIdx := firstChunk + 1; chunkIdx < lastChunk; chunkIdx++ {
			for _, x := range a.data[chunkIdx] {
				if !yield(i, x) {
					return
				}
				i++
			}
		}

		for idx := 0; idx <= lastIdx; idx++ {
			if !yield(i, a.data[lastChunk][idx]) {
				return
			}
			i++
		}

	}
}

func (a *DynSlice[T]) Clone() *DynSlice[T] {
	clone := NewDynSlice[T](a.chunk, a.Len(), a.Cap(), a.free, a.alloc)
	for i, x := range a.Iter() {
		clone.Set(i, x)
	}
	return clone
}

func (a *DynSlice[T]) ToSlice() Slice[T] {
	slice := NewSlice[T](a.Len(), a.Cap(), a.alloc)
	for i, x := range a.Iter() {
		slice[i] = x
	}
	return slice
}

func (a *DynSlice[T]) ToGoSlice() []T {
	slice := make([]T, a.Len(), a.Cap())
	for i, x := range a.Iter() {
		slice[i] = x
	}
	return slice
}

func (a *DynSlice[T]) String() string {
	var sb strings.Builder
	sb.WriteRune('[')

	for _, x := range a.Iter() {
		sb.WriteString(fmt.Sprint(x))
		sb.WriteRune(' ')
	}

	b := []byte(sb.String())
	b[len(b)-1] = ']'
	return string(b)
}
