package GrB

// #include "GraphBLAS.h"
import "C"
import "fmt"

// Info is the default error type that is returned by all GraphBLAS functions.
// There are two types of return codes: API error and execution error.
//
// API errors are returned to the caller of the corresponding GraphBLAS function,
// and need to be handled like other Go errors.
//
// Execution errors cause a panic, and can be handled by recover.
//
// The GraphBLAS C API specification also specifies informational codes as
// a third type of return code (GrB_SUCCESS and GrB_NO_VALUE). forGraphBLASGo
// functions do not return these informational codes, but instead return the nil
// error value to indicate successful execution of the function. In cases where it
// is important to have further information about the successful outcome,
// forGraphBLASGo functions return additional return values.
// (See for example [Matrix.ExtractElement].)
//
// This also applies to GxB_EXHAUSTED, which is an additional informational return
// code added by SuiteSparse:GraphBLAS as an extension: It is never returned by
// forGraphBLASGo, but instead translated into additional return values, where
// necessary. (GxB_EXHAUSTED is returned by some of the iterator functions in the
// C implementation.)
//
// In [Blocking] mode, when an operation returns nil, it completed successfully.
// In [NonBlocking] mode, when an operation returns nil, this indicates that basic checks for
// input arguments passed successfully, like the compatibility tests on dimensions and domains
// for scalars, vectors and matrices. Either way, result scalars, vectors, or matrices are ready
// to be used in subsequent function calls.
type Info int

// GraphBLAS informational return codes
var (
	success     = Info(C.GrB_SUCCESS)
	noValue     = Info(C.GrB_NO_VALUE)
	isExhausted = Info(C.GxB_EXHAUSTED) // is a SuiteSparse:GraphBLAS extension.
)

// GraphBLAS API errors that may be returned by GraphBLAS operations.
var (
	// UninitializedObject indicates that a GraphBLAS object is passed to a function before it was properly
	// initialized by a call to [BinaryOpNew], [ContextNew], [DescriptorNew], [IndexUnaryOpNew],
	// [Matrix.ColIteratorNew], [Matrix.Dup], [Matrix.IteratorNew], [Matrix.ReshapeDup], [Matrix.RowIteratorNew],
	// [MatrixDeserialize], [MatrixImport], [MatrixNew], [MonoidNew], [MonoidTerminalNew], [NamedBinaryOpNew],
	// [NamedIndexUnaryOpNew], [NamedTypeNew], [NamedUnaryOpNew], [Scalar.Dup], [ScalarNew], [SemiringNew], [TypeNew],
	// [UnaryOpNew], [Vector.Diag], [Vector.Dup], [Vector.IteratorNew], or [VectorNew].
	UninitializedObject = Info(C.GrB_UNINITIALIZED_OBJECT)

	// NullPointer indicates that a nil is passed for a pointer parameter.
	NullPointer = Info(C.GrB_NULL_POINTER)

	// InvalidValue indicates that an invalid value is passed as a parameter (for example, a value that is not
	// one of the predefined values of an enumeration type, or an index that is <= 0 or > [IndexMax]).
	InvalidValue = Info(C.GrB_INVALID_VALUE)

	// InvalidIndex indicates that one or more (single) index is passed to a function that is outside of
	// the dimensions of the corresponding vector or matrix.
	InvalidIndex = Info(C.GrB_INVALID_INDEX)

	// DomainMismatch indicates that the domains of the various input and/or output scalars, vectors and matrices
	// are incompatible with each other and/or with the provided accumulation or other operators; or that
	// the mask's domain is not compatible with bool (in the case where the [Structure] descriptor is not
	// set for the mask).
	DomainMismatch = Info(C.GrB_DOMAIN_MISMATCH)

	// DimensionMismatch indicates that the dimensions of input and/or output vectors, matrices or masks
	// are incompatible.
	DimensionMismatch = Info(C.GrB_DIMENSION_MISMATCH)

	// OutputNotEmpty indicates that the output vector or matrix already contains valid tuples (elements).
	// In other words, [Vector.Nvals] or [Matrix.Nvals] returns a positive value.
	OutputNotEmpty = Info(C.GrB_OUTPUT_NOT_EMPTY)

	// NotImplemented indicates that an attempt was made to call a GraphBLAS function for a combination of
	// input parameters that is not supported by a particular implementation.
	NotImplemented = Info(C.GrB_NOT_IMPLEMENTED)

	// SliceMismatch indicates that the lengths of different input slices that should be the same do not match.
	// SliceMismatch is a forGraphBLASGo extension.
	SliceMismatch = Info(-201)
)

// GraphBLAS execution errors that may cause a panic.
var (
	// Panic indicates an unknown internal error.
	Panic = Info(C.GrB_PANIC)

	// OutOfMemory indicates that not enough memory is available for an operation.
	OutOfMemory = Info(C.GrB_OUT_OF_MEMORY)

	// InsufficientSpace indicates that the slice provided is not large enough to hold output.
	InsufficientSpace = Info(C.GrB_INSUFFICIENT_SPACE)

	// InvalidObject indicates that one of the opaque GraphBLAS objects (input or output) is in an invalid state
	// caused by a previous execution error.  Call the Err methods on these GraphBLAS objects to access any error
	// messages generated by the implementation.
	InvalidObject = Info(C.GrB_INVALID_OBJECT)

	// IndexOutOfBounds indicates that a value in one or more of the input index slices is less than 0
	// or greater than the corresponding dimension of the output vector or matrix. In non-blocking mode,
	// this may be reported by call to a Wait method.
	IndexOutOfBounds = Info(C.GrB_INDEX_OUT_OF_BOUNDS)

	// EmptyObject indicates that the [Scalar] object used in the call is empty (nvals = 0) and therefore
	// a value cannot be passed to the provided operator.
	EmptyObject = Info(C.GrB_EMPTY_OBJECT)
)

var errorStrings = map[Info]string{
	success:             "success",
	noValue:             "no value",
	isExhausted:         "iterator exhausted",
	UninitializedObject: "uninitialized object",
	NullPointer:         "null pointer",
	InvalidValue:        "invalid value",
	InvalidIndex:        "invalid index",
	DomainMismatch:      "domain mismatch",
	DimensionMismatch:   "dimension mismatch",
	OutputNotEmpty:      "output not empty",
	NotImplemented:      "not implemented",
	SliceMismatch:       "slice mismatch",
	Panic:               "panic",
	OutOfMemory:         "out of memory",
	InsufficientSpace:   "insufficient space",
	InvalidObject:       "invalid object",
	IndexOutOfBounds:    "index out of bounds",
	EmptyObject:         "empty object",
}

func (info Info) String() string {
	if errorString, ok := errorStrings[info]; ok {
		return errorString
	}
	return fmt.Sprintf("unknown error code: %v", int(info))
}

func (info Info) Error() string {
	if errorString, ok := errorStrings[info]; ok {
		if isInformational(info) {
			return fmt.Sprintf("GraphBLAS information return code: %v", errorString)
		}
		if isExecutionError(info) {
			return fmt.Sprintf("GraphBLAS execution error: %v", errorString)
		}
		return fmt.Sprintf("GraphBLAS API error: %v", errorString)
	}
	return fmt.Sprintf("GraphBLAS API error: unknown error code %v", int(info))
}

func isInformational(info Info) bool {
	switch info {
	case success, noValue, isExhausted:
		return true
	default:
		return false
	}
}

func isExecutionError(info Info) bool {
	switch info {
	case Panic, OutOfMemory, InsufficientSpace, InvalidObject, IndexOutOfBounds, EmptyObject:
		return true
	default:
		return false
	}
}

var panicOnError = false

// GlobalSetPanicOnError changes how GraphBLAS API errors are returned:
//   - if onNotOff is true, then forGraphBLASGo functions panic for these
//     errors when they occur, rather than return them.
//   - if onNotOff is false (the default), then forGraphBLASGo functions
//     return these errors when they occur.
//
// The setting does not influence how GraphBLAS execution errors are reported;
// forGraphBLASGo functions always panic when they occur.
//
// GlobalSetPanicOnError is a forGraphBLASGo extension.
func GlobalSetPanicOnError(onNotOff bool) error {
	panicOnError = onNotOff
	return nil
}

func makeError(info Info) error {
	if isInformational(info) {
		panic(fmt.Errorf("informational return code %w must not be returned by forGraphBLASGo - this should not happen", info))
	}
	if panicOnError || isExecutionError(info) {
		panic(info.Error())
	}
	return info
}

type nopanic struct {
	wrapped error
}

// CheckErrors recovers a potential error handled by [OK] and assigns it to
// *err, unless *err != nil. Use CheckError with defer.
//
// [OK] panics on errors != nil, but CheckError only recovers
// errors handled by [OK]. CheckError will panic again on any
// other panics.
//
// Example:
//
//	func Example(...) (err error) {
//	   defer GrB.CheckError(&err)
//
//	   ...
//
//	   GrB.OK(GrB.SomeOperation(...))
//
//	   ...
//	}
//
// CheckErrors is a forGraphBLASGo extension.
func CheckErrors(err *error) {
	x := recover()
	if x == nil {
		return
	}
	if np, ok := x.(nopanic); ok {
		if *err == nil {
			*err = np.wrapped
		}
	} else {
		panic(x)
	}
}

// OK transfers control to the closest deferred [CheckErrors],
// if err != nil (by panicking using an internal marker type).
// If err == nil, OK just returns without doing anything else.
//
// This mechanism works on any error types, not only [Info].
//
// OK is an optional mechanism. It is also always correct to
// check errors manually, even when deferred [CheckErrors]
// calls are active.
//
// See [CheckErrors] for an example.
//
// OK is a forGraphBLASGo extension.
func OK(err error) {
	if err == nil {
		return
	}
	panic(nopanic{err})
}
