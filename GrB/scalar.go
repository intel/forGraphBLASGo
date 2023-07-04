package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

// A Scalar is defined by a domain D, and a set of zero or one scalar value.
type Scalar[D any] struct {
	grb C.GrB_Scalar
}

// ScalarView returns a view on the given scalar (with domain From) using a different domain To.
//
// In the GraphBLAS specification for the C programming language, collections (scalars, vectors and matrices) of
// [Predefined] domains can be arbitrarily intermixed. In SuiteSparse:GraphBLAS, this extends to collections of [Complex]
// domains. When entries of collections are accessed expecting a particular domain (type), then
// the entry values are typecast using the rules of the C programming language. (Collections of
// user-defined domains are not compatible with any other collections in this way.)
//
// In Go, generally only identical types are compatible with each other, and conversions are
// not implicit. To get around this restriction, ScalarView, [VectorView] and [MatrixView] can be used to view a
// collection using a different domain. These functions do not perform any conversion themselves, but are essentially
// NO-OPs.
//
// ScalarView is a forGraphBLASGo extension.
func ScalarView[To, From Predefined | Complex](scalar Scalar[From]) (view Scalar[To]) {
	view.grb = scalar.grb
	return
}

// Type returns the actual [Type] object representing the domain of the given scalar.
// This is not necessarily the [Type] object corresponding to D, if Type is called
// on a [ScalarView] of a matrix of some other domain.
//
// Type might return false as a second return value if the domain is not a [Predefined]
// or [Complex] domain, or if the type has not been registered with [TypeNew] or
// [NamedTypeNew].
//
// Type is a forGraphBLASGo extension. It can be used in place of GxB_Scalar_type_name
// and GxB_Type_from_name, which are SuiteSparse:GraphBLAS extensions.
func (scalar Scalar[D]) Type() (typ Type, ok bool, err error) {
	var ctypename [C.GxB_MAX_NAME_LEN]C.char
	info := Info(C.GxB_Scalar_type_name(&ctypename[0], scalar.grb))
	if info != success {
		err = makeError(info)
		return
	}
	var grb C.GrB_Type
	info = Info(C.GxB_Type_from_name(&grb, &ctypename[0]))
	if info != success {
		err = makeError(info)
	}
	typ, ok = goType[grb]
	return
}

// ScalarNew creates a new matrix with specified domain.
//
// Parameters:
//
//   - D: The type corresponding to the domain of the matrix being created.
//     Can be one of the [Predefined] or [Complex] types, or an existing
//     user-defined GraphBLAS type.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func ScalarNew[D any]() (scalar Scalar[D], err error) {
	var d D
	dt, ok := grbType[TypeOf(d)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	info := Info(C.GrB_Scalar_new(&scalar.grb, dt))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Dup creates a new scalar with the same domain and content as another scalar.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (scalar Scalar[D]) Dup() (dup Scalar[D], err error) {
	info := Info(C.GrB_Scalar_dup(&dup.grb, scalar.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Clear removes the stored element from a scalar.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (scalar Scalar[D]) Clear() error {
	info := Info(C.GrB_Scalar_clear(scalar.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Nvals retrieves the number of stored elements in a scalar (either zero or one).
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (scalar Scalar[D]) Nvals() (nvals int, err error) {
	var cnvals C.GrB_Index
	info := Info(C.GrB_Scalar_nvals(&cnvals, scalar.grb))
	if info == success {
		return int(cnvals), nil
	}
	err = makeError(info)
	return
}

// SetElement sets the single element of a scalar to a given value.
//
// Parameters:
//
//   - val (IN): Scalar to assign.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (scalar Scalar[D]) SetElement(val D) error {
	var info Info
	switch value := any(val).(type) {
	case bool:
		info = Info(C.GrB_Scalar_setElement_BOOL(scalar.grb, C.bool(value)))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Scalar_setElement_INT32(scalar.grb, C.int32_t(value)))
		} else {
			info = Info(C.GrB_Scalar_setElement_INT64(scalar.grb, C.int64_t(value)))
		}
	case int8:
		info = Info(C.GrB_Scalar_setElement_INT8(scalar.grb, C.int8_t(value)))
	case int16:
		info = Info(C.GrB_Scalar_setElement_INT16(scalar.grb, C.int16_t(value)))
	case int32:
		info = Info(C.GrB_Scalar_setElement_INT32(scalar.grb, C.int32_t(value)))
	case int64:
		info = Info(C.GrB_Scalar_setElement_INT64(scalar.grb, C.int64_t(value)))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Scalar_setElement_UINT32(scalar.grb, C.uint32_t(value)))
		} else {
			info = Info(C.GrB_Scalar_setElement_UINT64(scalar.grb, C.uint64_t(value)))
		}
	case uint8:
		info = Info(C.GrB_Scalar_setElement_UINT8(scalar.grb, C.uint8_t(value)))
	case uint16:
		info = Info(C.GrB_Scalar_setElement_UINT16(scalar.grb, C.uint16_t(value)))
	case uint32:
		info = Info(C.GrB_Scalar_setElement_UINT32(scalar.grb, C.uint32_t(value)))
	case uint64:
		info = Info(C.GrB_Scalar_setElement_UINT64(scalar.grb, C.uint64_t(value)))
	case float32:
		info = Info(C.GrB_Scalar_setElement_FP32(scalar.grb, C.float(value)))
	case float64:
		info = Info(C.GrB_Scalar_setElement_FP64(scalar.grb, C.double(value)))
	case complex64:
		info = Info(C.GxB_Scalar_setElement_FC32(scalar.grb, C.complexfloat(value)))
	case complex128:
		info = Info(C.GxB_Scalar_setElement_FC64(scalar.grb, C.complexdouble(value)))
	default:
		info = Info(C.GrB_Scalar_setElement_UDT(scalar.grb, unsafe.Pointer(&val)))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// ExtractElement extracts the single element of a scalar.
//
// When there is no stored value, ExtractElement returns
// ok == false. Otherwise, it returns ok == true.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (scalar Scalar[D]) ExtractElement() (result D, ok bool, err error) {
	var info Info
	switch res := any(&result).(type) {
	case *bool:
		var cresult C.bool
		info = Info(C.GrB_Scalar_extractElement_BOOL(&cresult, scalar.grb))
		if info == success {
			*res = bool(cresult)
			ok = true
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var cresult C.int32_t
			info = Info(C.GrB_Scalar_extractElement_INT32(&cresult, scalar.grb))
			if info == success {
				*res = int(cresult)
				ok = true
				return
			}
		} else {
			var cresult C.int64_t
			info = Info(C.GrB_Scalar_extractElement_INT64(&cresult, scalar.grb))
			if info == success {
				*res = int(cresult)
				ok = true
				return
			}
		}
	case *int8:
		var cresult C.int8_t
		info = Info(C.GrB_Scalar_extractElement_INT8(&cresult, scalar.grb))
		if info == success {
			*res = int8(cresult)
			ok = true
			return
		}
	case *int16:
		var cresult C.int16_t
		info = Info(C.GrB_Scalar_extractElement_INT16(&cresult, scalar.grb))
		if info == success {
			*res = int16(cresult)
			ok = true
			return
		}
	case *int32:
		var cresult C.int32_t
		info = Info(C.GrB_Scalar_extractElement_INT32(&cresult, scalar.grb))
		if info == success {
			*res = int32(cresult)
			ok = true
			return
		}
	case *int64:
		var cresult C.int64_t
		info = Info(C.GrB_Scalar_extractElement_INT64(&cresult, scalar.grb))
		if info == success {
			*res = int64(cresult)
			ok = true
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var cresult C.uint32_t
			info = Info(C.GrB_Scalar_extractElement_UINT32(&cresult, scalar.grb))
			if info == success {
				*res = uint(cresult)
				ok = true
				return
			}
		} else {
			var cresult C.uint64_t
			info = Info(C.GrB_Scalar_extractElement_UINT64(&cresult, scalar.grb))
			if info == success {
				*res = uint(cresult)
				ok = true
				return
			}
		}
	case *uint8:
		var cresult C.uint8_t
		info = Info(C.GrB_Scalar_extractElement_UINT8(&cresult, scalar.grb))
		if info == success {
			*res = uint8(cresult)
			ok = true
			return
		}
	case *uint16:
		var cresult C.uint16_t
		info = Info(C.GrB_Scalar_extractElement_UINT16(&cresult, scalar.grb))
		if info == success {
			*res = uint16(cresult)
			ok = true
			return
		}
	case *uint32:
		var cresult C.uint32_t
		info = Info(C.GrB_Scalar_extractElement_UINT32(&cresult, scalar.grb))
		if info == success {
			*res = uint32(cresult)
			ok = true
			return
		}
	case *uint64:
		var cresult C.uint64_t
		info = Info(C.GrB_Scalar_extractElement_UINT64(&cresult, scalar.grb))
		if info == success {
			*res = uint64(cresult)
			ok = true
			return
		}
	case *float32:
		var cresult C.float
		info = Info(C.GrB_Scalar_extractElement_FP32(&cresult, scalar.grb))
		if info == success {
			*res = float32(cresult)
			ok = true
			return
		}
	case *float64:
		var cresult C.double
		info = Info(C.GrB_Scalar_extractElement_FP64(&cresult, scalar.grb))
		if info == success {
			*res = float64(cresult)
			ok = true
			return
		}
	case *complex64:
		var cresult C.complexfloat
		info = Info(C.GxB_Scalar_extractElement_FC32(&cresult, scalar.grb))
		if info == success {
			*res = complex64(cresult)
			ok = true
			return
		}
	case *complex128:
		var cresult C.complexdouble
		info = Info(C.GxB_Scalar_extractElement_FC64(&cresult, scalar.grb))
		if info == success {
			*res = complex128(cresult)
			ok = true
			return
		}
	default:
		info = Info(C.GrB_Scalar_extractElement_UDT(unsafe.Pointer(&result), scalar.grb))
		if info == success {
			ok = true
			return
		}
	}
	if info == noValue {
		return
	}
	err = makeError(info)
	return
}

// MemoryUsage returns the memory space required for a scalar, in bytes.
//
// MemoryUsage is a SuiteSparse:GraphBLAS extension.
func (scalar Scalar[D]) MemoryUsage() (size int, err error) {
	var csize C.size_t
	info := Info(C.GxB_Scalar_memoryUsage(&csize, scalar.grb))
	if info == success {
		return int(csize), nil
	}
	err = makeError(info)
	return
}

// Valid returns true if scalar has been created by a successful call to [Scalar.Dup],
// or [ScalarNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (scalar Scalar[D]) Valid() bool {
	return scalar.grb != C.GrB_Scalar(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Scalar] and releases any resources associated with
// it. Calling Free on an object that is not [Scalar.Valid]() is legal.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (scalar *Scalar[D]) Free() error {
	info := Info(C.GrB_Scalar_free(&scalar.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the scalar into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (scalar Scalar[D]) Wait(mode WaitMode) error {
	info := Info(C.GrB_Scalar_wait(scalar.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the scalar.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (scalar Scalar[D]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_Scalar_error(&cerror, scalar.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the scalar to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: scalar is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (scalar Scalar[D]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Scalar_fprint(scalar.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}
