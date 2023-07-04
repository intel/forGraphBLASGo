package GrB

// PrintLevel is a SuiteSparse:GraphBLAS extension.
type PrintLevel int

// SuiteSparse:GraphBLAS extensions
const (
	// Silent means that nothing is printed, just check the object.
	Silent PrintLevel = iota

	// Summary means that a terse summary is printed.
	Summary

	// Short means that a short description is printed, with about 30 entries.
	Short

	// Completely means that the entire contents of the object is printed.
	Completely

	// ShortVerbose is like [Short], but with more precision for floating point numbers.
	ShortVerbose

	// CompletelyVerbose is like [Completely], but with more precision for floating point numbers.
	CompletelyVerbose
)

func (printLevel PrintLevel) String() string {
	switch printLevel {
	case Silent:
		return "silent"
	case Summary:
		return "summary"
	case Short:
		return "short"
	case Completely:
		return "complete"
	case ShortVerbose:
		return "short verbose"
	case CompletelyVerbose:
		return "complete verbose"
	}
	panic("invalid print level")
}
