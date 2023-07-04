package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// A Semiring is defined by three domains Dout, Din1 and Din2; an associative and commutative operator additive operation;
// a multiplicative operation; and an identity element.
//
// In a GraphBLAS semiring, the multiplicative operator does not have to distribute over the additive operator. This is
// unlike the conventional semiring algebraic structure.
type Semiring[Dout, Din1, Din2 any] struct {
	grb C.GrB_Semiring
}

// SemiringNew creates a new monoid with specified domains, operators, and elements.
//
// Parameters:
//
//   - addOp (IN): An existing GraphBLAS commutative monoid that specifies the addition operator and its identity.
//
//   - mulOp (IN): An existing GraphBLAS binary operator that specifies the semiring's multiplicative operator.
//
// GraphBLAS API errors that may be returned:
//   - [DomainMismatch], [NullPointer], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func SemiringNew[Dout, Din1, Din2 any](addOp Monoid[Dout], mulOp BinaryOp[Dout, Din1, Din2]) (semiring Semiring[Dout, Din1, Din2], err error) {
	info := Info(C.GrB_Semiring_new(&semiring.grb, addOp.grb, mulOp.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Add returns the additive monoid of the semiring.
//
// Add is a SuiteSparse:GraphBLAS extension.
func (semiring Semiring[Dout, Din1, Din2]) Add() (add Monoid[Dout], err error) {
	info := Info(C.GxB_Semiring_add(&add.grb, semiring.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Multiply returns the binary operator of the semiring.
//
// Multiply is a SuiteSparse:GraphBLAS extension.
func (semiring Semiring[Dout, Din1, Din2]) Multiply() (multiply BinaryOp[Dout, Din1, Din2], err error) {
	info := Info(C.GxB_Semiring_multiply(&multiply.grb, semiring.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Valid returns true if semiring has been created by a successful call to [SemiringNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (semiring Semiring[Dout, Din1, Din2]) Valid() bool {
	return semiring.grb != C.GrB_Semiring(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Semiring] and releases any resources associated with
// it. Calling Free on an object that is not [Semiring.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined semiring is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (semiring *Semiring[Dout, Din1, Din2]) Free() error {
	info := Info(C.GrB_Semiring_free(&semiring.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the semiring into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (semiring Semiring[Dout, Din1, Din2]) Wait(mode WaitMode) error {
	info := Info(C.GrB_Semiring_wait(semiring.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the semiring.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (semiring Semiring[Dout, Din1, Din2]) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_Semiring_error(&cerror, semiring.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the semiring to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: semiring is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (semiring Semiring[Dout, Din1, Din2]) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Semiring_fprint(semiring.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// PlusTimesSemiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Times].
func PlusTimesSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_PLUS_TIMES_SEMIRING_INT32
		} else {
			s.grb = C.GrB_PLUS_TIMES_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_PLUS_TIMES_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_PLUS_TIMES_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_PLUS_TIMES_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinPlusSemiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Plus].
func MinPlusSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_PLUS_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MIN_PLUS_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_PLUS_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MIN_PLUS_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MIN_PLUS_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxPlusSemiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Plus].
func MaxPlusSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_PLUS_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MAX_PLUS_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_PLUS_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MAX_PLUS_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MAX_PLUS_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinTimesSemiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Times].
func MinTimesSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_TIMES_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MIN_TIMES_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_TIMES_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MIN_TIMES_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MIN_TIMES_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinMaxSemiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Max].
func MinMaxSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_MAX_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MIN_MAX_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MIN_MAX_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MIN_MAX_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MIN_MAX_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MIN_MAX_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_MAX_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MIN_MAX_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MIN_MAX_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MIN_MAX_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MIN_MAX_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MIN_MAX_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MIN_MAX_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MIN_MAX_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxMinSemiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Min].
func MaxMinSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_MIN_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MAX_MIN_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MAX_MIN_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MAX_MIN_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MAX_MIN_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MAX_MIN_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_MIN_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MAX_MIN_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MAX_MIN_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MAX_MIN_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MAX_MIN_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MAX_MIN_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MAX_MIN_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MAX_MIN_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxTimesSemiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [times].
func MaxTimesSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_TIMES_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MAX_TIMES_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_TIMES_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MAX_TIMES_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MAX_TIMES_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusMinSemiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Min].
func PlusMinSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_PLUS_MIN_SEMIRING_INT32
		} else {
			s.grb = C.GrB_PLUS_MIN_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_PLUS_MIN_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_PLUS_MIN_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_PLUS_MIN_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

var (
	// LorLandSemiringBool with additive [Monoid] [LorMonoidBool] and [BinaryOp] [LandBool].
	LorLandSemiringBool = Semiring[bool, bool, bool]{C.GrB_LOR_LAND_SEMIRING_BOOL}

	// LandLorSemiringBool with additive [Monoid] [LandMonoidBool] and [BinaryOp] [LorBool].
	LandLorSemiringBool = Semiring[bool, bool, bool]{C.GrB_LAND_LOR_SEMIRING_BOOL}

	// LxorLandSemiringBool with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [LandBool].
	LxorLandSemiringBool = Semiring[bool, bool, bool]{C.GrB_LXOR_LAND_SEMIRING_BOOL}

	// LxnorLorSemiringBool with additive [Monoid] [LxnorMonoidBool] and [BinaryOp] [LorBool].
	LxnorLorSemiringBool = Semiring[bool, bool, bool]{C.GrB_LXNOR_LOR_SEMIRING_BOOL}
)

// MinFirstSemiring with additive [Monoid] [MinMonoid] and [BinaryOp] [First].
func MinFirstSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_FIRST_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MIN_FIRST_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_FIRST_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MIN_FIRST_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MIN_FIRST_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinSecondSemiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Second].
func MinSecondSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_SECOND_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MIN_SECOND_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MIN_SECOND_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MIN_SECOND_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MIN_SECOND_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxFirstSemiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [First].
func MaxFirstSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_FIRST_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MAX_FIRST_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_FIRST_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MAX_FIRST_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MAX_FIRST_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxSecondSemiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Second].
func MaxSecondSemiring[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_SECOND_SEMIRING_INT32
		} else {
			s.grb = C.GrB_MAX_SECOND_SEMIRING_INT64
		}
	case int8:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_INT8
	case int16:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_INT16
	case int32:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_INT32
	case int64:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GrB_MAX_SECOND_SEMIRING_UINT32
		} else {
			s.grb = C.GrB_MAX_SECOND_SEMIRING_UINT64
		}
	case uint8:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_UINT8
	case uint16:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_UINT16
	case uint32:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_UINT32
	case uint64:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_UINT64
	case float32:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_FP32
	case float64:
		s.grb = C.GrB_MAX_SECOND_SEMIRING_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusFirst semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [First].
//
// PlusFirst is a SuiteSparse:GraphBLAS extension.
func PlusFirst[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_FIRST_INT32
		} else {
			s.grb = C.GxB_PLUS_FIRST_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_FIRST_INT8
	case int16:
		s.grb = C.GxB_PLUS_FIRST_INT16
	case int32:
		s.grb = C.GxB_PLUS_FIRST_INT32
	case int64:
		s.grb = C.GxB_PLUS_FIRST_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_FIRST_UINT32
		} else {
			s.grb = C.GxB_PLUS_FIRST_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_FIRST_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_FIRST_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_FIRST_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_FIRST_UINT64
	case float32:
		s.grb = C.GxB_PLUS_FIRST_FP32
	case float64:
		s.grb = C.GxB_PLUS_FIRST_FP64
	case complex64:
		s.grb = C.GxB_PLUS_FIRST_FC32
	case complex128:
		s.grb = C.GxB_PLUS_FIRST_FC64
	default:
		panic("unreachable code")
	}
	return
}

// PlusSecond semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Second].
//
// PlusSecond is a SuiteSparse:GraphBLAS extension.
func PlusSecond[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_SECOND_INT32
		} else {
			s.grb = C.GxB_PLUS_SECOND_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_SECOND_INT8
	case int16:
		s.grb = C.GxB_PLUS_SECOND_INT16
	case int32:
		s.grb = C.GxB_PLUS_SECOND_INT32
	case int64:
		s.grb = C.GxB_PLUS_SECOND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_SECOND_UINT32
		} else {
			s.grb = C.GxB_PLUS_SECOND_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_SECOND_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_SECOND_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_SECOND_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_SECOND_UINT64
	case float32:
		s.grb = C.GxB_PLUS_SECOND_FP32
	case float64:
		s.grb = C.GxB_PLUS_SECOND_FP64
	case complex64:
		s.grb = C.GxB_PLUS_SECOND_FC32
	case complex128:
		s.grb = C.GxB_PLUS_SECOND_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesFirst semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [First].
//
// TimesFirst is a SuiteSparse:GraphBLAS extension.
func TimesFirst[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_FIRST_INT32
		} else {
			s.grb = C.GxB_TIMES_FIRST_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_FIRST_INT8
	case int16:
		s.grb = C.GxB_TIMES_FIRST_INT16
	case int32:
		s.grb = C.GxB_TIMES_FIRST_INT32
	case int64:
		s.grb = C.GxB_TIMES_FIRST_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_FIRST_UINT32
		} else {
			s.grb = C.GxB_TIMES_FIRST_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_FIRST_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_FIRST_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_FIRST_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_FIRST_UINT64
	case float32:
		s.grb = C.GxB_TIMES_FIRST_FP32
	case float64:
		s.grb = C.GxB_TIMES_FIRST_FP64
	case complex64:
		s.grb = C.GxB_TIMES_FIRST_FC32
	case complex128:
		s.grb = C.GxB_TIMES_FIRST_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesSecond semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Second].
//
// TimesSecond is a SuiteSparse:GraphBLAS extension.
func TimesSecond[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_SECOND_INT32
		} else {
			s.grb = C.GxB_TIMES_SECOND_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_SECOND_INT8
	case int16:
		s.grb = C.GxB_TIMES_SECOND_INT16
	case int32:
		s.grb = C.GxB_TIMES_SECOND_INT32
	case int64:
		s.grb = C.GxB_TIMES_SECOND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_SECOND_UINT32
		} else {
			s.grb = C.GxB_TIMES_SECOND_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_SECOND_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_SECOND_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_SECOND_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_SECOND_UINT64
	case float32:
		s.grb = C.GxB_TIMES_SECOND_FP32
	case float64:
		s.grb = C.GxB_TIMES_SECOND_FP64
	case complex64:
		s.grb = C.GxB_TIMES_SECOND_FC32
	case complex128:
		s.grb = C.GxB_TIMES_SECOND_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyFirst semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [First].
//
// AnyFirst is a SuiteSparse:GraphBLAS extension.
func AnyFirst[D Predefined | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_FIRST_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_FIRST_INT32
		} else {
			s.grb = C.GxB_ANY_FIRST_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_FIRST_INT8
	case int16:
		s.grb = C.GxB_ANY_FIRST_INT16
	case int32:
		s.grb = C.GxB_ANY_FIRST_INT32
	case int64:
		s.grb = C.GxB_ANY_FIRST_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_FIRST_UINT32
		} else {
			s.grb = C.GxB_ANY_FIRST_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_FIRST_UINT8
	case uint16:
		s.grb = C.GxB_ANY_FIRST_UINT16
	case uint32:
		s.grb = C.GxB_ANY_FIRST_UINT32
	case uint64:
		s.grb = C.GxB_ANY_FIRST_UINT64
	case float32:
		s.grb = C.GxB_ANY_FIRST_FP32
	case float64:
		s.grb = C.GxB_ANY_FIRST_FP64
	case complex64:
		s.grb = C.GxB_ANY_FIRST_FC32
	case complex128:
		s.grb = C.GxB_ANY_FIRST_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnySecond semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Second].
//
// AnySecond is a SuiteSparse:GraphBLAS extension.
func AnySecond[D Predefined | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_SECOND_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_SECOND_INT32
		} else {
			s.grb = C.GxB_ANY_SECOND_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_SECOND_INT8
	case int16:
		s.grb = C.GxB_ANY_SECOND_INT16
	case int32:
		s.grb = C.GxB_ANY_SECOND_INT32
	case int64:
		s.grb = C.GxB_ANY_SECOND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_SECOND_UINT32
		} else {
			s.grb = C.GxB_ANY_SECOND_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_SECOND_UINT8
	case uint16:
		s.grb = C.GxB_ANY_SECOND_UINT16
	case uint32:
		s.grb = C.GxB_ANY_SECOND_UINT32
	case uint64:
		s.grb = C.GxB_ANY_SECOND_UINT64
	case float32:
		s.grb = C.GxB_ANY_SECOND_FP32
	case float64:
		s.grb = C.GxB_ANY_SECOND_FP64
	case complex64:
		s.grb = C.GxB_ANY_SECOND_FC32
	case complex128:
		s.grb = C.GxB_ANY_SECOND_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinOneb semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Oneb].
//
// MinOneb is a SuiteSparse:GraphBLAS extension.
func MinOneb[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_PAIR_INT32
		} else {
			s.grb = C.GxB_MIN_PAIR_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_PAIR_INT8
	case int16:
		s.grb = C.GxB_MIN_PAIR_INT16
	case int32:
		s.grb = C.GxB_MIN_PAIR_INT32
	case int64:
		s.grb = C.GxB_MIN_PAIR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_PAIR_UINT32
		} else {
			s.grb = C.GxB_MIN_PAIR_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_PAIR_UINT8
	case uint16:
		s.grb = C.GxB_MIN_PAIR_UINT16
	case uint32:
		s.grb = C.GxB_MIN_PAIR_UINT32
	case uint64:
		s.grb = C.GxB_MIN_PAIR_UINT64
	case float32:
		s.grb = C.GxB_MIN_PAIR_FP32
	case float64:
		s.grb = C.GxB_MIN_PAIR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxOneb semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Oneb].
//
// MaxOneb is a SuiteSparse:GraphBLAS extension.
func MaxOneb[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_PAIR_INT32
		} else {
			s.grb = C.GxB_MAX_PAIR_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_PAIR_INT8
	case int16:
		s.grb = C.GxB_MAX_PAIR_INT16
	case int32:
		s.grb = C.GxB_MAX_PAIR_INT32
	case int64:
		s.grb = C.GxB_MAX_PAIR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_PAIR_UINT32
		} else {
			s.grb = C.GxB_MAX_PAIR_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_PAIR_UINT8
	case uint16:
		s.grb = C.GxB_MAX_PAIR_UINT16
	case uint32:
		s.grb = C.GxB_MAX_PAIR_UINT32
	case uint64:
		s.grb = C.GxB_MAX_PAIR_UINT64
	case float32:
		s.grb = C.GxB_MAX_PAIR_FP32
	case float64:
		s.grb = C.GxB_MAX_PAIR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusOneb semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Oneb].
//
// PlusOneb is a SuiteSparse:GraphBLAS extension.
func PlusOneb[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_PAIR_INT32
		} else {
			s.grb = C.GxB_PLUS_PAIR_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_PAIR_INT8
	case int16:
		s.grb = C.GxB_PLUS_PAIR_INT16
	case int32:
		s.grb = C.GxB_PLUS_PAIR_INT32
	case int64:
		s.grb = C.GxB_PLUS_PAIR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_PAIR_UINT32
		} else {
			s.grb = C.GxB_PLUS_PAIR_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_PAIR_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_PAIR_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_PAIR_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_PAIR_UINT64
	case float32:
		s.grb = C.GxB_PLUS_PAIR_FP32
	case float64:
		s.grb = C.GxB_PLUS_PAIR_FP64
	case complex64:
		s.grb = C.GxB_PLUS_PAIR_FC32
	case complex128:
		s.grb = C.GxB_PLUS_PAIR_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesOneb semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Oneb].
//
// TimesOneb is a SuiteSparse:GraphBLAS extension.
func TimesOneb[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_PAIR_INT32
		} else {
			s.grb = C.GxB_TIMES_PAIR_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_PAIR_INT8
	case int16:
		s.grb = C.GxB_TIMES_PAIR_INT16
	case int32:
		s.grb = C.GxB_TIMES_PAIR_INT32
	case int64:
		s.grb = C.GxB_TIMES_PAIR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_PAIR_UINT32
		} else {
			s.grb = C.GxB_TIMES_PAIR_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_PAIR_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_PAIR_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_PAIR_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_PAIR_UINT64
	case float32:
		s.grb = C.GxB_TIMES_PAIR_FP32
	case float64:
		s.grb = C.GxB_TIMES_PAIR_FP64
	case complex64:
		s.grb = C.GxB_TIMES_PAIR_FC32
	case complex128:
		s.grb = C.GxB_TIMES_PAIR_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyOneb semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Oneb].
//
// AnyOneb is a SuiteSparse:GraphBLAS extension.
func AnyOneb[D Predefined | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_PAIR_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_PAIR_INT32
		} else {
			s.grb = C.GxB_ANY_PAIR_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_PAIR_INT8
	case int16:
		s.grb = C.GxB_ANY_PAIR_INT16
	case int32:
		s.grb = C.GxB_ANY_PAIR_INT32
	case int64:
		s.grb = C.GxB_ANY_PAIR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_PAIR_UINT32
		} else {
			s.grb = C.GxB_ANY_PAIR_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_PAIR_UINT8
	case uint16:
		s.grb = C.GxB_ANY_PAIR_UINT16
	case uint32:
		s.grb = C.GxB_ANY_PAIR_UINT32
	case uint64:
		s.grb = C.GxB_ANY_PAIR_UINT64
	case float32:
		s.grb = C.GxB_ANY_PAIR_FP32
	case float64:
		s.grb = C.GxB_ANY_PAIR_FP64
	case complex64:
		s.grb = C.GxB_ANY_PAIR_FC32
	case complex128:
		s.grb = C.GxB_ANY_PAIR_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinMin semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Min].
//
// MinMin is a SuiteSparse:GraphBLAS extension.
func MinMin[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_MIN_INT32
		} else {
			s.grb = C.GxB_MIN_MIN_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_MIN_INT8
	case int16:
		s.grb = C.GxB_MIN_MIN_INT16
	case int32:
		s.grb = C.GxB_MIN_MIN_INT32
	case int64:
		s.grb = C.GxB_MIN_MIN_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_MIN_UINT32
		} else {
			s.grb = C.GxB_MIN_MIN_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_MIN_UINT8
	case uint16:
		s.grb = C.GxB_MIN_MIN_UINT16
	case uint32:
		s.grb = C.GxB_MIN_MIN_UINT32
	case uint64:
		s.grb = C.GxB_MIN_MIN_UINT64
	case float32:
		s.grb = C.GxB_MIN_MIN_FP32
	case float64:
		s.grb = C.GxB_MIN_MIN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxMin semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Min].
//
// MaxMin is a SuiteSparse:GraphBLAS extension.
func MaxMin[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_MIN_INT32
		} else {
			s.grb = C.GxB_MAX_MIN_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_MIN_INT8
	case int16:
		s.grb = C.GxB_MAX_MIN_INT16
	case int32:
		s.grb = C.GxB_MAX_MIN_INT32
	case int64:
		s.grb = C.GxB_MAX_MIN_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_MIN_UINT32
		} else {
			s.grb = C.GxB_MAX_MIN_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_MIN_UINT8
	case uint16:
		s.grb = C.GxB_MAX_MIN_UINT16
	case uint32:
		s.grb = C.GxB_MAX_MIN_UINT32
	case uint64:
		s.grb = C.GxB_MAX_MIN_UINT64
	case float32:
		s.grb = C.GxB_MAX_MIN_FP32
	case float64:
		s.grb = C.GxB_MAX_MIN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusMin semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Min].
//
// PlusMin is a SuiteSparse:GraphBLAS extension.
func PlusMin[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_MIN_INT32
		} else {
			s.grb = C.GxB_PLUS_MIN_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_MIN_INT8
	case int16:
		s.grb = C.GxB_PLUS_MIN_INT16
	case int32:
		s.grb = C.GxB_PLUS_MIN_INT32
	case int64:
		s.grb = C.GxB_PLUS_MIN_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_MIN_UINT32
		} else {
			s.grb = C.GxB_PLUS_MIN_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_MIN_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_MIN_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_MIN_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_MIN_UINT64
	case float32:
		s.grb = C.GxB_PLUS_MIN_FP32
	case float64:
		s.grb = C.GxB_PLUS_MIN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesMin semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Min].
//
// TimesMin is a SuiteSparse:GraphBLAS extension.
func TimesMin[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_MIN_INT32
		} else {
			s.grb = C.GxB_TIMES_MIN_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_MIN_INT8
	case int16:
		s.grb = C.GxB_TIMES_MIN_INT16
	case int32:
		s.grb = C.GxB_TIMES_MIN_INT32
	case int64:
		s.grb = C.GxB_TIMES_MIN_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_MIN_UINT32
		} else {
			s.grb = C.GxB_TIMES_MIN_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_MIN_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_MIN_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_MIN_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_MIN_UINT64
	case float32:
		s.grb = C.GxB_TIMES_MIN_FP32
	case float64:
		s.grb = C.GxB_TIMES_MIN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyMin semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Min].
//
// AnyMin is a SuiteSparse:GraphBLAS extension.
func AnyMin[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_MIN_INT32
		} else {
			s.grb = C.GxB_ANY_MIN_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_MIN_INT8
	case int16:
		s.grb = C.GxB_ANY_MIN_INT16
	case int32:
		s.grb = C.GxB_ANY_MIN_INT32
	case int64:
		s.grb = C.GxB_ANY_MIN_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_MIN_UINT32
		} else {
			s.grb = C.GxB_ANY_MIN_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_MIN_UINT8
	case uint16:
		s.grb = C.GxB_ANY_MIN_UINT16
	case uint32:
		s.grb = C.GxB_ANY_MIN_UINT32
	case uint64:
		s.grb = C.GxB_ANY_MIN_UINT64
	case float32:
		s.grb = C.GxB_ANY_MIN_FP32
	case float64:
		s.grb = C.GxB_ANY_MIN_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinMax semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Max].
//
// MinMax is a SuiteSparse:GraphBLAS extension.
func MinMax[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_MAX_INT32
		} else {
			s.grb = C.GxB_MIN_MAX_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_MAX_INT8
	case int16:
		s.grb = C.GxB_MIN_MAX_INT16
	case int32:
		s.grb = C.GxB_MIN_MAX_INT32
	case int64:
		s.grb = C.GxB_MIN_MAX_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_MAX_UINT32
		} else {
			s.grb = C.GxB_MIN_MAX_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_MAX_UINT8
	case uint16:
		s.grb = C.GxB_MIN_MAX_UINT16
	case uint32:
		s.grb = C.GxB_MIN_MAX_UINT32
	case uint64:
		s.grb = C.GxB_MIN_MAX_UINT64
	case float32:
		s.grb = C.GxB_MIN_MAX_FP32
	case float64:
		s.grb = C.GxB_MIN_MAX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxMax semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Max].
//
// MaxMax is a SuiteSparse:GraphBLAS extension.
func MaxMax[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_MAX_INT32
		} else {
			s.grb = C.GxB_MAX_MAX_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_MAX_INT8
	case int16:
		s.grb = C.GxB_MAX_MAX_INT16
	case int32:
		s.grb = C.GxB_MAX_MAX_INT32
	case int64:
		s.grb = C.GxB_MAX_MAX_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_MAX_UINT32
		} else {
			s.grb = C.GxB_MAX_MAX_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_MAX_UINT8
	case uint16:
		s.grb = C.GxB_MAX_MAX_UINT16
	case uint32:
		s.grb = C.GxB_MAX_MAX_UINT32
	case uint64:
		s.grb = C.GxB_MAX_MAX_UINT64
	case float32:
		s.grb = C.GxB_MAX_MAX_FP32
	case float64:
		s.grb = C.GxB_MAX_MAX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusMax semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Max].
//
// PlusMax is a SuiteSparse:GraphBLAS extension.
func PlusMax[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_MAX_INT32
		} else {
			s.grb = C.GxB_PLUS_MAX_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_MAX_INT8
	case int16:
		s.grb = C.GxB_PLUS_MAX_INT16
	case int32:
		s.grb = C.GxB_PLUS_MAX_INT32
	case int64:
		s.grb = C.GxB_PLUS_MAX_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_MAX_UINT32
		} else {
			s.grb = C.GxB_PLUS_MAX_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_MAX_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_MAX_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_MAX_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_MAX_UINT64
	case float32:
		s.grb = C.GxB_PLUS_MAX_FP32
	case float64:
		s.grb = C.GxB_PLUS_MAX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesMax semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Max].
//
// TimesMax is a SuiteSparse:GraphBLAS extension.
func TimesMax[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_MAX_INT32
		} else {
			s.grb = C.GxB_TIMES_MAX_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_MAX_INT8
	case int16:
		s.grb = C.GxB_TIMES_MAX_INT16
	case int32:
		s.grb = C.GxB_TIMES_MAX_INT32
	case int64:
		s.grb = C.GxB_TIMES_MAX_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_MAX_UINT32
		} else {
			s.grb = C.GxB_TIMES_MAX_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_MAX_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_MAX_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_MAX_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_MAX_UINT64
	case float32:
		s.grb = C.GxB_TIMES_MAX_FP32
	case float64:
		s.grb = C.GxB_TIMES_MAX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyMax semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Max].
//
// AnyMax is a SuiteSparse:GraphBLAS extension.
func AnyMax[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_MAX_INT32
		} else {
			s.grb = C.GxB_ANY_MAX_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_MAX_INT8
	case int16:
		s.grb = C.GxB_ANY_MAX_INT16
	case int32:
		s.grb = C.GxB_ANY_MAX_INT32
	case int64:
		s.grb = C.GxB_ANY_MAX_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_MAX_UINT32
		} else {
			s.grb = C.GxB_ANY_MAX_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_MAX_UINT8
	case uint16:
		s.grb = C.GxB_ANY_MAX_UINT16
	case uint32:
		s.grb = C.GxB_ANY_MAX_UINT32
	case uint64:
		s.grb = C.GxB_ANY_MAX_UINT64
	case float32:
		s.grb = C.GxB_ANY_MAX_FP32
	case float64:
		s.grb = C.GxB_ANY_MAX_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinPlus semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Plus].
//
// MinPlus is a SuiteSparse:GraphBLAS extension.
func MinPlus[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_PLUS_INT32
		} else {
			s.grb = C.GxB_MIN_PLUS_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_PLUS_INT8
	case int16:
		s.grb = C.GxB_MIN_PLUS_INT16
	case int32:
		s.grb = C.GxB_MIN_PLUS_INT32
	case int64:
		s.grb = C.GxB_MIN_PLUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_PLUS_UINT32
		} else {
			s.grb = C.GxB_MIN_PLUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_PLUS_UINT8
	case uint16:
		s.grb = C.GxB_MIN_PLUS_UINT16
	case uint32:
		s.grb = C.GxB_MIN_PLUS_UINT32
	case uint64:
		s.grb = C.GxB_MIN_PLUS_UINT64
	case float32:
		s.grb = C.GxB_MIN_PLUS_FP32
	case float64:
		s.grb = C.GxB_MIN_PLUS_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxPlus semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Plus].
//
// MaxPlus is a SuiteSparse:GraphBLAS extension.
func MaxPlus[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_PLUS_INT32
		} else {
			s.grb = C.GxB_MAX_PLUS_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_PLUS_INT8
	case int16:
		s.grb = C.GxB_MAX_PLUS_INT16
	case int32:
		s.grb = C.GxB_MAX_PLUS_INT32
	case int64:
		s.grb = C.GxB_MAX_PLUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_PLUS_UINT32
		} else {
			s.grb = C.GxB_MAX_PLUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_PLUS_UINT8
	case uint16:
		s.grb = C.GxB_MAX_PLUS_UINT16
	case uint32:
		s.grb = C.GxB_MAX_PLUS_UINT32
	case uint64:
		s.grb = C.GxB_MAX_PLUS_UINT64
	case float32:
		s.grb = C.GxB_MAX_PLUS_FP32
	case float64:
		s.grb = C.GxB_MAX_PLUS_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusPlus semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Plus].
//
// PlusPlus is a SuiteSparse:GraphBLAS extension.
func PlusPlus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_PLUS_INT32
		} else {
			s.grb = C.GxB_PLUS_PLUS_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_PLUS_INT8
	case int16:
		s.grb = C.GxB_PLUS_PLUS_INT16
	case int32:
		s.grb = C.GxB_PLUS_PLUS_INT32
	case int64:
		s.grb = C.GxB_PLUS_PLUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_PLUS_UINT32
		} else {
			s.grb = C.GxB_PLUS_PLUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_PLUS_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_PLUS_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_PLUS_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_PLUS_UINT64
	case float32:
		s.grb = C.GxB_PLUS_PLUS_FP32
	case float64:
		s.grb = C.GxB_PLUS_PLUS_FP64
	case complex64:
		s.grb = C.GxB_PLUS_PLUS_FC32
	case complex128:
		s.grb = C.GxB_PLUS_PLUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesPlus semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Plus].
//
// TimesPlus is a SuiteSparse:GraphBLAS extension.
func TimesPlus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_PLUS_INT32
		} else {
			s.grb = C.GxB_TIMES_PLUS_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_PLUS_INT8
	case int16:
		s.grb = C.GxB_TIMES_PLUS_INT16
	case int32:
		s.grb = C.GxB_TIMES_PLUS_INT32
	case int64:
		s.grb = C.GxB_TIMES_PLUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_PLUS_UINT32
		} else {
			s.grb = C.GxB_TIMES_PLUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_PLUS_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_PLUS_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_PLUS_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_PLUS_UINT64
	case float32:
		s.grb = C.GxB_TIMES_PLUS_FP32
	case float64:
		s.grb = C.GxB_TIMES_PLUS_FP64
	case complex64:
		s.grb = C.GxB_TIMES_PLUS_FC32
	case complex128:
		s.grb = C.GxB_TIMES_PLUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyPlus semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Plus].
//
// AnyPlus is a SuiteSparse:GraphBLAS extension.
func AnyPlus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_PLUS_INT32
		} else {
			s.grb = C.GxB_ANY_PLUS_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_PLUS_INT8
	case int16:
		s.grb = C.GxB_ANY_PLUS_INT16
	case int32:
		s.grb = C.GxB_ANY_PLUS_INT32
	case int64:
		s.grb = C.GxB_ANY_PLUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_PLUS_UINT32
		} else {
			s.grb = C.GxB_ANY_PLUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_PLUS_UINT8
	case uint16:
		s.grb = C.GxB_ANY_PLUS_UINT16
	case uint32:
		s.grb = C.GxB_ANY_PLUS_UINT32
	case uint64:
		s.grb = C.GxB_ANY_PLUS_UINT64
	case float32:
		s.grb = C.GxB_ANY_PLUS_FP32
	case float64:
		s.grb = C.GxB_ANY_PLUS_FP64
	case complex64:
		s.grb = C.GxB_ANY_PLUS_FC32
	case complex128:
		s.grb = C.GxB_ANY_PLUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinMinus semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Minus].
//
// MinMinus is a SuiteSparse:GraphBLAS extension.
func MinMinus[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_MINUS_INT32
		} else {
			s.grb = C.GxB_MIN_MINUS_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_MINUS_INT8
	case int16:
		s.grb = C.GxB_MIN_MINUS_INT16
	case int32:
		s.grb = C.GxB_MIN_MINUS_INT32
	case int64:
		s.grb = C.GxB_MIN_MINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_MINUS_UINT32
		} else {
			s.grb = C.GxB_MIN_MINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_MINUS_UINT8
	case uint16:
		s.grb = C.GxB_MIN_MINUS_UINT16
	case uint32:
		s.grb = C.GxB_MIN_MINUS_UINT32
	case uint64:
		s.grb = C.GxB_MIN_MINUS_UINT64
	case float32:
		s.grb = C.GxB_MIN_MINUS_FP32
	case float64:
		s.grb = C.GxB_MIN_MINUS_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxMinus semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Minus].
//
// MaxMinus is a SuiteSparse:GraphBLAS extension.
func MaxMinus[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_MINUS_INT32
		} else {
			s.grb = C.GxB_MAX_MINUS_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_MINUS_INT8
	case int16:
		s.grb = C.GxB_MAX_MINUS_INT16
	case int32:
		s.grb = C.GxB_MAX_MINUS_INT32
	case int64:
		s.grb = C.GxB_MAX_MINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_MINUS_UINT32
		} else {
			s.grb = C.GxB_MAX_MINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_MINUS_UINT8
	case uint16:
		s.grb = C.GxB_MAX_MINUS_UINT16
	case uint32:
		s.grb = C.GxB_MAX_MINUS_UINT32
	case uint64:
		s.grb = C.GxB_MAX_MINUS_UINT64
	case float32:
		s.grb = C.GxB_MAX_MINUS_FP32
	case float64:
		s.grb = C.GxB_MAX_MINUS_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusMinus semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Minus].
//
// PlusMinus is a SuiteSparse:GraphBLAS extension.
func PlusMinus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_MINUS_INT32
		} else {
			s.grb = C.GxB_PLUS_MINUS_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_MINUS_INT8
	case int16:
		s.grb = C.GxB_PLUS_MINUS_INT16
	case int32:
		s.grb = C.GxB_PLUS_MINUS_INT32
	case int64:
		s.grb = C.GxB_PLUS_MINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_MINUS_UINT32
		} else {
			s.grb = C.GxB_PLUS_MINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_MINUS_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_MINUS_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_MINUS_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_MINUS_UINT64
	case float32:
		s.grb = C.GxB_PLUS_MINUS_FP32
	case float64:
		s.grb = C.GxB_PLUS_MINUS_FP64
	case complex64:
		s.grb = C.GxB_PLUS_MINUS_FC32
	case complex128:
		s.grb = C.GxB_PLUS_MINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesMinus semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Minus].
//
// TimesMinus is a SuiteSparse:GraphBLAS extension.
func TimesMinus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_MINUS_INT32
		} else {
			s.grb = C.GxB_TIMES_MINUS_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_MINUS_INT8
	case int16:
		s.grb = C.GxB_TIMES_MINUS_INT16
	case int32:
		s.grb = C.GxB_TIMES_MINUS_INT32
	case int64:
		s.grb = C.GxB_TIMES_MINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_MINUS_UINT32
		} else {
			s.grb = C.GxB_TIMES_MINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_MINUS_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_MINUS_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_MINUS_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_MINUS_UINT64
	case float32:
		s.grb = C.GxB_TIMES_MINUS_FP32
	case float64:
		s.grb = C.GxB_TIMES_MINUS_FP64
	case complex64:
		s.grb = C.GxB_TIMES_MINUS_FC32
	case complex128:
		s.grb = C.GxB_TIMES_MINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyMinus semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Minus].
//
// AnyMinus is a SuiteSparse:GraphBLAS extension.
func AnyMinus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_MINUS_INT32
		} else {
			s.grb = C.GxB_ANY_MINUS_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_MINUS_INT8
	case int16:
		s.grb = C.GxB_ANY_MINUS_INT16
	case int32:
		s.grb = C.GxB_ANY_MINUS_INT32
	case int64:
		s.grb = C.GxB_ANY_MINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_MINUS_UINT32
		} else {
			s.grb = C.GxB_ANY_MINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_MINUS_UINT8
	case uint16:
		s.grb = C.GxB_ANY_MINUS_UINT16
	case uint32:
		s.grb = C.GxB_ANY_MINUS_UINT32
	case uint64:
		s.grb = C.GxB_ANY_MINUS_UINT64
	case float32:
		s.grb = C.GxB_ANY_MINUS_FP32
	case float64:
		s.grb = C.GxB_ANY_MINUS_FP64
	case complex64:
		s.grb = C.GxB_ANY_MINUS_FC32
	case complex128:
		s.grb = C.GxB_ANY_MINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinTimes semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Times].
//
// MinTimes is a SuiteSparse:GraphBLAS extension.
func MinTimes[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_TIMES_INT32
		} else {
			s.grb = C.GxB_MIN_TIMES_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_TIMES_INT8
	case int16:
		s.grb = C.GxB_MIN_TIMES_INT16
	case int32:
		s.grb = C.GxB_MIN_TIMES_INT32
	case int64:
		s.grb = C.GxB_MIN_TIMES_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_TIMES_UINT32
		} else {
			s.grb = C.GxB_MIN_TIMES_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_TIMES_UINT8
	case uint16:
		s.grb = C.GxB_MIN_TIMES_UINT16
	case uint32:
		s.grb = C.GxB_MIN_TIMES_UINT32
	case uint64:
		s.grb = C.GxB_MIN_TIMES_UINT64
	case float32:
		s.grb = C.GxB_MIN_TIMES_FP32
	case float64:
		s.grb = C.GxB_MIN_TIMES_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxTimes semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Times].
//
// MaxTimes is a SuiteSparse:GraphBLAS extension.
func MaxTimes[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_TIMES_INT32
		} else {
			s.grb = C.GxB_MAX_TIMES_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_TIMES_INT8
	case int16:
		s.grb = C.GxB_MAX_TIMES_INT16
	case int32:
		s.grb = C.GxB_MAX_TIMES_INT32
	case int64:
		s.grb = C.GxB_MAX_TIMES_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_TIMES_UINT32
		} else {
			s.grb = C.GxB_MAX_TIMES_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_TIMES_UINT8
	case uint16:
		s.grb = C.GxB_MAX_TIMES_UINT16
	case uint32:
		s.grb = C.GxB_MAX_TIMES_UINT32
	case uint64:
		s.grb = C.GxB_MAX_TIMES_UINT64
	case float32:
		s.grb = C.GxB_MAX_TIMES_FP32
	case float64:
		s.grb = C.GxB_MAX_TIMES_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusTimes semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Times].
//
// PlusTimes is a SuiteSparse:GraphBLAS extension.
func PlusTimes[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_TIMES_INT32
		} else {
			s.grb = C.GxB_PLUS_TIMES_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_TIMES_INT8
	case int16:
		s.grb = C.GxB_PLUS_TIMES_INT16
	case int32:
		s.grb = C.GxB_PLUS_TIMES_INT32
	case int64:
		s.grb = C.GxB_PLUS_TIMES_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_TIMES_UINT32
		} else {
			s.grb = C.GxB_PLUS_TIMES_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_TIMES_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_TIMES_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_TIMES_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_TIMES_UINT64
	case float32:
		s.grb = C.GxB_PLUS_TIMES_FP32
	case float64:
		s.grb = C.GxB_PLUS_TIMES_FP64
	case complex64:
		s.grb = C.GxB_PLUS_TIMES_FC32
	case complex128:
		s.grb = C.GxB_PLUS_TIMES_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesTimes semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Times].
//
// TimesTimes is a SuiteSparse:GraphBLAS extension.
func TimesTimes[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_TIMES_INT32
		} else {
			s.grb = C.GxB_TIMES_TIMES_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_TIMES_INT8
	case int16:
		s.grb = C.GxB_TIMES_TIMES_INT16
	case int32:
		s.grb = C.GxB_TIMES_TIMES_INT32
	case int64:
		s.grb = C.GxB_TIMES_TIMES_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_TIMES_UINT32
		} else {
			s.grb = C.GxB_TIMES_TIMES_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_TIMES_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_TIMES_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_TIMES_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_TIMES_UINT64
	case float32:
		s.grb = C.GxB_TIMES_TIMES_FP32
	case float64:
		s.grb = C.GxB_TIMES_TIMES_FP64
	case complex64:
		s.grb = C.GxB_TIMES_TIMES_FC32
	case complex128:
		s.grb = C.GxB_TIMES_TIMES_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyTimes semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Times].
//
// AnyTimes is a SuiteSparse:GraphBLAS extension.
func AnyTimes[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_TIMES_INT32
		} else {
			s.grb = C.GxB_ANY_TIMES_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_TIMES_INT8
	case int16:
		s.grb = C.GxB_ANY_TIMES_INT16
	case int32:
		s.grb = C.GxB_ANY_TIMES_INT32
	case int64:
		s.grb = C.GxB_ANY_TIMES_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_TIMES_UINT32
		} else {
			s.grb = C.GxB_ANY_TIMES_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_TIMES_UINT8
	case uint16:
		s.grb = C.GxB_ANY_TIMES_UINT16
	case uint32:
		s.grb = C.GxB_ANY_TIMES_UINT32
	case uint64:
		s.grb = C.GxB_ANY_TIMES_UINT64
	case float32:
		s.grb = C.GxB_ANY_TIMES_FP32
	case float64:
		s.grb = C.GxB_ANY_TIMES_FP64
	case complex64:
		s.grb = C.GxB_ANY_TIMES_FC32
	case complex128:
		s.grb = C.GxB_ANY_TIMES_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinDiv semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Div].
//
// MinDiv is a SuiteSparse:GraphBLAS extension.
func MinDiv[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_DIV_INT32
		} else {
			s.grb = C.GxB_MIN_DIV_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_DIV_INT8
	case int16:
		s.grb = C.GxB_MIN_DIV_INT16
	case int32:
		s.grb = C.GxB_MIN_DIV_INT32
	case int64:
		s.grb = C.GxB_MIN_DIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_DIV_UINT32
		} else {
			s.grb = C.GxB_MIN_DIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_DIV_UINT8
	case uint16:
		s.grb = C.GxB_MIN_DIV_UINT16
	case uint32:
		s.grb = C.GxB_MIN_DIV_UINT32
	case uint64:
		s.grb = C.GxB_MIN_DIV_UINT64
	case float32:
		s.grb = C.GxB_MIN_DIV_FP32
	case float64:
		s.grb = C.GxB_MIN_DIV_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxDiv semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Div].
//
// MaxDiv is a SuiteSparse:GraphBLAS extension.
func MaxDiv[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_DIV_INT32
		} else {
			s.grb = C.GxB_MAX_DIV_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_DIV_INT8
	case int16:
		s.grb = C.GxB_MAX_DIV_INT16
	case int32:
		s.grb = C.GxB_MAX_DIV_INT32
	case int64:
		s.grb = C.GxB_MAX_DIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_DIV_UINT32
		} else {
			s.grb = C.GxB_MAX_DIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_DIV_UINT8
	case uint16:
		s.grb = C.GxB_MAX_DIV_UINT16
	case uint32:
		s.grb = C.GxB_MAX_DIV_UINT32
	case uint64:
		s.grb = C.GxB_MAX_DIV_UINT64
	case float32:
		s.grb = C.GxB_MAX_DIV_FP32
	case float64:
		s.grb = C.GxB_MAX_DIV_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusDiv semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Div].
//
// PlusDiv is a SuiteSparse:GraphBLAS extension.
func PlusDiv[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_DIV_INT32
		} else {
			s.grb = C.GxB_PLUS_DIV_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_DIV_INT8
	case int16:
		s.grb = C.GxB_PLUS_DIV_INT16
	case int32:
		s.grb = C.GxB_PLUS_DIV_INT32
	case int64:
		s.grb = C.GxB_PLUS_DIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_DIV_UINT32
		} else {
			s.grb = C.GxB_PLUS_DIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_DIV_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_DIV_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_DIV_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_DIV_UINT64
	case float32:
		s.grb = C.GxB_PLUS_DIV_FP32
	case float64:
		s.grb = C.GxB_PLUS_DIV_FP64
	case complex64:
		s.grb = C.GxB_PLUS_DIV_FC32
	case complex128:
		s.grb = C.GxB_PLUS_DIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesDiv semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Div].
//
// TimesDiv is a SuiteSparse:GraphBLAS extension.
func TimesDiv[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_DIV_INT32
		} else {
			s.grb = C.GxB_TIMES_DIV_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_DIV_INT8
	case int16:
		s.grb = C.GxB_TIMES_DIV_INT16
	case int32:
		s.grb = C.GxB_TIMES_DIV_INT32
	case int64:
		s.grb = C.GxB_TIMES_DIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_DIV_UINT32
		} else {
			s.grb = C.GxB_TIMES_DIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_DIV_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_DIV_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_DIV_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_DIV_UINT64
	case float32:
		s.grb = C.GxB_TIMES_DIV_FP32
	case float64:
		s.grb = C.GxB_TIMES_DIV_FP64
	case complex64:
		s.grb = C.GxB_TIMES_DIV_FC32
	case complex128:
		s.grb = C.GxB_TIMES_DIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyDiv semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Div].
//
// AnyDiv is a SuiteSparse:GraphBLAS extension.
func AnyDiv[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_DIV_INT32
		} else {
			s.grb = C.GxB_ANY_DIV_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_DIV_INT8
	case int16:
		s.grb = C.GxB_ANY_DIV_INT16
	case int32:
		s.grb = C.GxB_ANY_DIV_INT32
	case int64:
		s.grb = C.GxB_ANY_DIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_DIV_UINT32
		} else {
			s.grb = C.GxB_ANY_DIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_DIV_UINT8
	case uint16:
		s.grb = C.GxB_ANY_DIV_UINT16
	case uint32:
		s.grb = C.GxB_ANY_DIV_UINT32
	case uint64:
		s.grb = C.GxB_ANY_DIV_UINT64
	case float32:
		s.grb = C.GxB_ANY_DIV_FP32
	case float64:
		s.grb = C.GxB_ANY_DIV_FP64
	case complex64:
		s.grb = C.GxB_ANY_DIV_FC32
	case complex128:
		s.grb = C.GxB_ANY_DIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinRdiv semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Rdiv].
//
// MinRdiv is a SuiteSparse:GraphBLAS extension.
func MinRdiv[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_RDIV_INT32
		} else {
			s.grb = C.GxB_MIN_RDIV_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_RDIV_INT8
	case int16:
		s.grb = C.GxB_MIN_RDIV_INT16
	case int32:
		s.grb = C.GxB_MIN_RDIV_INT32
	case int64:
		s.grb = C.GxB_MIN_RDIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_RDIV_UINT32
		} else {
			s.grb = C.GxB_MIN_RDIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_RDIV_UINT8
	case uint16:
		s.grb = C.GxB_MIN_RDIV_UINT16
	case uint32:
		s.grb = C.GxB_MIN_RDIV_UINT32
	case uint64:
		s.grb = C.GxB_MIN_RDIV_UINT64
	case float32:
		s.grb = C.GxB_MIN_RDIV_FP32
	case float64:
		s.grb = C.GxB_MIN_RDIV_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxRdiv semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Rdiv].
//
// MaxRdiv is a SuiteSparse:GraphBLAS extension.
func MaxRdiv[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_RDIV_INT32
		} else {
			s.grb = C.GxB_MAX_RDIV_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_RDIV_INT8
	case int16:
		s.grb = C.GxB_MAX_RDIV_INT16
	case int32:
		s.grb = C.GxB_MAX_RDIV_INT32
	case int64:
		s.grb = C.GxB_MAX_RDIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_RDIV_UINT32
		} else {
			s.grb = C.GxB_MAX_RDIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_RDIV_UINT8
	case uint16:
		s.grb = C.GxB_MAX_RDIV_UINT16
	case uint32:
		s.grb = C.GxB_MAX_RDIV_UINT32
	case uint64:
		s.grb = C.GxB_MAX_RDIV_UINT64
	case float32:
		s.grb = C.GxB_MAX_RDIV_FP32
	case float64:
		s.grb = C.GxB_MAX_RDIV_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusRdiv semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Rdiv].
//
// PlusRdiv is a SuiteSparse:GraphBLAS extension.
func PlusRdiv[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_RDIV_INT32
		} else {
			s.grb = C.GxB_PLUS_RDIV_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_RDIV_INT8
	case int16:
		s.grb = C.GxB_PLUS_RDIV_INT16
	case int32:
		s.grb = C.GxB_PLUS_RDIV_INT32
	case int64:
		s.grb = C.GxB_PLUS_RDIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_RDIV_UINT32
		} else {
			s.grb = C.GxB_PLUS_RDIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_RDIV_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_RDIV_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_RDIV_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_RDIV_UINT64
	case float32:
		s.grb = C.GxB_PLUS_RDIV_FP32
	case float64:
		s.grb = C.GxB_PLUS_RDIV_FP64
	case complex64:
		s.grb = C.GxB_PLUS_RDIV_FC32
	case complex128:
		s.grb = C.GxB_PLUS_RDIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesRdiv semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Rdiv].
//
// TimesRdiv is a SuiteSparse:GraphBLAS extension.
func TimesRdiv[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_RDIV_INT32
		} else {
			s.grb = C.GxB_TIMES_RDIV_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_RDIV_INT8
	case int16:
		s.grb = C.GxB_TIMES_RDIV_INT16
	case int32:
		s.grb = C.GxB_TIMES_RDIV_INT32
	case int64:
		s.grb = C.GxB_TIMES_RDIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_RDIV_UINT32
		} else {
			s.grb = C.GxB_TIMES_RDIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_RDIV_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_RDIV_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_RDIV_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_RDIV_UINT64
	case float32:
		s.grb = C.GxB_TIMES_RDIV_FP32
	case float64:
		s.grb = C.GxB_TIMES_RDIV_FP64
	case complex64:
		s.grb = C.GxB_TIMES_RDIV_FC32
	case complex128:
		s.grb = C.GxB_TIMES_RDIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyRdiv semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Rdiv].
//
// AnyRdiv is a SuiteSparse:GraphBLAS extension.
func AnyRdiv[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_RDIV_INT32
		} else {
			s.grb = C.GxB_ANY_RDIV_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_RDIV_INT8
	case int16:
		s.grb = C.GxB_ANY_RDIV_INT16
	case int32:
		s.grb = C.GxB_ANY_RDIV_INT32
	case int64:
		s.grb = C.GxB_ANY_RDIV_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_RDIV_UINT32
		} else {
			s.grb = C.GxB_ANY_RDIV_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_RDIV_UINT8
	case uint16:
		s.grb = C.GxB_ANY_RDIV_UINT16
	case uint32:
		s.grb = C.GxB_ANY_RDIV_UINT32
	case uint64:
		s.grb = C.GxB_ANY_RDIV_UINT64
	case float32:
		s.grb = C.GxB_ANY_RDIV_FP32
	case float64:
		s.grb = C.GxB_ANY_RDIV_FP64
	case complex64:
		s.grb = C.GxB_ANY_RDIV_FC32
	case complex128:
		s.grb = C.GxB_ANY_RDIV_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinRminus semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Rminus].
//
// MinRminus is a SuiteSparse:GraphBLAS extension.
func MinRminus[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_RMINUS_INT32
		} else {
			s.grb = C.GxB_MIN_RMINUS_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_RMINUS_INT8
	case int16:
		s.grb = C.GxB_MIN_RMINUS_INT16
	case int32:
		s.grb = C.GxB_MIN_RMINUS_INT32
	case int64:
		s.grb = C.GxB_MIN_RMINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_RMINUS_UINT32
		} else {
			s.grb = C.GxB_MIN_RMINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_RMINUS_UINT8
	case uint16:
		s.grb = C.GxB_MIN_RMINUS_UINT16
	case uint32:
		s.grb = C.GxB_MIN_RMINUS_UINT32
	case uint64:
		s.grb = C.GxB_MIN_RMINUS_UINT64
	case float32:
		s.grb = C.GxB_MIN_RMINUS_FP32
	case float64:
		s.grb = C.GxB_MIN_RMINUS_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxRminus semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Rminus].
//
// MaxRminus is a SuiteSparse:GraphBLAS extension.
func MaxRminus[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_RMINUS_INT32
		} else {
			s.grb = C.GxB_MAX_RMINUS_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_RMINUS_INT8
	case int16:
		s.grb = C.GxB_MAX_RMINUS_INT16
	case int32:
		s.grb = C.GxB_MAX_RMINUS_INT32
	case int64:
		s.grb = C.GxB_MAX_RMINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_RMINUS_UINT32
		} else {
			s.grb = C.GxB_MAX_RMINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_RMINUS_UINT8
	case uint16:
		s.grb = C.GxB_MAX_RMINUS_UINT16
	case uint32:
		s.grb = C.GxB_MAX_RMINUS_UINT32
	case uint64:
		s.grb = C.GxB_MAX_RMINUS_UINT64
	case float32:
		s.grb = C.GxB_MAX_RMINUS_FP32
	case float64:
		s.grb = C.GxB_MAX_RMINUS_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusRminus semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Rminus].
//
// PlusRminus is a SuiteSparse:GraphBLAS extension.
func PlusRminus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_RMINUS_INT32
		} else {
			s.grb = C.GxB_PLUS_RMINUS_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_RMINUS_INT8
	case int16:
		s.grb = C.GxB_PLUS_RMINUS_INT16
	case int32:
		s.grb = C.GxB_PLUS_RMINUS_INT32
	case int64:
		s.grb = C.GxB_PLUS_RMINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_RMINUS_UINT32
		} else {
			s.grb = C.GxB_PLUS_RMINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_RMINUS_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_RMINUS_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_RMINUS_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_RMINUS_UINT64
	case float32:
		s.grb = C.GxB_PLUS_RMINUS_FP32
	case float64:
		s.grb = C.GxB_PLUS_RMINUS_FP64
	case complex64:
		s.grb = C.GxB_PLUS_RMINUS_FC32
	case complex128:
		s.grb = C.GxB_PLUS_RMINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// TimesRminus semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Rminus].
//
// TimesRminus is a SuiteSparse:GraphBLAS extension.
func TimesRminus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_RMINUS_INT32
		} else {
			s.grb = C.GxB_TIMES_RMINUS_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_RMINUS_INT8
	case int16:
		s.grb = C.GxB_TIMES_RMINUS_INT16
	case int32:
		s.grb = C.GxB_TIMES_RMINUS_INT32
	case int64:
		s.grb = C.GxB_TIMES_RMINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_RMINUS_UINT32
		} else {
			s.grb = C.GxB_TIMES_RMINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_RMINUS_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_RMINUS_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_RMINUS_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_RMINUS_UINT64
	case float32:
		s.grb = C.GxB_TIMES_RMINUS_FP32
	case float64:
		s.grb = C.GxB_TIMES_RMINUS_FP64
	case complex64:
		s.grb = C.GxB_TIMES_RMINUS_FC32
	case complex128:
		s.grb = C.GxB_TIMES_RMINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// AnyRminus semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Rminus].
//
// AnyRminus is a SuiteSparse:GraphBLAS extension.
func AnyRminus[D Number | Complex]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_RMINUS_INT32
		} else {
			s.grb = C.GxB_ANY_RMINUS_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_RMINUS_INT8
	case int16:
		s.grb = C.GxB_ANY_RMINUS_INT16
	case int32:
		s.grb = C.GxB_ANY_RMINUS_INT32
	case int64:
		s.grb = C.GxB_ANY_RMINUS_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_RMINUS_UINT32
		} else {
			s.grb = C.GxB_ANY_RMINUS_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_RMINUS_UINT8
	case uint16:
		s.grb = C.GxB_ANY_RMINUS_UINT16
	case uint32:
		s.grb = C.GxB_ANY_RMINUS_UINT32
	case uint64:
		s.grb = C.GxB_ANY_RMINUS_UINT64
	case float32:
		s.grb = C.GxB_ANY_RMINUS_FP32
	case float64:
		s.grb = C.GxB_ANY_RMINUS_FP64
	case complex64:
		s.grb = C.GxB_ANY_RMINUS_FC32
	case complex128:
		s.grb = C.GxB_ANY_RMINUS_FC64
	default:
		panic("unreachable code")
	}
	return
}

// MinIseq semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Iseq].
//
// MinIseq is a SuiteSparse:GraphBLAS extension.
func MinIseq[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISEQ_INT32
		} else {
			s.grb = C.GxB_MIN_ISEQ_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_ISEQ_INT8
	case int16:
		s.grb = C.GxB_MIN_ISEQ_INT16
	case int32:
		s.grb = C.GxB_MIN_ISEQ_INT32
	case int64:
		s.grb = C.GxB_MIN_ISEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISEQ_UINT32
		} else {
			s.grb = C.GxB_MIN_ISEQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_ISEQ_UINT8
	case uint16:
		s.grb = C.GxB_MIN_ISEQ_UINT16
	case uint32:
		s.grb = C.GxB_MIN_ISEQ_UINT32
	case uint64:
		s.grb = C.GxB_MIN_ISEQ_UINT64
	case float32:
		s.grb = C.GxB_MIN_ISEQ_FP32
	case float64:
		s.grb = C.GxB_MIN_ISEQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxIseq semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Iseq].
//
// MaxIseq is a SuiteSparse:GraphBLAS extension.
func MaxIseq[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISEQ_INT32
		} else {
			s.grb = C.GxB_MAX_ISEQ_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_ISEQ_INT8
	case int16:
		s.grb = C.GxB_MAX_ISEQ_INT16
	case int32:
		s.grb = C.GxB_MAX_ISEQ_INT32
	case int64:
		s.grb = C.GxB_MAX_ISEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISEQ_UINT32
		} else {
			s.grb = C.GxB_MAX_ISEQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_ISEQ_UINT8
	case uint16:
		s.grb = C.GxB_MAX_ISEQ_UINT16
	case uint32:
		s.grb = C.GxB_MAX_ISEQ_UINT32
	case uint64:
		s.grb = C.GxB_MAX_ISEQ_UINT64
	case float32:
		s.grb = C.GxB_MAX_ISEQ_FP32
	case float64:
		s.grb = C.GxB_MAX_ISEQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusIseq semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Iseq].
//
// PlusIseq is a SuiteSparse:GraphBLAS extension.
func PlusIseq[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISEQ_INT32
		} else {
			s.grb = C.GxB_PLUS_ISEQ_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_ISEQ_INT8
	case int16:
		s.grb = C.GxB_PLUS_ISEQ_INT16
	case int32:
		s.grb = C.GxB_PLUS_ISEQ_INT32
	case int64:
		s.grb = C.GxB_PLUS_ISEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISEQ_UINT32
		} else {
			s.grb = C.GxB_PLUS_ISEQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_ISEQ_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_ISEQ_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_ISEQ_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_ISEQ_UINT64
	case float32:
		s.grb = C.GxB_PLUS_ISEQ_FP32
	case float64:
		s.grb = C.GxB_PLUS_ISEQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesIseq semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Iseq].
//
// TimesIseq is a SuiteSparse:GraphBLAS extension.
func TimesIseq[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISEQ_INT32
		} else {
			s.grb = C.GxB_TIMES_ISEQ_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_ISEQ_INT8
	case int16:
		s.grb = C.GxB_TIMES_ISEQ_INT16
	case int32:
		s.grb = C.GxB_TIMES_ISEQ_INT32
	case int64:
		s.grb = C.GxB_TIMES_ISEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISEQ_UINT32
		} else {
			s.grb = C.GxB_TIMES_ISEQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_ISEQ_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_ISEQ_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_ISEQ_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_ISEQ_UINT64
	case float32:
		s.grb = C.GxB_TIMES_ISEQ_FP32
	case float64:
		s.grb = C.GxB_TIMES_ISEQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyIseq semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Iseq].
//
// AnyIseq is a SuiteSparse:GraphBLAS extension.
func AnyIseq[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISEQ_INT32
		} else {
			s.grb = C.GxB_ANY_ISEQ_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_ISEQ_INT8
	case int16:
		s.grb = C.GxB_ANY_ISEQ_INT16
	case int32:
		s.grb = C.GxB_ANY_ISEQ_INT32
	case int64:
		s.grb = C.GxB_ANY_ISEQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISEQ_UINT32
		} else {
			s.grb = C.GxB_ANY_ISEQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_ISEQ_UINT8
	case uint16:
		s.grb = C.GxB_ANY_ISEQ_UINT16
	case uint32:
		s.grb = C.GxB_ANY_ISEQ_UINT32
	case uint64:
		s.grb = C.GxB_ANY_ISEQ_UINT64
	case float32:
		s.grb = C.GxB_ANY_ISEQ_FP32
	case float64:
		s.grb = C.GxB_ANY_ISEQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinIsne semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Isne].
//
// MinIsne is a SuiteSparse:GraphBLAS extension.
func MinIsne[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISNE_INT32
		} else {
			s.grb = C.GxB_MIN_ISNE_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_ISNE_INT8
	case int16:
		s.grb = C.GxB_MIN_ISNE_INT16
	case int32:
		s.grb = C.GxB_MIN_ISNE_INT32
	case int64:
		s.grb = C.GxB_MIN_ISNE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISNE_UINT32
		} else {
			s.grb = C.GxB_MIN_ISNE_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_ISNE_UINT8
	case uint16:
		s.grb = C.GxB_MIN_ISNE_UINT16
	case uint32:
		s.grb = C.GxB_MIN_ISNE_UINT32
	case uint64:
		s.grb = C.GxB_MIN_ISNE_UINT64
	case float32:
		s.grb = C.GxB_MIN_ISNE_FP32
	case float64:
		s.grb = C.GxB_MIN_ISNE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxIsne semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Isne].
//
// MaxIsne is a SuiteSparse:GraphBLAS extension.
func MaxIsne[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISNE_INT32
		} else {
			s.grb = C.GxB_MAX_ISNE_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_ISNE_INT8
	case int16:
		s.grb = C.GxB_MAX_ISNE_INT16
	case int32:
		s.grb = C.GxB_MAX_ISNE_INT32
	case int64:
		s.grb = C.GxB_MAX_ISNE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISNE_UINT32
		} else {
			s.grb = C.GxB_MAX_ISNE_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_ISNE_UINT8
	case uint16:
		s.grb = C.GxB_MAX_ISNE_UINT16
	case uint32:
		s.grb = C.GxB_MAX_ISNE_UINT32
	case uint64:
		s.grb = C.GxB_MAX_ISNE_UINT64
	case float32:
		s.grb = C.GxB_MAX_ISNE_FP32
	case float64:
		s.grb = C.GxB_MAX_ISNE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusIsne semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Isne].
//
// PlusIsne is a SuiteSparse:GraphBLAS extension.
func PlusIsne[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISNE_INT32
		} else {
			s.grb = C.GxB_PLUS_ISNE_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_ISNE_INT8
	case int16:
		s.grb = C.GxB_PLUS_ISNE_INT16
	case int32:
		s.grb = C.GxB_PLUS_ISNE_INT32
	case int64:
		s.grb = C.GxB_PLUS_ISNE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISNE_UINT32
		} else {
			s.grb = C.GxB_PLUS_ISNE_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_ISNE_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_ISNE_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_ISNE_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_ISNE_UINT64
	case float32:
		s.grb = C.GxB_PLUS_ISNE_FP32
	case float64:
		s.grb = C.GxB_PLUS_ISNE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesIsne semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Isne].
//
// TimesIsne is a SuiteSparse:GraphBLAS extension.
func TimesIsne[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISNE_INT32
		} else {
			s.grb = C.GxB_TIMES_ISNE_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_ISNE_INT8
	case int16:
		s.grb = C.GxB_TIMES_ISNE_INT16
	case int32:
		s.grb = C.GxB_TIMES_ISNE_INT32
	case int64:
		s.grb = C.GxB_TIMES_ISNE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISNE_UINT32
		} else {
			s.grb = C.GxB_TIMES_ISNE_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_ISNE_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_ISNE_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_ISNE_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_ISNE_UINT64
	case float32:
		s.grb = C.GxB_TIMES_ISNE_FP32
	case float64:
		s.grb = C.GxB_TIMES_ISNE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyIsne semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Isne].
//
// AnyIsne is a SuiteSparse:GraphBLAS extension.
func AnyIsne[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISNE_INT32
		} else {
			s.grb = C.GxB_ANY_ISNE_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_ISNE_INT8
	case int16:
		s.grb = C.GxB_ANY_ISNE_INT16
	case int32:
		s.grb = C.GxB_ANY_ISNE_INT32
	case int64:
		s.grb = C.GxB_ANY_ISNE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISNE_UINT32
		} else {
			s.grb = C.GxB_ANY_ISNE_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_ISNE_UINT8
	case uint16:
		s.grb = C.GxB_ANY_ISNE_UINT16
	case uint32:
		s.grb = C.GxB_ANY_ISNE_UINT32
	case uint64:
		s.grb = C.GxB_ANY_ISNE_UINT64
	case float32:
		s.grb = C.GxB_ANY_ISNE_FP32
	case float64:
		s.grb = C.GxB_ANY_ISNE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinIsgt semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Isgt].
//
// MinIsgt is a SuiteSparse:GraphBLAS extension.
func MinIsgt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISGT_INT32
		} else {
			s.grb = C.GxB_MIN_ISGT_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_ISGT_INT8
	case int16:
		s.grb = C.GxB_MIN_ISGT_INT16
	case int32:
		s.grb = C.GxB_MIN_ISGT_INT32
	case int64:
		s.grb = C.GxB_MIN_ISGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISGT_UINT32
		} else {
			s.grb = C.GxB_MIN_ISGT_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_ISGT_UINT8
	case uint16:
		s.grb = C.GxB_MIN_ISGT_UINT16
	case uint32:
		s.grb = C.GxB_MIN_ISGT_UINT32
	case uint64:
		s.grb = C.GxB_MIN_ISGT_UINT64
	case float32:
		s.grb = C.GxB_MIN_ISGT_FP32
	case float64:
		s.grb = C.GxB_MIN_ISGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxIsgt semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Isgt].
//
// MaxIsgt is a SuiteSparse:GraphBLAS extension.
func MaxIsgt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISGT_INT32
		} else {
			s.grb = C.GxB_MAX_ISGT_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_ISGT_INT8
	case int16:
		s.grb = C.GxB_MAX_ISGT_INT16
	case int32:
		s.grb = C.GxB_MAX_ISGT_INT32
	case int64:
		s.grb = C.GxB_MAX_ISGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISGT_UINT32
		} else {
			s.grb = C.GxB_MAX_ISGT_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_ISGT_UINT8
	case uint16:
		s.grb = C.GxB_MAX_ISGT_UINT16
	case uint32:
		s.grb = C.GxB_MAX_ISGT_UINT32
	case uint64:
		s.grb = C.GxB_MAX_ISGT_UINT64
	case float32:
		s.grb = C.GxB_MAX_ISGT_FP32
	case float64:
		s.grb = C.GxB_MAX_ISGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusIsgt semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Isgt].
//
// PlusIsgt is a SuiteSparse:GraphBLAS extension.
func PlusIsgt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISGT_INT32
		} else {
			s.grb = C.GxB_PLUS_ISGT_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_ISGT_INT8
	case int16:
		s.grb = C.GxB_PLUS_ISGT_INT16
	case int32:
		s.grb = C.GxB_PLUS_ISGT_INT32
	case int64:
		s.grb = C.GxB_PLUS_ISGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISGT_UINT32
		} else {
			s.grb = C.GxB_PLUS_ISGT_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_ISGT_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_ISGT_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_ISGT_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_ISGT_UINT64
	case float32:
		s.grb = C.GxB_PLUS_ISGT_FP32
	case float64:
		s.grb = C.GxB_PLUS_ISGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesIsgt semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Isgt].
//
// TimesIsgt is a SuiteSparse:GraphBLAS extension.
func TimesIsgt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISGT_INT32
		} else {
			s.grb = C.GxB_TIMES_ISGT_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_ISGT_INT8
	case int16:
		s.grb = C.GxB_TIMES_ISGT_INT16
	case int32:
		s.grb = C.GxB_TIMES_ISGT_INT32
	case int64:
		s.grb = C.GxB_TIMES_ISGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISGT_UINT32
		} else {
			s.grb = C.GxB_TIMES_ISGT_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_ISGT_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_ISGT_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_ISGT_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_ISGT_UINT64
	case float32:
		s.grb = C.GxB_TIMES_ISGT_FP32
	case float64:
		s.grb = C.GxB_TIMES_ISGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyIsgt semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Isgt].
//
// AnyIsgt is a SuiteSparse:GraphBLAS extension.
func AnyIsgt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISGT_INT32
		} else {
			s.grb = C.GxB_ANY_ISGT_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_ISGT_INT8
	case int16:
		s.grb = C.GxB_ANY_ISGT_INT16
	case int32:
		s.grb = C.GxB_ANY_ISGT_INT32
	case int64:
		s.grb = C.GxB_ANY_ISGT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISGT_UINT32
		} else {
			s.grb = C.GxB_ANY_ISGT_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_ISGT_UINT8
	case uint16:
		s.grb = C.GxB_ANY_ISGT_UINT16
	case uint32:
		s.grb = C.GxB_ANY_ISGT_UINT32
	case uint64:
		s.grb = C.GxB_ANY_ISGT_UINT64
	case float32:
		s.grb = C.GxB_ANY_ISGT_FP32
	case float64:
		s.grb = C.GxB_ANY_ISGT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinIslt semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Islt].
//
// MinIslt is a SuiteSparse:GraphBLAS extension.
func MinIslt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISLT_INT32
		} else {
			s.grb = C.GxB_MIN_ISLT_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_ISLT_INT8
	case int16:
		s.grb = C.GxB_MIN_ISLT_INT16
	case int32:
		s.grb = C.GxB_MIN_ISLT_INT32
	case int64:
		s.grb = C.GxB_MIN_ISLT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISLT_UINT32
		} else {
			s.grb = C.GxB_MIN_ISLT_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_ISLT_UINT8
	case uint16:
		s.grb = C.GxB_MIN_ISLT_UINT16
	case uint32:
		s.grb = C.GxB_MIN_ISLT_UINT32
	case uint64:
		s.grb = C.GxB_MIN_ISLT_UINT64
	case float32:
		s.grb = C.GxB_MIN_ISLT_FP32
	case float64:
		s.grb = C.GxB_MIN_ISLT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxIslt semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Islt].
//
// MaxIslt is a SuiteSparse:GraphBLAS extension.
func MaxIslt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISLT_INT32
		} else {
			s.grb = C.GxB_MAX_ISLT_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_ISLT_INT8
	case int16:
		s.grb = C.GxB_MAX_ISLT_INT16
	case int32:
		s.grb = C.GxB_MAX_ISLT_INT32
	case int64:
		s.grb = C.GxB_MAX_ISLT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISLT_UINT32
		} else {
			s.grb = C.GxB_MAX_ISLT_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_ISLT_UINT8
	case uint16:
		s.grb = C.GxB_MAX_ISLT_UINT16
	case uint32:
		s.grb = C.GxB_MAX_ISLT_UINT32
	case uint64:
		s.grb = C.GxB_MAX_ISLT_UINT64
	case float32:
		s.grb = C.GxB_MAX_ISLT_FP32
	case float64:
		s.grb = C.GxB_MAX_ISLT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusIslt semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Islt].
//
// PlusIslt is a SuiteSparse:GraphBLAS extension.
func PlusIslt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISLT_INT32
		} else {
			s.grb = C.GxB_PLUS_ISLT_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_ISLT_INT8
	case int16:
		s.grb = C.GxB_PLUS_ISLT_INT16
	case int32:
		s.grb = C.GxB_PLUS_ISLT_INT32
	case int64:
		s.grb = C.GxB_PLUS_ISLT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISLT_UINT32
		} else {
			s.grb = C.GxB_PLUS_ISLT_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_ISLT_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_ISLT_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_ISLT_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_ISLT_UINT64
	case float32:
		s.grb = C.GxB_PLUS_ISLT_FP32
	case float64:
		s.grb = C.GxB_PLUS_ISLT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesIslt semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Islt].
//
// TimesIslt is a SuiteSparse:GraphBLAS extension.
func TimesIslt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISLT_INT32
		} else {
			s.grb = C.GxB_TIMES_ISLT_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_ISLT_INT8
	case int16:
		s.grb = C.GxB_TIMES_ISLT_INT16
	case int32:
		s.grb = C.GxB_TIMES_ISLT_INT32
	case int64:
		s.grb = C.GxB_TIMES_ISLT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISLT_UINT32
		} else {
			s.grb = C.GxB_TIMES_ISLT_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_ISLT_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_ISLT_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_ISLT_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_ISLT_UINT64
	case float32:
		s.grb = C.GxB_TIMES_ISLT_FP32
	case float64:
		s.grb = C.GxB_TIMES_ISLT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyIslt semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Islt].
//
// AnyIslt is a SuiteSparse:GraphBLAS extension.
func AnyIslt[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISLT_INT32
		} else {
			s.grb = C.GxB_ANY_ISLT_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_ISLT_INT8
	case int16:
		s.grb = C.GxB_ANY_ISLT_INT16
	case int32:
		s.grb = C.GxB_ANY_ISLT_INT32
	case int64:
		s.grb = C.GxB_ANY_ISLT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISLT_UINT32
		} else {
			s.grb = C.GxB_ANY_ISLT_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_ISLT_UINT8
	case uint16:
		s.grb = C.GxB_ANY_ISLT_UINT16
	case uint32:
		s.grb = C.GxB_ANY_ISLT_UINT32
	case uint64:
		s.grb = C.GxB_ANY_ISLT_UINT64
	case float32:
		s.grb = C.GxB_ANY_ISLT_FP32
	case float64:
		s.grb = C.GxB_ANY_ISLT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinIsge semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Isge].
//
// MinIsge is a SuiteSparse:GraphBLAS extension.
func MinIsge[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISGE_INT32
		} else {
			s.grb = C.GxB_MIN_ISGE_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_ISGE_INT8
	case int16:
		s.grb = C.GxB_MIN_ISGE_INT16
	case int32:
		s.grb = C.GxB_MIN_ISGE_INT32
	case int64:
		s.grb = C.GxB_MIN_ISGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISGE_UINT32
		} else {
			s.grb = C.GxB_MIN_ISGE_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_ISGE_UINT8
	case uint16:
		s.grb = C.GxB_MIN_ISGE_UINT16
	case uint32:
		s.grb = C.GxB_MIN_ISGE_UINT32
	case uint64:
		s.grb = C.GxB_MIN_ISGE_UINT64
	case float32:
		s.grb = C.GxB_MIN_ISGE_FP32
	case float64:
		s.grb = C.GxB_MIN_ISGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxIsge semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Isge].
//
// MaxIsge is a SuiteSparse:GraphBLAS extension.
func MaxIsge[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISGE_INT32
		} else {
			s.grb = C.GxB_MAX_ISGE_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_ISGE_INT8
	case int16:
		s.grb = C.GxB_MAX_ISGE_INT16
	case int32:
		s.grb = C.GxB_MAX_ISGE_INT32
	case int64:
		s.grb = C.GxB_MAX_ISGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISGE_UINT32
		} else {
			s.grb = C.GxB_MAX_ISGE_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_ISGE_UINT8
	case uint16:
		s.grb = C.GxB_MAX_ISGE_UINT16
	case uint32:
		s.grb = C.GxB_MAX_ISGE_UINT32
	case uint64:
		s.grb = C.GxB_MAX_ISGE_UINT64
	case float32:
		s.grb = C.GxB_MAX_ISGE_FP32
	case float64:
		s.grb = C.GxB_MAX_ISGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusIsge semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Isge].
//
// PlusIsge is a SuiteSparse:GraphBLAS extension.
func PlusIsge[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISGE_INT32
		} else {
			s.grb = C.GxB_PLUS_ISGE_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_ISGE_INT8
	case int16:
		s.grb = C.GxB_PLUS_ISGE_INT16
	case int32:
		s.grb = C.GxB_PLUS_ISGE_INT32
	case int64:
		s.grb = C.GxB_PLUS_ISGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISGE_UINT32
		} else {
			s.grb = C.GxB_PLUS_ISGE_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_ISGE_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_ISGE_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_ISGE_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_ISGE_UINT64
	case float32:
		s.grb = C.GxB_PLUS_ISGE_FP32
	case float64:
		s.grb = C.GxB_PLUS_ISGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesIsge semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Isge].
//
// TimesIsge is a SuiteSparse:GraphBLAS extension.
func TimesIsge[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISGE_INT32
		} else {
			s.grb = C.GxB_TIMES_ISGE_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_ISGE_INT8
	case int16:
		s.grb = C.GxB_TIMES_ISGE_INT16
	case int32:
		s.grb = C.GxB_TIMES_ISGE_INT32
	case int64:
		s.grb = C.GxB_TIMES_ISGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISGE_UINT32
		} else {
			s.grb = C.GxB_TIMES_ISGE_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_ISGE_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_ISGE_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_ISGE_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_ISGE_UINT64
	case float32:
		s.grb = C.GxB_TIMES_ISGE_FP32
	case float64:
		s.grb = C.GxB_TIMES_ISGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyIsge semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Isge].
//
// AnyIsge is a SuiteSparse:GraphBLAS extension.
func AnyIsge[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISGE_INT32
		} else {
			s.grb = C.GxB_ANY_ISGE_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_ISGE_INT8
	case int16:
		s.grb = C.GxB_ANY_ISGE_INT16
	case int32:
		s.grb = C.GxB_ANY_ISGE_INT32
	case int64:
		s.grb = C.GxB_ANY_ISGE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISGE_UINT32
		} else {
			s.grb = C.GxB_ANY_ISGE_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_ISGE_UINT8
	case uint16:
		s.grb = C.GxB_ANY_ISGE_UINT16
	case uint32:
		s.grb = C.GxB_ANY_ISGE_UINT32
	case uint64:
		s.grb = C.GxB_ANY_ISGE_UINT64
	case float32:
		s.grb = C.GxB_ANY_ISGE_FP32
	case float64:
		s.grb = C.GxB_ANY_ISGE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinIsle semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Isle].
//
// MinIsle is a SuiteSparse:GraphBLAS extension.
func MinIsle[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISLE_INT32
		} else {
			s.grb = C.GxB_MIN_ISLE_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_ISLE_INT8
	case int16:
		s.grb = C.GxB_MIN_ISLE_INT16
	case int32:
		s.grb = C.GxB_MIN_ISLE_INT32
	case int64:
		s.grb = C.GxB_MIN_ISLE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_ISLE_UINT32
		} else {
			s.grb = C.GxB_MIN_ISLE_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_ISLE_UINT8
	case uint16:
		s.grb = C.GxB_MIN_ISLE_UINT16
	case uint32:
		s.grb = C.GxB_MIN_ISLE_UINT32
	case uint64:
		s.grb = C.GxB_MIN_ISLE_UINT64
	case float32:
		s.grb = C.GxB_MIN_ISLE_FP32
	case float64:
		s.grb = C.GxB_MIN_ISLE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxIsle semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Isle].
//
// MaxIsle is a SuiteSparse:GraphBLAS extension.
func MaxIsle[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISLE_INT32
		} else {
			s.grb = C.GxB_MAX_ISLE_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_ISLE_INT8
	case int16:
		s.grb = C.GxB_MAX_ISLE_INT16
	case int32:
		s.grb = C.GxB_MAX_ISLE_INT32
	case int64:
		s.grb = C.GxB_MAX_ISLE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_ISLE_UINT32
		} else {
			s.grb = C.GxB_MAX_ISLE_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_ISLE_UINT8
	case uint16:
		s.grb = C.GxB_MAX_ISLE_UINT16
	case uint32:
		s.grb = C.GxB_MAX_ISLE_UINT32
	case uint64:
		s.grb = C.GxB_MAX_ISLE_UINT64
	case float32:
		s.grb = C.GxB_MAX_ISLE_FP32
	case float64:
		s.grb = C.GxB_MAX_ISLE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusIsle semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Isle].
//
// PlusIsle is a SuiteSparse:GraphBLAS extension.
func PlusIsle[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISLE_INT32
		} else {
			s.grb = C.GxB_PLUS_ISLE_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_ISLE_INT8
	case int16:
		s.grb = C.GxB_PLUS_ISLE_INT16
	case int32:
		s.grb = C.GxB_PLUS_ISLE_INT32
	case int64:
		s.grb = C.GxB_PLUS_ISLE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_ISLE_UINT32
		} else {
			s.grb = C.GxB_PLUS_ISLE_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_ISLE_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_ISLE_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_ISLE_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_ISLE_UINT64
	case float32:
		s.grb = C.GxB_PLUS_ISLE_FP32
	case float64:
		s.grb = C.GxB_PLUS_ISLE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesIsle semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Isle].
//
// TimesIsle is a SuiteSparse:GraphBLAS extension.
func TimesIsle[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISLE_INT32
		} else {
			s.grb = C.GxB_TIMES_ISLE_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_ISLE_INT8
	case int16:
		s.grb = C.GxB_TIMES_ISLE_INT16
	case int32:
		s.grb = C.GxB_TIMES_ISLE_INT32
	case int64:
		s.grb = C.GxB_TIMES_ISLE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_ISLE_UINT32
		} else {
			s.grb = C.GxB_TIMES_ISLE_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_ISLE_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_ISLE_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_ISLE_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_ISLE_UINT64
	case float32:
		s.grb = C.GxB_TIMES_ISLE_FP32
	case float64:
		s.grb = C.GxB_TIMES_ISLE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyIsle semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Isle].
//
// AnyIsle is a SuiteSparse:GraphBLAS extension.
func AnyIsle[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISLE_INT32
		} else {
			s.grb = C.GxB_ANY_ISLE_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_ISLE_INT8
	case int16:
		s.grb = C.GxB_ANY_ISLE_INT16
	case int32:
		s.grb = C.GxB_ANY_ISLE_INT32
	case int64:
		s.grb = C.GxB_ANY_ISLE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_ISLE_UINT32
		} else {
			s.grb = C.GxB_ANY_ISLE_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_ISLE_UINT8
	case uint16:
		s.grb = C.GxB_ANY_ISLE_UINT16
	case uint32:
		s.grb = C.GxB_ANY_ISLE_UINT32
	case uint64:
		s.grb = C.GxB_ANY_ISLE_UINT64
	case float32:
		s.grb = C.GxB_ANY_ISLE_FP32
	case float64:
		s.grb = C.GxB_ANY_ISLE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinLor semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Lor].
//
// MinLor is a SuiteSparse:GraphBLAS extension.
func MinLor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_LOR_INT32
		} else {
			s.grb = C.GxB_MIN_LOR_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_LOR_INT8
	case int16:
		s.grb = C.GxB_MIN_LOR_INT16
	case int32:
		s.grb = C.GxB_MIN_LOR_INT32
	case int64:
		s.grb = C.GxB_MIN_LOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_LOR_UINT32
		} else {
			s.grb = C.GxB_MIN_LOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_LOR_UINT8
	case uint16:
		s.grb = C.GxB_MIN_LOR_UINT16
	case uint32:
		s.grb = C.GxB_MIN_LOR_UINT32
	case uint64:
		s.grb = C.GxB_MIN_LOR_UINT64
	case float32:
		s.grb = C.GxB_MIN_LOR_FP32
	case float64:
		s.grb = C.GxB_MIN_LOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxLor semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Lor].
//
// MaxLor is a SuiteSparse:GraphBLAS extension.
func MaxLor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_LOR_INT32
		} else {
			s.grb = C.GxB_MAX_LOR_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_LOR_INT8
	case int16:
		s.grb = C.GxB_MAX_LOR_INT16
	case int32:
		s.grb = C.GxB_MAX_LOR_INT32
	case int64:
		s.grb = C.GxB_MAX_LOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_LOR_UINT32
		} else {
			s.grb = C.GxB_MAX_LOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_LOR_UINT8
	case uint16:
		s.grb = C.GxB_MAX_LOR_UINT16
	case uint32:
		s.grb = C.GxB_MAX_LOR_UINT32
	case uint64:
		s.grb = C.GxB_MAX_LOR_UINT64
	case float32:
		s.grb = C.GxB_MAX_LOR_FP32
	case float64:
		s.grb = C.GxB_MAX_LOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusLor semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Lor].
//
// PlusLor is a SuiteSparse:GraphBLAS extension.
func PlusLor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_LOR_INT32
		} else {
			s.grb = C.GxB_PLUS_LOR_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_LOR_INT8
	case int16:
		s.grb = C.GxB_PLUS_LOR_INT16
	case int32:
		s.grb = C.GxB_PLUS_LOR_INT32
	case int64:
		s.grb = C.GxB_PLUS_LOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_LOR_UINT32
		} else {
			s.grb = C.GxB_PLUS_LOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_LOR_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_LOR_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_LOR_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_LOR_UINT64
	case float32:
		s.grb = C.GxB_PLUS_LOR_FP32
	case float64:
		s.grb = C.GxB_PLUS_LOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesLor semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Lor].
//
// TimesLor is a SuiteSparse:GraphBLAS extension.
func TimesLor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_LOR_INT32
		} else {
			s.grb = C.GxB_TIMES_LOR_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_LOR_INT8
	case int16:
		s.grb = C.GxB_TIMES_LOR_INT16
	case int32:
		s.grb = C.GxB_TIMES_LOR_INT32
	case int64:
		s.grb = C.GxB_TIMES_LOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_LOR_UINT32
		} else {
			s.grb = C.GxB_TIMES_LOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_LOR_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_LOR_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_LOR_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_LOR_UINT64
	case float32:
		s.grb = C.GxB_TIMES_LOR_FP32
	case float64:
		s.grb = C.GxB_TIMES_LOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyLor semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Lor].
//
// AnyLor is a SuiteSparse:GraphBLAS extension.
func AnyLor[D Predefined]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_LOR_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LOR_INT32
		} else {
			s.grb = C.GxB_ANY_LOR_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_LOR_INT8
	case int16:
		s.grb = C.GxB_ANY_LOR_INT16
	case int32:
		s.grb = C.GxB_ANY_LOR_INT32
	case int64:
		s.grb = C.GxB_ANY_LOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LOR_UINT32
		} else {
			s.grb = C.GxB_ANY_LOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_LOR_UINT8
	case uint16:
		s.grb = C.GxB_ANY_LOR_UINT16
	case uint32:
		s.grb = C.GxB_ANY_LOR_UINT32
	case uint64:
		s.grb = C.GxB_ANY_LOR_UINT64
	case float32:
		s.grb = C.GxB_ANY_LOR_FP32
	case float64:
		s.grb = C.GxB_ANY_LOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinLand semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Land].
//
// MinLand is a SuiteSparse:GraphBLAS extension.
func MinLand[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_LAND_INT32
		} else {
			s.grb = C.GxB_MIN_LAND_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_LAND_INT8
	case int16:
		s.grb = C.GxB_MIN_LAND_INT16
	case int32:
		s.grb = C.GxB_MIN_LAND_INT32
	case int64:
		s.grb = C.GxB_MIN_LAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_LAND_UINT32
		} else {
			s.grb = C.GxB_MIN_LAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_LAND_UINT8
	case uint16:
		s.grb = C.GxB_MIN_LAND_UINT16
	case uint32:
		s.grb = C.GxB_MIN_LAND_UINT32
	case uint64:
		s.grb = C.GxB_MIN_LAND_UINT64
	case float32:
		s.grb = C.GxB_MIN_LAND_FP32
	case float64:
		s.grb = C.GxB_MIN_LAND_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxLand semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Land].
//
// MaxLand is a SuiteSparse:GraphBLAS extension.
func MaxLand[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_LAND_INT32
		} else {
			s.grb = C.GxB_MAX_LAND_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_LAND_INT8
	case int16:
		s.grb = C.GxB_MAX_LAND_INT16
	case int32:
		s.grb = C.GxB_MAX_LAND_INT32
	case int64:
		s.grb = C.GxB_MAX_LAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_LAND_UINT32
		} else {
			s.grb = C.GxB_MAX_LAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_LAND_UINT8
	case uint16:
		s.grb = C.GxB_MAX_LAND_UINT16
	case uint32:
		s.grb = C.GxB_MAX_LAND_UINT32
	case uint64:
		s.grb = C.GxB_MAX_LAND_UINT64
	case float32:
		s.grb = C.GxB_MAX_LAND_FP32
	case float64:
		s.grb = C.GxB_MAX_LAND_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusLand semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Land].
//
// PlusLand is a SuiteSparse:GraphBLAS extension.
func PlusLand[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_LAND_INT32
		} else {
			s.grb = C.GxB_PLUS_LAND_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_LAND_INT8
	case int16:
		s.grb = C.GxB_PLUS_LAND_INT16
	case int32:
		s.grb = C.GxB_PLUS_LAND_INT32
	case int64:
		s.grb = C.GxB_PLUS_LAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_LAND_UINT32
		} else {
			s.grb = C.GxB_PLUS_LAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_LAND_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_LAND_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_LAND_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_LAND_UINT64
	case float32:
		s.grb = C.GxB_PLUS_LAND_FP32
	case float64:
		s.grb = C.GxB_PLUS_LAND_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesLand semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Land].
//
// TimesLand is a SuiteSparse:GraphBLAS extension.
func TimesLand[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_LAND_INT32
		} else {
			s.grb = C.GxB_TIMES_LAND_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_LAND_INT8
	case int16:
		s.grb = C.GxB_TIMES_LAND_INT16
	case int32:
		s.grb = C.GxB_TIMES_LAND_INT32
	case int64:
		s.grb = C.GxB_TIMES_LAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_LAND_UINT32
		} else {
			s.grb = C.GxB_TIMES_LAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_LAND_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_LAND_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_LAND_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_LAND_UINT64
	case float32:
		s.grb = C.GxB_TIMES_LAND_FP32
	case float64:
		s.grb = C.GxB_TIMES_LAND_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyLand semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Land].
//
// AnyLand is a SuiteSparse:GraphBLAS extension.
func AnyLand[D Predefined]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_LAND_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LAND_INT32
		} else {
			s.grb = C.GxB_ANY_LAND_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_LAND_INT8
	case int16:
		s.grb = C.GxB_ANY_LAND_INT16
	case int32:
		s.grb = C.GxB_ANY_LAND_INT32
	case int64:
		s.grb = C.GxB_ANY_LAND_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LAND_UINT32
		} else {
			s.grb = C.GxB_ANY_LAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_LAND_UINT8
	case uint16:
		s.grb = C.GxB_ANY_LAND_UINT16
	case uint32:
		s.grb = C.GxB_ANY_LAND_UINT32
	case uint64:
		s.grb = C.GxB_ANY_LAND_UINT64
	case float32:
		s.grb = C.GxB_ANY_LAND_FP32
	case float64:
		s.grb = C.GxB_ANY_LAND_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MinLxor semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Lxor].
//
// MinLxor is a SuiteSparse:GraphBLAS extension.
func MinLxor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_LXOR_INT32
		} else {
			s.grb = C.GxB_MIN_LXOR_INT64
		}
	case int8:
		s.grb = C.GxB_MIN_LXOR_INT8
	case int16:
		s.grb = C.GxB_MIN_LXOR_INT16
	case int32:
		s.grb = C.GxB_MIN_LXOR_INT32
	case int64:
		s.grb = C.GxB_MIN_LXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_LXOR_UINT32
		} else {
			s.grb = C.GxB_MIN_LXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_MIN_LXOR_UINT8
	case uint16:
		s.grb = C.GxB_MIN_LXOR_UINT16
	case uint32:
		s.grb = C.GxB_MIN_LXOR_UINT32
	case uint64:
		s.grb = C.GxB_MIN_LXOR_UINT64
	case float32:
		s.grb = C.GxB_MIN_LXOR_FP32
	case float64:
		s.grb = C.GxB_MIN_LXOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// MaxLxor semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Lxor].
//
// MaxLxor is a SuiteSparse:GraphBLAS extension.
func MaxLxor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_LXOR_INT32
		} else {
			s.grb = C.GxB_MAX_LXOR_INT64
		}
	case int8:
		s.grb = C.GxB_MAX_LXOR_INT8
	case int16:
		s.grb = C.GxB_MAX_LXOR_INT16
	case int32:
		s.grb = C.GxB_MAX_LXOR_INT32
	case int64:
		s.grb = C.GxB_MAX_LXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_LXOR_UINT32
		} else {
			s.grb = C.GxB_MAX_LXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_MAX_LXOR_UINT8
	case uint16:
		s.grb = C.GxB_MAX_LXOR_UINT16
	case uint32:
		s.grb = C.GxB_MAX_LXOR_UINT32
	case uint64:
		s.grb = C.GxB_MAX_LXOR_UINT64
	case float32:
		s.grb = C.GxB_MAX_LXOR_FP32
	case float64:
		s.grb = C.GxB_MAX_LXOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// PlusLxor semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Lxor].
//
// PlusLxor is a SuiteSparse:GraphBLAS extension.
func PlusLxor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_LXOR_INT32
		} else {
			s.grb = C.GxB_PLUS_LXOR_INT64
		}
	case int8:
		s.grb = C.GxB_PLUS_LXOR_INT8
	case int16:
		s.grb = C.GxB_PLUS_LXOR_INT16
	case int32:
		s.grb = C.GxB_PLUS_LXOR_INT32
	case int64:
		s.grb = C.GxB_PLUS_LXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_LXOR_UINT32
		} else {
			s.grb = C.GxB_PLUS_LXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_PLUS_LXOR_UINT8
	case uint16:
		s.grb = C.GxB_PLUS_LXOR_UINT16
	case uint32:
		s.grb = C.GxB_PLUS_LXOR_UINT32
	case uint64:
		s.grb = C.GxB_PLUS_LXOR_UINT64
	case float32:
		s.grb = C.GxB_PLUS_LXOR_FP32
	case float64:
		s.grb = C.GxB_PLUS_LXOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// TimesLxor semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Lxor].
//
// TimesLxor is a SuiteSparse:GraphBLAS extension.
func TimesLxor[D Number]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_LXOR_INT32
		} else {
			s.grb = C.GxB_TIMES_LXOR_INT64
		}
	case int8:
		s.grb = C.GxB_TIMES_LXOR_INT8
	case int16:
		s.grb = C.GxB_TIMES_LXOR_INT16
	case int32:
		s.grb = C.GxB_TIMES_LXOR_INT32
	case int64:
		s.grb = C.GxB_TIMES_LXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_LXOR_UINT32
		} else {
			s.grb = C.GxB_TIMES_LXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_TIMES_LXOR_UINT8
	case uint16:
		s.grb = C.GxB_TIMES_LXOR_UINT16
	case uint32:
		s.grb = C.GxB_TIMES_LXOR_UINT32
	case uint64:
		s.grb = C.GxB_TIMES_LXOR_UINT64
	case float32:
		s.grb = C.GxB_TIMES_LXOR_FP32
	case float64:
		s.grb = C.GxB_TIMES_LXOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyLxor semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Lxor].
//
// AnyLxor is a SuiteSparse:GraphBLAS extension.
func AnyLxor[D Predefined]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_LXOR_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LXOR_INT32
		} else {
			s.grb = C.GxB_ANY_LXOR_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_LXOR_INT8
	case int16:
		s.grb = C.GxB_ANY_LXOR_INT16
	case int32:
		s.grb = C.GxB_ANY_LXOR_INT32
	case int64:
		s.grb = C.GxB_ANY_LXOR_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LXOR_UINT32
		} else {
			s.grb = C.GxB_ANY_LXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_LXOR_UINT8
	case uint16:
		s.grb = C.GxB_ANY_LXOR_UINT16
	case uint32:
		s.grb = C.GxB_ANY_LXOR_UINT32
	case uint64:
		s.grb = C.GxB_ANY_LXOR_UINT64
	case float32:
		s.grb = C.GxB_ANY_LXOR_FP32
	case float64:
		s.grb = C.GxB_ANY_LXOR_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LorEq semiring with additive [Monoid] [LorMonoid] and [BinaryOp] [Eq].
//
// LorEq is a SuiteSparse:GraphBLAS extension.
func LorEq[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_EQ_INT32
		} else {
			s.grb = C.GxB_LOR_EQ_INT64
		}
	case int8:
		s.grb = C.GxB_LOR_EQ_INT8
	case int16:
		s.grb = C.GxB_LOR_EQ_INT16
	case int32:
		s.grb = C.GxB_LOR_EQ_INT32
	case int64:
		s.grb = C.GxB_LOR_EQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_EQ_UINT32
		} else {
			s.grb = C.GxB_LOR_EQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_LOR_EQ_UINT8
	case uint16:
		s.grb = C.GxB_LOR_EQ_UINT16
	case uint32:
		s.grb = C.GxB_LOR_EQ_UINT32
	case uint64:
		s.grb = C.GxB_LOR_EQ_UINT64
	case float32:
		s.grb = C.GxB_LOR_EQ_FP32
	case float64:
		s.grb = C.GxB_LOR_EQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LandEq semiring with additive [Monoid] [LandMonoid] and [BinaryOp] [Eq].
//
// LandEq is a SuiteSparse:GraphBLAS extension.
func LandEq[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_EQ_INT32
		} else {
			s.grb = C.GxB_LAND_EQ_INT64
		}
	case int8:
		s.grb = C.GxB_LAND_EQ_INT8
	case int16:
		s.grb = C.GxB_LAND_EQ_INT16
	case int32:
		s.grb = C.GxB_LAND_EQ_INT32
	case int64:
		s.grb = C.GxB_LAND_EQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_EQ_UINT32
		} else {
			s.grb = C.GxB_LAND_EQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_LAND_EQ_UINT8
	case uint16:
		s.grb = C.GxB_LAND_EQ_UINT16
	case uint32:
		s.grb = C.GxB_LAND_EQ_UINT32
	case uint64:
		s.grb = C.GxB_LAND_EQ_UINT64
	case float32:
		s.grb = C.GxB_LAND_EQ_FP32
	case float64:
		s.grb = C.GxB_LAND_EQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LxorEq semiring with additive [Monoid] [LxorMonoid] and [BinaryOp] [Eq].
//
// LxorEq is a SuiteSparse:GraphBLAS extension.
func LxorEq[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_EQ_INT32
		} else {
			s.grb = C.GxB_LXOR_EQ_INT64
		}
	case int8:
		s.grb = C.GxB_LXOR_EQ_INT8
	case int16:
		s.grb = C.GxB_LXOR_EQ_INT16
	case int32:
		s.grb = C.GxB_LXOR_EQ_INT32
	case int64:
		s.grb = C.GxB_LXOR_EQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_EQ_UINT32
		} else {
			s.grb = C.GxB_LXOR_EQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_LXOR_EQ_UINT8
	case uint16:
		s.grb = C.GxB_LXOR_EQ_UINT16
	case uint32:
		s.grb = C.GxB_LXOR_EQ_UINT32
	case uint64:
		s.grb = C.GxB_LXOR_EQ_UINT64
	case float32:
		s.grb = C.GxB_LXOR_EQ_FP32
	case float64:
		s.grb = C.GxB_LXOR_EQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// EqEq semiring with additive [Monoid] [Eq] and [BinaryOp] [Eq].
//
// EqEq is a SuiteSparse:GraphBLAS extension.
func EqEq[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_EQ_INT32
		} else {
			s.grb = C.GxB_EQ_EQ_INT64
		}
	case int8:
		s.grb = C.GxB_EQ_EQ_INT8
	case int16:
		s.grb = C.GxB_EQ_EQ_INT16
	case int32:
		s.grb = C.GxB_EQ_EQ_INT32
	case int64:
		s.grb = C.GxB_EQ_EQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_EQ_UINT32
		} else {
			s.grb = C.GxB_EQ_EQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_EQ_EQ_UINT8
	case uint16:
		s.grb = C.GxB_EQ_EQ_UINT16
	case uint32:
		s.grb = C.GxB_EQ_EQ_UINT32
	case uint64:
		s.grb = C.GxB_EQ_EQ_UINT64
	case float32:
		s.grb = C.GxB_EQ_EQ_FP32
	case float64:
		s.grb = C.GxB_EQ_EQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyEq semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Eq].
//
// AnyEq is a SuiteSparse:GraphBLAS extension.
func AnyEq[D Predefined]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_EQ_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_EQ_INT32
		} else {
			s.grb = C.GxB_ANY_EQ_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_EQ_INT8
	case int16:
		s.grb = C.GxB_ANY_EQ_INT16
	case int32:
		s.grb = C.GxB_ANY_EQ_INT32
	case int64:
		s.grb = C.GxB_ANY_EQ_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_EQ_UINT32
		} else {
			s.grb = C.GxB_ANY_EQ_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_EQ_UINT8
	case uint16:
		s.grb = C.GxB_ANY_EQ_UINT16
	case uint32:
		s.grb = C.GxB_ANY_EQ_UINT32
	case uint64:
		s.grb = C.GxB_ANY_EQ_UINT64
	case float32:
		s.grb = C.GxB_ANY_EQ_FP32
	case float64:
		s.grb = C.GxB_ANY_EQ_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LorNe semiring with additive [Monoid] [LorMonoid] and [BinaryOp] [Ne].
//
// LorNe is a SuiteSparse:GraphBLAS extension.
func LorNe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_NE_INT32
		} else {
			s.grb = C.GxB_LOR_NE_INT64
		}
	case int8:
		s.grb = C.GxB_LOR_NE_INT8
	case int16:
		s.grb = C.GxB_LOR_NE_INT16
	case int32:
		s.grb = C.GxB_LOR_NE_INT32
	case int64:
		s.grb = C.GxB_LOR_NE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_NE_UINT32
		} else {
			s.grb = C.GxB_LOR_NE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LOR_NE_UINT8
	case uint16:
		s.grb = C.GxB_LOR_NE_UINT16
	case uint32:
		s.grb = C.GxB_LOR_NE_UINT32
	case uint64:
		s.grb = C.GxB_LOR_NE_UINT64
	case float32:
		s.grb = C.GxB_LOR_NE_FP32
	case float64:
		s.grb = C.GxB_LOR_NE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LandNe semiring with additive [Monoid] [LandMonoid] and [BinaryOp] [Ne].
//
// LandNe is a SuiteSparse:GraphBLAS extension.
func LandNe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_NE_INT32
		} else {
			s.grb = C.GxB_LAND_NE_INT64
		}
	case int8:
		s.grb = C.GxB_LAND_NE_INT8
	case int16:
		s.grb = C.GxB_LAND_NE_INT16
	case int32:
		s.grb = C.GxB_LAND_NE_INT32
	case int64:
		s.grb = C.GxB_LAND_NE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_NE_UINT32
		} else {
			s.grb = C.GxB_LAND_NE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LAND_NE_UINT8
	case uint16:
		s.grb = C.GxB_LAND_NE_UINT16
	case uint32:
		s.grb = C.GxB_LAND_NE_UINT32
	case uint64:
		s.grb = C.GxB_LAND_NE_UINT64
	case float32:
		s.grb = C.GxB_LAND_NE_FP32
	case float64:
		s.grb = C.GxB_LAND_NE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LxorNe semiring with additive [Monoid] [LxorMonoid] and [BinaryOp] [Ne].
//
// LxorNe is a SuiteSparse:GraphBLAS extension.
func LxorNe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_NE_INT32
		} else {
			s.grb = C.GxB_LXOR_NE_INT64
		}
	case int8:
		s.grb = C.GxB_LXOR_NE_INT8
	case int16:
		s.grb = C.GxB_LXOR_NE_INT16
	case int32:
		s.grb = C.GxB_LXOR_NE_INT32
	case int64:
		s.grb = C.GxB_LXOR_NE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_NE_UINT32
		} else {
			s.grb = C.GxB_LXOR_NE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LXOR_NE_UINT8
	case uint16:
		s.grb = C.GxB_LXOR_NE_UINT16
	case uint32:
		s.grb = C.GxB_LXOR_NE_UINT32
	case uint64:
		s.grb = C.GxB_LXOR_NE_UINT64
	case float32:
		s.grb = C.GxB_LXOR_NE_FP32
	case float64:
		s.grb = C.GxB_LXOR_NE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// EqNe semiring with additive [Monoid] [Eq] and [BinaryOp] [Ne].
//
// EqNe is a SuiteSparse:GraphBLAS extension.
func EqNe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_NE_INT32
		} else {
			s.grb = C.GxB_EQ_NE_INT64
		}
	case int8:
		s.grb = C.GxB_EQ_NE_INT8
	case int16:
		s.grb = C.GxB_EQ_NE_INT16
	case int32:
		s.grb = C.GxB_EQ_NE_INT32
	case int64:
		s.grb = C.GxB_EQ_NE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_NE_UINT32
		} else {
			s.grb = C.GxB_EQ_NE_UINT64
		}
	case uint8:
		s.grb = C.GxB_EQ_NE_UINT8
	case uint16:
		s.grb = C.GxB_EQ_NE_UINT16
	case uint32:
		s.grb = C.GxB_EQ_NE_UINT32
	case uint64:
		s.grb = C.GxB_EQ_NE_UINT64
	case float32:
		s.grb = C.GxB_EQ_NE_FP32
	case float64:
		s.grb = C.GxB_EQ_NE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyNe semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Ne].
//
// AnyNe is a SuiteSparse:GraphBLAS extension.
func AnyNe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_NE_INT32
		} else {
			s.grb = C.GxB_ANY_NE_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_NE_INT8
	case int16:
		s.grb = C.GxB_ANY_NE_INT16
	case int32:
		s.grb = C.GxB_ANY_NE_INT32
	case int64:
		s.grb = C.GxB_ANY_NE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_NE_UINT32
		} else {
			s.grb = C.GxB_ANY_NE_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_NE_UINT8
	case uint16:
		s.grb = C.GxB_ANY_NE_UINT16
	case uint32:
		s.grb = C.GxB_ANY_NE_UINT32
	case uint64:
		s.grb = C.GxB_ANY_NE_UINT64
	case float32:
		s.grb = C.GxB_ANY_NE_FP32
	case float64:
		s.grb = C.GxB_ANY_NE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LorGt semiring with additive [Monoid] [LorMonoid] and [BinaryOp] [Gt].
//
// LorGt is a SuiteSparse:GraphBLAS extension.
func LorGt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_GT_INT32
		} else {
			s.grb = C.GxB_LOR_GT_INT64
		}
	case int8:
		s.grb = C.GxB_LOR_GT_INT8
	case int16:
		s.grb = C.GxB_LOR_GT_INT16
	case int32:
		s.grb = C.GxB_LOR_GT_INT32
	case int64:
		s.grb = C.GxB_LOR_GT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_GT_UINT32
		} else {
			s.grb = C.GxB_LOR_GT_UINT64
		}
	case uint8:
		s.grb = C.GxB_LOR_GT_UINT8
	case uint16:
		s.grb = C.GxB_LOR_GT_UINT16
	case uint32:
		s.grb = C.GxB_LOR_GT_UINT32
	case uint64:
		s.grb = C.GxB_LOR_GT_UINT64
	case float32:
		s.grb = C.GxB_LOR_GT_FP32
	case float64:
		s.grb = C.GxB_LOR_GT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LandGt semiring with additive [Monoid] [LandMonoid] and [BinaryOp] [Gt].
//
// LandGt is a SuiteSparse:GraphBLAS extension.
func LandGt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_GT_INT32
		} else {
			s.grb = C.GxB_LAND_GT_INT64
		}
	case int8:
		s.grb = C.GxB_LAND_GT_INT8
	case int16:
		s.grb = C.GxB_LAND_GT_INT16
	case int32:
		s.grb = C.GxB_LAND_GT_INT32
	case int64:
		s.grb = C.GxB_LAND_GT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_GT_UINT32
		} else {
			s.grb = C.GxB_LAND_GT_UINT64
		}
	case uint8:
		s.grb = C.GxB_LAND_GT_UINT8
	case uint16:
		s.grb = C.GxB_LAND_GT_UINT16
	case uint32:
		s.grb = C.GxB_LAND_GT_UINT32
	case uint64:
		s.grb = C.GxB_LAND_GT_UINT64
	case float32:
		s.grb = C.GxB_LAND_GT_FP32
	case float64:
		s.grb = C.GxB_LAND_GT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LxorGt semiring with additive [Monoid] [LxorMonoid] and [BinaryOp] [Gt].
//
// LxorGt is a SuiteSparse:GraphBLAS extension.
func LxorGt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_GT_INT32
		} else {
			s.grb = C.GxB_LXOR_GT_INT64
		}
	case int8:
		s.grb = C.GxB_LXOR_GT_INT8
	case int16:
		s.grb = C.GxB_LXOR_GT_INT16
	case int32:
		s.grb = C.GxB_LXOR_GT_INT32
	case int64:
		s.grb = C.GxB_LXOR_GT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_GT_UINT32
		} else {
			s.grb = C.GxB_LXOR_GT_UINT64
		}
	case uint8:
		s.grb = C.GxB_LXOR_GT_UINT8
	case uint16:
		s.grb = C.GxB_LXOR_GT_UINT16
	case uint32:
		s.grb = C.GxB_LXOR_GT_UINT32
	case uint64:
		s.grb = C.GxB_LXOR_GT_UINT64
	case float32:
		s.grb = C.GxB_LXOR_GT_FP32
	case float64:
		s.grb = C.GxB_LXOR_GT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// EqGt semiring with additive [Monoid] [Eq] and [BinaryOp] [Gt].
//
// EqGt is a SuiteSparse:GraphBLAS extension.
func EqGt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_GT_INT32
		} else {
			s.grb = C.GxB_EQ_GT_INT64
		}
	case int8:
		s.grb = C.GxB_EQ_GT_INT8
	case int16:
		s.grb = C.GxB_EQ_GT_INT16
	case int32:
		s.grb = C.GxB_EQ_GT_INT32
	case int64:
		s.grb = C.GxB_EQ_GT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_GT_UINT32
		} else {
			s.grb = C.GxB_EQ_GT_UINT64
		}
	case uint8:
		s.grb = C.GxB_EQ_GT_UINT8
	case uint16:
		s.grb = C.GxB_EQ_GT_UINT16
	case uint32:
		s.grb = C.GxB_EQ_GT_UINT32
	case uint64:
		s.grb = C.GxB_EQ_GT_UINT64
	case float32:
		s.grb = C.GxB_EQ_GT_FP32
	case float64:
		s.grb = C.GxB_EQ_GT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyGt semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Gt].
//
// AnyGt is a SuiteSparse:GraphBLAS extension.
func AnyGt[D Predefined]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_GT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_GT_INT32
		} else {
			s.grb = C.GxB_ANY_GT_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_GT_INT8
	case int16:
		s.grb = C.GxB_ANY_GT_INT16
	case int32:
		s.grb = C.GxB_ANY_GT_INT32
	case int64:
		s.grb = C.GxB_ANY_GT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_GT_UINT32
		} else {
			s.grb = C.GxB_ANY_GT_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_GT_UINT8
	case uint16:
		s.grb = C.GxB_ANY_GT_UINT16
	case uint32:
		s.grb = C.GxB_ANY_GT_UINT32
	case uint64:
		s.grb = C.GxB_ANY_GT_UINT64
	case float32:
		s.grb = C.GxB_ANY_GT_FP32
	case float64:
		s.grb = C.GxB_ANY_GT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LorLt semiring with additive [Monoid] [LorMonoid] and [BinaryOp] [Lt].
//
// LorLt is a SuiteSparse:GraphBLAS extension.
func LorLt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_LT_INT32
		} else {
			s.grb = C.GxB_LOR_LT_INT64
		}
	case int8:
		s.grb = C.GxB_LOR_LT_INT8
	case int16:
		s.grb = C.GxB_LOR_LT_INT16
	case int32:
		s.grb = C.GxB_LOR_LT_INT32
	case int64:
		s.grb = C.GxB_LOR_LT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_LT_UINT32
		} else {
			s.grb = C.GxB_LOR_LT_UINT64
		}
	case uint8:
		s.grb = C.GxB_LOR_LT_UINT8
	case uint16:
		s.grb = C.GxB_LOR_LT_UINT16
	case uint32:
		s.grb = C.GxB_LOR_LT_UINT32
	case uint64:
		s.grb = C.GxB_LOR_LT_UINT64
	case float32:
		s.grb = C.GxB_LOR_LT_FP32
	case float64:
		s.grb = C.GxB_LOR_LT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LandLt semiring with additive [Monoid] [LandMonoid] and [BinaryOp] [Lt].
//
// LandLt is a SuiteSparse:GraphBLAS extension.
func LandLt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_LT_INT32
		} else {
			s.grb = C.GxB_LAND_LT_INT64
		}
	case int8:
		s.grb = C.GxB_LAND_LT_INT8
	case int16:
		s.grb = C.GxB_LAND_LT_INT16
	case int32:
		s.grb = C.GxB_LAND_LT_INT32
	case int64:
		s.grb = C.GxB_LAND_LT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_LT_UINT32
		} else {
			s.grb = C.GxB_LAND_LT_UINT64
		}
	case uint8:
		s.grb = C.GxB_LAND_LT_UINT8
	case uint16:
		s.grb = C.GxB_LAND_LT_UINT16
	case uint32:
		s.grb = C.GxB_LAND_LT_UINT32
	case uint64:
		s.grb = C.GxB_LAND_LT_UINT64
	case float32:
		s.grb = C.GxB_LAND_LT_FP32
	case float64:
		s.grb = C.GxB_LAND_LT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LxorLt semiring with additive [Monoid] [LxorMonoid] and [BinaryOp] [Lt].
//
// LxorLt is a SuiteSparse:GraphBLAS extension.
func LxorLt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_LT_INT32
		} else {
			s.grb = C.GxB_LXOR_LT_INT64
		}
	case int8:
		s.grb = C.GxB_LXOR_LT_INT8
	case int16:
		s.grb = C.GxB_LXOR_LT_INT16
	case int32:
		s.grb = C.GxB_LXOR_LT_INT32
	case int64:
		s.grb = C.GxB_LXOR_LT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_LT_UINT32
		} else {
			s.grb = C.GxB_LXOR_LT_UINT64
		}
	case uint8:
		s.grb = C.GxB_LXOR_LT_UINT8
	case uint16:
		s.grb = C.GxB_LXOR_LT_UINT16
	case uint32:
		s.grb = C.GxB_LXOR_LT_UINT32
	case uint64:
		s.grb = C.GxB_LXOR_LT_UINT64
	case float32:
		s.grb = C.GxB_LXOR_LT_FP32
	case float64:
		s.grb = C.GxB_LXOR_LT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// EqLt semiring with additive [Monoid] [Eq] and [BinaryOp] [Lt].
//
// EqLt is a SuiteSparse:GraphBLAS extension.
func EqLt[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_LT_INT32
		} else {
			s.grb = C.GxB_EQ_LT_INT64
		}
	case int8:
		s.grb = C.GxB_EQ_LT_INT8
	case int16:
		s.grb = C.GxB_EQ_LT_INT16
	case int32:
		s.grb = C.GxB_EQ_LT_INT32
	case int64:
		s.grb = C.GxB_EQ_LT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_LT_UINT32
		} else {
			s.grb = C.GxB_EQ_LT_UINT64
		}
	case uint8:
		s.grb = C.GxB_EQ_LT_UINT8
	case uint16:
		s.grb = C.GxB_EQ_LT_UINT16
	case uint32:
		s.grb = C.GxB_EQ_LT_UINT32
	case uint64:
		s.grb = C.GxB_EQ_LT_UINT64
	case float32:
		s.grb = C.GxB_EQ_LT_FP32
	case float64:
		s.grb = C.GxB_EQ_LT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyLt semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Lt].
//
// AnyLt is a SuiteSparse:GraphBLAS extension.
func AnyLt[D Predefined]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_LT_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LT_INT32
		} else {
			s.grb = C.GxB_ANY_LT_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_LT_INT8
	case int16:
		s.grb = C.GxB_ANY_LT_INT16
	case int32:
		s.grb = C.GxB_ANY_LT_INT32
	case int64:
		s.grb = C.GxB_ANY_LT_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LT_UINT32
		} else {
			s.grb = C.GxB_ANY_LT_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_LT_UINT8
	case uint16:
		s.grb = C.GxB_ANY_LT_UINT16
	case uint32:
		s.grb = C.GxB_ANY_LT_UINT32
	case uint64:
		s.grb = C.GxB_ANY_LT_UINT64
	case float32:
		s.grb = C.GxB_ANY_LT_FP32
	case float64:
		s.grb = C.GxB_ANY_LT_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LorGe semiring with additive [Monoid] [LorMonoid] and [BinaryOp] [Ge].
//
// LorGe is a SuiteSparse:GraphBLAS extension.
func LorGe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_GE_INT32
		} else {
			s.grb = C.GxB_LOR_GE_INT64
		}
	case int8:
		s.grb = C.GxB_LOR_GE_INT8
	case int16:
		s.grb = C.GxB_LOR_GE_INT16
	case int32:
		s.grb = C.GxB_LOR_GE_INT32
	case int64:
		s.grb = C.GxB_LOR_GE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_GE_UINT32
		} else {
			s.grb = C.GxB_LOR_GE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LOR_GE_UINT8
	case uint16:
		s.grb = C.GxB_LOR_GE_UINT16
	case uint32:
		s.grb = C.GxB_LOR_GE_UINT32
	case uint64:
		s.grb = C.GxB_LOR_GE_UINT64
	case float32:
		s.grb = C.GxB_LOR_GE_FP32
	case float64:
		s.grb = C.GxB_LOR_GE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LandGe semiring with additive [Monoid] [LandMonoid] and [BinaryOp] [Ge].
//
// LandGe is a SuiteSparse:GraphBLAS extension.
func LandGe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_GE_INT32
		} else {
			s.grb = C.GxB_LAND_GE_INT64
		}
	case int8:
		s.grb = C.GxB_LAND_GE_INT8
	case int16:
		s.grb = C.GxB_LAND_GE_INT16
	case int32:
		s.grb = C.GxB_LAND_GE_INT32
	case int64:
		s.grb = C.GxB_LAND_GE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_GE_UINT32
		} else {
			s.grb = C.GxB_LAND_GE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LAND_GE_UINT8
	case uint16:
		s.grb = C.GxB_LAND_GE_UINT16
	case uint32:
		s.grb = C.GxB_LAND_GE_UINT32
	case uint64:
		s.grb = C.GxB_LAND_GE_UINT64
	case float32:
		s.grb = C.GxB_LAND_GE_FP32
	case float64:
		s.grb = C.GxB_LAND_GE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LxorGe semiring with additive [Monoid] [LxorMonoid] and [BinaryOp] [Ge].
//
// LxorGe is a SuiteSparse:GraphBLAS extension.
func LxorGe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_GE_INT32
		} else {
			s.grb = C.GxB_LXOR_GE_INT64
		}
	case int8:
		s.grb = C.GxB_LXOR_GE_INT8
	case int16:
		s.grb = C.GxB_LXOR_GE_INT16
	case int32:
		s.grb = C.GxB_LXOR_GE_INT32
	case int64:
		s.grb = C.GxB_LXOR_GE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_GE_UINT32
		} else {
			s.grb = C.GxB_LXOR_GE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LXOR_GE_UINT8
	case uint16:
		s.grb = C.GxB_LXOR_GE_UINT16
	case uint32:
		s.grb = C.GxB_LXOR_GE_UINT32
	case uint64:
		s.grb = C.GxB_LXOR_GE_UINT64
	case float32:
		s.grb = C.GxB_LXOR_GE_FP32
	case float64:
		s.grb = C.GxB_LXOR_GE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// EqGe semiring with additive [Monoid] [EqMonoid] and [BinaryOp] [Ge].
//
// EqGe is a SuiteSparse:GraphBLAS extension.
func EqGe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_GE_INT32
		} else {
			s.grb = C.GxB_EQ_GE_INT64
		}
	case int8:
		s.grb = C.GxB_EQ_GE_INT8
	case int16:
		s.grb = C.GxB_EQ_GE_INT16
	case int32:
		s.grb = C.GxB_EQ_GE_INT32
	case int64:
		s.grb = C.GxB_EQ_GE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_GE_UINT32
		} else {
			s.grb = C.GxB_EQ_GE_UINT64
		}
	case uint8:
		s.grb = C.GxB_EQ_GE_UINT8
	case uint16:
		s.grb = C.GxB_EQ_GE_UINT16
	case uint32:
		s.grb = C.GxB_EQ_GE_UINT32
	case uint64:
		s.grb = C.GxB_EQ_GE_UINT64
	case float32:
		s.grb = C.GxB_EQ_GE_FP32
	case float64:
		s.grb = C.GxB_EQ_GE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyGe semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Ge].
//
// AnyGe is a SuiteSparse:GraphBLAS extension.
func AnyGe[D Predefined]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_GE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_GE_INT32
		} else {
			s.grb = C.GxB_ANY_GE_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_GE_INT8
	case int16:
		s.grb = C.GxB_ANY_GE_INT16
	case int32:
		s.grb = C.GxB_ANY_GE_INT32
	case int64:
		s.grb = C.GxB_ANY_GE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_GE_UINT32
		} else {
			s.grb = C.GxB_ANY_GE_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_GE_UINT8
	case uint16:
		s.grb = C.GxB_ANY_GE_UINT16
	case uint32:
		s.grb = C.GxB_ANY_GE_UINT32
	case uint64:
		s.grb = C.GxB_ANY_GE_UINT64
	case float32:
		s.grb = C.GxB_ANY_GE_FP32
	case float64:
		s.grb = C.GxB_ANY_GE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LorLe semiring with additive [Monoid] [LorMonoid] and [BinaryOp] [Le].
//
// LorLe is a SuiteSparse:GraphBLAS extension.
func LorLe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_LE_INT32
		} else {
			s.grb = C.GxB_LOR_LE_INT64
		}
	case int8:
		s.grb = C.GxB_LOR_LE_INT8
	case int16:
		s.grb = C.GxB_LOR_LE_INT16
	case int32:
		s.grb = C.GxB_LOR_LE_INT32
	case int64:
		s.grb = C.GxB_LOR_LE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LOR_LE_UINT32
		} else {
			s.grb = C.GxB_LOR_LE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LOR_LE_UINT8
	case uint16:
		s.grb = C.GxB_LOR_LE_UINT16
	case uint32:
		s.grb = C.GxB_LOR_LE_UINT32
	case uint64:
		s.grb = C.GxB_LOR_LE_UINT64
	case float32:
		s.grb = C.GxB_LOR_LE_FP32
	case float64:
		s.grb = C.GxB_LOR_LE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LandLe semiring with additive [Monoid] [LandMonoid] and [BinaryOp] [Le].
//
// LandLe is a SuiteSparse:GraphBLAS extension.
func LandLe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_LE_INT32
		} else {
			s.grb = C.GxB_LAND_LE_INT64
		}
	case int8:
		s.grb = C.GxB_LAND_LE_INT8
	case int16:
		s.grb = C.GxB_LAND_LE_INT16
	case int32:
		s.grb = C.GxB_LAND_LE_INT32
	case int64:
		s.grb = C.GxB_LAND_LE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LAND_LE_UINT32
		} else {
			s.grb = C.GxB_LAND_LE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LAND_LE_UINT8
	case uint16:
		s.grb = C.GxB_LAND_LE_UINT16
	case uint32:
		s.grb = C.GxB_LAND_LE_UINT32
	case uint64:
		s.grb = C.GxB_LAND_LE_UINT64
	case float32:
		s.grb = C.GxB_LAND_LE_FP32
	case float64:
		s.grb = C.GxB_LAND_LE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// LxorLe semiring with additive [Monoid] [LxorMonoid] and [BinaryOp] [Le].
//
// LxorLe is a SuiteSparse:GraphBLAS extension.
func LxorLe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_LE_INT32
		} else {
			s.grb = C.GxB_LXOR_LE_INT64
		}
	case int8:
		s.grb = C.GxB_LXOR_LE_INT8
	case int16:
		s.grb = C.GxB_LXOR_LE_INT16
	case int32:
		s.grb = C.GxB_LXOR_LE_INT32
	case int64:
		s.grb = C.GxB_LXOR_LE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_LXOR_LE_UINT32
		} else {
			s.grb = C.GxB_LXOR_LE_UINT64
		}
	case uint8:
		s.grb = C.GxB_LXOR_LE_UINT8
	case uint16:
		s.grb = C.GxB_LXOR_LE_UINT16
	case uint32:
		s.grb = C.GxB_LXOR_LE_UINT32
	case uint64:
		s.grb = C.GxB_LXOR_LE_UINT64
	case float32:
		s.grb = C.GxB_LXOR_LE_FP32
	case float64:
		s.grb = C.GxB_LXOR_LE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// EqLe semiring with additive [Monoid] [Eq] and [BinaryOp] [Le].
//
// EqLe is a SuiteSparse:GraphBLAS extension.
func EqLe[D Number]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_LE_INT32
		} else {
			s.grb = C.GxB_EQ_LE_INT64
		}
	case int8:
		s.grb = C.GxB_EQ_LE_INT8
	case int16:
		s.grb = C.GxB_EQ_LE_INT16
	case int32:
		s.grb = C.GxB_EQ_LE_INT32
	case int64:
		s.grb = C.GxB_EQ_LE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_EQ_LE_UINT32
		} else {
			s.grb = C.GxB_EQ_LE_UINT64
		}
	case uint8:
		s.grb = C.GxB_EQ_LE_UINT8
	case uint16:
		s.grb = C.GxB_EQ_LE_UINT16
	case uint32:
		s.grb = C.GxB_EQ_LE_UINT32
	case uint64:
		s.grb = C.GxB_EQ_LE_UINT64
	case float32:
		s.grb = C.GxB_EQ_LE_FP32
	case float64:
		s.grb = C.GxB_EQ_LE_FP64
	default:
		panic("unreachable code")
	}
	return
}

// AnyLe semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Le].
//
// AnyLe is a SuiteSparse:GraphBLAS extension.
func AnyLe[D Predefined]() (s Semiring[bool, D, D]) {
	var d D
	switch any(d).(type) {
	case bool:
		s.grb = C.GxB_ANY_LE_BOOL
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LE_INT32
		} else {
			s.grb = C.GxB_ANY_LE_INT64
		}
	case int8:
		s.grb = C.GxB_ANY_LE_INT8
	case int16:
		s.grb = C.GxB_ANY_LE_INT16
	case int32:
		s.grb = C.GxB_ANY_LE_INT32
	case int64:
		s.grb = C.GxB_ANY_LE_INT64
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_LE_UINT32
		} else {
			s.grb = C.GxB_ANY_LE_UINT64
		}
	case uint8:
		s.grb = C.GxB_ANY_LE_UINT8
	case uint16:
		s.grb = C.GxB_ANY_LE_UINT16
	case uint32:
		s.grb = C.GxB_ANY_LE_UINT32
	case uint64:
		s.grb = C.GxB_ANY_LE_UINT64
	case float32:
		s.grb = C.GxB_ANY_LE_FP32
	case float64:
		s.grb = C.GxB_ANY_LE_FP64
	default:
		panic("unreachable code")
	}
	return
}

var (
	// LorFirstBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [First].
	//
	// LorFirstBool is a SuiteSparse:GraphBLAS extension.
	LorFirstBool = Semiring[bool, bool, bool]{C.GxB_LOR_FIRST_BOOL}

	// LandFirstBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [First].
	//
	// LandFirstBool is a SuiteSparse:GraphBLAS extension.
	LandFirstBool = Semiring[bool, bool, bool]{C.GxB_LAND_FIRST_BOOL}

	// LxorFirstBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [First].
	//
	// LxorFirstBool is a SuiteSparse:GraphBLAS extension.
	LxorFirstBool = Semiring[bool, bool, bool]{C.GxB_LXOR_FIRST_BOOL}

	// EqFirstBool semiring with additive [Monoid] [Eq] and [BinaryOp] [First].
	//
	// EqFirstBool is a SuiteSparse:GraphBLAS extension.
	EqFirstBool = Semiring[bool, bool, bool]{C.GxB_EQ_FIRST_BOOL}

	// LorSecondBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Second].
	//
	// LorSecondBool is a SuiteSparse:GraphBLAS extension.
	LorSecondBool = Semiring[bool, bool, bool]{C.GxB_LOR_SECOND_BOOL}

	// LandSecondBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Second].
	//
	// LandSecondBool is a SuiteSparse:GraphBLAS extension.
	LandSecondBool = Semiring[bool, bool, bool]{C.GxB_LAND_SECOND_BOOL}

	// LxorSecondBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Second].
	//
	// LxorSecondBool is a SuiteSparse:GraphBLAS extension.
	LxorSecondBool = Semiring[bool, bool, bool]{C.GxB_LXOR_SECOND_BOOL}

	// EqSecondBool semiring with additive [Monoid] [Eq] and [BinaryOp] [Second].
	//
	// EqSecondBool is a SuiteSparse:GraphBLAS extension.
	EqSecondBool = Semiring[bool, bool, bool]{C.GxB_EQ_SECOND_BOOL}

	// LorOnebBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Oneb].
	//
	// LorOnebBool is a SuiteSparse:GraphBLAS extension.
	LorOnebBool = Semiring[bool, bool, bool]{C.GxB_LOR_PAIR_BOOL}

	// LandOnebBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Oneb].
	//
	// LandOnebBool is a SuiteSparse:GraphBLAS extension.
	LandOnebBool = Semiring[bool, bool, bool]{C.GxB_LAND_PAIR_BOOL}

	// LxorOnebBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Oneb].
	//
	// LxorOnebBool is a SuiteSparse:GraphBLAS extension.
	LxorOnebBool = Semiring[bool, bool, bool]{C.GxB_LXOR_PAIR_BOOL}

	// EqOnebBool semiring with additive [Monoid] [EqMonoidBool] and [BinaryOp] [Oneb].
	//
	// EqOnebBool is a SuiteSparse:GraphBLAS extension.
	EqOnebBool = Semiring[bool, bool, bool]{C.GxB_EQ_PAIR_BOOL}

	// LorLorBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [LorBool].
	//
	// LorLorBool is a SuiteSparse:GraphBLAS extension.
	LorLorBool = Semiring[bool, bool, bool]{C.GxB_LOR_LOR_BOOL}

	// LandLorBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [LorBool].
	//
	// LandLorBool is a SuiteSparse:GraphBLAS extension.
	LandLorBool = Semiring[bool, bool, bool]{C.GxB_LAND_LOR_BOOL}

	// LxorLorBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [LorBool].
	//
	// LxorLorBool is a SuiteSparse:GraphBLAS extension.
	LxorLorBool = Semiring[bool, bool, bool]{C.GxB_LXOR_LOR_BOOL}

	// EqLorBool semiring with additive [Monoid] [Eq] and [BinaryOp] [LorBool].
	//
	// EqLorBool is a SuiteSparse:GraphBLAS extension.
	EqLorBool = Semiring[bool, bool, bool]{C.GxB_EQ_LOR_BOOL}

	// LorLandBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [LandBool].
	//
	// LorLandBool is a SuiteSparse:GraphBLAS extension.
	LorLandBool = Semiring[bool, bool, bool]{C.GxB_LOR_LAND_BOOL}

	// LandLandBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [LandBool].
	//
	// LandLandBool is a SuiteSparse:GraphBLAS extension.
	LandLandBool = Semiring[bool, bool, bool]{C.GxB_LAND_LAND_BOOL}

	// LxorLandBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [LandBool].
	//
	// LxorLandBool is a SuiteSparse:GraphBLAS extension.
	LxorLandBool = Semiring[bool, bool, bool]{C.GxB_LXOR_LAND_BOOL}

	// EqLandBool semiring with additive [Monoid] [Eq] and [BinaryOp] [LandBool].
	//
	// EqLandBool is a SuiteSparse:GraphBLAS extension.
	EqLandBool = Semiring[bool, bool, bool]{C.GxB_EQ_LAND_BOOL}

	// LorLxorBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [LxorBool].
	//
	// LorLxorBool is a SuiteSparse:GraphBLAS extension.
	LorLxorBool = Semiring[bool, bool, bool]{C.GxB_LOR_LXOR_BOOL}

	// LandLxorBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [LxorBool].
	//
	// LandLxorBool is a SuiteSparse:GraphBLAS extension.
	LandLxorBool = Semiring[bool, bool, bool]{C.GxB_LAND_LXOR_BOOL}

	// LxorLxorBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [LxorBool].
	//
	// LxorLxorBool is a SuiteSparse:GraphBLAS extension.
	LxorLxorBool = Semiring[bool, bool, bool]{C.GxB_LXOR_LXOR_BOOL}

	// EqLxorBool semiring with additive [Monoid] [Eq] and [BinaryOp] [LxorBool].
	//
	// EqLxorBool is a SuiteSparse:GraphBLAS extension.
	EqLxorBool = Semiring[bool, bool, bool]{C.GxB_EQ_LXOR_BOOL}

	// LorEqBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Eq].
	//
	// LorEqBool is a SuiteSparse:GraphBLAS extension.
	LorEqBool = Semiring[bool, bool, bool]{C.GxB_LOR_EQ_BOOL}

	// LandEqBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Eq].
	//
	// LandEqBool is a SuiteSparse:GraphBLAS extension.
	LandEqBool = Semiring[bool, bool, bool]{C.GxB_LAND_EQ_BOOL}

	// LxorEqBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Eq].
	//
	// LxorEqBool is a SuiteSparse:GraphBLAS extension.
	LxorEqBool = Semiring[bool, bool, bool]{C.GxB_LXOR_EQ_BOOL}

	// EqEqBool semiring with additive [Monoid] [Eq] and [BinaryOp] [Eq].
	//
	// EqEqBool is a SuiteSparse:GraphBLAS extension.
	EqEqBool = Semiring[bool, bool, bool]{C.GxB_EQ_EQ_BOOL}

	// LorGtBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Gt].
	//
	// LorGtBool is a SuiteSparse:GraphBLAS extension.
	LorGtBool = Semiring[bool, bool, bool]{C.GxB_LOR_GT_BOOL}

	// LandGtBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Gt].
	//
	// LandGtBool is a SuiteSparse:GraphBLAS extension.
	LandGtBool = Semiring[bool, bool, bool]{C.GxB_LAND_GT_BOOL}

	// LxorGtBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Gt].
	//
	// LxorGtBool is a SuiteSparse:GraphBLAS extension.
	LxorGtBool = Semiring[bool, bool, bool]{C.GxB_LXOR_GT_BOOL}

	// EqGtBool semiring with additive [Monoid] [EqMonoidBool] and [BinaryOp] [Gt].
	//
	// EqGtBool is a SuiteSparse:GraphBLAS extension.
	EqGtBool = Semiring[bool, bool, bool]{C.GxB_EQ_GT_BOOL}

	// LorLtBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Lt].
	//
	// LorLtBool is a SuiteSparse:GraphBLAS extension.
	LorLtBool = Semiring[bool, bool, bool]{C.GxB_LOR_LT_BOOL}

	// LandLtBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Lt].
	//
	// LandLtBool is a SuiteSparse:GraphBLAS extension.
	LandLtBool = Semiring[bool, bool, bool]{C.GxB_LAND_LT_BOOL}

	// LxorLtBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Lt].
	//
	// LxorLtBool is a SuiteSparse:GraphBLAS extension.
	LxorLtBool = Semiring[bool, bool, bool]{C.GxB_LXOR_LT_BOOL}

	// EqLtBool semiring with additive [Monoid] [Eq] and [BinaryOp] [Lt].
	//
	// EqLtBool is a SuiteSparse:GraphBLAS extension.
	EqLtBool = Semiring[bool, bool, bool]{C.GxB_EQ_LT_BOOL}

	// LorGeBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Ge].
	//
	// LorGeBool is a SuiteSparse:GraphBLAS extension.
	LorGeBool = Semiring[bool, bool, bool]{C.GxB_LOR_GE_BOOL}

	// LandGeBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Ge].
	//
	// LandGeBool is a SuiteSparse:GraphBLAS extension.
	LandGeBool = Semiring[bool, bool, bool]{C.GxB_LAND_GE_BOOL}

	// LxorGeBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Ge].
	//
	// LxorGeBool is a SuiteSparse:GraphBLAS extension.
	LxorGeBool = Semiring[bool, bool, bool]{C.GxB_LXOR_GE_BOOL}

	// EqGeBool semiring with additive [Monoid] [Eq] and [BinaryOp] [Ge].
	//
	// EqGeBool is a SuiteSparse:GraphBLAS extension.
	EqGeBool = Semiring[bool, bool, bool]{C.GxB_EQ_GE_BOOL}

	// LorLeBool semiring with additive [Monoid] [LorMonoidBool] and [BinaryOp] [Le].
	//
	// LorLeBool is a SuiteSparse:GraphBLAS extension.
	LorLeBool = Semiring[bool, bool, bool]{C.GxB_LOR_LE_BOOL}

	// LandLeBool semiring with additive [Monoid] [LandMonoidBool] and [BinaryOp] [Le].
	//
	// LandLeBool is a SuiteSparse:GraphBLAS extension.
	LandLeBool = Semiring[bool, bool, bool]{C.GxB_LAND_LE_BOOL}

	// LxorLeBool semiring with additive [Monoid] [LxorMonoidBool] and [BinaryOp] [Le].
	//
	// LxorLeBool is a SuiteSparse:GraphBLAS extension.
	LxorLeBool = Semiring[bool, bool, bool]{C.GxB_LXOR_LE_BOOL}

	// EqLeBool semiring with additive [Monoid] [Eq] and [BinaryOp] [Le].
	//
	// EqLeBool is a SuiteSparse:GraphBLAS extension.
	EqLeBool = Semiring[bool, bool, bool]{C.GxB_EQ_LE_BOOL}
)

// BorBor semiring with additive [Monoid] [BorMonoid] and [BinaryOp] [Bor].
//
// BorBor is a SuiteSparse:GraphBLAS extension.
func BorBor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BOR_BOR_UINT32
		} else {
			s.grb = C.GxB_BOR_BOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BOR_BOR_UINT8
	case uint16:
		s.grb = C.GxB_BOR_BOR_UINT16
	case uint32:
		s.grb = C.GxB_BOR_BOR_UINT32
	case uint64:
		s.grb = C.GxB_BOR_BOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BorBand semiring with additive [Monoid] [BorMonoid] and [BinaryOp] [Band].
//
// BorBand is a SuiteSparse:GraphBLAS extension.
func BorBand[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BOR_BAND_UINT32
		} else {
			s.grb = C.GxB_BOR_BAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_BOR_BAND_UINT8
	case uint16:
		s.grb = C.GxB_BOR_BAND_UINT16
	case uint32:
		s.grb = C.GxB_BOR_BAND_UINT32
	case uint64:
		s.grb = C.GxB_BOR_BAND_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BorBxor semiring with additive [Monoid] [BorMonoid] and [BinaryOp] [Bxor].
//
// BorBxor is a SuiteSparse:GraphBLAS extension.
func BorBxor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BOR_BXOR_UINT32
		} else {
			s.grb = C.GxB_BOR_BXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BOR_BXOR_UINT8
	case uint16:
		s.grb = C.GxB_BOR_BXOR_UINT16
	case uint32:
		s.grb = C.GxB_BOR_BXOR_UINT32
	case uint64:
		s.grb = C.GxB_BOR_BXOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BorBxnor semiring with additive [Monoid] [BorMonoid] and [BinaryOp] [Bxnor].
//
// BorBxnor is a SuiteSparse:GraphBLAS extension.
func BorBxnor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BOR_BXNOR_UINT32
		} else {
			s.grb = C.GxB_BOR_BXNOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BOR_BXNOR_UINT8
	case uint16:
		s.grb = C.GxB_BOR_BXNOR_UINT16
	case uint32:
		s.grb = C.GxB_BOR_BXNOR_UINT32
	case uint64:
		s.grb = C.GxB_BOR_BXNOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BandBor semiring with additive [Monoid] [BandMonoid] and [BinaryOp] [Bor].
//
// BandBor is a SuiteSparse:GraphBLAS extension.
func BandBor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BAND_BOR_UINT32
		} else {
			s.grb = C.GxB_BAND_BOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BAND_BOR_UINT8
	case uint16:
		s.grb = C.GxB_BAND_BOR_UINT16
	case uint32:
		s.grb = C.GxB_BAND_BOR_UINT32
	case uint64:
		s.grb = C.GxB_BAND_BOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BandBand semiring with additive [Monoid] [BandMonoid] and [BinaryOp] [Band].
//
// BandBand is a SuiteSparse:GraphBLAS extension.
func BandBand[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BAND_BAND_UINT32
		} else {
			s.grb = C.GxB_BAND_BAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_BAND_BAND_UINT8
	case uint16:
		s.grb = C.GxB_BAND_BAND_UINT16
	case uint32:
		s.grb = C.GxB_BAND_BAND_UINT32
	case uint64:
		s.grb = C.GxB_BAND_BAND_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BandBxor semiring with additive [Monoid] [BandMonoid] and [BinaryOp] [Bxor].
//
// BandBxor is a SuiteSparse:GraphBLAS extension.
func BandBxor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BAND_BXOR_UINT32
		} else {
			s.grb = C.GxB_BAND_BXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BAND_BXOR_UINT8
	case uint16:
		s.grb = C.GxB_BAND_BXOR_UINT16
	case uint32:
		s.grb = C.GxB_BAND_BXOR_UINT32
	case uint64:
		s.grb = C.GxB_BAND_BXOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BandBxnor semiring with additive [Monoid] [BandMonoid] and [BinaryOp] [Bxnor].
//
// BandBxnor is a SuiteSparse:GraphBLAS extension.
func BandBxnor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BAND_BXNOR_UINT32
		} else {
			s.grb = C.GxB_BAND_BXNOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BAND_BXNOR_UINT8
	case uint16:
		s.grb = C.GxB_BAND_BXNOR_UINT16
	case uint32:
		s.grb = C.GxB_BAND_BXNOR_UINT32
	case uint64:
		s.grb = C.GxB_BAND_BXNOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxorBor semiring with additive [Monoid] [BxorMonoid] and [BinaryOp] [Bor].
//
// BxorBor is a SuiteSparse:GraphBLAS extension.
func BxorBor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXOR_BOR_UINT32
		} else {
			s.grb = C.GxB_BXOR_BOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXOR_BOR_UINT8
	case uint16:
		s.grb = C.GxB_BXOR_BOR_UINT16
	case uint32:
		s.grb = C.GxB_BXOR_BOR_UINT32
	case uint64:
		s.grb = C.GxB_BXOR_BOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxorBand semiring with additive [Monoid] [BxorMonoid] and [BinaryOp] [Band].
//
// BxorBand is a SuiteSparse:GraphBLAS extension.
// is a SuiteSparse:GraphBLAS extension.
func BxorBand[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXOR_BAND_UINT32
		} else {
			s.grb = C.GxB_BXOR_BAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXOR_BAND_UINT8
	case uint16:
		s.grb = C.GxB_BXOR_BAND_UINT16
	case uint32:
		s.grb = C.GxB_BXOR_BAND_UINT32
	case uint64:
		s.grb = C.GxB_BXOR_BAND_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxorBxor semiring with additive [Monoid] [BxorMonoid] and [BinaryOp] [Bxor].
//
// BxorBxor is a SuiteSparse:GraphBLAS extension.
func BxorBxor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXOR_BXOR_UINT32
		} else {
			s.grb = C.GxB_BXOR_BXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXOR_BXOR_UINT8
	case uint16:
		s.grb = C.GxB_BXOR_BXOR_UINT16
	case uint32:
		s.grb = C.GxB_BXOR_BXOR_UINT32
	case uint64:
		s.grb = C.GxB_BXOR_BXOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxorBxnor semiring with additive [Monoid] [BxorMonoid] and [BinaryOp] [Bxnor].
//
// BxorBxnor is a SuiteSparse:GraphBLAS extension.
func BxorBxnor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXOR_BXNOR_UINT32
		} else {
			s.grb = C.GxB_BXOR_BXNOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXOR_BXNOR_UINT8
	case uint16:
		s.grb = C.GxB_BXOR_BXNOR_UINT16
	case uint32:
		s.grb = C.GxB_BXOR_BXNOR_UINT32
	case uint64:
		s.grb = C.GxB_BXOR_BXNOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxnorBor semiring with additive [Monoid] [BxnorMonoid] and [BinaryOp] [Bor].
//
// BxnorBor is a SuiteSparse:GraphBLAS extension.
func BxnorBor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXNOR_BOR_UINT32
		} else {
			s.grb = C.GxB_BXNOR_BOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXNOR_BOR_UINT8
	case uint16:
		s.grb = C.GxB_BXNOR_BOR_UINT16
	case uint32:
		s.grb = C.GxB_BXNOR_BOR_UINT32
	case uint64:
		s.grb = C.GxB_BXNOR_BOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxnorBand semiring with additive [Monoid] [BxnorMonoid] and [BinaryOp] [Band].
//
// BxnorBand is a SuiteSparse:GraphBLAS extension.
func BxnorBand[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXNOR_BAND_UINT32
		} else {
			s.grb = C.GxB_BXNOR_BAND_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXNOR_BAND_UINT8
	case uint16:
		s.grb = C.GxB_BXNOR_BAND_UINT16
	case uint32:
		s.grb = C.GxB_BXNOR_BAND_UINT32
	case uint64:
		s.grb = C.GxB_BXNOR_BAND_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxnorBxor semiring with additive [Monoid] [BxnorMonoid] and [BinaryOp] [Bxor].
//
// BxnorBxor is a SuiteSparse:GraphBLAS extension.
func BxnorBxor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXNOR_BXOR_UINT32
		} else {
			s.grb = C.GxB_BXNOR_BXOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXNOR_BXOR_UINT8
	case uint16:
		s.grb = C.GxB_BXNOR_BXOR_UINT16
	case uint32:
		s.grb = C.GxB_BXNOR_BXOR_UINT32
	case uint64:
		s.grb = C.GxB_BXNOR_BXOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// BxnorBxnor semiring with additive [Monoid] [BxnorMonoid] and [BinaryOp] [Bxnor].
//
// BxnorBxnor is a SuiteSparse:GraphBLAS extension.
func BxnorBxnor[D Unsigned]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case uint:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_BXNOR_BXNOR_UINT32
		} else {
			s.grb = C.GxB_BXNOR_BXNOR_UINT64
		}
	case uint8:
		s.grb = C.GxB_BXNOR_BXNOR_UINT8
	case uint16:
		s.grb = C.GxB_BXNOR_BXNOR_UINT16
	case uint32:
		s.grb = C.GxB_BXNOR_BXNOR_UINT32
	case uint64:
		s.grb = C.GxB_BXNOR_BXNOR_UINT64
	default:
		panic("unreachable code")
	}
	return
}

// MinFirsti semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Firsti].
//
// MinFirsti is a SuiteSparse:GraphBLAS extension.
func MinFirsti[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_FIRSTI_INT32
		} else {
			s.grb = C.GxB_MIN_FIRSTI_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_FIRSTI_INT32
	case int64:
		s.grb = C.GxB_MIN_FIRSTI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxFirsti semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Firsti].
//
// MaxFirsti is a SuiteSparse:GraphBLAS extension.
func MaxFirsti[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_FIRSTI_INT32
		} else {
			s.grb = C.GxB_MAX_FIRSTI_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_FIRSTI_INT32
	case int64:
		s.grb = C.GxB_MAX_FIRSTI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnyFirsti semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Firsti].
//
// AnyFirsti is a SuiteSparse:GraphBLAS extension.
func AnyFirsti[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_FIRSTI_INT32
		} else {
			s.grb = C.GxB_ANY_FIRSTI_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_FIRSTI_INT32
	case int64:
		s.grb = C.GxB_ANY_FIRSTI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusFirsti semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Firsti].
//
// PlusFirsti is a SuiteSparse:GraphBLAS extension.
func PlusFirsti[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_FIRSTI_INT32
		} else {
			s.grb = C.GxB_PLUS_FIRSTI_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_FIRSTI_INT32
	case int64:
		s.grb = C.GxB_PLUS_FIRSTI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesFirsti semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Firsti].
//
// TimesFirsti is a SuiteSparse:GraphBLAS extension.
func TimesFirsti[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_FIRSTI_INT32
		} else {
			s.grb = C.GxB_TIMES_FIRSTI_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_FIRSTI_INT32
	case int64:
		s.grb = C.GxB_TIMES_FIRSTI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinFirsti1 semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Firsti1].
//
// MinFirsti1 is a SuiteSparse:GraphBLAS extension.
func MinFirsti1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_FIRSTI1_INT32
		} else {
			s.grb = C.GxB_MIN_FIRSTI1_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_FIRSTI1_INT32
	case int64:
		s.grb = C.GxB_MIN_FIRSTI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxFirsti1 semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Firsti1].
//
// MaxFirsti1 is a SuiteSparse:GraphBLAS extension.
func MaxFirsti1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_FIRSTI1_INT32
		} else {
			s.grb = C.GxB_MAX_FIRSTI1_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_FIRSTI1_INT32
	case int64:
		s.grb = C.GxB_MAX_FIRSTI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnyFirsti1 semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Firsti1].
//
// AnyFirsti1 is a SuiteSparse:GraphBLAS extension.
func AnyFirsti1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_FIRSTI1_INT32
		} else {
			s.grb = C.GxB_ANY_FIRSTI1_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_FIRSTI1_INT32
	case int64:
		s.grb = C.GxB_ANY_FIRSTI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusFirsti1 semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Firsti1].
//
// PlusFirsti1 is a SuiteSparse:GraphBLAS extension.
func PlusFirsti1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_FIRSTI1_INT32
		} else {
			s.grb = C.GxB_PLUS_FIRSTI1_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_FIRSTI1_INT32
	case int64:
		s.grb = C.GxB_PLUS_FIRSTI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesFirsti1 semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Firsti1].
//
// TimesFirsti1 is a SuiteSparse:GraphBLAS extension.
func TimesFirsti1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_FIRSTI1_INT32
		} else {
			s.grb = C.GxB_TIMES_FIRSTI1_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_FIRSTI1_INT32
	case int64:
		s.grb = C.GxB_TIMES_FIRSTI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinFirstj semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Firstj].
//
// MinFirstj is a SuiteSparse:GraphBLAS extension.
func MinFirstj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_FIRSTJ_INT32
		} else {
			s.grb = C.GxB_MIN_FIRSTJ_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_FIRSTJ_INT32
	case int64:
		s.grb = C.GxB_MIN_FIRSTJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxFirstj semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Firstj].
//
// MaxFirstj is a SuiteSparse:GraphBLAS extension.
func MaxFirstj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_FIRSTJ_INT32
		} else {
			s.grb = C.GxB_MAX_FIRSTJ_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_FIRSTJ_INT32
	case int64:
		s.grb = C.GxB_MAX_FIRSTJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnyFirstj semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Firstj].
//
// AnyFirstj is a SuiteSparse:GraphBLAS extension.
func AnyFirstj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_FIRSTJ_INT32
		} else {
			s.grb = C.GxB_ANY_FIRSTJ_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_FIRSTJ_INT32
	case int64:
		s.grb = C.GxB_ANY_FIRSTJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusFirstj semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Firstj].
//
// PlusFirstj is a SuiteSparse:GraphBLAS extension.
func PlusFirstj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_FIRSTJ_INT32
		} else {
			s.grb = C.GxB_PLUS_FIRSTJ_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_FIRSTJ_INT32
	case int64:
		s.grb = C.GxB_PLUS_FIRSTJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesFirstj semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Firstj].
//
// TimesFirstj is a SuiteSparse:GraphBLAS extension.
func TimesFirstj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_FIRSTJ_INT32
		} else {
			s.grb = C.GxB_TIMES_FIRSTJ_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_FIRSTJ_INT32
	case int64:
		s.grb = C.GxB_TIMES_FIRSTJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinFirstj1 semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Firstj1].
//
// MinFirstj1 is a SuiteSparse:GraphBLAS extension.
func MinFirstj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_FIRSTJ1_INT32
		} else {
			s.grb = C.GxB_MIN_FIRSTJ1_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_FIRSTJ1_INT32
	case int64:
		s.grb = C.GxB_MIN_FIRSTJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxFirstj1 semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Firstj1].
//
// MaxFirstj1 is a SuiteSparse:GraphBLAS extension.
func MaxFirstj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_FIRSTJ1_INT32
		} else {
			s.grb = C.GxB_MAX_FIRSTJ1_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_FIRSTJ1_INT32
	case int64:
		s.grb = C.GxB_MAX_FIRSTJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnyFirstj1 semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Firstj1].
//
// AnyFirstj1 is a SuiteSparse:GraphBLAS extension.
func AnyFirstj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_FIRSTJ1_INT32
		} else {
			s.grb = C.GxB_ANY_FIRSTJ1_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_FIRSTJ1_INT32
	case int64:
		s.grb = C.GxB_ANY_FIRSTJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusFirstj1 semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Firstj1].
//
// PlusFirstj1 is a SuiteSparse:GraphBLAS extension.
func PlusFirstj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_FIRSTJ1_INT32
		} else {
			s.grb = C.GxB_PLUS_FIRSTJ1_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_FIRSTJ1_INT32
	case int64:
		s.grb = C.GxB_PLUS_FIRSTJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesFirstj1 semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Firstj1].
//
// TimesFirstj1 is a SuiteSparse:GraphBLAS extension.
func TimesFirstj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_FIRSTJ1_INT32
		} else {
			s.grb = C.GxB_TIMES_FIRSTJ1_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_FIRSTJ1_INT32
	case int64:
		s.grb = C.GxB_TIMES_FIRSTJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinSecondi semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Secondi].
//
// MinSecondi is a SuiteSparse:GraphBLAS extension.
func MinSecondi[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_SECONDI_INT32
		} else {
			s.grb = C.GxB_MIN_SECONDI_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_SECONDI_INT32
	case int64:
		s.grb = C.GxB_MIN_SECONDI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxSecondi semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Secondi].
//
// MaxSecondi is a SuiteSparse:GraphBLAS extension.
func MaxSecondi[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_SECONDI_INT32
		} else {
			s.grb = C.GxB_MAX_SECONDI_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_SECONDI_INT32
	case int64:
		s.grb = C.GxB_MAX_SECONDI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnySecondi semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Secondi].
//
// AnySecondi is a SuiteSparse:GraphBLAS extension.
func AnySecondi[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_SECONDI_INT32
		} else {
			s.grb = C.GxB_ANY_SECONDI_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_SECONDI_INT32
	case int64:
		s.grb = C.GxB_ANY_SECONDI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusSecondi semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Secondi].
//
// PlusSecondi is a SuiteSparse:GraphBLAS extension.
func PlusSecondi[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_SECONDI_INT32
		} else {
			s.grb = C.GxB_PLUS_SECONDI_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_SECONDI_INT32
	case int64:
		s.grb = C.GxB_PLUS_SECONDI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesSecondi semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Secondi].
//
// TimesSecondi is a SuiteSparse:GraphBLAS extension.
func TimesSecondi[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_SECONDI_INT32
		} else {
			s.grb = C.GxB_TIMES_SECONDI_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_SECONDI_INT32
	case int64:
		s.grb = C.GxB_TIMES_SECONDI_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinSecondi1 semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Secondi1].
//
// MinSecondi1 is a SuiteSparse:GraphBLAS extension.
func MinSecondi1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_SECONDI1_INT32
		} else {
			s.grb = C.GxB_MIN_SECONDI1_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_SECONDI1_INT32
	case int64:
		s.grb = C.GxB_MIN_SECONDI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxSecondi1 semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Secondi1].
//
// MaxSecondi1 is a SuiteSparse:GraphBLAS extension.
func MaxSecondi1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_SECONDI1_INT32
		} else {
			s.grb = C.GxB_MAX_SECONDI1_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_SECONDI1_INT32
	case int64:
		s.grb = C.GxB_MAX_SECONDI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnySecondi1 semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Secondi1].
//
// AnySecondi1 is a SuiteSparse:GraphBLAS extension.
func AnySecondi1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_SECONDI1_INT32
		} else {
			s.grb = C.GxB_ANY_SECONDI1_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_SECONDI1_INT32
	case int64:
		s.grb = C.GxB_ANY_SECONDI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusSecondi1 semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Secondi1].
//
// PlusSecondi1 is a SuiteSparse:GraphBLAS extension.
func PlusSecondi1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_SECONDI1_INT32
		} else {
			s.grb = C.GxB_PLUS_SECONDI1_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_SECONDI1_INT32
	case int64:
		s.grb = C.GxB_PLUS_SECONDI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesSecondi1 semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Secondi1].
//
// TimesSecondi1 is a SuiteSparse:GraphBLAS extension.
func TimesSecondi1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_SECONDI1_INT32
		} else {
			s.grb = C.GxB_TIMES_SECONDI1_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_SECONDI1_INT32
	case int64:
		s.grb = C.GxB_TIMES_SECONDI1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinSecondj semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Secondj].
//
// MinSecondj is a SuiteSparse:GraphBLAS extension.
func MinSecondj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_SECONDJ_INT32
		} else {
			s.grb = C.GxB_MIN_SECONDJ_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_SECONDJ_INT32
	case int64:
		s.grb = C.GxB_MIN_SECONDJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxSecondj semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Secondj].
//
// MaxSecondj is a SuiteSparse:GraphBLAS extension.
func MaxSecondj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_SECONDJ_INT32
		} else {
			s.grb = C.GxB_MAX_SECONDJ_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_SECONDJ_INT32
	case int64:
		s.grb = C.GxB_MAX_SECONDJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnySecondj semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Secondj].
//
// AnySecondj is a SuiteSparse:GraphBLAS extension.
func AnySecondj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_SECONDJ_INT32
		} else {
			s.grb = C.GxB_ANY_SECONDJ_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_SECONDJ_INT32
	case int64:
		s.grb = C.GxB_ANY_SECONDJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusSecondj semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Secondj].
//
// PlusSecondj is a SuiteSparse:GraphBLAS extension.
func PlusSecondj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_SECONDJ_INT32
		} else {
			s.grb = C.GxB_PLUS_SECONDJ_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_SECONDJ_INT32
	case int64:
		s.grb = C.GxB_PLUS_SECONDJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesSecondj semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Secondj].
//
// TimesSecondj is a SuiteSparse:GraphBLAS extension.
func TimesSecondj[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_SECONDJ_INT32
		} else {
			s.grb = C.GxB_TIMES_SECONDJ_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_SECONDJ_INT32
	case int64:
		s.grb = C.GxB_TIMES_SECONDJ_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MinSecondj1 semiring with additive [Monoid] [MinMonoid] and [BinaryOp] [Secondj1].
//
// MinSecondj1 is a SuiteSparse:GraphBLAS extension.
func MinSecondj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MIN_SECONDJ1_INT32
		} else {
			s.grb = C.GxB_MIN_SECONDJ1_INT64
		}
	case int32:
		s.grb = C.GxB_MIN_SECONDJ1_INT32
	case int64:
		s.grb = C.GxB_MIN_SECONDJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// MaxSecondj1 semiring with additive [Monoid] [MaxMonoid] and [BinaryOp] [Secondj1].
//
// MaxSecondj1 is a SuiteSparse:GraphBLAS extension.
func MaxSecondj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_MAX_SECONDJ1_INT32
		} else {
			s.grb = C.GxB_MAX_SECONDJ1_INT64
		}
	case int32:
		s.grb = C.GxB_MAX_SECONDJ1_INT32
	case int64:
		s.grb = C.GxB_MAX_SECONDJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// AnySecondj1 semiring with additive [Monoid] [AnyMonoid] and [BinaryOp] [Secondj1].
//
// AnySecondj1 is a SuiteSparse:GraphBLAS extension.
func AnySecondj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_ANY_SECONDJ1_INT32
		} else {
			s.grb = C.GxB_ANY_SECONDJ1_INT64
		}
	case int32:
		s.grb = C.GxB_ANY_SECONDJ1_INT32
	case int64:
		s.grb = C.GxB_ANY_SECONDJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// PlusSecondj1 semiring with additive [Monoid] [PlusMonoid] and [BinaryOp] [Secondj1].
//
// PlusSecondj1 is a SuiteSparse:GraphBLAS extension.
func PlusSecondj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_PLUS_SECONDJ1_INT32
		} else {
			s.grb = C.GxB_PLUS_SECONDJ1_INT64
		}
	case int32:
		s.grb = C.GxB_PLUS_SECONDJ1_INT32
	case int64:
		s.grb = C.GxB_PLUS_SECONDJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}

// TimesSecondj1 semiring with additive [Monoid] [TimesMonoid] and [BinaryOp] [Secondj1].
//
// TimesSecondj1 is a SuiteSparse:GraphBLAS extension.
func TimesSecondj1[D int32 | int64 | int]() (s Semiring[D, D, D]) {
	var d D
	switch any(d).(type) {
	case int:
		if unsafe.Sizeof(0) == 4 {
			s.grb = C.GxB_TIMES_SECONDJ1_INT32
		} else {
			s.grb = C.GxB_TIMES_SECONDJ1_INT64
		}
	case int32:
		s.grb = C.GxB_TIMES_SECONDJ1_INT32
	case int64:
		s.grb = C.GxB_TIMES_SECONDJ1_INT64
	default:
		panic("unreachable code")
	}
	return
}
