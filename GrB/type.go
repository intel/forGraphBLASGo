package GrB

// #include "GraphBLAS.h"
import "C"
import (
	"math"
	"reflect"
	"unsafe"
)

// A Type represents a domain for elements that can be stored in collections and operated on through
// GraphBLAS functions.
type Type struct {
	typ reflect.Type
}

// TypeOf returns the [Type] for a given Go value. The return value is only safe to
// use if it corresponds to a predefined GraphBLAS type, or if it has been created
// by [TypeNew] or [NamedTypeNew] as a user-defined type.
//
// TypeOf is a forGraphBLASGo extension.
func TypeOf(x any) Type {
	return Type{reflect.TypeOf(x)}
}

// Predefined GraphBLAS types (domains).
var (
	Bool       = TypeOf(false)
	Int        = TypeOf(int(0))
	Int8       = TypeOf(int8(0))
	Int16      = TypeOf(int16(0))
	Int32      = TypeOf(int32(0))
	Int64      = TypeOf(int64(0))
	Uint       = TypeOf(uint(0))
	Uint8      = TypeOf(uint8(0))
	Uint16     = TypeOf(uint16(0))
	Uint32     = TypeOf(uint32(0))
	Uint64     = TypeOf(uint64(0))
	Float32    = TypeOf(float32(0))
	Float64    = TypeOf(float64(0))
	Complex64  = TypeOf(complex64(0))  // a SuiteSparse:GraphBLAS extension
	Complex128 = TypeOf(complex128(0)) // a SuiteSparse:GraphBLAS extension
)

var (
	grbType = map[Type]C.GrB_Type{
		Bool:       C.GrB_BOOL,
		Int:        C.GrB_INT64,
		Int8:       C.GrB_INT8,
		Int16:      C.GrB_INT16,
		Int32:      C.GrB_INT32,
		Int64:      C.GrB_INT64,
		Uint:       C.GrB_UINT64,
		Uint8:      C.GrB_UINT8,
		Uint16:     C.GrB_UINT16,
		Uint32:     C.GrB_UINT32,
		Uint64:     C.GrB_UINT64,
		Float32:    C.GrB_FP32,
		Float64:    C.GrB_FP64,
		Complex64:  C.GxB_FC32,
		Complex128: C.GxB_FC64,
	}

	goType = map[C.GrB_Type]Type{
		C.GrB_BOOL:   Bool,
		C.GrB_INT8:   Int8,
		C.GrB_INT16:  Int16,
		C.GrB_INT32:  Int32,
		C.GrB_INT64:  Int64,
		C.GrB_UINT8:  Uint8,
		C.GrB_UINT16: Uint16,
		C.GrB_UINT32: Uint32,
		C.GrB_UINT64: Uint64,
		C.GrB_FP32:   Float32,
		C.GrB_FP64:   Float64,
		C.GxB_FC32:   Complex64,
		C.GxB_FC64:   Complex128,
	}
)

func init() {
	if unsafe.Sizeof(0) == 4 {
		grbType[TypeOf(int(0))] = C.GrB_INT32
	}
	if unsafe.Sizeof(0) == 4 {
		grbType[TypeOf(uint(0))] = C.GrB_INT32
	}
}

// Meaningful type constraints for the predefined GraphBLAS types (domains).
//
// These are forGraphBLASGo extension.
type (
	Signed interface {
		int | int8 | int16 | int32 | int64
	}

	Unsigned interface {
		uint | uint8 | uint16 | uint32 | uint64
	}

	Integer interface {
		Signed | Unsigned
	}

	Float interface {
		float32 | float64
	}

	Complex interface {
		complex64 | complex128
	}

	Number interface {
		Integer | Float
	}

	Predefined interface {
		bool | Number
	}
)

// Minimum values for the predefined GraphBLAS [Number] types (domains).
//
// This returns the smallest negative value for signed integers,
// 0 for unsigned integers, and negative infinity for floating point
// numbers.
//
// Minimum is a forGraphBLAS Go extension.
func Minimum[D Number]() (result D) {
	switch x := any(&result).(type) {
	case *int:
		*x = math.MinInt
	case *int8:
		*x = math.MinInt8
	case *int16:
		*x = math.MinInt16
	case *int32:
		*x = math.MinInt32
	case *int64:
		*x = math.MinInt64
	case *uint, *uint8, *uint16, *uint32, *uint64:
		result = 0
	case *float32:
		*x = float32(math.Inf(-1))
	case *float64:
		*x = math.Inf(-1)
	default:
		panic("unreachable code")
	}
	return
}

// Maximum values for the predefined GraphBLAS [Number] types (domains).
//
// This returns the largest positive value for integers,
// and positive infinity for floating point numbers.
//
// Maximum is a forGraphBLAS Go extension.
func Maximum[D Number]() (result D) {
	switch x := any(&result).(type) {
	case *int:
		*x = math.MaxInt
	case *int8:
		*x = math.MaxInt8
	case *int16:
		*x = math.MaxInt16
	case *int32:
		*x = math.MaxInt32
	case *int64:
		*x = math.MaxInt64
	case *uint:
		*x = math.MaxUint
	case *uint8:
		*x = math.MaxUint8
	case *uint16:
		*x = math.MaxUint16
	case *uint32:
		*x = math.MaxUint32
	case *uint64:
		*x = math.MaxUint64
	case *float32:
		*x = float32(math.Inf(+1))
	case *float64:
		*x = math.Inf(+1)
	default:
		panic("unreachable code")
	}
	return
}

func hasGoPointer(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Invalid:
		panic("Encountered an invalid type kind.")
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return false
	case reflect.Array:
		return hasGoPointer(t.Elem())
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.String:
		return true
	case reflect.Struct:
		for f, n := 0, t.NumField(); f < n; f++ {
			if hasGoPointer(t.Field(f).Type) {
				return true
			}
		}
		return false
	case reflect.UnsafePointer:
		return true
	}
	panic("unreachable code")
}

// TypeNew creates a new user-defined GraphBLAS type. This type can then be used
// to create new operators, monoids, semirings, vectors and matrices.
//
// Variables of this type must be a struct or (fixed-size) array. In particular, given two variables, src and dst,
// of type D, the following operation must be a valid way to copy the contents of src to dst:
//
//	memcpy(&dst, &src, sizeof(D))
//
// Parameters:
//   - D: A Go type. Values of this type will be copied by the functions of the
//     SuiteSparse:GraphBLAS implementation in C and stored in C memory. Therefore, this type
//     must adhere to the restrictions of [cgo]. Specifically, D must not contain any Go pointers.
//   - size: The size of a D value when stored in C memory.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [InvalidObject], [OutOfMemory], [Panic]
//
// [cgo]: https://pkg.go.dev/cmd/cgo#hdr-Passing_pointers
func TypeNew[D any](size int) (typ Type, err error) {
	var d D
	t := reflect.TypeOf(d)
	if hasGoPointer(t) {
		err = makeError(InvalidValue)
		return
	}
	var grb C.GrB_Type
	info := Info(C.GrB_Type_new(&grb, C.size_t(size)))
	if info == success {
		typ = Type{t}
		grbType[typ] = grb
		goType[grb] = typ
		return
	}
	err = makeError(info)
	return
}

// MaxNameLen is the maximum length for strings naming types.
//
// In C, a null byte has to be added to the end of a string,
// so this constant is one smaller than the corresponding GxB_MAX_NAME_LEN
// in SuiteSparse:GraphBLAS.
//
// MaxNameLen is a SuiteSparse:GraphBLAS extension.
const MaxNameLen = 127

// NamedTypeNew creates a type with a name and definition that are known to GraphBLAS, as strings.
//
// The typename is any valid string (max length of 127 characters) that may appear as the name of
// a C type created by a C typedef statement. It must not contain any whitespace characters. For example,
// to create a type with a 4-by-4 dense float array and a 32-bit integer:
//
//	type myquaternion struct {
//		x [4*4]float32
//		color int32
//	}
//
//	typ, err := GrB.NamedTypeNew[myquaternion](0, "myquaternion", `typedef struct {
//		float x [4][4] ;
//		int color ;
//	} myquaternion ;`)
//
// The two strings are optional, but are required to enable the JIT compilation of kernels that use this type.
//
// If the size parameter is zero, and the strings are valid, a JIT kernel is compiled just to determine the size
// of the type. The Go type D has to be compatible with the corresponding C type.
//
// NamedTypeNew is a SuiteSparse:GraphBLAS extension.
func NamedTypeNew[D any](size int, typename string, typedefn string) (typ Type, err error) {
	var d D
	t := reflect.TypeOf(d)
	if hasGoPointer(t) {
		err = makeError(InvalidValue)
		return
	}
	ctypename := C.CString(typename)
	defer C.free(unsafe.Pointer(ctypename))
	ctypedefn := C.CString(typedefn)
	defer C.free(unsafe.Pointer(ctypedefn))
	var grb C.GrB_Type
	info := Info(C.GxB_Type_new(&grb, C.size_t(size), ctypename, ctypedefn))
	if info == success {
		typ = Type{t}
		grbType[typ] = grb
		goType[grb] = typ
		return
	}
	err = makeError(info)
	return
}

// Size returns the size of a type.
//
// This functions acts just like sizeof(type) in the C language. For example,
// GrB.Int32.Size() returns 4, the same as sizeof(int32_t).
//
// Size is a SuiteSparse:GraphBLAS extension.
func (typ Type) Size() (int, error) {
	var csize C.size_t
	info := Info(C.GxB_Type_size(&csize, grbType[typ]))
	if info == success {
		return int(csize), nil
	}
	return 0, makeError(info)
}

// Valid returns true if typ has been created by a successful call to [TypeNew] or [NamedTypeNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (typ Type) Valid() bool {
	return grbType[typ] != C.GrB_Type(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Type] and releases any resources associated with
// it. Calling Free on an object that is not [Type.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined type is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (typ *Type) Free() error {
	grb := grbType[*typ]
	grbCopy := grb
	info := Info(C.GrB_Type_free(&grb))
	if info == success {
		delete(grbType, goType[grbCopy])
		delete(goType, grbCopy)
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the type into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (typ Type) Wait(mode WaitMode) error {
	info := Info(C.GrB_Type_wait(grbType[typ], C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the type.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (typ Type) Err() (errorString string, err error) {
	var cerror *C.char
	info := Info(C.GrB_Type_error(&cerror, grbType[typ]))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the type to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: typ is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (typ Type) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Type_fprint(grbType[typ], cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}
