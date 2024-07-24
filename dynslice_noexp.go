//go:build !goexperiment.rangefunc

package mem

import (
	"fmt"
	"strings"
)

func (a *DynSlice[T]) Iter() (<-chan func() (int, T), func()) {
	yieldCh := make(chan func() (int, T))
	breakCh := make(chan struct{})

	go func() {
		defer close(yieldCh)

		for i := 0; i < a.Len(); i++ {
			x := a.Get(i)

			select {
			case <-breakCh:
				return
			case yieldCh <- func() (int, T) { return i, x }:
			}
		}
	}()

	return yieldCh, func() { close(breakCh) }
}

func (a *DynSlice[T]) Clone() *DynSlice[T] {
	clone := NewDynSlice[T](a.chunk, a.Len(), a.Cap(), a.free, a.alloc)

	yield, done := a.Iter()
	defer done()

	for f := range yield {
		i, x := f()
		clone.Set(i, x)
	}
	return clone
}

func (a *DynSlice[T]) ToSlice() Slice[T] {
	slice := NewSlice[T](a.Len(), a.Cap(), a.alloc)
	
	yield, done := a.Iter()
	defer done()

	for f := range yield {
		i, x := f()
		slice[i] = x
	}
	
	return slice
}

func (a *DynSlice[T]) ToGoSlice() []T {
	slice := make([]T, a.Len(), a.Cap())

	yield, done := a.Iter()
	defer done()

	for f := range yield {
		i, x := f()
		slice[i] = x
	}

	return slice
}

func (a *DynSlice[T]) String() string {
	var sb strings.Builder
	sb.WriteRune('[')

	yield, done := a.Iter()
	defer done()

	for f := range yield {
		_, x := f()
		sb.WriteString(fmt.Sprint(x))
		sb.WriteRune(' ')
	}

	b := []byte(sb.String())
	b[len(b)-1] = ']'
	return string(b)
}
