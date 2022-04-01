package forGraphBLASGo

type functionalScalar[T any] interface {
	extractElement(ref *scalarReference[T]) (T, bool)
	valid() bool
	optimized() bool
	optimize() functionalScalar[T]
}
