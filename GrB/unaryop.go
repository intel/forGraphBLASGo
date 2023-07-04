package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

type (
	// UnaryFunction is the C type for unary functions:
	//    typedef void (*GxB_unary_function)  (void *, const void *) ;
	UnaryFunction C.GxB_unary_function

	// UnaryOp represents a GraphBLAS function that takes one argument of type Din,
	// and returns an argument of type Dout.
	UnaryOp[Dout, Din any] struct {
		grb C.GrB_UnaryOp
	}
)

// UnaryOpNew returns a new GraphBLAS unary operator with a specified user-defined function in C and its types (domains).
//
// Parameters:
//   - unaryFunc (IN): A pointer to a user-defined function in C that takes an input parameter of type Din and returns
//     a value of type Dout, all passed as void pointers. Dout and Din should be one of the [Predefined] GraphBLAS types, one
//     of the [Complex] GraphBLAS types, or a user-defined GraphBLAS type.
//
// GraphBLAS API errors that may be returned:
//   - [NullPointer], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func UnaryOpNew[Dout, Din any](unaryFunc UnaryFunction) (unaryOp UnaryOp[Dout, Din], err error) {
	var dout Dout
	doutt, ok := grbType[TypeOf(dout)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	var din Din
	dint, ok := grbType[TypeOf(din)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	info := Info(C.GrB_UnaryOp_new(&unaryOp.grb, unaryFunc, doutt, dint))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// NamedUnaryOpNew creates a named unary function. It is like [UnaryOpNew], except:
//   - unopname is the name for the GraphBLAS unary operator. Only the first 127 characters are used.
//   - unopdefn is a string containing the entire function itself.
//
// The two strings unopname and unopdefn are optional, but are required to enable the JIT compilation
// of kernels that use this operator.
//
// If the JIT is enabled, or if the corresponding JIT kernel has been copied into the PreJIT folder,
// the function may be nil. In this case, a JIT kernel is compiled that contains just the user-defined
// function. If the JIT is disabled and the function is nil, this method returns a [NullPointer] error.
//
// NamedUnaryOpNew is a SuiteSparse:GraphBLAS extension.
func NamedUnaryOpNew[Dout, Din any](unaryFunc UnaryFunction, unopname string, unopdefn string) (unaryOp UnaryOp[Dout, Din], err error) {
	var dout Dout
	doutt, ok := grbType[TypeOf(dout)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	var din Din
	dint, ok := grbType[TypeOf(din)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	cunopname := C.CString(unopname)
	defer C.free(unsafe.Pointer(cunopname))
	cunopdefn := C.CString(unopdefn)
	defer C.free(unsafe.Pointer(cunopdefn))
	info := Info(C.GxB_UnaryOp_new(&unaryOp.grb, unaryFunc, doutt, dint, cunopname, cunopdefn))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Valid returns true if unaryOp has been created by a successful call to [UnaryOpNew] or [NamedUnaryOpNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (unaryOp UnaryOp[Dout, Din]) Valid() bool {
	return unaryOp.grb != C.GrB_UnaryOp(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [UnaryOp] and releases any resources associated with
// it. Calling Free on an object that is not [UnaryOp.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined unary operator is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (unaryOp *UnaryOp[Dout, Din]) Free() error {
	info := Info(C.GrB_UnaryOp_free(&unaryOp.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the unary operator into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (unaryOp UnaryOp[Dout, Din]) Wait(mode WaitMode) error {
	info := Info(C.GrB_UnaryOp_wait(unaryOp.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the unary operator.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (unaryOp UnaryOp[Dout, Din]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_UnaryOp_error(&cerror, unaryOp.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the binary operator to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: unaryOp is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (unaryOp UnaryOp[Dout, Din]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_UnaryOp_fprint(unaryOp.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// One is f(x) = 1
//
// One is a SuiteSparse:GraphBLAS extension.
func One[D Predefined | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ONE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ONE_INT32
		} else {
			f.grb = C.GxB_ONE_INT64
		}
	case int8:
		f.grb = C.GxB_ONE_INT8
	case int16:
		f.grb = C.GxB_ONE_INT16
	case int32:
		f.grb = C.GxB_ONE_INT32
	case int64:
		f.grb = C.GxB_ONE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ONE_UINT32
		} else {
			f.grb = C.GxB_ONE_UINT64
		}
	case uint8:
		f.grb = C.GxB_ONE_UINT8
	case uint16:
		f.grb = C.GxB_ONE_UINT16
	case uint32:
		f.grb = C.GxB_ONE_UINT32
	case uint64:
		f.grb = C.GxB_ONE_UINT64
	case float32:
		f.grb = C.GxB_ONE_FP32
	case float64:
		f.grb = C.GxB_ONE_FP64
	case complex64:
		f.grb = C.GxB_ONE_FC32
	case complex128:
		f.grb = C.GxB_ONE_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Identity is f(x) = x
func Identity[D Predefined | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_IDENTITY_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_IDENTITY_INT32
		} else {
			f.grb = C.GrB_IDENTITY_INT64
		}
	case int8:
		f.grb = C.GrB_IDENTITY_INT8
	case int16:
		f.grb = C.GrB_IDENTITY_INT16
	case int32:
		f.grb = C.GrB_IDENTITY_INT32
	case int64:
		f.grb = C.GrB_IDENTITY_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_IDENTITY_UINT32
		} else {
			f.grb = C.GrB_IDENTITY_UINT64
		}
	case uint8:
		f.grb = C.GrB_IDENTITY_UINT8
	case uint16:
		f.grb = C.GrB_IDENTITY_UINT16
	case uint32:
		f.grb = C.GrB_IDENTITY_UINT32
	case uint64:
		f.grb = C.GrB_IDENTITY_UINT64
	case float32:
		f.grb = C.GrB_IDENTITY_FP32
	case float64:
		f.grb = C.GrB_IDENTITY_FP64
	case complex64:
		f.grb = C.GxB_IDENTITY_FC32
	case complex128:
		f.grb = C.GxB_IDENTITY_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Ainv is f(x) = -x
func Ainv[D Predefined | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_AINV_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_AINV_INT32
		} else {
			f.grb = C.GrB_AINV_INT64
		}
	case int8:
		f.grb = C.GrB_AINV_INT8
	case int16:
		f.grb = C.GrB_AINV_INT16
	case int32:
		f.grb = C.GrB_AINV_INT32
	case int64:
		f.grb = C.GrB_AINV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_AINV_UINT32
		} else {
			f.grb = C.GrB_AINV_UINT64
		}
	case uint8:
		f.grb = C.GrB_AINV_UINT8
	case uint16:
		f.grb = C.GrB_AINV_UINT16
	case uint32:
		f.grb = C.GrB_AINV_UINT32
	case uint64:
		f.grb = C.GrB_AINV_UINT64
	case float32:
		f.grb = C.GrB_AINV_FP32
	case float64:
		f.grb = C.GrB_AINV_FP64
	case complex64:
		f.grb = C.GxB_AINV_FC32
	case complex128:
		f.grb = C.GxB_AINV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Minv is f(x) = 1/x
func Minv[D Predefined | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_MINV_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MINV_INT32
		} else {
			f.grb = C.GrB_MINV_INT64
		}
	case int8:
		f.grb = C.GrB_MINV_INT8
	case int16:
		f.grb = C.GrB_MINV_INT16
	case int32:
		f.grb = C.GrB_MINV_INT32
	case int64:
		f.grb = C.GrB_MINV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MINV_UINT32
		} else {
			f.grb = C.GrB_MINV_UINT64
		}
	case uint8:
		f.grb = C.GrB_MINV_UINT8
	case uint16:
		f.grb = C.GrB_MINV_UINT16
	case uint32:
		f.grb = C.GrB_MINV_UINT32
	case uint64:
		f.grb = C.GrB_MINV_UINT64
	case float32:
		f.grb = C.GrB_MINV_FP32
	case float64:
		f.grb = C.GrB_MINV_FP64
	case complex64:
		f.grb = C.GxB_MINV_FC32
	case complex128:
		f.grb = C.GxB_MINV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Abs is f(x) = |x|
func Abs[D Predefined | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_ABS_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_ABS_INT32
		} else {
			f.grb = C.GrB_ABS_INT64
		}
	case int8:
		f.grb = C.GrB_ABS_INT8
	case int16:
		f.grb = C.GrB_ABS_INT16
	case int32:
		f.grb = C.GrB_ABS_INT32
	case int64:
		f.grb = C.GrB_ABS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_ABS_UINT32
		} else {
			f.grb = C.GrB_ABS_UINT64
		}
	case uint8:
		f.grb = C.GrB_ABS_UINT8
	case uint16:
		f.grb = C.GrB_ABS_UINT16
	case uint32:
		f.grb = C.GrB_ABS_UINT32
	case uint64:
		f.grb = C.GrB_ABS_UINT64
	case float32:
		f.grb = C.GrB_ABS_FP32
	case float64:
		f.grb = C.GrB_ABS_FP64
	case complex64:
		f.grb = C.GxB_ABS_FC32
	case complex128:
		f.grb = C.GxB_ABS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// LnotBool is f(x) = !x
var LnotBool = UnaryOp[bool, bool]{C.GrB_LNOT}

// Lnot is f(x) = !x (C semantics)
//
// Lnot is a SuiteSparse:GraphBLAS extension.
func Lnot[D Predefined]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_LNOT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LNOT_INT32
		} else {
			f.grb = C.GxB_LNOT_INT64
		}
	case int8:
		f.grb = C.GxB_LNOT_INT8
	case int16:
		f.grb = C.GxB_LNOT_INT16
	case int32:
		f.grb = C.GxB_LNOT_INT32
	case int64:
		f.grb = C.GxB_LNOT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LNOT_UINT32
		} else {
			f.grb = C.GxB_LNOT_UINT64
		}
	case uint8:
		f.grb = C.GxB_LNOT_UINT8
	case uint16:
		f.grb = C.GxB_LNOT_UINT16
	case uint32:
		f.grb = C.GxB_LNOT_UINT32
	case uint64:
		f.grb = C.GxB_LNOT_UINT64
	case float32:
		f.grb = C.GxB_LNOT_FP32
	case float64:
		f.grb = C.GxB_LNOT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Bnot is f(x) = ^x
func Bnot[D Integer]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BNOT_INT32
		} else {
			f.grb = C.GrB_BNOT_INT64
		}
	case int8:
		f.grb = C.GrB_BNOT_INT8
	case int16:
		f.grb = C.GrB_BNOT_INT16
	case int32:
		f.grb = C.GrB_BNOT_INT32
	case int64:
		f.grb = C.GrB_BNOT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BNOT_UINT32
		} else {
			f.grb = C.GrB_BNOT_UINT64
		}
	case uint8:
		f.grb = C.GrB_BNOT_UINT8
	case uint16:
		f.grb = C.GrB_BNOT_UINT16
	case uint32:
		f.grb = C.GrB_BNOT_UINT32
	case uint64:
		f.grb = C.GrB_BNOT_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Positioni is f(x) = i (0-based row index)
//
// Positioni is a SuiteSparse:GraphBLAS extension.
func Positioni[D int32 | int64 | int, Din any]() (f UnaryOp[D, Din]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_POSITIONI_INT32
		} else {
			f.grb = C.GxB_POSITIONI_INT64
		}
	case int32:
		f.grb = C.GxB_POSITIONI_INT32
	case int64:
		f.grb = C.GxB_POSITIONI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Positioni1 is f(x) = i (1-based row index)
//
// Positioni1 is a SuiteSparse:GraphBLAS extension.
func Positioni1[D int32 | int64 | int, Din any]() (f UnaryOp[D, Din]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_POSITIONI1_INT32
		} else {
			f.grb = C.GxB_POSITIONI1_INT64
		}
	case int32:
		f.grb = C.GxB_POSITIONI1_INT32
	case int64:
		f.grb = C.GxB_POSITIONI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Positionj is f(x) = j (0-based column index)
//
// Positionj is a SuiteSparse:GraphBLAS extension.
func Positionj[D int32 | int64 | int, Din any]() (f UnaryOp[D, Din]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_POSITIONJ_INT32
		} else {
			f.grb = C.GxB_POSITIONJ_INT64
		}
	case int32:
		f.grb = C.GxB_POSITIONJ_INT32
	case int64:
		f.grb = C.GxB_POSITIONJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Positionj1 is f(x) = j (1-based column index)
//
// Positionj1 is a SuiteSparse:GraphBLAS extension.
func Positionj1[D int32 | int64 | int, Din any]() (f UnaryOp[D, Din]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_POSITIONJ1_INT32
		} else {
			f.grb = C.GxB_POSITIONJ1_INT64
		}
	case int32:
		f.grb = C.GxB_POSITIONJ1_INT32
	case int64:
		f.grb = C.GxB_POSITIONJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Sqrt is f(x) = sqrt(x) (square root)
//
// Sqrt is a SuiteSparse:GraphBLAS extension.
func Sqrt[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_SQRT_FP32
	case float64:
		f.grb = C.GxB_SQRT_FP64
	case complex64:
		f.grb = C.GxB_SQRT_FC32
	case complex128:
		f.grb = C.GxB_SQRT_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Log is f(x) = log(x) (natural logarithm)
//
// Log is a SuiteSparse:GraphBLAS extension.
func Log[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_LOG_FP32
	case float64:
		f.grb = C.GxB_LOG_FP64
	case complex64:
		f.grb = C.GxB_LOG_FC32
	case complex128:
		f.grb = C.GxB_LOG_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Exp is f(x) = exp(x) (natural exponent)
//
// Exp is a SuiteSparse:GraphBLAS extension.
func Exp[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_EXP_FP32
	case float64:
		f.grb = C.GxB_EXP_FP64
	case complex64:
		f.grb = C.GxB_EXP_FC32
	case complex128:
		f.grb = C.GxB_EXP_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Log10 is f(x) = log10(x) (base-10 logarithm)
//
// Log10 is a SuiteSparse:GraphBLAS extension.
func Log10[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_LOG10_FP32
	case float64:
		f.grb = C.GxB_LOG10_FP64
	case complex64:
		f.grb = C.GxB_LOG10_FC32
	case complex128:
		f.grb = C.GxB_LOG10_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Log2 is f(x) = log2(x) (base-2 logarithm)
//
// Log2 is a SuiteSparse:GraphBLAS extension.
func Log2[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_LOG2_FP32
	case float64:
		f.grb = C.GxB_LOG2_FP64
	case complex64:
		f.grb = C.GxB_LOG2_FC32
	case complex128:
		f.grb = C.GxB_LOG2_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Exp2 is f(x) = exp2(x) (base-2 exponent)
//
// Exp2 is a SuiteSparse:GraphBLAS extension.
func Exp2[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_EXP2_FP32
	case float64:
		f.grb = C.GxB_EXP2_FP64
	case complex64:
		f.grb = C.GxB_EXP2_FC32
	case complex128:
		f.grb = C.GxB_EXP2_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Expm1 is f(x) = expm1(x) (natural exponent - 1)
//
// Expm1 is a SuiteSparse:GraphBLAS extension.
func Expm1[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_EXPM1_FP32
	case float64:
		f.grb = C.GxB_EXPM1_FP64
	case complex64:
		f.grb = C.GxB_EXPM1_FC32
	case complex128:
		f.grb = C.GxB_EXPM1_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Log1p is f(x) = log1p(x) (natural logarithm + 1)
//
// Log1p is a SuiteSparse:GraphBLAS extension.
func Log1p[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_LOG1P_FP32
	case float64:
		f.grb = C.GxB_LOG1P_FP64
	case complex64:
		f.grb = C.GxB_LOG1P_FC32
	case complex128:
		f.grb = C.GxB_LOG1P_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Sin is f(x) = sin(x) (sine)
//
// Sin is a SuiteSparse:GraphBLAS extension.
func Sin[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_SIN_FP32
	case float64:
		f.grb = C.GxB_SIN_FP64
	case complex64:
		f.grb = C.GxB_SIN_FC32
	case complex128:
		f.grb = C.GxB_SIN_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Cos is f(x) = cos(x) (cosine)
//
// Cos is a SuiteSparse:GraphBLAS extension.
func Cos[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_COS_FP32
	case float64:
		f.grb = C.GxB_COS_FP64
	case complex64:
		f.grb = C.GxB_COS_FC32
	case complex128:
		f.grb = C.GxB_COS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Tan is f(x) = tan(x) (tangent)
//
// Tan is a SuiteSparse:GraphBLAS extension.
func Tan[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_TAN_FP32
	case float64:
		f.grb = C.GxB_TAN_FP64
	case complex64:
		f.grb = C.GxB_TAN_FC32
	case complex128:
		f.grb = C.GxB_TAN_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Asin is f(x) = asin(x) (inverse sine)
//
// Asin is a SuiteSparse:GraphBLAS extension.
func Asin[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ASIN_FP32
	case float64:
		f.grb = C.GxB_ASIN_FP64
	case complex64:
		f.grb = C.GxB_ASIN_FC32
	case complex128:
		f.grb = C.GxB_ASIN_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Acos is f(x) = acos(x) (inverse cosine)
//
// Acos is a SuiteSparse:GraphBLAS extension.
func Acos[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ACOS_FP32
	case float64:
		f.grb = C.GxB_ACOS_FP64
	case complex64:
		f.grb = C.GxB_ACOS_FC32
	case complex128:
		f.grb = C.GxB_ACOS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Atan is f(x) = atan(x) (inverse tangent)
//
// Atan is a SuiteSparse:GraphBLAS extension.
func Atan[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ATAN_FP32
	case float64:
		f.grb = C.GxB_ATAN_FP64
	case complex64:
		f.grb = C.GxB_ATAN_FC32
	case complex128:
		f.grb = C.GxB_ATAN_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Sinh is f(x) = sinh(x) (hyperbolic sine)
//
// Sinh is a SuiteSparse:GraphBLAS extension.
func Sinh[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_SINH_FP32
	case float64:
		f.grb = C.GxB_SINH_FP64
	case complex64:
		f.grb = C.GxB_SINH_FC32
	case complex128:
		f.grb = C.GxB_SINH_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Cosh is f(x) = cosh(x) (hyperbolic cosine)
//
// Cosh is a SuiteSparse:GraphBLAS extension.
func Cosh[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_COSH_FP32
	case float64:
		f.grb = C.GxB_COSH_FP64
	case complex64:
		f.grb = C.GxB_COSH_FC32
	case complex128:
		f.grb = C.GxB_COSH_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Tanh is f(x) = tanh(x) (hyperbolic tangent)
//
// Tanh is a SuiteSparse:GraphBLAS extension.
func Tanh[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_TANH_FP32
	case float64:
		f.grb = C.GxB_TANH_FP64
	case complex64:
		f.grb = C.GxB_TANH_FC32
	case complex128:
		f.grb = C.GxB_TANH_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Asinh is f(x) = asinh(x) (inverse hyperbolic sine)
//
// Asinh is a SuiteSparse:GraphBLAS extension.
func Asinh[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ASINH_FP32
	case float64:
		f.grb = C.GxB_ASINH_FP64
	case complex64:
		f.grb = C.GxB_ASINH_FC32
	case complex128:
		f.grb = C.GxB_ASINH_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Acosh is f(x) = acosh(x) (inverse hyperbolic cosine)
//
// Acosh is a SuiteSparse:GraphBLAS extension.
func Acosh[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ACOSH_FP32
	case float64:
		f.grb = C.GxB_ACOSH_FP64
	case complex64:
		f.grb = C.GxB_ACOSH_FC32
	case complex128:
		f.grb = C.GxB_ACOSH_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Atanh is f(x) = atanh(x) (inverse hyperbolic tangent)
//
// Atanh is a SuiteSparse:GraphBLAS extension.
func Atanh[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ATANH_FP32
	case float64:
		f.grb = C.GxB_ATANH_FP64
	case complex64:
		f.grb = C.GxB_ATANH_FC32
	case complex128:
		f.grb = C.GxB_ATANH_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Signum is f(x) = sgn(x) (sign, or signum)
//
// Signum is a SuiteSparse:GraphBLAS extension.
func Signum[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_SIGNUM_FP32
	case float64:
		f.grb = C.GxB_SIGNUM_FP64
	case complex64:
		f.grb = C.GxB_SIGNUM_FC32
	case complex128:
		f.grb = C.GxB_SIGNUM_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Ceil is f(x) = ceil(x) (ceiling)
//
// Ceil is a SuiteSparse:GraphBLAS extension.
func Ceil[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_CEIL_FP32
	case float64:
		f.grb = C.GxB_CEIL_FP64
	case complex64:
		f.grb = C.GxB_CEIL_FC32
	case complex128:
		f.grb = C.GxB_CEIL_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Floor is f(x) = floor(x) (floor)
//
// Floor is a SuiteSparse:GraphBLAS extension.
func Floor[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_FLOOR_FP32
	case float64:
		f.grb = C.GxB_FLOOR_FP64
	case complex64:
		f.grb = C.GxB_FLOOR_FC32
	case complex128:
		f.grb = C.GxB_FLOOR_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Round is f(x) = round(x) (round to nearest)
//
// Round is a SuiteSparse:GraphBLAS extension.
func Round[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ROUND_FP32
	case float64:
		f.grb = C.GxB_ROUND_FP64
	case complex64:
		f.grb = C.GxB_ROUND_FC32
	case complex128:
		f.grb = C.GxB_ROUND_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Trunc is f(x) = trunc(x) (round towards zero)
//
// Trunc is a SuiteSparse:GraphBLAS extension.
func Trunc[D Float | Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_TRUNC_FP32
	case float64:
		f.grb = C.GxB_TRUNC_FP64
	case complex64:
		f.grb = C.GxB_TRUNC_FC32
	case complex128:
		f.grb = C.GxB_TRUNC_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Isinf is f(x) = true if +/- infinity
//
// Isinf is a SuiteSparse:GraphBLAS extension.
func Isinf[D Float | Complex]() (f UnaryOp[bool, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ISINF_FP32
	case float64:
		f.grb = C.GxB_ISINF_FP64
	case complex64:
		f.grb = C.GxB_ISINF_FC32
	case complex128:
		f.grb = C.GxB_ISINF_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Isnan is f(x) = true if not a number
//
// Isnan is a SuiteSparse:GraphBLAS extension.
func Isnan[D Float | Complex]() (f UnaryOp[bool, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ISNAN_FP32
	case float64:
		f.grb = C.GxB_ISNAN_FP64
	case complex64:
		f.grb = C.GxB_ISNAN_FC32
	case complex128:
		f.grb = C.GxB_ISNAN_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Isfinite is f(x) = true if finite
//
// Isfinite is a SuiteSparse:GraphBLAS extension.
func Isfinite[D Float | Complex]() (f UnaryOp[bool, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ISFINITE_FP32
	case float64:
		f.grb = C.GxB_ISFINITE_FP64
	case complex64:
		f.grb = C.GxB_ISFINITE_FC32
	case complex128:
		f.grb = C.GxB_ISFINITE_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Lgamma is Logarithm of gamma function
//
// Lgamma is a SuiteSparse:GraphBLAS extension.
func Lgamma[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_LGAMMA_FP32
	case float64:
		f.grb = C.GxB_LGAMMA_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Tgamma is gamma function
//
// Tgamma is a SuiteSparse:GraphBLAS extension.
func Tgamma[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_TGAMMA_FP32
	case float64:
		f.grb = C.GxB_TGAMMA_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Erf is error function
//
// Erf is a SuiteSparse:GraphBLAS extension.
func Erf[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ERF_FP32
	case float64:
		f.grb = C.GxB_ERF_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Erfc is complimentary error function
//
// Erfc is a SuiteSparse:GraphBLAS extension.
func Erfc[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ERFC_FP32
	case float64:
		f.grb = C.GxB_ERFC_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Cbrt is cube root
//
// Cbrt is a SuiteSparse:GraphBLAS extension.
func Cbrt[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_CBRT_FP32
	case float64:
		f.grb = C.GxB_CBRT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Frexpx is normalized fraction
//
// Frexpx is a SuiteSparse:GraphBLAS extension.
func Frexpx[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_FREXPX_FP32
	case float64:
		f.grb = C.GxB_FREXPX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Frexpe is normalized exponent
//
// Frexpe is a SuiteSparse:GraphBLAS extension.
func Frexpe[D Float]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_FREXPE_FP32
	case float64:
		f.grb = C.GxB_FREXPE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Conj is complex conjugate
//
// Conj is a SuiteSparse:GraphBLAS extension.
func Conj[D Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case complex64:
		f.grb = C.GxB_CONJ_FC32
	case complex128:
		f.grb = C.GxB_CONJ_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Creal is real part
//
// Creal is a SuiteSparse:GraphBLAS extension.
func Creal[D Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case complex64:
		f.grb = C.GxB_CREAL_FC32
	case complex128:
		f.grb = C.GxB_CREAL_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Cimag is imaginary part
//
// Cimag is a SuiteSparse:GraphBLAS extension.
func Cimag[D Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case complex64:
		f.grb = C.GxB_CIMAG_FC32
	case complex128:
		f.grb = C.GxB_CIMAG_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Carg is angle
//
// Carg is a SuiteSparse:GraphBLAS extension.
func Carg[D Complex]() (f UnaryOp[D, D]) {
	var d D
	switch any(d).(type) {
	case complex64:
		f.grb = C.GxB_CARG_FC32
	case complex128:
		f.grb = C.GxB_CARG_FC64
	default:
		panic("unreachable code")
	}
	return
}
