package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

type (
	// BinaryFunction is the C type for binary functions:
	//    typedef void (*GxB_binary_function) (void *, const void *, const void *) ;
	BinaryFunction C.GxB_binary_function

	// BinaryOp represents a GraphBLAS function that takes one argument of type Din1,
	// one argument of type Din2, and returns an argument of type Dout.
	BinaryOp[Dout, Din1, Din2 any] struct {
		grb C.GrB_BinaryOp
	}
)

// BinaryOpNew returns a new GraphBLAS binary operator with a specified user-defined function in C and its types (domains).
//
// Parameters:
//   - binaryFunc (IN): A pointer to a user-defined function in C that takes two input parameters of types Din1 and Din2 and returns
//     a value of type Dout, all passed as void pointers. Dout, Din1, and Din2 should be one of the [Predefined] GraphBLAS types, one
//     of the [Complex] GraphBLAS types, or a user-defined GraphBLAS type.
//
// GraphBLAS API errors that may be returned:
//   - [NullPointer], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func BinaryOpNew[Dout, Din1, Din2 any](binaryFunc BinaryFunction) (binaryOp BinaryOp[Dout, Din1, Din2], err error) {
	var dout Dout
	doutt, ok := grbType[TypeOf(dout)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	var din1 Din1
	din1t, ok := grbType[TypeOf(din1)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	var din2 Din1
	din2t, ok := grbType[TypeOf(din2)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	info := Info(C.GrB_BinaryOp_new(&binaryOp.grb, binaryFunc, doutt, din1t, din2t))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// NamedBinaryOpNew creates a named binary function. It is like [BinaryOpNew], except:
//   - binopname is the name for the GraphBLAS binary operator. Only the first 127 characters are used.
//   - binopdefn is a string containing the entire function itself.
//
// The two strings binopname and binopdefn are optional, but are required to enable the JIT compilation
// of kernels that use this operator.
//
// If the JIT is enabled, or if the corresponding JIT kernel has been copied into the PreJIT folder,
// the function may be nil. In this case, a JIT kernel is compiled that contains just the user-defined
// function. If the JIT is disabled and the function is nil, this method returns a [NullPointer] error.
//
// NamedBinaryOpNew is a SuiteSparse:GraphBLAS extension.
func NamedBinaryOpNew[Dout, Din1, Din2 any](binaryFunc BinaryFunction, binopname string, binopdefn string) (binaryOp BinaryOp[Dout, Din1, Din2], err error) {
	var dout Dout
	doutt, ok := grbType[TypeOf(dout)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	var din1 Din1
	din1t, ok := grbType[TypeOf(din1)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	var din2 Din1
	din2t, ok := grbType[TypeOf(din2)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	cbinopname := C.CString(binopname)
	defer C.free(unsafe.Pointer(cbinopname))
	cbinopdefn := C.CString(binopdefn)
	defer C.free(unsafe.Pointer(cbinopdefn))
	info := Info(C.GxB_BinaryOp_new(&binaryOp.grb, binaryFunc, doutt, din1t, din2t, cbinopname, cbinopdefn))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Valid returns true if binaryOp has been created by a successful call to [BinaryOpNew] or [NamedBinaryOpNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (binaryOp BinaryOp[Dout, Din1, Din2]) Valid() bool {
	return binaryOp.grb != C.GrB_BinaryOp(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [BinaryOp] and releases any resources associated with
// it. Calling Free on an object that is not [BinaryOp.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined binary operator is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (binaryOp *BinaryOp[Dout, Din1, Din2]) Free() error {
	info := Info(C.GrB_BinaryOp_free(&binaryOp.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the binary operator into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (binaryOp BinaryOp[Dout, Din1, Din2]) Wait(mode WaitMode) error {
	info := Info(C.GrB_BinaryOp_wait(binaryOp.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the binary operator.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (binaryOp BinaryOp[Dout, Din1, Din2]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_BinaryOp_error(&cerror, binaryOp.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the binary operator to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: binaryOp is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (binaryOp BinaryOp[Dout, Din1, Din2]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_BinaryOp_fprint(binaryOp.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// First is f(x, y) = x
func First[D Predefined | Complex, Any any]() (f BinaryOp[D, D, Any]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_FIRST_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_FIRST_INT32
		} else {
			f.grb = C.GrB_FIRST_INT64
		}
	case int8:
		f.grb = C.GrB_FIRST_INT8
	case int16:
		f.grb = C.GrB_FIRST_INT16
	case int32:
		f.grb = C.GrB_FIRST_INT32
	case int64:
		f.grb = C.GrB_FIRST_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_FIRST_UINT32
		} else {
			f.grb = C.GrB_FIRST_UINT64
		}
	case uint8:
		f.grb = C.GrB_FIRST_UINT8
	case uint16:
		f.grb = C.GrB_FIRST_UINT16
	case uint32:
		f.grb = C.GrB_FIRST_UINT32
	case uint64:
		f.grb = C.GrB_FIRST_UINT64
	case float32:
		f.grb = C.GrB_FIRST_FP32
	case float64:
		f.grb = C.GrB_FIRST_FP64
	case complex64:
		f.grb = C.GxB_FIRST_FC32
	case complex128:
		f.grb = C.GxB_FIRST_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Second is f(x, y) = y
func Second[Any any, D Predefined | Complex]() (f BinaryOp[D, Any, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_SECOND_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_SECOND_INT32
		} else {
			f.grb = C.GrB_SECOND_INT64
		}
	case int8:
		f.grb = C.GrB_SECOND_INT8
	case int16:
		f.grb = C.GrB_SECOND_INT16
	case int32:
		f.grb = C.GrB_SECOND_INT32
	case int64:
		f.grb = C.GrB_SECOND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_SECOND_UINT32
		} else {
			f.grb = C.GrB_SECOND_UINT64
		}
	case uint8:
		f.grb = C.GrB_SECOND_UINT8
	case uint16:
		f.grb = C.GrB_SECOND_UINT16
	case uint32:
		f.grb = C.GrB_SECOND_UINT32
	case uint64:
		f.grb = C.GrB_SECOND_UINT64
	case float32:
		f.grb = C.GrB_SECOND_FP32
	case float64:
		f.grb = C.GrB_SECOND_FP64
	case complex64:
		f.grb = C.GxB_SECOND_FC32
	case complex128:
		f.grb = C.GxB_SECOND_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Any is f(x, y) = x or y, picked arbitrarily
//
// Any is a a SuiteSparse:GraphBLAS extension.
func Any[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ANY_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ANY_INT32
		} else {
			f.grb = C.GxB_ANY_INT64
		}
	case int8:
		f.grb = C.GxB_ANY_INT8
	case int16:
		f.grb = C.GxB_ANY_INT16
	case int32:
		f.grb = C.GxB_ANY_INT32
	case int64:
		f.grb = C.GxB_ANY_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ANY_UINT32
		} else {
			f.grb = C.GxB_ANY_UINT64
		}
	case uint8:
		f.grb = C.GxB_ANY_UINT8
	case uint16:
		f.grb = C.GxB_ANY_UINT16
	case uint32:
		f.grb = C.GxB_ANY_UINT32
	case uint64:
		f.grb = C.GxB_ANY_UINT64
	case float32:
		f.grb = C.GxB_ANY_FP32
	case float64:
		f.grb = C.GxB_ANY_FP64
	case complex64:
		f.grb = C.GxB_ANY_FC32
	case complex128:
		f.grb = C.GxB_ANY_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Oneb is f(x, y) = 1 when D is { Number | Complex }. f(x, y) = true when D is bool.
func Oneb[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_ONEB_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_ONEB_INT32
		} else {
			f.grb = C.GrB_ONEB_INT64
		}
	case int8:
		f.grb = C.GrB_ONEB_INT8
	case int16:
		f.grb = C.GrB_ONEB_INT16
	case int32:
		f.grb = C.GrB_ONEB_INT32
	case int64:
		f.grb = C.GrB_ONEB_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_ONEB_UINT32
		} else {
			f.grb = C.GrB_ONEB_UINT64
		}
	case uint8:
		f.grb = C.GrB_ONEB_UINT8
	case uint16:
		f.grb = C.GrB_ONEB_UINT16
	case uint32:
		f.grb = C.GrB_ONEB_UINT32
	case uint64:
		f.grb = C.GrB_ONEB_UINT64
	case float32:
		f.grb = C.GrB_ONEB_FP32
	case float64:
		f.grb = C.GrB_ONEB_FP64
	case complex64:
		f.grb = C.GxB_ONEB_FC32
	case complex128:
		f.grb = C.GxB_ONEB_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Plus is f(x, y) = x + y
func Plus[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_PLUS_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_PLUS_INT32
		} else {
			f.grb = C.GrB_PLUS_INT64
		}
	case int8:
		f.grb = C.GrB_PLUS_INT8
	case int16:
		f.grb = C.GrB_PLUS_INT16
	case int32:
		f.grb = C.GrB_PLUS_INT32
	case int64:
		f.grb = C.GrB_PLUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_PLUS_UINT32
		} else {
			f.grb = C.GrB_PLUS_UINT64
		}
	case uint8:
		f.grb = C.GrB_PLUS_UINT8
	case uint16:
		f.grb = C.GrB_PLUS_UINT16
	case uint32:
		f.grb = C.GrB_PLUS_UINT32
	case uint64:
		f.grb = C.GrB_PLUS_UINT64
	case float32:
		f.grb = C.GrB_PLUS_FP32
	case float64:
		f.grb = C.GrB_PLUS_FP64
	case complex64:
		f.grb = C.GxB_PLUS_FC32
	case complex128:
		f.grb = C.GxB_PLUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Minus is f(x, y) = x - y
func Minus[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_MINUS_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MINUS_INT32
		} else {
			f.grb = C.GrB_MINUS_INT64
		}
	case int8:
		f.grb = C.GrB_MINUS_INT8
	case int16:
		f.grb = C.GrB_MINUS_INT16
	case int32:
		f.grb = C.GrB_MINUS_INT32
	case int64:
		f.grb = C.GrB_MINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MINUS_UINT32
		} else {
			f.grb = C.GrB_MINUS_UINT64
		}
	case uint8:
		f.grb = C.GrB_MINUS_UINT8
	case uint16:
		f.grb = C.GrB_MINUS_UINT16
	case uint32:
		f.grb = C.GrB_MINUS_UINT32
	case uint64:
		f.grb = C.GrB_MINUS_UINT64
	case float32:
		f.grb = C.GrB_MINUS_FP32
	case float64:
		f.grb = C.GrB_MINUS_FP64
	case complex64:
		f.grb = C.GxB_MINUS_FC32
	case complex128:
		f.grb = C.GxB_MINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Rminus is f(x, y) = y - x
//
// Rminus is a SuiteSparse:GraphBLAS extension.
func Rminus[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_RMINUS_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_RMINUS_INT32
		} else {
			f.grb = C.GxB_RMINUS_INT64
		}
	case int8:
		f.grb = C.GxB_RMINUS_INT8
	case int16:
		f.grb = C.GxB_RMINUS_INT16
	case int32:
		f.grb = C.GxB_RMINUS_INT32
	case int64:
		f.grb = C.GxB_RMINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_RMINUS_UINT32
		} else {
			f.grb = C.GxB_RMINUS_UINT64
		}
	case uint8:
		f.grb = C.GxB_RMINUS_UINT8
	case uint16:
		f.grb = C.GxB_RMINUS_UINT16
	case uint32:
		f.grb = C.GxB_RMINUS_UINT32
	case uint64:
		f.grb = C.GxB_RMINUS_UINT64
	case float32:
		f.grb = C.GxB_RMINUS_FP32
	case float64:
		f.grb = C.GxB_RMINUS_FP64
	case complex64:
		f.grb = C.GxB_RMINUS_FC32
	case complex128:
		f.grb = C.GxB_RMINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Times is f(x, y) = x * y
func Times[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_TIMES_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_TIMES_INT32
		} else {
			f.grb = C.GrB_TIMES_INT64
		}
	case int8:
		f.grb = C.GrB_TIMES_INT8
	case int16:
		f.grb = C.GrB_TIMES_INT16
	case int32:
		f.grb = C.GrB_TIMES_INT32
	case int64:
		f.grb = C.GrB_TIMES_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_TIMES_UINT32
		} else {
			f.grb = C.GrB_TIMES_UINT64
		}
	case uint8:
		f.grb = C.GrB_TIMES_UINT8
	case uint16:
		f.grb = C.GrB_TIMES_UINT16
	case uint32:
		f.grb = C.GrB_TIMES_UINT32
	case uint64:
		f.grb = C.GrB_TIMES_UINT64
	case float32:
		f.grb = C.GrB_TIMES_FP32
	case float64:
		f.grb = C.GrB_TIMES_FP64
	case complex64:
		f.grb = C.GxB_TIMES_FC32
	case complex128:
		f.grb = C.GxB_TIMES_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Div is f(x, y) = x / y
func Div[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_DIV_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_DIV_INT32
		} else {
			f.grb = C.GrB_DIV_INT64
		}
	case int8:
		f.grb = C.GrB_DIV_INT8
	case int16:
		f.grb = C.GrB_DIV_INT16
	case int32:
		f.grb = C.GrB_DIV_INT32
	case int64:
		f.grb = C.GrB_DIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_DIV_UINT32
		} else {
			f.grb = C.GrB_DIV_UINT64
		}
	case uint8:
		f.grb = C.GrB_DIV_UINT8
	case uint16:
		f.grb = C.GrB_DIV_UINT16
	case uint32:
		f.grb = C.GrB_DIV_UINT32
	case uint64:
		f.grb = C.GrB_DIV_UINT64
	case float32:
		f.grb = C.GrB_DIV_FP32
	case float64:
		f.grb = C.GrB_DIV_FP64
	case complex64:
		f.grb = C.GxB_DIV_FC32
	case complex128:
		f.grb = C.GxB_DIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Rdiv is f(x, y) = y / x
//
// Rdiv is a a SuiteSparse:GraphBLAS extension.
func Rdiv[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_RDIV_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_RDIV_INT32
		} else {
			f.grb = C.GxB_RDIV_INT64
		}
	case int8:
		f.grb = C.GxB_RDIV_INT8
	case int16:
		f.grb = C.GxB_RDIV_INT16
	case int32:
		f.grb = C.GxB_RDIV_INT32
	case int64:
		f.grb = C.GxB_RDIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_RDIV_UINT32
		} else {
			f.grb = C.GxB_RDIV_UINT64
		}
	case uint8:
		f.grb = C.GxB_RDIV_UINT8
	case uint16:
		f.grb = C.GxB_RDIV_UINT16
	case uint32:
		f.grb = C.GxB_RDIV_UINT32
	case uint64:
		f.grb = C.GxB_RDIV_UINT64
	case float32:
		f.grb = C.GxB_RDIV_FP32
	case float64:
		f.grb = C.GxB_RDIV_FP64
	case complex64:
		f.grb = C.GxB_RDIV_FC32
	case complex128:
		f.grb = C.GxB_RDIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Pow is f(x, y) = x raised to the power of y
//
// Pow is a SuiteSparse:GraphBLAS extension.
func Pow[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_POW_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_POW_INT32
		} else {
			f.grb = C.GxB_POW_INT64
		}
	case int8:
		f.grb = C.GxB_POW_INT8
	case int16:
		f.grb = C.GxB_POW_INT16
	case int32:
		f.grb = C.GxB_POW_INT32
	case int64:
		f.grb = C.GxB_POW_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_POW_UINT32
		} else {
			f.grb = C.GxB_POW_UINT64
		}
	case uint8:
		f.grb = C.GxB_POW_UINT8
	case uint16:
		f.grb = C.GxB_POW_UINT16
	case uint32:
		f.grb = C.GxB_POW_UINT32
	case uint64:
		f.grb = C.GxB_POW_UINT64
	case float32:
		f.grb = C.GxB_POW_FP32
	case float64:
		f.grb = C.GxB_POW_FP64
	case complex64:
		f.grb = C.GxB_POW_FC32
	case complex128:
		f.grb = C.GxB_POW_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Iseq is f(x, y) = x == y
//
// Iseq is a SuiteSparse:GraphBLAS extension.
func Iseq[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ISEQ_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISEQ_INT32
		} else {
			f.grb = C.GxB_ISEQ_INT64
		}
	case int8:
		f.grb = C.GxB_ISEQ_INT8
	case int16:
		f.grb = C.GxB_ISEQ_INT16
	case int32:
		f.grb = C.GxB_ISEQ_INT32
	case int64:
		f.grb = C.GxB_ISEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISEQ_UINT32
		} else {
			f.grb = C.GxB_ISEQ_UINT64
		}
	case uint8:
		f.grb = C.GxB_ISEQ_UINT8
	case uint16:
		f.grb = C.GxB_ISEQ_UINT16
	case uint32:
		f.grb = C.GxB_ISEQ_UINT32
	case uint64:
		f.grb = C.GxB_ISEQ_UINT64
	case float32:
		f.grb = C.GxB_ISEQ_FP32
	case float64:
		f.grb = C.GxB_ISEQ_FP64
	case complex64:
		f.grb = C.GxB_ISEQ_FC32
	case complex128:
		f.grb = C.GxB_ISEQ_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Isne is f(x, y) = x != y
//
// Isne is a SuiteSparse:GraphBLAS extension.
func Isne[D Predefined | Complex]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ISNE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISNE_INT32
		} else {
			f.grb = C.GxB_ISNE_INT64
		}
	case int8:
		f.grb = C.GxB_ISNE_INT8
	case int16:
		f.grb = C.GxB_ISNE_INT16
	case int32:
		f.grb = C.GxB_ISNE_INT32
	case int64:
		f.grb = C.GxB_ISNE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISNE_UINT32
		} else {
			f.grb = C.GxB_ISNE_UINT64
		}
	case uint8:
		f.grb = C.GxB_ISNE_UINT8
	case uint16:
		f.grb = C.GxB_ISNE_UINT16
	case uint32:
		f.grb = C.GxB_ISNE_UINT32
	case uint64:
		f.grb = C.GxB_ISNE_UINT64
	case float32:
		f.grb = C.GxB_ISNE_FP32
	case float64:
		f.grb = C.GxB_ISNE_FP64
	case complex64:
		f.grb = C.GxB_ISNE_FC32
	case complex128:
		f.grb = C.GxB_ISNE_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Min is f(x, y) = minimum of x and y
func Min[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_MIN_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MIN_INT32
		} else {
			f.grb = C.GrB_MIN_INT64
		}
	case int8:
		f.grb = C.GrB_MIN_INT8
	case int16:
		f.grb = C.GrB_MIN_INT16
	case int32:
		f.grb = C.GrB_MIN_INT32
	case int64:
		f.grb = C.GrB_MIN_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MIN_UINT32
		} else {
			f.grb = C.GrB_MIN_UINT64
		}
	case uint8:
		f.grb = C.GrB_MIN_UINT8
	case uint16:
		f.grb = C.GrB_MIN_UINT16
	case uint32:
		f.grb = C.GrB_MIN_UINT32
	case uint64:
		f.grb = C.GrB_MIN_UINT64
	case float32:
		f.grb = C.GrB_MIN_FP32
	case float64:
		f.grb = C.GrB_MIN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Max is f(x, y) = maximum of x and y
func Max[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_MAX_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MAX_INT32
		} else {
			f.grb = C.GrB_MAX_INT64
		}
	case int8:
		f.grb = C.GrB_MAX_INT8
	case int16:
		f.grb = C.GrB_MAX_INT16
	case int32:
		f.grb = C.GrB_MAX_INT32
	case int64:
		f.grb = C.GrB_MAX_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_MAX_UINT32
		} else {
			f.grb = C.GrB_MAX_UINT64
		}
	case uint8:
		f.grb = C.GrB_MAX_UINT8
	case uint16:
		f.grb = C.GrB_MAX_UINT16
	case uint32:
		f.grb = C.GrB_MAX_UINT32
	case uint64:
		f.grb = C.GrB_MAX_UINT64
	case float32:
		f.grb = C.GrB_MAX_FP32
	case float64:
		f.grb = C.GrB_MAX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Isgt is f(x, y) = x > y
//
// Isgt is a SuiteSparse:GraphBLAS extension.
func Isgt[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ISGT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISGT_INT32
		} else {
			f.grb = C.GxB_ISGT_INT64
		}
	case int8:
		f.grb = C.GxB_ISGT_INT8
	case int16:
		f.grb = C.GxB_ISGT_INT16
	case int32:
		f.grb = C.GxB_ISGT_INT32
	case int64:
		f.grb = C.GxB_ISGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISGT_UINT32
		} else {
			f.grb = C.GxB_ISGT_UINT64
		}
	case uint8:
		f.grb = C.GxB_ISGT_UINT8
	case uint16:
		f.grb = C.GxB_ISGT_UINT16
	case uint32:
		f.grb = C.GxB_ISGT_UINT32
	case uint64:
		f.grb = C.GxB_ISGT_UINT64
	case float32:
		f.grb = C.GxB_ISGT_FP32
	case float64:
		f.grb = C.GxB_ISGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Islt is f(x, y) = x < y
//
// Islt is a SuiteSparse:GraphBLAS extension.
func Islt[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ISLT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISLT_INT32
		} else {
			f.grb = C.GxB_ISLT_INT64
		}
	case int8:
		f.grb = C.GxB_ISLT_INT8
	case int16:
		f.grb = C.GxB_ISLT_INT16
	case int32:
		f.grb = C.GxB_ISLT_INT32
	case int64:
		f.grb = C.GxB_ISLT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISLT_UINT32
		} else {
			f.grb = C.GxB_ISLT_UINT64
		}
	case uint8:
		f.grb = C.GxB_ISLT_UINT8
	case uint16:
		f.grb = C.GxB_ISLT_UINT16
	case uint32:
		f.grb = C.GxB_ISLT_UINT32
	case uint64:
		f.grb = C.GxB_ISLT_UINT64
	case float32:
		f.grb = C.GxB_ISLT_FP32
	case float64:
		f.grb = C.GxB_ISLT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Isge is f(x, y) = x >= y
//
// Isge is a SuiteSparse:GraphBLAS extension.
func Isge[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ISGE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISGE_INT32
		} else {
			f.grb = C.GxB_ISGE_INT64
		}
	case int8:
		f.grb = C.GxB_ISGE_INT8
	case int16:
		f.grb = C.GxB_ISGE_INT16
	case int32:
		f.grb = C.GxB_ISGE_INT32
	case int64:
		f.grb = C.GxB_ISGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISGE_UINT32
		} else {
			f.grb = C.GxB_ISGE_UINT64
		}
	case uint8:
		f.grb = C.GxB_ISGE_UINT8
	case uint16:
		f.grb = C.GxB_ISGE_UINT16
	case uint32:
		f.grb = C.GxB_ISGE_UINT32
	case uint64:
		f.grb = C.GxB_ISGE_UINT64
	case float32:
		f.grb = C.GxB_ISGE_FP32
	case float64:
		f.grb = C.GxB_ISGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Isle is f(x, y) = x <= y
//
// Isle is a SuiteSparse:GraphBLAS extension.
func Isle[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_ISLE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISLE_INT32
		} else {
			f.grb = C.GxB_ISLE_INT64
		}
	case int8:
		f.grb = C.GxB_ISLE_INT8
	case int16:
		f.grb = C.GxB_ISLE_INT16
	case int32:
		f.grb = C.GxB_ISLE_INT32
	case int64:
		f.grb = C.GxB_ISLE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_ISLE_UINT32
		} else {
			f.grb = C.GxB_ISLE_UINT64
		}
	case uint8:
		f.grb = C.GxB_ISLE_UINT8
	case uint16:
		f.grb = C.GxB_ISLE_UINT16
	case uint32:
		f.grb = C.GxB_ISLE_UINT32
	case uint64:
		f.grb = C.GxB_ISLE_UINT64
	case float32:
		f.grb = C.GxB_ISLE_FP32
	case float64:
		f.grb = C.GxB_ISLE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Lor is f(x, y) = (x != 0) || (y != 0)
//
// Lor is a SuiteSparse:GraphBLAS extension.
func Lor[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_LOR_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LOR_INT32
		} else {
			f.grb = C.GxB_LOR_INT64
		}
	case int8:
		f.grb = C.GxB_LOR_INT8
	case int16:
		f.grb = C.GxB_LOR_INT16
	case int32:
		f.grb = C.GxB_LOR_INT32
	case int64:
		f.grb = C.GxB_LOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LOR_UINT32
		} else {
			f.grb = C.GxB_LOR_UINT64
		}
	case uint8:
		f.grb = C.GxB_LOR_UINT8
	case uint16:
		f.grb = C.GxB_LOR_UINT16
	case uint32:
		f.grb = C.GxB_LOR_UINT32
	case uint64:
		f.grb = C.GxB_LOR_UINT64
	case float32:
		f.grb = C.GxB_LOR_FP32
	case float64:
		f.grb = C.GxB_LOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Land is f(x, y) = (x != 0) && (y != 0)
//
// Land is a SuiteSparse:GraphBLAS extension.
func Land[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_LAND_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LAND_INT32
		} else {
			f.grb = C.GxB_LAND_INT64
		}
	case int8:
		f.grb = C.GxB_LAND_INT8
	case int16:
		f.grb = C.GxB_LAND_INT16
	case int32:
		f.grb = C.GxB_LAND_INT32
	case int64:
		f.grb = C.GxB_LAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LAND_UINT32
		} else {
			f.grb = C.GxB_LAND_UINT64
		}
	case uint8:
		f.grb = C.GxB_LAND_UINT8
	case uint16:
		f.grb = C.GxB_LAND_UINT16
	case uint32:
		f.grb = C.GxB_LAND_UINT32
	case uint64:
		f.grb = C.GxB_LAND_UINT64
	case float32:
		f.grb = C.GxB_LAND_FP32
	case float64:
		f.grb = C.GxB_LAND_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Lxor is f(x, y) = (x != 0) != (y != 0)
//
// Lxor is a SuiteSparse:GraphBLAS extension.
func Lxor[D Predefined]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GxB_LXOR_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LXOR_INT32
		} else {
			f.grb = C.GxB_LXOR_INT64
		}
	case int8:
		f.grb = C.GxB_LXOR_INT8
	case int16:
		f.grb = C.GxB_LXOR_INT16
	case int32:
		f.grb = C.GxB_LXOR_INT32
	case int64:
		f.grb = C.GxB_LXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_LXOR_UINT32
		} else {
			f.grb = C.GxB_LXOR_UINT64
		}
	case uint8:
		f.grb = C.GxB_LXOR_UINT8
	case uint16:
		f.grb = C.GxB_LXOR_UINT16
	case uint32:
		f.grb = C.GxB_LXOR_UINT32
	case uint64:
		f.grb = C.GxB_LXOR_UINT64
	case float32:
		f.grb = C.GxB_LXOR_FP32
	case float64:
		f.grb = C.GxB_LXOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Eq is f(x, y) = x == y
func Eq[D Predefined | Complex]() (f BinaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_EQ_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_EQ_INT32
		} else {
			f.grb = C.GrB_EQ_INT64
		}
	case int8:
		f.grb = C.GrB_EQ_INT8
	case int16:
		f.grb = C.GrB_EQ_INT16
	case int32:
		f.grb = C.GrB_EQ_INT32
	case int64:
		f.grb = C.GrB_EQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_EQ_UINT32
		} else {
			f.grb = C.GrB_EQ_UINT64
		}
	case uint8:
		f.grb = C.GrB_EQ_UINT8
	case uint16:
		f.grb = C.GrB_EQ_UINT16
	case uint32:
		f.grb = C.GrB_EQ_UINT32
	case uint64:
		f.grb = C.GrB_EQ_UINT64
	case float32:
		f.grb = C.GrB_EQ_FP32
	case float64:
		f.grb = C.GrB_EQ_FP64
	case complex64:
		f.grb = C.GxB_EQ_FC32
	case complex128:
		f.grb = C.GxB_EQ_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Ne is f(x, y) = x != y
func Ne[D Predefined | Complex]() (f BinaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_NE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_NE_INT32
		} else {
			f.grb = C.GrB_NE_INT64
		}
	case int8:
		f.grb = C.GrB_NE_INT8
	case int16:
		f.grb = C.GrB_NE_INT16
	case int32:
		f.grb = C.GrB_NE_INT32
	case int64:
		f.grb = C.GrB_NE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_NE_UINT32
		} else {
			f.grb = C.GrB_NE_UINT64
		}
	case uint8:
		f.grb = C.GrB_NE_UINT8
	case uint16:
		f.grb = C.GrB_NE_UINT16
	case uint32:
		f.grb = C.GrB_NE_UINT32
	case uint64:
		f.grb = C.GrB_NE_UINT64
	case float32:
		f.grb = C.GrB_NE_FP32
	case float64:
		f.grb = C.GrB_NE_FP64
	case complex64:
		f.grb = C.GxB_NE_FC32
	case complex128:
		f.grb = C.GxB_NE_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Gt is f(x, y) = x > y
func Gt[D Predefined]() (f BinaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_GT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_GT_INT32
		} else {
			f.grb = C.GrB_GT_INT64
		}
	case int8:
		f.grb = C.GrB_GT_INT8
	case int16:
		f.grb = C.GrB_GT_INT16
	case int32:
		f.grb = C.GrB_GT_INT32
	case int64:
		f.grb = C.GrB_GT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_GT_UINT32
		} else {
			f.grb = C.GrB_GT_UINT64
		}
	case uint8:
		f.grb = C.GrB_GT_UINT8
	case uint16:
		f.grb = C.GrB_GT_UINT16
	case uint32:
		f.grb = C.GrB_GT_UINT32
	case uint64:
		f.grb = C.GrB_GT_UINT64
	case float32:
		f.grb = C.GrB_GT_FP32
	case float64:
		f.grb = C.GrB_GT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Lt is f(x, y) = x < y
func Lt[D Predefined]() (f BinaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_LT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_LT_INT32
		} else {
			f.grb = C.GrB_LT_INT64
		}
	case int8:
		f.grb = C.GrB_LT_INT8
	case int16:
		f.grb = C.GrB_LT_INT16
	case int32:
		f.grb = C.GrB_LT_INT32
	case int64:
		f.grb = C.GrB_LT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_LT_UINT32
		} else {
			f.grb = C.GrB_LT_UINT64
		}
	case uint8:
		f.grb = C.GrB_LT_UINT8
	case uint16:
		f.grb = C.GrB_LT_UINT16
	case uint32:
		f.grb = C.GrB_LT_UINT32
	case uint64:
		f.grb = C.GrB_LT_UINT64
	case float32:
		f.grb = C.GrB_LT_FP32
	case float64:
		f.grb = C.GrB_LT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Ge is f(x, y) = x >= y
func Ge[D Predefined]() (f BinaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_GE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_GE_INT32
		} else {
			f.grb = C.GrB_GE_INT64
		}
	case int8:
		f.grb = C.GrB_GE_INT8
	case int16:
		f.grb = C.GrB_GE_INT16
	case int32:
		f.grb = C.GrB_GE_INT32
	case int64:
		f.grb = C.GrB_GE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_GE_UINT32
		} else {
			f.grb = C.GrB_GE_UINT64
		}
	case uint8:
		f.grb = C.GrB_GE_UINT8
	case uint16:
		f.grb = C.GrB_GE_UINT16
	case uint32:
		f.grb = C.GrB_GE_UINT32
	case uint64:
		f.grb = C.GrB_GE_UINT64
	case float32:
		f.grb = C.GrB_GE_FP32
	case float64:
		f.grb = C.GrB_GE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Le is f(x, y) = x <= y
func Le[D Predefined]() (f BinaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_LE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_LE_INT32
		} else {
			f.grb = C.GrB_LE_INT64
		}
	case int8:
		f.grb = C.GrB_LE_INT8
	case int16:
		f.grb = C.GrB_LE_INT16
	case int32:
		f.grb = C.GrB_LE_INT32
	case int64:
		f.grb = C.GrB_LE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_LE_UINT32
		} else {
			f.grb = C.GrB_LE_UINT64
		}
	case uint8:
		f.grb = C.GrB_LE_UINT8
	case uint16:
		f.grb = C.GrB_LE_UINT16
	case uint32:
		f.grb = C.GrB_LE_UINT32
	case uint64:
		f.grb = C.GrB_LE_UINT64
	case float32:
		f.grb = C.GrB_LE_FP32
	case float64:
		f.grb = C.GrB_LE_FP64
	default:
		panic("unreachable code")
	}
	return
}

var (
	// LorBool is f(x, y) = (x != 0) || (y != 0)
	LorBool = BinaryOp[bool, bool, bool]{C.GrB_LOR}

	// LandBool is f(x, y) = (x != 0) && (y != 0)
	LandBool = BinaryOp[bool, bool, bool]{C.GrB_LAND}

	// LxorBool is f(x, y) = (x != 0) != (y != 0)
	LxorBool = BinaryOp[bool, bool, bool]{C.GrB_LXOR}

	// LxnorBool is f(x, y) = (x != 0) == (y != 0)
	LxnorBool = BinaryOp[bool, bool, bool]{C.GrB_LXNOR}
)

// Atan2 is 4-quadrant arc tangent
//
// Atan2 is a SuiteSparse:GraphBLAS extension.
func Atan2[D Float]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_ATAN2_FP32
	case float64:
		f.grb = C.GxB_ATAN2_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Hypot is hypotenuse
//
// Hypot is a SuiteSparse:GraphBLAS extension.
func Hypot[D Float]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_HYPOT_FP32
	case float64:
		f.grb = C.GxB_HYPOT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Fmod is ANSI C11 fmod
//
// Fmod is a SuiteSparse:GraphBLAS extension.
func Fmod[D Float]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_FMOD_FP32
	case float64:
		f.grb = C.GxB_FMOD_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Remainder is ANSI C11 remainder
//
// Remainder is a SuiteSparse:GraphBLAS extension.
func Remainder[D Float]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_REMAINDER_FP32
	case float64:
		f.grb = C.GxB_REMAINDER_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Ldexp is ANSI C11 ldexp
//
// Ldexp is a SuiteSparse:GraphBLAS extension.
func Ldexp[D Float]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_LDEXP_FP32
	case float64:
		f.grb = C.GxB_LDEXP_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Copysign is ANSI C11 copysign.
//
// Copysign is a SuiteSparse extension.
func Copysign[D Float]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case float32:
		f.grb = C.GxB_COPYSIGN_FP32
	case float64:
		f.grb = C.GxB_COPYSIGN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Cmplx is f(x, y) = x + y * i
//
// If C is complex64, F must be float32; if C is complex128, F must be float64.
//
// GraphBLAS API errors that may be returned:
//   - DomainMismatch: C and F are not compatible with each other.
//
// Cmplx is a SuiteSparse:GraphBLAS extension.
func Cmplx[C Complex, F Float]() (f BinaryOp[C, F, F], err error) {
	var c C
	switch any(c).(type) {
	case complex64:
		var x F
		switch any(x).(type) {
		case float32:
			f.grb = C.GxB_CMPLX_FP32
		default:
			err = DomainMismatch
		}
	case complex128:
		var x F
		switch any(x).(type) {
		case float64:
			f.grb = C.GxB_CMPLX_FP64
		default:
			err = DomainMismatch
		}
	default:
		panic("unreachable code")
	}
	return
}

// Bor is f(x, y) = x | y
func Bor[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BOR_INT32
		} else {
			f.grb = C.GrB_BOR_INT64
		}
	case int8:
		f.grb = C.GrB_BOR_INT8
	case int16:
		f.grb = C.GrB_BOR_INT16
	case int32:
		f.grb = C.GrB_BOR_INT32
	case int64:
		f.grb = C.GrB_BOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BOR_UINT32
		} else {
			f.grb = C.GrB_BOR_UINT64
		}
	case uint8:
		f.grb = C.GrB_BOR_UINT8
	case uint16:
		f.grb = C.GrB_BOR_UINT16
	case uint32:
		f.grb = C.GrB_BOR_UINT32
	case uint64:
		f.grb = C.GrB_BOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Band is f(x, y) = x & y
func Band[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BAND_INT32
		} else {
			f.grb = C.GrB_BAND_INT64
		}
	case int8:
		f.grb = C.GrB_BAND_INT8
	case int16:
		f.grb = C.GrB_BAND_INT16
	case int32:
		f.grb = C.GrB_BAND_INT32
	case int64:
		f.grb = C.GrB_BAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BAND_UINT32
		} else {
			f.grb = C.GrB_BAND_UINT64
		}
	case uint8:
		f.grb = C.GrB_BAND_UINT8
	case uint16:
		f.grb = C.GrB_BAND_UINT16
	case uint32:
		f.grb = C.GrB_BAND_UINT32
	case uint64:
		f.grb = C.GrB_BAND_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Bxor is f(x, y) = x ^ y
func Bxor[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BXOR_INT32
		} else {
			f.grb = C.GrB_BXOR_INT64
		}
	case int8:
		f.grb = C.GrB_BXOR_INT8
	case int16:
		f.grb = C.GrB_BXOR_INT16
	case int32:
		f.grb = C.GrB_BXOR_INT32
	case int64:
		f.grb = C.GrB_BXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BXOR_UINT32
		} else {
			f.grb = C.GrB_BXOR_UINT64
		}
	case uint8:
		f.grb = C.GrB_BXOR_UINT8
	case uint16:
		f.grb = C.GrB_BXOR_UINT16
	case uint32:
		f.grb = C.GrB_BXOR_UINT32
	case uint64:
		f.grb = C.GrB_BXOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Bxnor is f(x, y) = ^(x ^ y)
func Bxnor[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BXNOR_INT32
		} else {
			f.grb = C.GrB_BXNOR_INT64
		}
	case int8:
		f.grb = C.GrB_BXNOR_INT8
	case int16:
		f.grb = C.GrB_BXNOR_INT16
	case int32:
		f.grb = C.GrB_BXNOR_INT32
	case int64:
		f.grb = C.GrB_BXNOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_BXNOR_UINT32
		} else {
			f.grb = C.GrB_BXNOR_UINT64
		}
	case uint8:
		f.grb = C.GrB_BXNOR_UINT8
	case uint16:
		f.grb = C.GrB_BXNOR_UINT16
	case uint32:
		f.grb = C.GrB_BXNOR_UINT32
	case uint64:
		f.grb = C.GrB_BXNOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Bget is f(x, y) = get bit y of x
//
// Bget is a SuiteSparse:GraphBLAS extension.
func Bget[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BGET_INT32
		} else {
			f.grb = C.GxB_BGET_INT64
		}
	case int8:
		f.grb = C.GxB_BGET_INT8
	case int16:
		f.grb = C.GxB_BGET_INT16
	case int32:
		f.grb = C.GxB_BGET_INT32
	case int64:
		f.grb = C.GxB_BGET_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BGET_UINT32
		} else {
			f.grb = C.GxB_BGET_UINT64
		}
	case uint8:
		f.grb = C.GxB_BGET_UINT8
	case uint16:
		f.grb = C.GxB_BGET_UINT16
	case uint32:
		f.grb = C.GxB_BGET_UINT32
	case uint64:
		f.grb = C.GxB_BGET_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Bset is f(x, y) = set bit y of x
//
// Bset is a SuiteSparse:GraphBLAS extension.
func Bset[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BSET_INT32
		} else {
			f.grb = C.GxB_BSET_INT64
		}
	case int8:
		f.grb = C.GxB_BSET_INT8
	case int16:
		f.grb = C.GxB_BSET_INT16
	case int32:
		f.grb = C.GxB_BSET_INT32
	case int64:
		f.grb = C.GxB_BSET_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BSET_UINT32
		} else {
			f.grb = C.GxB_BSET_UINT64
		}
	case uint8:
		f.grb = C.GxB_BSET_UINT8
	case uint16:
		f.grb = C.GxB_BSET_UINT16
	case uint32:
		f.grb = C.GxB_BSET_UINT32
	case uint64:
		f.grb = C.GxB_BSET_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Bclr is f(x, y) = clear bit y of x
//
// Bclr is a SuiteSparse:GraphBLAS extension.
func Bclr[D Integer]() (f BinaryOp[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BCLR_INT32
		} else {
			f.grb = C.GxB_BCLR_INT64
		}
	case int8:
		f.grb = C.GxB_BCLR_INT8
	case int16:
		f.grb = C.GxB_BCLR_INT16
	case int32:
		f.grb = C.GxB_BCLR_INT32
	case int64:
		f.grb = C.GxB_BCLR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BCLR_UINT32
		} else {
			f.grb = C.GxB_BCLR_UINT64
		}
	case uint8:
		f.grb = C.GxB_BCLR_UINT8
	case uint16:
		f.grb = C.GxB_BCLR_UINT16
	case uint32:
		f.grb = C.GxB_BCLR_UINT32
	case uint64:
		f.grb = C.GxB_BCLR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Bshift is bit shift
//
// Bshift is a SuiteSparse:GraphBLAS extension.
func Bshift[D Integer]() (f BinaryOp[D, D, int8]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BSHIFT_INT32
		} else {
			f.grb = C.GxB_BSHIFT_INT64
		}
	case int8:
		f.grb = C.GxB_BSHIFT_INT8
	case int16:
		f.grb = C.GxB_BSHIFT_INT16
	case int32:
		f.grb = C.GxB_BSHIFT_INT32
	case int64:
		f.grb = C.GxB_BSHIFT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_BSHIFT_UINT32
		} else {
			f.grb = C.GxB_BSHIFT_UINT64
		}
	case uint8:
		f.grb = C.GxB_BSHIFT_UINT8
	case uint16:
		f.grb = C.GxB_BSHIFT_UINT16
	case uint32:
		f.grb = C.GxB_BSHIFT_UINT32
	case uint64:
		f.grb = C.GxB_BSHIFT_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// Firsti is f(x, y) = row index of x (0-based)
//
// Firsti is a SuiteSparse:GraphBLAS extension.
func Firsti[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_FIRSTI_INT32
		} else {
			f.grb = C.GxB_FIRSTI_INT64
		}
	case int32:
		f.grb = C.GxB_FIRSTI_INT32
	case int64:
		f.grb = C.GxB_FIRSTI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Firsti1 is f(x, y) = row index of x (1-based)
//
// Firsti1 is a SuiteSparse:GraphBLAS extension.
func Firsti1[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_FIRSTI1_INT32
		} else {
			f.grb = C.GxB_FIRSTI1_INT64
		}
	case int32:
		f.grb = C.GxB_FIRSTI1_INT32
	case int64:
		f.grb = C.GxB_FIRSTI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Firstj is f(x, y) = column index of x (0-based)
//
// Firstj is a SuiteSparse:GraphBLAS extension.
func Firstj[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_FIRSTJ_INT32
		} else {
			f.grb = C.GxB_FIRSTJ_INT64
		}
	case int32:
		f.grb = C.GxB_FIRSTJ_INT32
	case int64:
		f.grb = C.GxB_FIRSTJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Firstj1 is f(x, y) = column index of x (1-based)
//
// Firstj1 is a SuiteSparse:GraphBLAS extension.
func Firstj1[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_FIRSTJ1_INT32
		} else {
			f.grb = C.GxB_FIRSTJ1_INT64
		}
	case int32:
		f.grb = C.GxB_FIRSTJ1_INT32
	case int64:
		f.grb = C.GxB_FIRSTJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Secondi is f(x, y) = row index of y (0-based)
//
// Secondi is a SuiteSparse:GraphBLAS extension.
func Secondi[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_SECONDI_INT32
		} else {
			f.grb = C.GxB_SECONDI_INT64
		}
	case int32:
		f.grb = C.GxB_SECONDI_INT32
	case int64:
		f.grb = C.GxB_SECONDI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Secondi1 is f(x, y) = row index of y (1-based)
//
// Secondi1 is a SuiteSparse:GraphBLAS extension.
func Secondi1[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_SECONDI1_INT32
		} else {
			f.grb = C.GxB_SECONDI1_INT64
		}
	case int32:
		f.grb = C.GxB_SECONDI1_INT32
	case int64:
		f.grb = C.GxB_SECONDI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Secondj is f(x, y) = column index of y (0-based)
//
// Secondj is a SuiteSparse:GraphBLAS extension.
func Secondj[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_SECONDJ_INT32
		} else {
			f.grb = C.GxB_SECONDJ_INT64
		}
	case int32:
		f.grb = C.GxB_SECONDJ_INT32
	case int64:
		f.grb = C.GxB_SECONDJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Secondj1 is f(x, y) = column index of y (1-based)
//
// Secondj1 is a SuiteSparse:GraphBLAS extension.
func Secondj1[D int32 | int64 | int, Din1, Din2 any]() (f BinaryOp[D, Din1, Din2]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GxB_SECONDJ1_INT32
		} else {
			f.grb = C.GxB_SECONDJ1_INT64
		}
	case int32:
		f.grb = C.GxB_SECONDJ1_INT32
	case int64:
		f.grb = C.GxB_SECONDJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// IgnoreDup is ignore duplicates during [Vector.Build] or [Matrix.Build].
//
// IgnoreDup is a SuiteSparse:GraphBLAS extension.
func IgnoreDup[Dout, Din1, Din2 any]() BinaryOp[Dout, Din1, Din2] {
	return BinaryOp[Dout, Din1, Din2]{C.GxB_IGNORE_DUP}
}
