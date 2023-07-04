package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// A Descriptor is used to modify the behavior of a GraphBLAS method. When present in the
// signature of a method, they appear as the last argument in the method. Descriptors specify how
// the other input arguments corresponding to GraphBLAS collections – vectors, matrices, and masks
// – should be processed (modified) before the main operation of a method is performed.
//
// If a default descriptor is desired, nil should be specified. Non-default field/value pairs are
// as follows: If the [Replace] descriptor is set for [Outp], then the output vector or matrix is
// cleared (all elements removed) before the result is stored in it. If the [Structure] descriptor
// is set for [Mask], then the write mask is constructed from the structure (pattern of stored
// values) of the input mask vector or matrix. The stored values are not examined. If the [Comp]
// descriptor is set for [Mask], then the complement of mask is used. If the [Tran] descriptor is
// set for [Inp0], then the transpose of the first non-mask input matrix is used for the operation.
// If the [Tran] descriptor is set for [Inp1], then the transpose of the second non-mask input matrix
// is used for the operation.
//
// The descriptor is a lightweight object. It is composed of (field, value) pairs where the field selects
// one of the GraphBLAS objects from the argument list of a method and the value defines the
// indicated modification associated with that object.
//
// For the purpose of constructing descriptors, the arguments of a method that can be modified
// are identified by specific field names. The output parameter (typically the first parameter in a
// GraphBLAS method) is indicated by the field name, [Outp]. The mask is indicated by the
// [Mask] field name. The input parameters corresponding to the input vectors and matrices are
// indicated by [Inp0] and [Inp1] in the order they appear in the signature of the GraphBLAS
// method. The descriptor is an opaque object and hence we do not define how objects of this type
// should be implemented.
//
// In the definitions of the GraphBLAS methods, we often refer to the default behavior of a method
// with respect to the action of a descriptor. If a descriptor is not provided or if the value associated
// with a particular field in a descriptor is not set, the default behavior of a GraphBLAS method is
// defined as follows:
//
//   - Input matrices are not transposed.
//   - The mask is used, as is, without complementing, and stored values are examined to determine
//     whether they evaluate to true or false.
//   - Values of the output object that are not directly modified by the operation are preserved.
//
// GraphBLAS specifies all of the valid combinations of (field, value) pairs as predefined descriptors.
type Descriptor struct {
	grb C.GrB_Descriptor
}

type (
	// DescField is the descriptor field enumeration.
	DescField int

	// DescValue is the descriptor value enumeration.
	DescValue int
)

// [Descriptor] field names.
const (
	Outp DescField = iota
	Mask
	Inp0
	Inp1

	// AxBMethod is a descriptor for selecting the C=A*B algorithm.
	// AxBMethod is a SuiteSparse:GraphBLAS extension.
	AxBMethod DescField = 1000

	// SortHint controls sorting in [MxM], [MxV], [VxM], and reduction functions.
	// SortHint is a SuiteSparse:GraphBLAS extension.
	SortHint DescField = 35

	// Compression selects the compression for [Matrix.Serialize].
	// Compression is a SuiteSparse:GraphBLAS extension.
	Compression DescField = 36

	// Import selects between secure and fast packing.
	// Import is a SuiteSparse:GraphBLAS extension.
	Import DescField = 37 // a SuiteSparse:GraphBLAS extension
)

func (field DescField) String() string {
	switch field {
	case Outp:
		return "output"
	case Mask:
		return "mask"
	case Inp0:
		return "first input"
	case Inp1:
		return "second input"
	case AxBMethod:
		return "multiply algorithm"
	case SortHint:
		return "sort control"
	case Compression:
		return "compression"
	case Import:
		return "import security"
	}
	panic("invalid descriptor field")
}

// Descriptor field values.
const (
	Default DescValue = iota // a SuiteSparse:GraphBLAS extension
	Replace
	Comp
	Tran
	Structure

	// AxBGustavson selects an extended version of Gustavson's method for [AxBMethod].
	// AxBGustavson is a SuiteSparse:GraphBLAS extension.
	AxBGustavson DescValue = 1001

	// AxBDot selects a very specialized method for [AxBMethod] that works well only if the
	// mask is present, very sparse, and not complemented, when the output matrix is very
	// small, or when the output matrix is bitmap or full.
	// AxBDot is a SuiteSparse:GraphBLAS extension.
	AxBDot DescValue = 1003

	// AxBHash selects a hash-based method for [AxBMethod].
	// AxBHash is a SuiteSparse:GraphBLAS extension.
	AxBHash DescValue = 1004

	// AxBSaxpy selects a saxpy-based method for [AxBMethod].
	// AxBSaxpy is a SuiteSparse:GraphBLAS extension.
	AxBSaxpy DescValue = 1005

	// SecureImport informs the pack functions that the data is being packed
	// from an untrusted source, so additional checks will be made.
	// SecureImport is a SuiteSparse:GraphBLAS extension.
	SecureImport DescValue = 502

	// PreferSorted provides a hint to [MxM], [MxV], [VxM], and reduction functions
	// to sort the output result. (This can be any value other than 0.)
	// PreferSorted is a SuiteSparse:GraphBLAS extension.
	PreferSorted DescValue = 999

	// CompressionNone selects no compression for [Matrix.Serialize].
	// CompressionNone is a SuiteSparse:GraphBLAS extension.
	CompressionNone DescValue = -1

	// CompressionLZ4 selects LZ4 compression for [Matrix.Serialize].
	// CompressionLZ4 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4 DescValue = 1000

	// CompressionLZ4HC selects LZ4HC with default level 9 for [Matrix.Serialize].
	// CompressionLZ4HC is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC DescValue = 2000

	// CompressionLZ4HC1 selects LZ4HC with level 1 for [Matrix.Serialize].
	// CompressionLZ4HC1 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC1 DescValue = 2001

	// CompressionLZ4HC2 selects LZ4HC with level 2 for [Matrix.Serialize].
	// CompressionLZ4HC2 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC2 DescValue = 2002

	// CompressionLZ4HC3 selects LZ4HC with level 3 for [Matrix.Serialize].
	// CompressionLZ4HC3 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC3 DescValue = 2003

	// CompressionLZ4HC4 selects LZ4HC with level 4 for [Matrix.Serialize].
	// CompressionLZ4HC4 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC4 DescValue = 2004

	// CompressionLZ4HC5 selects LZ4HC with level 5 for [Matrix.Serialize].
	// CompressionLZ4HC5 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC5 DescValue = 2005

	// CompressionLZ4HC6 selects LZ4HC with level 6 for [Matrix.Serialize].
	// CompressionLZ4HC6 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC6 DescValue = 2006

	// CompressionLZ4HC7 selects LZ4HC with level 7 for [Matrix.Serialize].
	// CompressionLZ4HC7 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC7 DescValue = 2007

	// CompressionLZ4HC8 selects LZ4HC with level 8 for [Matrix.Serialize].
	// CompressionLZ4HC8 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC8 DescValue = 2008

	// CompressionLZ4HC9 selects LZ4HC with level 9 for [Matrix.Serialize].
	// CompressionLZ4HC9 is a SuiteSparse:GraphBLAS extension.
	CompressionLZ4HC9 DescValue = 2009

	// CompressionZSTD selects ZSTD with default level 1 for [Matrix.Serialize].
	// CompressionZSTD is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD DescValue = 3000

	// CompressionZSTD1 selects ZSTD with level 1 for [Matrix.Serialize].
	// CompressionZSTD1 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD1 DescValue = 3001

	// CompressionZSTD2 selects ZSTD with level 2 for [Matrix.Serialize].
	// CompressionZSTD2 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD2 DescValue = 3002

	// CompressionZSTD3 selects ZSTD with level 3 for [Matrix.Serialize].
	// CompressionZSTD3 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD3 DescValue = 3003

	// CompressionZSTD4 selects ZSTD with level 4 for [Matrix.Serialize].
	// CompressionZSTD4 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD4 DescValue = 3004

	// CompressionZSTD5 selects ZSTD with level 5 for [Matrix.Serialize].
	// CompressionZSTD5 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD5 DescValue = 3005

	// CompressionZSTD6 selects ZSTD with level 6 for [Matrix.Serialize].
	// CompressionZSTD6 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD6 DescValue = 3006

	// CompressionZSTD7 selects ZSTD with level 7 for [Matrix.Serialize].
	// CompressionZSTD7 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD7 DescValue = 3007

	// CompressionZSTD8 selects ZSTD with level 8 for [Matrix.Serialize].
	// CompressionZSTD8 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD8 DescValue = 3008

	// CompressionZSTD9 selects ZSTD with level 9 for [Matrix.Serialize].
	// CompressionZSTD9 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD9 DescValue = 3009

	// CompressionZSTD10 selects ZSTD with level 10 for [Matrix.Serialize].
	// CompressionZSTD10 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD10 DescValue = 3010

	// CompressionZSTD11 selects ZSTD with level 11 for [Matrix.Serialize].
	// CompressionZSTD11 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD11 DescValue = 3011

	// CompressionZSTD12 selects ZSTD with level 12 for [Matrix.Serialize].
	// CompressionZSTD12 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD12 DescValue = 3012

	// CompressionZSTD13 selects ZSTD with level 13 for [Matrix.Serialize].
	// CompressionZSTD13 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD13 DescValue = 3013

	// CompressionZSTD14 selects ZSTD with level 14 for [Matrix.Serialize].
	// CompressionZSTD14 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD14 DescValue = 3014

	// CompressionZSTD15 selects ZSTD with level 15 for [Matrix.Serialize].
	// CompressionZSTD15 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD15 DescValue = 3015

	// CompressionZSTD16 selects ZSTD with level 16 for [Matrix.Serialize].
	// CompressionZSTD16 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD16 DescValue = 3016

	// CompressionZSTD17 selects ZSTD with level 17 for [Matrix.Serialize].
	// CompressionZSTD17 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD17 DescValue = 3017

	// CompressionZSTD18 selects ZSTD with level 18 for [Matrix.Serialize].
	// CompressionZSTD18 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD18 DescValue = 3018

	// CompressionZSTD19 selects ZSTD with level 19 for [Matrix.Serialize].
	// CompressionZSTD19 is a SuiteSparse:GraphBLAS extension.
	CompressionZSTD19 DescValue = 3019
)

func (value DescValue) String() string {
	switch value {
	case Default:
		return "default"
	case Replace:
		return "replace"
	case Comp:
		return "mask complement"
	case Tran:
		return "input transpose"
	case Structure:
		return "mask structure"
	case AxBGustavson:
		return "gather-scatter saxpy method"
	case AxBDot:
		return "dot product"
	case AxBHash:
		return "hash-based saxpy method"
	case AxBSaxpy:
		return "any saxpy method"
	case CompressionNone:
		return "no compression"
	case CompressionLZ4:
		return "LZ4"
	case CompressionLZ4HC:
		return "LZ4HC default level (9)"
	case CompressionLZ4HC1:
		return "LZ4HC level 1"
	case CompressionLZ4HC2:
		return "LZ4HC level 2"
	case CompressionLZ4HC3:
		return "LZ4HC level 3"
	case CompressionLZ4HC4:
		return "LZ4HC level 4"
	case CompressionLZ4HC5:
		return "LZ4HC level 5"
	case CompressionLZ4HC6:
		return "LZ4HC level 6"
	case CompressionLZ4HC7:
		return "LZ4HC level 7"
	case CompressionLZ4HC8:
		return "LZ4HC level 8"
	case CompressionLZ4HC9:
		return "LZ4HC level 9"
	case CompressionZSTD:
		return "ZSTD default level (1)"
	case CompressionZSTD1:
		return "ZSTD level 1"
	case CompressionZSTD2:
		return "ZSTD level 2"
	case CompressionZSTD3:
		return "ZSTD level 3"
	case CompressionZSTD4:
		return "ZSTD level 4"
	case CompressionZSTD5:
		return "ZSTD level 5"
	case CompressionZSTD6:
		return "ZSTD level 6"
	case CompressionZSTD7:
		return "ZSTD level 7"
	case CompressionZSTD8:
		return "ZSTD level 8"
	case CompressionZSTD9:
		return "ZSTD level 9"
	case CompressionZSTD10:
		return "ZSTD level 10"
	case CompressionZSTD11:
		return "ZSTD level 11"
	case CompressionZSTD12:
		return "ZSTD level 12"
	case CompressionZSTD13:
		return "ZSTD level 13"
	case CompressionZSTD14:
		return "ZSTD level 14"
	case CompressionZSTD15:
		return "ZSTD level 15"
	case CompressionZSTD16:
		return "ZSTD level 16"
	case CompressionZSTD17:
		return "ZSTD level 17"
	case CompressionZSTD18:
		return "ZSTD level 18"
	case CompressionZSTD19:
		return "ZSTD level 19"
	case SecureImport:
		return "secure import"
	case PreferSorted:
		return "prefer sorted result"
	}
	panic("invalid descriptor value")
}

// DescriptorNew creates a new (empty or default) descriptor.
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func DescriptorNew() (descriptor Descriptor, err error) {
	info := Info(C.GrB_Descriptor_new(&descriptor.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Set the content for a field for an existing descriptor.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [OutOfMemory], [Panic]
func (descriptor Descriptor) Set(field DescField, value DescValue) error {
	info := Info(C.GrB_Descriptor_set(descriptor.grb, C.GrB_Desc_Field(field), C.GrB_Desc_Value(value)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Get the content for a field of an existing descriptor.
//
// Get is a SuiteSparse:GraphBLAS extension.
func (descriptor Descriptor) Get(field DescField) (DescValue, error) {
	var cvalue C.GrB_Desc_Value
	info := Info(C.GxB_Descriptor_get(&cvalue, descriptor.grb, C.GrB_Desc_Field(field)))
	if info == success {
		return DescValue(cvalue), nil
	}
	return 0, makeError(info)
}

// Valid returns true if descriptor has been created by a successful call to [DescriptorNew].
//
// Valid is a forGraphBLASGo extension. It is used in place of comparing against GrB_INVALID_HANDLE.
func (descriptor Descriptor) Valid() bool {
	return descriptor.grb != C.GrB_Descriptor(C.GrB_INVALID_HANDLE)
}

// Free destroys a previously created [Descriptor] and releases any resources associated with
// it. Calling Free on an object that is not [Descriptor.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined descriptor is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (descriptor *Descriptor) Free() error {
	info := Info(C.GrB_Descriptor_free(&descriptor.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Wait until function calls in a sequence put the descriptor into a state of completion or
// materialization.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [IndexOutOfBounds], [OutOfMemory], [Panic]
func (descriptor Descriptor) Wait(mode WaitMode) error {
	info := Info(C.GrB_Descriptor_wait(descriptor.grb, C.GrB_WaitMode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Err returns an error message about any errors encountered during the processing associated with
// the descriptor.
//
// GraphBLAS API errors that may be returned:
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (descriptor Descriptor) Err() (string, error) {
	var cerror *C.char
	info := Info(C.GrB_Descriptor_error(&cerror, descriptor.grb))
	if info == success {
		return C.GoString(cerror), nil
	}
	return "", makeError(info)
}

// Print the contents of the descriptor to stdout.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]: The underlying print routine returned an I/O error.
//   - [NullPointer]: descriptor is a nil pointer.
//   - [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject], [Panic]
//
// Print is a SuiteSparse:GraphBLAS extension.
func (descriptor Descriptor) Print(name string, pr PrintLevel) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Descriptor_fprint(descriptor.grb, cname, C.GxB_Print_Level(pr), (*C.FILE)(C.NULL)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// Predefined GraphBLAS descriptors. The list includes all possible descriptors, according to the current
// GraphBLAS standard (without SuiteSparse:GraphBLAS-specific extensions).
//
// nil is also a valid descriptor value, indicating default behavior.
var (
	DescT1      = &Descriptor{C.GrB_DESC_T1}
	DescT0      = &Descriptor{C.GrB_DESC_T0}
	DescT0T1    = &Descriptor{C.GrB_DESC_T0T1}
	DescC       = &Descriptor{C.GrB_DESC_C}
	DescS       = &Descriptor{C.GrB_DESC_S}
	DescCT1     = &Descriptor{C.GrB_DESC_CT1}
	DescST1     = &Descriptor{C.GrB_DESC_ST1}
	DescCT0     = &Descriptor{C.GrB_DESC_CT0}
	DescST0     = &Descriptor{C.GrB_DESC_ST0}
	DescCT0T1   = &Descriptor{C.GrB_DESC_CT0T1}
	DescST0T1   = &Descriptor{C.GrB_DESC_ST0T1}
	DescSC      = &Descriptor{C.GrB_DESC_SC}
	DescSCT1    = &Descriptor{C.GrB_DESC_SCT1}
	DescSCT0    = &Descriptor{C.GrB_DESC_SCT0}
	DescSCT0T1  = &Descriptor{C.GrB_DESC_SCT0T1}
	DescR       = &Descriptor{C.GrB_DESC_R}
	DescRT1     = &Descriptor{C.GrB_DESC_RT1}
	DescRT0     = &Descriptor{C.GrB_DESC_RT0}
	DescRT0T1   = &Descriptor{C.GrB_DESC_RT0T1}
	DescRC      = &Descriptor{C.GrB_DESC_RC}
	DescRS      = &Descriptor{C.GrB_DESC_RS}
	DescRCT1    = &Descriptor{C.GrB_DESC_RCT1}
	DescRST1    = &Descriptor{C.GrB_DESC_RST1}
	DescRCT0    = &Descriptor{C.GrB_DESC_RCT0}
	DescRST0    = &Descriptor{C.GrB_DESC_RST0}
	DescRCT0T1  = &Descriptor{C.GrB_DESC_RCT0T1}
	DescRST0T1  = &Descriptor{C.GrB_DESC_RST0T1}
	DescRSC     = &Descriptor{C.GrB_DESC_RSC}
	DescRSCT1   = &Descriptor{C.GrB_DESC_RSCT1}
	DescRSCT0   = &Descriptor{C.GrB_DESC_RSCT0}
	DescRSCT0T1 = &Descriptor{C.GrB_DESC_RSCT0T1}
)
