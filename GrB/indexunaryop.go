package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

type (
	// IndexUnaryFunction is the C type for index unary function:
	//    typedef void (*GxB_index_unary_function) (void *, const void *, GrB_Index, GrB_Index, const void *) ;
	IndexUnaryFunction C.GxB_index_unary_function

	// IndexUnaryOp represents a GraphBLAS function that takes arguments of type Din1, GrB_Index, GrB_Index, and Din2,
	// and returns an argument of type Dout.
	IndexUnaryOp[Dout, Din1, Din2 any] struct {
		grb C.GrB_IndexUnaryOp
	}
)

// IndexUnaryOpNew returns a new GraphBLAS index unary operator with a specified user-defined function in C and its types (domains).
//
// Parameters:
//   - indexUnaryFunc (IN): A pointer to a user-defined function in C that takes an input parameters of type Din1,
//     GrB_Index, GrB_Index and Din2, and returns a value of type Dout. Except for the GrB_Index parameters,
//     all are passed as void pointers. Dout, Din1, and Din2 should be one of the [Predefined] GraphBLAS types,
//     one of the [Complex] GraphBLAS types, or a user-defined GraphBLAS type.
//
// GraphBLAS API errors that may be returned:
//   - [NullPointer], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func IndexUnaryOpNew[Dout, Din1, Din2 any](indexUnaryFunc IndexUnaryFunction) (indexUnaryOp IndexUnaryOp[Dout, Din1, Din2], err error) {
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
	info := Info(C.GrB_IndexUnaryOp_new(&indexUnaryOp.grb, indexUnaryFunc, doutt, din1t, din2t))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// NamedIndexUnaryOpNew creates a named index unary function. It is like [IndexUnaryOpNew], except:
//   - idxopname is the name for the GraphBLAS index unary operator. Only the first 127 characters are used.
//   - idxopdefn is a string containing the entire function itself.
//
// The two strings idxopname and idxopdefn are optional, but are required to enable the JIT compilation
// of kernels that use this operator.
//
// If the JIT is enabled, or if the corresponding JIT kernel has been copied into the PreJIT folder,
// the function may be nil. In this case, a JIT kernel is compiled that contains just the user-defined
// function. If the JIT is disabled and the function is nil, this method panics with a [NullPointer] error.
//
// NamedIndexUnaryOpNew is a SuiteSparse:GraphBLAS extension.
func NamedIndexUnaryOpNew[Dout, Din1, Din2 any](indexUnaryFunc IndexUnaryFunction, idxopname string, idxopdefn string) (indexUnaryOp IndexUnaryOp[Dout, Din1, Din2], err error) {
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
	cidxopname := C.CString(idxopname)
	defer C.free(unsafe.Pointer(cidxopname))
	cidxopdefn := C.CString(idxopdefn)
	defer C.free(unsafe.Pointer(cidxopdefn))
	info := Info(C.GxB_IndexUnaryOp_new(&indexUnaryOp.grb, indexUnaryFunc, doutt, din1t, din2t, cidxopname, cidxopdefn))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Valid returns true if indexUnaryOp has been created by a successful call to [IndexUnaryOpNew] or [NamedIndexUnaryOpNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (indexUnaryOp IndexUnaryOp[Dout, Din1, Din2]) Valid() bool {
	return indexUnaryOp.grb != C.GrB_IndexUnaryOp(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [IndexUnaryOp] and releases any resources associated with
// it. Calling Free on an object that is not [IndexUnaryOp.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined index unary operator is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (indexUnaryOp *IndexUnaryOp[Dout, Din1, Din2]) Free() error {
	info := Info(C.GrB_IndexUnaryOp_free(&indexUnaryOp.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the index unary operator into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (indexUnaryOp IndexUnaryOp[Dout, Din1, Din2]) Wait(mode WaitMode) error {
	info := Info(C.GrB_IndexUnaryOp_wait(indexUnaryOp.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the index unary operator.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (indexUnaryOp IndexUnaryOp[Dout, Din1, Din2]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_IndexUnaryOp_error(&cerror, indexUnaryOp.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the index unary operator to stdout.
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
func (indexUnaryOp IndexUnaryOp[Dout, Din1, Din2]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_IndexUnaryOp_fprint(indexUnaryOp.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// RowIndex is
//   - for matrices: f(a(i, j), i, j, s) = i + s
//   - for vectors:  f(u(i), i, 0, s) = i + s
func RowIndex[D int32 | int64 | int, Any any]() (f IndexUnaryOp[D, Any, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_ROWINDEX_INT32
		} else {
			f.grb = C.GrB_ROWINDEX_INT64
		}
	case int32:
		f.grb = C.GrB_ROWINDEX_INT32
	case int64:
		f.grb = C.GrB_ROWINDEX_INT64
	default:
		panic("unreachable code")
	}
	return
}

// ColIndex is f(a(i, j), i, j, s) = j + s
func ColIndex[D int32 | int64 | int, Any any]() (f IndexUnaryOp[D, Any, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_COLINDEX_INT32
		} else {
			f.grb = C.GrB_COLINDEX_INT64
		}
	case int32:
		f.grb = C.GrB_COLINDEX_INT32
	case int64:
		f.grb = C.GrB_COLINDEX_INT64
	default:
		panic("unreachable code")
	}
	return
}

// DiagIndex is f(a(i, j), i, j, s) = j - i + s
func DiagIndex[D int32 | int64 | int, Any any]() (f IndexUnaryOp[D, Any, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_DIAGINDEX_INT32
		} else {
			f.grb = C.GrB_DIAGINDEX_INT64
		}
	case int32:
		f.grb = C.GrB_DIAGINDEX_INT32
	case int64:
		f.grb = C.GrB_DIAGINDEX_INT64
	default:
		panic("unreachable code")
	}
	return
}

// Tril is f(a(i, j), i, j, s) = j <= i + s
func Tril[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_TRIL}
}

// Triu is f(a(i, j), i, j, s) = j >= i + s
func Triu[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_TRIU}
}

// Diag is f(a(i, j), i, j, s) = j == i + s
func Diag[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_DIAG}
}

// Offdiag is f(a(i, j), i, j, s) = j != i + s
func Offdiag[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_OFFDIAG}
}

// Colle is f(a(i, j), i, j, s) = j <= s
func Colle[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_COLLE}
}

// Colgt is f(a(i, j), i, j, s) = j > s
func Colgt[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_COLGT}
}

// Rowle is
//   - for matrices: f(a(i, j), i, j, s) = i <= s
//   - for vectors:  f(u(i), i, 0, s) = i <= s
func Rowle[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_ROWLE}
}

// Rowgt is
//   - for matrices: f(a(i, j), i, j, s) = i > s
//   - for vectors:  f(u(i), i, 0, s) = i > s
func Rowgt[Any any]() IndexUnaryOp[bool, Any, int64] {
	return IndexUnaryOp[bool, Any, int64]{C.GrB_ROWGT}
}

// Valuene is
//   - for matrices: f(a(i, j), i, j, s) = a(i, j) != s
//   - for vectors:  f(u(i), i, 0, s) = u(i) != s
func Valuene[D Predefined | Complex]() (f IndexUnaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_VALUENE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUENE_INT32
		} else {
			f.grb = C.GrB_VALUENE_INT64
		}
	case int8:
		f.grb = C.GrB_VALUENE_INT8
	case int16:
		f.grb = C.GrB_VALUENE_INT16
	case int32:
		f.grb = C.GrB_VALUENE_INT32
	case int64:
		f.grb = C.GrB_VALUENE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUENE_UINT32
		} else {
			f.grb = C.GrB_VALUENE_UINT64
		}
	case uint8:
		f.grb = C.GrB_VALUENE_UINT8
	case uint16:
		f.grb = C.GrB_VALUENE_UINT16
	case uint32:
		f.grb = C.GrB_VALUENE_UINT32
	case uint64:
		f.grb = C.GrB_VALUENE_UINT64
	case float32:
		f.grb = C.GrB_VALUENE_FP32
	case float64:
		f.grb = C.GrB_VALUENE_FP64
	case complex64:
		f.grb = C.GxB_VALUENE_FC32
	case complex128:
		f.grb = C.GxB_VALUENE_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Valueeq is
//   - for matrices: f(a(i, j), i, j, s) = a(i, j) == s
//   - for vectors:  f(u(i), i, 0, s) = u(i) == s
func Valueeq[D Predefined | Complex]() (f IndexUnaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_VALUEEQ_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUEEQ_INT32
		} else {
			f.grb = C.GrB_VALUEEQ_INT64
		}
	case int8:
		f.grb = C.GrB_VALUEEQ_INT8
	case int16:
		f.grb = C.GrB_VALUEEQ_INT16
	case int32:
		f.grb = C.GrB_VALUEEQ_INT32
	case int64:
		f.grb = C.GrB_VALUEEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUEEQ_UINT32
		} else {
			f.grb = C.GrB_VALUEEQ_UINT64
		}
	case uint8:
		f.grb = C.GrB_VALUEEQ_UINT8
	case uint16:
		f.grb = C.GrB_VALUEEQ_UINT16
	case uint32:
		f.grb = C.GrB_VALUEEQ_UINT32
	case uint64:
		f.grb = C.GrB_VALUEEQ_UINT64
	case float32:
		f.grb = C.GrB_VALUEEQ_FP32
	case float64:
		f.grb = C.GrB_VALUEEQ_FP64
	case complex64:
		f.grb = C.GxB_VALUEEQ_FC32
	case complex128:
		f.grb = C.GxB_VALUEEQ_FC64
	default:
		panic("unreachable code")
	}
	return
}

// Valuegt is
//   - for matrices: f(a(i, j), i, j, s) = a(i, j) > s
//   - for vectors:  f(u(i), i, 0, s) = u(i) > s
func Valuegt[D Predefined]() (f IndexUnaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_VALUEGT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUEGT_INT32
		} else {
			f.grb = C.GrB_VALUEGT_INT64
		}
	case int8:
		f.grb = C.GrB_VALUEGT_INT8
	case int16:
		f.grb = C.GrB_VALUEGT_INT16
	case int32:
		f.grb = C.GrB_VALUEGT_INT32
	case int64:
		f.grb = C.GrB_VALUEGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUEGT_UINT32
		} else {
			f.grb = C.GrB_VALUEGT_UINT64
		}
	case uint8:
		f.grb = C.GrB_VALUEGT_UINT8
	case uint16:
		f.grb = C.GrB_VALUEGT_UINT16
	case uint32:
		f.grb = C.GrB_VALUEGT_UINT32
	case uint64:
		f.grb = C.GrB_VALUEGT_UINT64
	case float32:
		f.grb = C.GrB_VALUEGT_FP32
	case float64:
		f.grb = C.GrB_VALUEGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Valuege is
//   - for matrices: f(a(i, j), i, j, s) = a(i, j) >= s
//   - for vectors:  f(u(i), i, 0, s) = u(i) >= s
func Valuege[D Predefined]() (f IndexUnaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_VALUEGE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUEGE_INT32
		} else {
			f.grb = C.GrB_VALUEGE_INT64
		}
	case int8:
		f.grb = C.GrB_VALUEGE_INT8
	case int16:
		f.grb = C.GrB_VALUEGE_INT16
	case int32:
		f.grb = C.GrB_VALUEGE_INT32
	case int64:
		f.grb = C.GrB_VALUEGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUEGE_UINT32
		} else {
			f.grb = C.GrB_VALUEGE_UINT64
		}
	case uint8:
		f.grb = C.GrB_VALUEGE_UINT8
	case uint16:
		f.grb = C.GrB_VALUEGE_UINT16
	case uint32:
		f.grb = C.GrB_VALUEGE_UINT32
	case uint64:
		f.grb = C.GrB_VALUEGE_UINT64
	case float32:
		f.grb = C.GrB_VALUEGE_FP32
	case float64:
		f.grb = C.GrB_VALUEGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Valuelt is
//   - for matrices: f(a(i, j), i, j, s) = a(i, j) < s
//   - for vectors:  f(u(i), i, 0, s) = u(i) < s
func Valuelt[D Predefined]() (f IndexUnaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_VALUELT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUELT_INT32
		} else {
			f.grb = C.GrB_VALUELT_INT64
		}
	case int8:
		f.grb = C.GrB_VALUELT_INT8
	case int16:
		f.grb = C.GrB_VALUELT_INT16
	case int32:
		f.grb = C.GrB_VALUELT_INT32
	case int64:
		f.grb = C.GrB_VALUELT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUELT_UINT32
		} else {
			f.grb = C.GrB_VALUELT_UINT64
		}
	case uint8:
		f.grb = C.GrB_VALUELT_UINT8
	case uint16:
		f.grb = C.GrB_VALUELT_UINT16
	case uint32:
		f.grb = C.GrB_VALUELT_UINT32
	case uint64:
		f.grb = C.GrB_VALUELT_UINT64
	case float32:
		f.grb = C.GrB_VALUELT_FP32
	case float64:
		f.grb = C.GrB_VALUELT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// Valuele is
//   - for matrices: f(a(i, j), i, j, s) = a(i, j) < s
//   - for vectors:  f(u(i), i, 0, s) = u(i) < s
func Valuele[D Predefined]() (f IndexUnaryOp[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		f.grb = C.GrB_VALUELE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUELE_INT32
		} else {
			f.grb = C.GrB_VALUELE_INT64
		}
	case int8:
		f.grb = C.GrB_VALUELE_INT8
	case int16:
		f.grb = C.GrB_VALUELE_INT16
	case int32:
		f.grb = C.GrB_VALUELE_INT32
	case int64:
		f.grb = C.GrB_VALUELE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			f.grb = C.GrB_VALUELE_UINT32
		} else {
			f.grb = C.GrB_VALUELE_UINT64
		}
	case uint8:
		f.grb = C.GrB_VALUELE_UINT8
	case uint16:
		f.grb = C.GrB_VALUELE_UINT16
	case uint32:
		f.grb = C.GrB_VALUELE_UINT32
	case uint64:
		f.grb = C.GrB_VALUELE_UINT64
	case float32:
		f.grb = C.GrB_VALUELE_FP32
	case float64:
		f.grb = C.GrB_VALUELE_FP64
	default:
		panic("unreachable code")
	}
	return
}
