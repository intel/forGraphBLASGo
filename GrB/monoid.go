package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// A Monoid is defined by a single domain D, an associative operation and an identity element.
// A GraphBLAS monoid is equivalent to the conventional monoid algebraic structure.
type Monoid[D any] struct {
	grb C.GrB_Monoid
}

// MonoidNew creates a new monoid with specified binary operator and identity value.
//
// Parameters:
//
//   - binaryOp (IN): An existing GraphBLAS associative binary operator whose input and
//     output types are the same.
//
//   - identity (IN): The value of the identity element of the monoid.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [NullPointer], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func MonoidNew[D any](binaryOp BinaryOp[D, D, D], identity D) (monoid Monoid[D], err error) {
	var info Info
	switch id := any(identity).(type) {
	case bool:
		info = Info(C.GrB_Monoid_new_BOOL(&monoid.grb, binaryOp.grb, C.bool(id)))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Monoid_new_INT32(&monoid.grb, binaryOp.grb, C.int32_t(id)))
		} else {
			info = Info(C.GrB_Monoid_new_INT64(&monoid.grb, binaryOp.grb, C.int64_t(id)))
		}
	case int8:
		info = Info(C.GrB_Monoid_new_INT8(&monoid.grb, binaryOp.grb, C.int8_t(id)))
	case int16:
		info = Info(C.GrB_Monoid_new_INT16(&monoid.grb, binaryOp.grb, C.int16_t(id)))
	case int32:
		info = Info(C.GrB_Monoid_new_INT32(&monoid.grb, binaryOp.grb, C.int32_t(id)))
	case int64:
		info = Info(C.GrB_Monoid_new_INT64(&monoid.grb, binaryOp.grb, C.int64_t(id)))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Monoid_new_UINT32(&monoid.grb, binaryOp.grb, C.uint32_t(id)))
		} else {
			info = Info(C.GrB_Monoid_new_UINT64(&monoid.grb, binaryOp.grb, C.uint64_t(id)))
		}
	case uint8:
		info = Info(C.GrB_Monoid_new_UINT8(&monoid.grb, binaryOp.grb, C.uint8_t(id)))
	case uint16:
		info = Info(C.GrB_Monoid_new_UINT16(&monoid.grb, binaryOp.grb, C.uint16_t(id)))
	case uint32:
		info = Info(C.GrB_Monoid_new_UINT32(&monoid.grb, binaryOp.grb, C.uint32_t(id)))
	case uint64:
		info = Info(C.GrB_Monoid_new_UINT64(&monoid.grb, binaryOp.grb, C.uint64_t(id)))
	case float32:
		info = Info(C.GrB_Monoid_new_FP32(&monoid.grb, binaryOp.grb, C.float(id)))
	case float64:
		info = Info(C.GrB_Monoid_new_FP64(&monoid.grb, binaryOp.grb, C.double(id)))
	case complex64:
		info = Info(C.GxB_Monoid_new_FC32(&monoid.grb, binaryOp.grb, C.complexfloat(id)))
	case complex128:
		info = Info(C.GxB_Monoid_new_FC64(&monoid.grb, binaryOp.grb, C.complexdouble(id)))
	default:
		info = Info(C.GrB_Monoid_new_UDT(&monoid.grb, binaryOp.grb, unsafe.Pointer(&identity)))
	}
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// MonoidTerminalNew is identical to [MonoidNew], except that it allows for the specification of
// a terminal value.
//
// The terminal value of a monoid is the value z for which z = f(z, y) for any y, where f is the
// binary operator of the monoid. If the terminal value is encountered during computation,
// the rest of the computations are skipped. This can greatly improve the performance of
// reductions, and matrix multiply in specific cases.
//
// MonoidTerminalNew is a SuiteSparse:GraphBLAS extension.
func MonoidTerminalNew[D any](binaryOp BinaryOp[D, D, D], identity, terminal D) (monoid Monoid[D], err error) {
	var info Info
	switch id := any(identity).(type) {
	case bool:
		info = Info(C.GxB_Monoid_terminal_new_BOOL(&monoid.grb, binaryOp.grb, C.bool(id), C.bool(any(terminal).(bool))))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GxB_Monoid_terminal_new_INT32(&monoid.grb, binaryOp.grb, C.int32_t(id), C.int32_t(any(terminal).(int32))))
		} else {
			info = Info(C.GxB_Monoid_terminal_new_INT64(&monoid.grb, binaryOp.grb, C.int64_t(id), C.int64_t(any(terminal).(int64))))
		}
	case int8:
		info = Info(C.GxB_Monoid_terminal_new_INT8(&monoid.grb, binaryOp.grb, C.int8_t(id), C.int8_t(any(terminal).(int8))))
	case int16:
		info = Info(C.GxB_Monoid_terminal_new_INT16(&monoid.grb, binaryOp.grb, C.int16_t(id), C.int16_t(any(terminal).(int16))))
	case int32:
		info = Info(C.GxB_Monoid_terminal_new_INT32(&monoid.grb, binaryOp.grb, C.int32_t(id), C.int32_t(any(terminal).(int32))))
	case int64:
		info = Info(C.GxB_Monoid_terminal_new_INT64(&monoid.grb, binaryOp.grb, C.int64_t(id), C.int64_t(any(terminal).(int64))))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GxB_Monoid_terminal_new_UINT32(&monoid.grb, binaryOp.grb, C.uint32_t(id), C.uint32_t(any(terminal).(uint32))))
		} else {
			info = Info(C.GxB_Monoid_terminal_new_UINT64(&monoid.grb, binaryOp.grb, C.uint64_t(id), C.uint64_t(any(terminal).(uint64))))
		}
	case uint8:
		info = Info(C.GxB_Monoid_terminal_new_UINT8(&monoid.grb, binaryOp.grb, C.uint8_t(id), C.uint8_t(any(terminal).(uint8))))
	case uint16:
		info = Info(C.GxB_Monoid_terminal_new_UINT16(&monoid.grb, binaryOp.grb, C.uint16_t(id), C.uint16_t(any(terminal).(uint16))))
	case uint32:
		info = Info(C.GxB_Monoid_terminal_new_UINT32(&monoid.grb, binaryOp.grb, C.uint32_t(id), C.uint32_t(any(terminal).(uint32))))
	case uint64:
		info = Info(C.GxB_Monoid_terminal_new_UINT64(&monoid.grb, binaryOp.grb, C.uint64_t(id), C.uint64_t(any(terminal).(uint64))))
	case float32:
		info = Info(C.GxB_Monoid_terminal_new_FP32(&monoid.grb, binaryOp.grb, C.float(id), C.float(any(terminal).(float32))))
	case float64:
		info = Info(C.GxB_Monoid_terminal_new_FP64(&monoid.grb, binaryOp.grb, C.double(id), C.double(any(terminal).(float64))))
	case complex64:
		info = Info(C.GxB_Monoid_terminal_new_FC32(&monoid.grb, binaryOp.grb, C.complexfloat(id), C.complexfloat(any(terminal).(complex64))))
	case complex128:
		info = Info(C.GxB_Monoid_terminal_new_FC64(&monoid.grb, binaryOp.grb, C.complexdouble(id), C.complexdouble(any(terminal).(complex128))))
	default:
		info = Info(C.GxB_Monoid_terminal_new_UDT(&monoid.grb, binaryOp.grb, unsafe.Pointer(&identity), unsafe.Pointer(&terminal)))
	}
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Operator returns the binary operator of the monoid.
//
// Operator is a SuiteSparse:GraphBLAS extension.
func (monoid Monoid[D]) Operator() (binaryOp BinaryOp[D, D, D], err error) {
	info := Info(C.GxB_Monoid_operator(&binaryOp.grb, monoid.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Identity returns the identity value of the monoid.
//
// Identity is a SuiteSparse:GraphBLAS extension.
func (monoid Monoid[D]) Identity() (identity D, err error) {
	var info Info
	switch id := any(&identity).(type) {
	case *bool:
		var x C.bool
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = bool(x)
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var x C.int32_t
			info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*id = int(x)
				return
			}
		} else {
			var x C.int64_t
			info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*id = int(x)
				return
			}
		}
	case *int8:
		var x C.int8_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = int8(x)
			return
		}
	case *int16:
		var x C.int16_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = int16(x)
			return
		}
	case *int32:
		var x C.int32_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = int32(x)
			return
		}
	case *int64:
		var x C.int64_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = int64(x)
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var x C.uint32_t
			info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*id = uint(x)
				return
			}
		} else {
			var x C.uint64_t
			info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*id = uint(x)
				return
			}
		}
	case *uint8:
		var x C.uint8_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = uint8(x)
			return
		}
	case *uint16:
		var x C.uint16_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = uint16(x)
			return
		}
	case *uint32:
		var x C.uint32_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = uint32(x)
			return
		}
	case *uint64:
		var x C.uint64_t
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = uint64(x)
			return
		}
	case *float32:
		var x C.float
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = float32(x)
			return
		}
	case *float64:
		var x C.double
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = float64(x)
			return
		}
	case *complex64:
		var x C.complexfloat
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = complex64(x)
			return
		}
	case *complex128:
		var x C.complexdouble
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*id = complex128(x)
			return
		}
	default:
		info = Info(C.GxB_Monoid_identity(unsafe.Pointer(&identity), monoid.grb))
		if info == success {
			return
		}
	}
	err = makeError(info)
	return
}

// Terminal returns the terminal value of the monoid, if any.
//
// If the monoid has a terminal value, then ok == true. If it has
// no terminal value, then ok == false.
//
// Terminal is a SuiteSparse:GraphBLAS extension.
func (monoid Monoid[D]) Terminal() (terminal D, ok bool, err error) {
	var hasTerminal C.bool
	var info Info
	switch term := any(&terminal).(type) {
	case *bool:
		var x C.bool
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = bool(x)
			ok = bool(hasTerminal)
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var x C.int32_t
			info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*term = int(x)
				ok = bool(hasTerminal)
				return
			}
		} else {
			var x C.int64_t
			info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*term = int(x)
				ok = bool(hasTerminal)
				return
			}
		}
	case *int8:
		var x C.int8_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = int8(x)
			ok = bool(hasTerminal)
			return
		}
	case *int16:
		var x C.int16_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = int16(x)
			ok = bool(hasTerminal)
			return
		}
	case *int32:
		var x C.int32_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = int32(x)
			ok = bool(hasTerminal)
			return
		}
	case *int64:
		var x C.int64_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = int64(x)
			ok = bool(hasTerminal)
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var x C.uint32_t
			info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*term = uint(x)
				ok = bool(hasTerminal)
				return
			}
		} else {
			var x C.uint64_t
			info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
			if info == success {
				*term = uint(x)
				ok = bool(hasTerminal)
				return
			}
		}
	case *uint8:
		var x C.uint8_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = uint8(x)
			ok = bool(hasTerminal)
			return
		}
	case *uint16:
		var x C.uint16_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = uint16(x)
			ok = bool(hasTerminal)
			return
		}
	case *uint32:
		var x C.uint32_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = uint32(x)
			ok = bool(hasTerminal)
			return
		}
	case *uint64:
		var x C.uint64_t
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = uint64(x)
			ok = bool(hasTerminal)
			return
		}
	case *float32:
		var x C.float
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = float32(x)
			ok = bool(hasTerminal)
			return
		}
	case *float64:
		var x C.double
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = float64(x)
			ok = bool(hasTerminal)
			return
		}
	case *complex64:
		var x C.complexfloat
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = complex64(x)
			ok = bool(hasTerminal)
			return
		}
	case *complex128:
		var x C.complexdouble
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&x), monoid.grb))
		if info == success {
			*term = complex128(x)
			ok = bool(hasTerminal)
			return
		}
	default:
		info = Info(C.GxB_Monoid_terminal(&hasTerminal, unsafe.Pointer(&terminal), monoid.grb))
		if info == success {
			ok = bool(hasTerminal)
			return
		}
	}
	err = makeError(info)
	return
}

// Valid returns true if monoid has been created by a successful call to [MonoidNew] or [MonoidTerminalNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (monoid Monoid[D]) Valid() bool {
	return monoid.grb != C.GrB_Monoid(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Monoid] and releases any resources associated with
// it. Calling Free on an object that is not [Monoid.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined monoid is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (monoid *Monoid[D]) Free() error {
	info := Info(C.GrB_Monoid_free(&monoid.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the monoid into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (monoid Monoid[D]) Wait(mode WaitMode) error {
	info := Info(C.GrB_Monoid_wait(monoid.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the monoid.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (monoid Monoid[D]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_Monoid_error(&cerror, monoid.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the monoid to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: monoid is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (monoid Monoid[D]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Monoid_fprint(monoid.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// PlusMonoid is addition with identity 0
func PlusMonoid[D Number | Complex]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_PLUS_MONOID_INT32
		} else {
			m.grb = C.GrB_PLUS_MONOID_INT64
		}
	case int8:
		m.grb = C.GrB_PLUS_MONOID_INT8
	case int16:
		m.grb = C.GrB_PLUS_MONOID_INT16
	case int32:
		m.grb = C.GrB_PLUS_MONOID_INT32
	case int64:
		m.grb = C.GrB_PLUS_MONOID_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_PLUS_MONOID_UINT32
		} else {
			m.grb = C.GrB_PLUS_MONOID_UINT64
		}
	case uint8:
		m.grb = C.GrB_PLUS_MONOID_UINT8
	case uint16:
		m.grb = C.GrB_PLUS_MONOID_UINT16
	case uint32:
		m.grb = C.GrB_PLUS_MONOID_UINT32
	case uint64:
		m.grb = C.GrB_PLUS_MONOID_UINT64
	case float32:
		m.grb = C.GrB_PLUS_MONOID_FP32
	case float64:
		m.grb = C.GrB_PLUS_MONOID_FP64
	case complex64:
		m.grb = C.GxB_PLUS_FC32_MONOID
	case complex128:
		m.grb = C.GxB_PLUS_FC64_MONOID
	default:
		panic("unreachable code")
	}
	return
}

// TimesMonoid is multiplication with identity 1
func TimesMonoid[D Number | Complex]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_TIMES_MONOID_INT32
		} else {
			m.grb = C.GrB_TIMES_MONOID_INT64
		}
	case int8:
		m.grb = C.GrB_TIMES_MONOID_INT8
	case int16:
		m.grb = C.GrB_TIMES_MONOID_INT16
	case int32:
		m.grb = C.GrB_TIMES_MONOID_INT32
	case int64:
		m.grb = C.GrB_TIMES_MONOID_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_TIMES_MONOID_UINT32
		} else {
			m.grb = C.GrB_TIMES_MONOID_UINT64
		}
	case uint8:
		m.grb = C.GrB_TIMES_MONOID_UINT8
	case uint16:
		m.grb = C.GrB_TIMES_MONOID_UINT16
	case uint32:
		m.grb = C.GrB_TIMES_MONOID_UINT32
	case uint64:
		m.grb = C.GrB_TIMES_MONOID_UINT64
	case float32:
		m.grb = C.GrB_TIMES_MONOID_FP32
	case float64:
		m.grb = C.GrB_TIMES_MONOID_FP64
	case complex64:
		m.grb = C.GxB_TIMES_FC32_MONOID
	case complex128:
		m.grb = C.GxB_TIMES_FC64_MONOID
	default:
		panic("unreachable code")
	}
	return
}

// MinMonoid is minimum value with identity being the maximum value of the domain
func MinMonoid[D Number]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_MIN_MONOID_INT32
		} else {
			m.grb = C.GrB_MIN_MONOID_INT64
		}
	case int8:
		m.grb = C.GrB_MIN_MONOID_INT8
	case int16:
		m.grb = C.GrB_MIN_MONOID_INT16
	case int32:
		m.grb = C.GrB_MIN_MONOID_INT32
	case int64:
		m.grb = C.GrB_MIN_MONOID_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_MIN_MONOID_UINT32
		} else {
			m.grb = C.GrB_MIN_MONOID_UINT64
		}
	case uint8:
		m.grb = C.GrB_MIN_MONOID_UINT8
	case uint16:
		m.grb = C.GrB_MIN_MONOID_UINT16
	case uint32:
		m.grb = C.GrB_MIN_MONOID_UINT32
	case uint64:
		m.grb = C.GrB_MIN_MONOID_UINT64
	case float32:
		m.grb = C.GrB_MIN_MONOID_FP32
	case float64:
		m.grb = C.GrB_MIN_MONOID_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxMonoid is maximum value with identity being the minimum value of the domain
func MaxMonoid[D Number]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_MAX_MONOID_INT32
		} else {
			m.grb = C.GrB_MAX_MONOID_INT64
		}
	case int8:
		m.grb = C.GrB_MAX_MONOID_INT8
	case int16:
		m.grb = C.GrB_MAX_MONOID_INT16
	case int32:
		m.grb = C.GrB_MAX_MONOID_INT32
	case int64:
		m.grb = C.GrB_MAX_MONOID_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GrB_MAX_MONOID_UINT32
		} else {
			m.grb = C.GrB_MAX_MONOID_UINT64
		}
	case uint8:
		m.grb = C.GrB_MAX_MONOID_UINT8
	case uint16:
		m.grb = C.GrB_MAX_MONOID_UINT16
	case uint32:
		m.grb = C.GrB_MAX_MONOID_UINT32
	case uint64:
		m.grb = C.GrB_MAX_MONOID_UINT64
	case float32:
		m.grb = C.GrB_MAX_MONOID_FP32
	case float64:
		m.grb = C.GrB_MAX_MONOID_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyMonoid is any value with identity being any value
//
// AnyMonoid is a SuiteSparse:GraphBLAS extension.
func AnyMonoid[D Number | Complex]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GxB_ANY_INT32_MONOID
		} else {
			m.grb = C.GxB_ANY_INT64_MONOID
		}
	case int8:
		m.grb = C.GxB_ANY_INT8_MONOID
	case int16:
		m.grb = C.GxB_ANY_INT16_MONOID
	case int32:
		m.grb = C.GxB_ANY_INT32_MONOID
	case int64:
		m.grb = C.GxB_ANY_INT64_MONOID
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GxB_ANY_UINT32_MONOID
		} else {
			m.grb = C.GxB_ANY_UINT64_MONOID
		}
	case uint8:
		m.grb = C.GxB_ANY_UINT8_MONOID
	case uint16:
		m.grb = C.GxB_ANY_UINT16_MONOID
	case uint32:
		m.grb = C.GxB_ANY_UINT32_MONOID
	case uint64:
		m.grb = C.GxB_ANY_UINT64_MONOID
	case float32:
		m.grb = C.GxB_ANY_FP32_MONOID
	case float64:
		m.grb = C.GxB_ANY_FP64_MONOID
	case complex64:
		m.grb = C.GxB_ANY_FC32_MONOID
	case complex128:
		m.grb = C.GxB_ANY_FC64_MONOID
	default:
		panic("unreachable code")
	}
	return
}

var (
	// LorMonoidBool is logical or with identity false
	LorMonoidBool = Monoid[bool]{C.GrB_LOR_MONOID_BOOL}

	// LandMonoidBool is logical and with identity true
	LandMonoidBool = Monoid[bool]{C.GrB_LAND_MONOID_BOOL}

	// LxorMonoidBool is logical xor (not equal) with identity false
	LxorMonoidBool = Monoid[bool]{C.GrB_LXOR_MONOID_BOOL}

	// LxnorMonoidBool is logical xnor (equal) with identity true
	LxnorMonoidBool = Monoid[bool]{C.GrB_LXNOR_MONOID_BOOL}
)

// BorMonoid is bitwise or with identity being all bits zero
//
// BorMonoid is a SuiteSparse:GraphBLAS extension.
func BorMonoid[D Unsigned]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GxB_BOR_UINT32_MONOID
		} else {
			m.grb = C.GxB_BOR_UINT64_MONOID
		}
	case uint8:
		m.grb = C.GxB_BOR_UINT8_MONOID
	case uint16:
		m.grb = C.GxB_BOR_UINT16_MONOID
	case uint32:
		m.grb = C.GxB_BOR_UINT32_MONOID
	case uint64:
		m.grb = C.GxB_BOR_UINT64_MONOID
	default:
		panic("unreachable code")
	}
	return
}

// BandMonoid is bitwise and with identity being all bits one
//
// BandMonoid is a SuiteSparse:GraphBLAS extension.
func BandMonoid[D Unsigned]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GxB_BAND_UINT32_MONOID
		} else {
			m.grb = C.GxB_BAND_UINT64_MONOID
		}
	case uint8:
		m.grb = C.GxB_BAND_UINT8_MONOID
	case uint16:
		m.grb = C.GxB_BAND_UINT16_MONOID
	case uint32:
		m.grb = C.GxB_BAND_UINT32_MONOID
	case uint64:
		m.grb = C.GxB_BAND_UINT64_MONOID
	default:
		panic("unreachable code")
	}
	return
}

// BxorMonoid is bitwise xor with identity being all bits zero
//
// BxorMonoid is a SuiteSparse:GraphBLAS extension.
func BxorMonoid[D Unsigned]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GxB_BXOR_UINT32_MONOID
		} else {
			m.grb = C.GxB_BXOR_UINT64_MONOID
		}
	case uint8:
		m.grb = C.GxB_BXOR_UINT8_MONOID
	case uint16:
		m.grb = C.GxB_BXOR_UINT16_MONOID
	case uint32:
		m.grb = C.GxB_BXOR_UINT32_MONOID
	case uint64:
		m.grb = C.GxB_BXOR_UINT64_MONOID
	default:
		panic("unreachable code")
	}
	return
}

// BxnorMonoid is bitwise xnor with identity being all bits one
//
// BxnorMonoid is a SuiteSparse:GraphBLAS extension.
func BxnorMonoid[D Unsigned]() (m Monoid[D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			m.grb = C.GxB_BXNOR_UINT32_MONOID
		} else {
			m.grb = C.GxB_BXNOR_UINT64_MONOID
		}
	case uint8:
		m.grb = C.GxB_BXNOR_UINT8_MONOID
	case uint16:
		m.grb = C.GxB_BXNOR_UINT16_MONOID
	case uint32:
		m.grb = C.GxB_BXNOR_UINT32_MONOID
	case uint64:
		m.grb = C.GxB_BXNOR_UINT64_MONOID
	default:
		panic("unreachable code")
	}
	return
}
