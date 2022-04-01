package forGraphBLASGo

type (
	emptyScalar[T any] struct{}
	fullScalar[T any]  struct{ value T }
)

func newEmptyScalar[T any]() emptyScalar[T] {
	return emptyScalar[T]{}
}

func (s emptyScalar[T]) extractElement(_ *scalarReference[T]) (result T, ok bool) {
	return
}

func (_ emptyScalar[T]) optimized() bool {
	return true
}

func (_ emptyScalar[T]) valid() bool {
	return false
}

func (s emptyScalar[T]) optimize() functionalScalar[T] {
	return s
}

func newFullScalar[T any](value T) fullScalar[T] {
	return fullScalar[T]{value}
}

func (s fullScalar[T]) extractElement(_ *scalarReference[T]) (result T, ok bool) {
	return s.value, true
}

func (_ fullScalar[T]) optimized() bool {
	return true
}

func (_ fullScalar[T]) valid() bool {
	return true
}

func (s fullScalar[T]) optimize() functionalScalar[T] {
	return s
}
