package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"unsafe"
)

// A Vector is defined by a domain D, a size N > 0,
// and a set up tuples (i, v(i)), where 0 <= i < N, and v(i) ∈ D. A
// particular value of i can occur at most once in v.
type Vector[D any] struct {
	grb C.GrB_Vector
}

// A VectorMask can be used to optionally control which results from a GraphBLAS operation
// are stored into an output vector. The vector size must match that of the
// output. If the [Structure] descriptor is not set for the mask, the domain
// of the mask vector must be of type bool, any of the [Predefined] "built-in" types,
// or any of the [Complex] "built-in" types. Use [Vector.AsMask] to convert
// to the required parameter type. If the default mask is desired (i.e., a mask that is all
// true with the size of the output vector), nil should be specified.
//
// The forGraphBLASGo API does not use the VectorMask type, but directly uses *Vector[bool] instead.
type VectorMask = *Vector[bool]

// VectorView returns a view on the given vector (with domain From) using a different domain To.
//
// In the GraphBLAS specification for the C programming language, collections (scalars, vectors and matrices) of
// [Predefined] domains can be arbitrarily intermixed. In SuiteSparse:GraphBLAS, this extends to collections of [Complex]
// domains. When entries of collections are accessed expecting a particular domain (type), then
// the entry values are typecast using the rules of the C programming language. (Collections of
// user-defined domains are not compatible with any other collections in this way.)
//
// In Go, generally only identical types are compatible with each other, and conversions are
// not implicit. To get around this restriction, [ScalarView], VectorView and [MatrixView] can be used to view a
// collection using a different domain. These functions do not perform any conversion themselves, but are essentially
// NO-OPs.
//
// VectorView is a forGraphBLASGo extension.
func VectorView[To, From Predefined | Complex](vector Vector[From]) (view Vector[To]) {
	view.grb = vector.grb
	return
}

// AsMask returns a view on the given vector using the domain bool.
//
// In GraphBLAS, whenever a mask is required as an input parameter for a GraphBLAS operation,
// a vector of any domain can be passed, and depending on whether [Structure] is set or not in the
// [Descriptor] passed to that operation, the only requirement is that the domain is compatible
// with bool. In the C programming language, this holds for any of the [Predefined] domains.
// In SuiteSparse:GraphBLAS, this extends to any of the [Complex] domains.
//
// In Go, generally only identical types are compatible with each other, and conversions are
// not implicit. To get around this restriction, AsMask can be used to view a vector as a bool
// mask. AsMask does not perform any conversion itself, but is essentially a NO-OP.
//
// AsMask is a forGraphBLASGo extension.
func (vector Vector[D]) AsMask() *Vector[bool] {
	return &Vector[bool]{vector.grb}
}

// Type returns the actual [Type] object representing the domain of the given vector.
// This is not necessarily the [Type] object corresponding to D, if Type is called
// on a [VectorView] of a vector of some other domain.
//
// Type might return false as a second return value if the domain is not a [Predefined]
// or [Complex] domain, or if the type has not been registered with [TypeNew] or
// [NamedTypeNew].
//
// Type is a forGraphBLASGo extension. It can be used in place of GxB_Vector_type_name
// and GxB_Type_from_name, which are SuiteSparse:GraphBLAS extensions.
func (vector Vector[D]) Type() (typ Type, ok bool, err error) {
	var ctypename [C.GxB_MAX_NAME_LEN]C.char
	info := Info(C.GxB_Vector_type_name(&ctypename[0], vector.grb))
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

// VectorNew creates a new vector with specified domain and size.
//
// Parameters:
//
//   - D: The type corresponding to the domain of the matrix being created.
//     Can be one of the [Predefined] or [Complex] types, or an existing
//     user-defined GraphBLAS type.
//
//   - size (IN): The size of the vector being created.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func VectorNew[D any](size int) (vector Vector[D], err error) {
	if size < 0 {
		err = makeError(InvalidValue)
		return
	}
	var d D
	dt, ok := grbType[TypeOf(d)]
	if !ok {
		err = makeError(UninitializedObject)
		return
	}
	info := Info(C.GrB_Vector_new(&vector.grb, dt, C.GrB_Index(size)))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Dup creates a new vector with the same domain, size, and contents as another vector.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) Dup() (dup Vector[D], err error) {
	info := Info(C.GrB_Vector_dup(&dup.grb, vector.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Resize changes the size of an existing vector.
//
// Parameters:
//
//   - size (IN): The new size of the vector. It can be smaller or larger than the current size.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) Resize(size int) error {
	if size < 0 {
		return makeError(InvalidValue)
	}
	info := Info(C.GrB_Vector_resize(vector.grb, C.GrB_Index(size)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Clear removes all elements (tuples) from the vector.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) Clear() error {
	info := Info(C.GrB_Vector_clear(vector.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Size retrieves the size of a vector.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
func (vector Vector[D]) Size() (size int, err error) {
	var csize C.GrB_Index
	info := Info(C.GrB_Vector_size(&csize, vector.grb))
	if info == success {
		return int(csize), nil
	}
	err = makeError(info)
	return
}

// Nvals retrieves the number of stored elements (tuples) in a vector.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) Nvals() (nvals int, err error) {
	var cnvals C.GrB_Index
	info := Info(C.GrB_Vector_nvals(&cnvals, vector.grb))
	if info == success {
		return int(cnvals), nil
	}
	err = makeError(info)
	return
}

// Build stores elements from tuples in a vector.
//
// Parameters:
//
//   - indices: A slice of indices.
//
//   - values: A slice of scalars of type D.
//
//   - dup: An associative and commutative binary operator to apply when duplicate
//     values for the same index are present in the input slices. All three domains
//     of dup must be D. If dup is nil, then duplicate indices will result in an [InvalidValue] error.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidValue], [SliceMismatch], [OutputNotEmpty], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) Build(indices []int, values []D, dup *BinaryOp[D, D, D]) error {
	if len(indices) != len(values) {
		return makeError(SliceMismatch)
	}
	for _, index := range indices {
		if index < 0 {
			return makeError(InvalidIndex)
		}
	}
	var cdup C.GrB_BinaryOp
	if dup == nil {
		cdup = C.GrB_BinaryOp(C.GrB_NULL)
	} else {
		cdup = dup.grb
	}
	var info Info
	switch vals := any(values).(type) {
	case []bool:
		info = Info(C.GrB_Vector_build_BOOL(
			vector.grb, grbIndices(indices),
			cSlice[C.bool, bool](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_build_INT32(
				vector.grb, grbIndices(indices),
				cSlice[C.int32_t, int](vals),
				C.GrB_Index(len(indices)), cdup,
			))
		} else {
			info = Info(C.GrB_Vector_build_INT64(
				vector.grb, grbIndices(indices),
				cSlice[C.int64_t, int](vals),
				C.GrB_Index(len(indices)), cdup,
			))
		}
	case []int8:
		info = Info(C.GrB_Vector_build_INT8(
			vector.grb, grbIndices(indices),
			cSlice[C.int8_t, int8](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []int16:
		info = Info(C.GrB_Vector_build_INT16(
			vector.grb, grbIndices(indices),
			cSlice[C.int16_t, int16](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []int32:
		info = Info(C.GrB_Vector_build_INT32(
			vector.grb, grbIndices(indices),
			cSlice[C.int32_t, int32](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []int64:
		info = Info(C.GrB_Vector_build_INT64(
			vector.grb, grbIndices(indices),
			cSlice[C.int64_t, int64](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_build_UINT32(
				vector.grb, grbIndices(indices),
				cSlice[C.uint32_t, uint](vals),
				C.GrB_Index(len(indices)), cdup,
			))
		} else {
			info = Info(C.GrB_Vector_build_UINT64(
				vector.grb, grbIndices(indices),
				cSlice[C.uint64_t, uint](vals),
				C.GrB_Index(len(indices)), cdup,
			))
		}
	case []uint8:
		info = Info(C.GrB_Vector_build_UINT8(
			vector.grb, grbIndices(indices),
			cSlice[C.uint8_t, uint8](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []uint16:
		info = Info(C.GrB_Vector_build_UINT16(
			vector.grb, grbIndices(indices),
			cSlice[C.uint16_t, uint16](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []uint32:
		info = Info(C.GrB_Vector_build_UINT32(
			vector.grb, grbIndices(indices),
			cSlice[C.uint32_t, uint32](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []uint64:
		info = Info(C.GrB_Vector_build_UINT64(
			vector.grb, grbIndices(indices),
			cSlice[C.uint64_t, uint64](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []float32:
		info = Info(C.GrB_Vector_build_FP32(
			vector.grb, grbIndices(indices),
			cSlice[C.float, float32](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []float64:
		info = Info(C.GrB_Vector_build_FP64(
			vector.grb, grbIndices(indices),
			cSlice[C.double, float64](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []complex64:
		info = Info(C.GxB_Vector_build_FC32(
			vector.grb, grbIndices(indices),
			cSlice[C.complexfloat, complex64](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	case []complex128:
		info = Info(C.GxB_Vector_build_FC64(
			vector.grb, grbIndices(indices),
			cSlice[C.complexdouble, complex128](vals),
			C.GrB_Index(len(indices)), cdup,
		))
	default:
		info = Info(C.GrB_Vector_build_UDT(
			vector.grb, grbIndices(indices),
			unsafe.Pointer(unsafe.SliceData(values)),
			C.GrB_Index(len(indices)), cdup,
		))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// BuildScalar is like [Vector.Build], except that the scalar is the value of all the tuples.
//
// Unlike [Vector.Build], there is no dup operator to handle duplicate entries. Instead, any
// duplicates are silently ignored.
//
// BuildScalar is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) BuildScalar(indices []int, scalar Scalar[D]) error {
	for _, index := range indices {
		if index < 0 {
			return makeError(InvalidIndex)
		}
	}
	info := Info(C.GxB_Vector_build_Scalar(
		vector.grb, grbIndices(indices),
		scalar.grb, C.GrB_Index(len(indices)),
	))
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetElement sets one element of a vector to a given value.
//
// To pass a [Scalar] object instead of a non-opaque variable, use [Vector.SetElementScalar].
//
// Parameters:
//
//   - val (IN): Scalar to assign.
//
//   - index (IN): Index of element to be assigned.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) SetElement(val D, index int) error {
	if index < 0 {
		return makeError(InvalidIndex)
	}
	var info Info
	switch value := any(val).(type) {
	case bool:
		info = Info(C.GrB_Vector_setElement_BOOL(vector.grb, C.bool(value), C.GrB_Index(index)))
	case int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_setElement_INT32(vector.grb, C.int32_t(value), C.GrB_Index(index)))
		} else {
			info = Info(C.GrB_Vector_setElement_INT64(vector.grb, C.int64_t(value), C.GrB_Index(index)))
		}
	case int8:
		info = Info(C.GrB_Vector_setElement_INT8(vector.grb, C.int8_t(value), C.GrB_Index(index)))
	case int16:
		info = Info(C.GrB_Vector_setElement_INT16(vector.grb, C.int16_t(value), C.GrB_Index(index)))
	case int32:
		info = Info(C.GrB_Vector_setElement_INT32(vector.grb, C.int32_t(value), C.GrB_Index(index)))
	case int64:
		info = Info(C.GrB_Vector_setElement_INT64(vector.grb, C.int64_t(value), C.GrB_Index(index)))
	case uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_setElement_UINT32(vector.grb, C.uint32_t(value), C.GrB_Index(index)))
		} else {
			info = Info(C.GrB_Vector_setElement_UINT64(vector.grb, C.uint64_t(value), C.GrB_Index(index)))
		}
	case uint8:
		info = Info(C.GrB_Vector_setElement_UINT8(vector.grb, C.uint8_t(value), C.GrB_Index(index)))
	case uint16:
		info = Info(C.GrB_Vector_setElement_UINT16(vector.grb, C.uint16_t(value), C.GrB_Index(index)))
	case uint32:
		info = Info(C.GrB_Vector_setElement_UINT32(vector.grb, C.uint32_t(value), C.GrB_Index(index)))
	case uint64:
		info = Info(C.GrB_Vector_setElement_UINT64(vector.grb, C.uint64_t(value), C.GrB_Index(index)))
	case float32:
		info = Info(C.GrB_Vector_setElement_FP32(vector.grb, C.float(value), C.GrB_Index(index)))
	case float64:
		info = Info(C.GrB_Vector_setElement_FP64(vector.grb, C.double(value), C.GrB_Index(index)))
	case complex64:
		info = Info(C.GxB_Vector_setElement_FC32(vector.grb, C.complexfloat(value), C.GrB_Index(index)))
	case complex128:
		info = Info(C.GxB_Vector_setElement_FC64(vector.grb, C.complexdouble(value), C.GrB_Index(index)))
	default:
		info = Info(C.GrB_Vector_setElement_UDT(vector.grb, unsafe.Pointer(&val), C.GrB_Index(index)))
	}
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetElementScalar is like [Vector.SetElement], except that the scalar value is passed as a [Scalar]
// object. It may be empty.
func (vector Vector[D]) SetElementScalar(val Scalar[D], index int) error {
	if index < 0 {
		return makeError(InvalidIndex)
	}
	info := Info(C.GrB_Vector_setElement_Scalar(vector.grb, val.grb, C.GrB_Index(index)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// RemoveElement removes (annihilates) one stored element from a vector.
//
// Parameters:
//
//   - index (IN): Index of element to be removed.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) RemoveElement(index int) error {
	if index < 0 {
		return makeError(InvalidIndex)
	}
	info := Info(C.GrB_Vector_removeElement(vector.grb, C.GrB_Index(index)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// ExtractElement extracts one element of a vector.
//
// When there is no stored value at the specified index, ExtractElement returns
// ok == false. Otherwise, it returns ok == true.
//
// To store the element in a [Scalar] object instead of returning a non-opaque value,
// use [Vector.ExtractElementScalar].
//
// Parameters:
//
//   - index (IN): Index of element to be assigned.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) ExtractElement(index int) (result D, ok bool, err error) {
	if index < 0 {
		err = makeError(InvalidIndex)
		return
	}
	var info Info
	switch res := any(&result).(type) {
	case *bool:
		var cresult C.bool
		info = Info(C.GrB_Vector_extractElement_BOOL(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = bool(cresult)
			ok = true
			return
		}
	case *int:
		if unsafe.Sizeof(0) == 4 {
			var cresult C.int32_t
			info = Info(C.GrB_Vector_extractElement_INT32(&cresult, vector.grb, C.GrB_Index(index)))
			if info == success {
				*res = int(cresult)
				ok = true
				return
			}
		} else {
			var cresult C.int64_t
			info = Info(C.GrB_Vector_extractElement_INT64(&cresult, vector.grb, C.GrB_Index(index)))
			if info == success {
				*res = int(cresult)
				ok = true
				return
			}
		}
	case *int8:
		var cresult C.int8_t
		info = Info(C.GrB_Vector_extractElement_INT8(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = int8(cresult)
			ok = true
			return
		}
	case *int16:
		var cresult C.int16_t
		info = Info(C.GrB_Vector_extractElement_INT16(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = int16(cresult)
			ok = true
			return
		}
	case *int32:
		var cresult C.int32_t
		info = Info(C.GrB_Vector_extractElement_INT32(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = int32(cresult)
			ok = true
			return
		}
	case *int64:
		var cresult C.int64_t
		info = Info(C.GrB_Vector_extractElement_INT64(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = int64(cresult)
			ok = true
			return
		}
	case *uint:
		if unsafe.Sizeof(0) == 4 {
			var cresult C.uint32_t
			info = Info(C.GrB_Vector_extractElement_UINT32(&cresult, vector.grb, C.GrB_Index(index)))
			if info == success {
				*res = uint(cresult)
				ok = true
				return
			}
		} else {
			var cresult C.uint64_t
			info = Info(C.GrB_Vector_extractElement_UINT64(&cresult, vector.grb, C.GrB_Index(index)))
			if info == success {
				*res = uint(cresult)
				ok = true
				return
			}
		}
	case *uint8:
		var cresult C.uint8_t
		info = Info(C.GrB_Vector_extractElement_UINT8(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = uint8(cresult)
			ok = true
			return
		}
	case *uint16:
		var cresult C.uint16_t
		info = Info(C.GrB_Vector_extractElement_UINT16(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = uint16(cresult)
			ok = true
			return
		}
	case *uint32:
		var cresult C.uint32_t
		info = Info(C.GrB_Vector_extractElement_UINT32(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = uint32(cresult)
			ok = true
			return
		}
	case *uint64:
		var cresult C.uint64_t
		info = Info(C.GrB_Vector_extractElement_UINT64(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = uint64(cresult)
			ok = true
			return
		}
	case *float32:
		var cresult C.float
		info = Info(C.GrB_Vector_extractElement_FP32(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = float32(cresult)
			ok = true
			return
		}
	case *float64:
		var cresult C.double
		info = Info(C.GrB_Vector_extractElement_FP64(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = float64(cresult)
			ok = true
			return
		}
	case *complex64:
		var cresult C.complexfloat
		info = Info(C.GxB_Vector_extractElement_FC32(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = complex64(cresult)
			ok = true
			return
		}
	case *complex128:
		var cresult C.complexdouble
		info = Info(C.GxB_Vector_extractElement_FC64(&cresult, vector.grb, C.GrB_Index(index)))
		if info == success {
			*res = complex128(cresult)
			ok = true
			return
		}
	default:
		info = Info(C.GrB_Vector_extractElement_UDT(unsafe.Pointer(&result), vector.grb, C.GrB_Index(index)))
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

// ExtractElementScalar is like [Vector.ExtractElement], except that the element is stored in a [Scalar]
// object.
//
// When there is no stored value at the specified location, the result becomes empty.
func (vector Vector[D]) ExtractElementScalar(result Scalar[D], index int) error {
	if index < 0 {
		return makeError(InvalidIndex)
	}
	info := Info(C.GrB_Vector_extractElement_Scalar(result.grb, vector.grb, C.GrB_Index(index)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// IsStoredElement determines whether there is a stored value at the specified
// index or not.
//
// Parameters:
//
//   - Index (IN): Index of element to be assigned.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidIndex], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
//
// IsStoredElement is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) IsStoredElement(index int) (ok bool, err error) {
	if index < 0 {
		err = makeError(InvalidIndex)
		return
	}
	switch info := Info(C.GxB_Vector_isStoredElement(vector.grb, C.GrB_Index(index))); info {
	case success:
		return true, nil
	case noValue:
		return false, nil
	default:
		err = makeError(info)
		return
	}
}

// ExtractTuples extracts the contents of a GraphBLAS vector into non-opaque slices,
// by appending the indices and values to the slices
// passed to this function (by using Go's built-in append function).
//
// Parameters:
//
//   - indices (INOUT): Pointer to a slice of indices. If nil, ExtractTuples does not
//     produces the indices of the vector.
//
//   - values (INOUT): Pointer to a slice of indices. If nil, ExtractTuples does not
//     produces the values of the vector.
//
// It is valid to pass pointers to nil slices, and ExtractTuples then produces the
// corresponding indices or values.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) ExtractTuples(indices *[]int, values *[]D) error {
	nvals, err := vector.Nvals()
	if err != nil {
		return err
	}
	targetIndices, finalizeTargetIndices := growIndices(indices, nvals)
	targetValues := growslice(values, nvals)
	var info Info
	cnvals := C.GrB_Index(nvals)
	switch vals := any(targetValues).(type) {
	case []bool:
		info = Info(C.GrB_Vector_extractTuples_BOOL(
			targetIndices,
			cSlice[C.bool, bool](vals),
			&cnvals, vector.grb,
		))
	case []int:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_extractTuples_INT32(
				targetIndices,
				cSlice[C.int32_t, int](vals),
				&cnvals, vector.grb,
			))
		} else {
			info = Info(C.GrB_Vector_extractTuples_INT64(
				targetIndices,
				cSlice[C.int64_t, int](vals),
				&cnvals, vector.grb,
			))
		}
	case []int8:
		info = Info(C.GrB_Vector_extractTuples_INT8(
			targetIndices,
			cSlice[C.int8_t, int8](vals),
			&cnvals, vector.grb,
		))
	case []int16:
		info = Info(C.GrB_Vector_extractTuples_INT16(
			targetIndices,
			cSlice[C.int16_t, int16](vals),
			&cnvals, vector.grb,
		))
	case []int32:
		info = Info(C.GrB_Vector_extractTuples_INT32(
			targetIndices,
			cSlice[C.int32_t, int32](vals),
			&cnvals, vector.grb,
		))
	case []int64:
		info = Info(C.GrB_Vector_extractTuples_INT64(
			targetIndices,
			cSlice[C.int64_t, int64](vals),
			&cnvals, vector.grb,
		))
	case []uint:
		if unsafe.Sizeof(0) == 4 {
			info = Info(C.GrB_Vector_extractTuples_UINT32(
				targetIndices,
				cSlice[C.uint32_t, uint](vals),
				&cnvals, vector.grb,
			))
		} else {
			info = Info(C.GrB_Vector_extractTuples_UINT64(
				targetIndices,
				cSlice[C.uint64_t, uint](vals),
				&cnvals, vector.grb,
			))
		}
	case []uint8:
		info = Info(C.GrB_Vector_extractTuples_UINT8(
			targetIndices,
			cSlice[C.uint8_t, uint8](vals),
			&cnvals, vector.grb,
		))
	case []uint16:
		info = Info(C.GrB_Vector_extractTuples_UINT16(
			targetIndices,
			cSlice[C.uint16_t, uint16](vals),
			&cnvals, vector.grb,
		))
	case []uint32:
		info = Info(C.GrB_Vector_extractTuples_UINT32(
			targetIndices,
			cSlice[C.uint32_t, uint32](vals),
			&cnvals, vector.grb,
		))
	case []uint64:
		info = Info(C.GrB_Vector_extractTuples_UINT64(
			targetIndices,
			cSlice[C.uint64_t, uint64](vals),
			&cnvals, vector.grb,
		))
	case []float32:
		info = Info(C.GrB_Vector_extractTuples_FP32(
			targetIndices,
			cSlice[C.float, float32](vals),
			&cnvals, vector.grb,
		))
	case []float64:
		info = Info(C.GrB_Vector_extractTuples_FP64(
			targetIndices,
			cSlice[C.double, float64](vals),
			&cnvals, vector.grb,
		))
	case []complex64:
		info = Info(C.GxB_Vector_extractTuples_FC32(
			targetIndices,
			cSlice[C.complexfloat, complex64](vals),
			&cnvals, vector.grb,
		))
	case []complex128:
		info = Info(C.GxB_Vector_extractTuples_FC64(
			targetIndices,
			cSlice[C.complexdouble, complex128](vals),
			&cnvals, vector.grb,
		))
	default:
		info = Info(C.GrB_Vector_extractTuples_UDT(
			targetIndices,
			unsafe.Pointer(unsafe.SliceData(targetValues)),
			&cnvals, vector.grb,
		))
	}
	if info == success {
		if nvals != int(cnvals) {
			return makeError(InvalidObject)
		}
		finalizeTargetIndices()
		return nil
	}
	return makeError(info)
}

// Diag constructs a diagonal GraphBLAS matrix.
//
// Parameters:
//
//   - vector (IN): The GraphBLAS vector whose contents will be copied to the
//     diagonal of the matrix. The matrix is square with each dimension equal
//     to size(vector) + |k|.
//
//   - k (IN): The diagonal to which the vector is assigned. k == 0 represents
//     the main diagonal, k > 0 is above the main dioganal, and k < 0 is below.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) Diag(k int) (diag Matrix[D], err error) {
	info := Info(C.GrB_Matrix_diag(&diag.grb, vector.grb, C.int64_t(int64(k))))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// IteratorNew creates a vector iterator and attaches it to the vector.
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
//
// IteratorNew is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) IteratorNew(desc *Descriptor) (it VectorIterator[D], err error) {
	info := Info(C.GxB_Iterator_new(&it.grb))
	if info != success {
		err = makeError(info)
		return
	}
	cdesc := processDescriptor(desc)
	info = Info(C.GxB_Vector_Iterator_attach(it.grb, vector.grb, cdesc))
	if info == success {
		it.init()
		return
	}
	err = makeError(info)
	return
}

// Sort a vector.
//
// Parameters:
//
//   - into (OUT): Contains the vector of sorted values. If nil, this output is not produced.
//
//   - p (OUT): Contains the permutations of the sorted values. If nil, this output is not produced.
//
//   - op (IN): The comparator operation.
//
// Sort is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) Sort(
	into *Vector[D],
	p *Vector[int],
	op BinaryOp[bool, D, D],
	desc *Descriptor,
) error {
	var cinto, cp C.GrB_Vector
	if into == nil {
		cinto = C.GrB_Vector(C.NULL)
	} else {
		cinto = into.grb
	}
	if p == nil {
		cp = C.GrB_Vector(C.NULL)
	} else {
		cp = p.grb
	}
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_sort(cinto, cp, op.grb, vector.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MemoryUsage returns the memory space required for a vector, in bytes.
//
// MemoryUsage is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) MemoryUsage() (size int, err error) {
	var csize C.size_t
	info := Info(C.GxB_Vector_memoryUsage(&csize, vector.grb))
	if info == success {
		return int(csize), nil
	}
	err = makeError(info)
	return
}

// Iso returns true if the vector is iso-valued, false otherwise.
//
// Iso is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) Iso() (iso bool, err error) {
	var ciso C.bool
	info := Info(C.GxB_Vector_iso(&ciso, vector.grb))
	if info == success {
		return bool(ciso), nil
	}
	err = makeError(info)
	return
}

// ExtractDiag extracts a diagonal from a GraphBLAS matrix.
//
// Parameters:
//
//   - vector (OUT): The GraphBLAS vector whose contents will be a copy of the
//     diagonal of the matrix.
//
//   - k (IN): The diagonal from which the vector is assigned. k == 0 represents
//     the main diagonal, k > 0 is above the main dioganal, and k < 0 is below.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [OutOfMemory], [Panic]
func (vector Vector[D]) ExtractDiag(a Matrix[D], k int, desc *Descriptor) error {
	cdesc := processDescriptor(desc)
	info := Info(C.GxB_Vector_diag(vector.grb, a.grb, C.int64_t(int64(k)), cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetBitmapSwitch determines how the vector is converted to the bitmap format.
//
// Parameters:
//
//   - bitmapSwitch (IN): A value between 0 and 1.
//
// SetBitmapSwitch is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) SetBitmapSwitch(bitmapSwitch float64) error {
	info := Info(C.GxB_Vector_Option_set_FP64(vector.grb, C.GxB_BITMAP_SWITCH, C.double(bitmapSwitch)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetBitmapSwitch retrieves the current switch to bitmap. See [Vector.SetBitmapSwitch].
//
// GetBitmapSwitch is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) GetBitmapSwitch() (bitmapSwitch float64, err error) {
	var cBitmapSwitch C.double
	info := Info(C.GxB_Vector_Option_get_FP64(vector.grb, C.GxB_BITMAP_SWITCH, &cBitmapSwitch))
	if info == success {
		return float64(cBitmapSwitch), nil
	}
	err = makeError(info)
	return
}

// SetSparsityControl determines the valid [Sparsity] format(s) for the vector.
//
// SetSparsityControl is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) SetSparsityControl(sparsity Sparsity) error {
	info := Info(C.GxB_Vector_Option_set_INT32(vector.grb, C.GxB_SPARSITY_CONTROL, C.int32_t(sparsity)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetSparsityControl retrieves the valid [Sparsity] format(s) of the vector.
//
// GetSparsityControl is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) GetSparsityControl() (sparsity Sparsity, err error) {
	var csparsity C.int32_t
	info := Info(C.GxB_Vector_Option_get_INT32(vector.grb, C.GxB_SPARSITY_CONTROL, &csparsity))
	if info == success {
		return Sparsity(csparsity), nil
	}
	err = makeError(info)
	return
}

// GetSparsityStatus retrieves the current [Sparsity] format of the vector.
//
// GetSparsityStatus is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) GetSparsityStatus() (status Sparsity, err error) {
	var cstatus C.int32_t
	info := Info(C.GxB_Vector_Option_get_INT32(vector.grb, C.GxB_SPARSITY_STATUS, &cstatus))
	if info == success {
		return Sparsity(cstatus), nil
	}
	err = makeError(info)
	return
}

// Valid returns true if matrix has been created by a successful call to [Vector.Dup] or [VectorNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (vector Vector[D]) Valid() bool {
	return vector.grb != C.GrB_Vector(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Vector] and releases any resources associated with
// it. Calling Free on an object that is not [Vector.Valid]() is legal.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (vector *Vector[D]) Free() error {
	info := Info(C.GrB_Vector_free(&vector.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the vector into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (vector Vector[D]) Wait(mode WaitMode) error {
	info := Info(C.GrB_Vector_wait(vector.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the matrix.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (vector Vector[D]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_Vector_error(&cerror, vector.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the vector to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: vector is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (vector Vector[D]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Vector_fprint(vector.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}
