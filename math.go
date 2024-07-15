package mem

const (
	MaxUintptr = ^uintptr(0)
	PtrSize = 4 << (MaxUintptr >> 63)
)

func mulUintptr(a uintptr, b uintptr) (uintptr, bool) {
	if a|b < 1<<(4*PtrSize) || a == 0 {
		return a * b, false
	}
	overflow := b > MaxUintptr/a
	return a * b, overflow
}