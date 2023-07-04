package GrB

/*
#cgo LDFLAGS: -lgraphblas
#include <assert.h>
#include "GraphBLAS.h"

typedef void * (* user_malloc_function  ) (size_t);
typedef void * (* user_calloc_function  ) (size_t, size_t);
typedef void * (* user_realloc_function ) (void *, size_t);
typedef void   (* user_free_function    ) (void *);

void* user_malloc(user_malloc_function umalloc, size_t size) {
	if (umalloc != 0) {
		return umalloc(size);
	}
	return malloc(size);
}

void* user_calloc(user_calloc_function ucalloc, size_t num, size_t size) {
	if (ucalloc != 0) {
		return ucalloc(num, size);
	}
	return calloc(num, size);
}

void user_free(user_free_function ufree, void* ptr) {
	if (ufree != 0) {
		ufree(ptr);
	} else {
		free(ptr);
	}
}

void assertions() {
	static_assert(GRB_VERSION==2 && GRB_SUBVERSION==0, "This version of GraphBLAS is not supported.");
	static_assert(sizeof(bool) == 1, "The size of the C compiler's bool type is not 1.");
}
*/
import "C"
import "unsafe"

func init() {
	if unsafe.Sizeof(false) != 1 {
		panic("The size of Go's bool type is not 1.")
	}
}

// VersionToInt returns a single integer for comparing spec and version levels.
//
// VersionToInt is a SuiteSparse:GraphBLAS extension.
func VersionToInt(major, minor, sub int) int {
	return (major*1000+minor)*1000 + sub
}

// SuiteSparse:GraphBLAS extensions
const (
	SuiteSparseImplementationName    = C.GxB_IMPLEMENTATION_NAME
	SuiteSparseImplementationDate    = C.GxB_IMPLEMENTATION_DATE
	SuiteSparseImplementationMajor   = C.GxB_IMPLEMENTATION_MAJOR
	SuiteSparseImplementationMinor   = C.GxB_IMPLEMENTATION_MINOR
	SuiteSparseImplementationSub     = C.GxB_IMPLEMENTATION_SUB
	SuiteSparseImplementation        = C.GxB_IMPLEMENTATION
	SuiteSparseImplementationAbout   = C.GxB_IMPLEMENTATION_ABOUT
	SuiteSparseImplementationLicense = C.GxB_IMPLEMENTATION_LICENSE
)

// forGraphBLASGo extensions
const (
	ImplementationName    = "Intel® Generic Implementation of GraphBLAS* for Go*"
	ImplementationDate    = "March 1, 2023"
	ImplementationMajor   = 0
	ImplementationMinor   = 0
	ImplementationSub     = 0
	Implementation        = (ImplementationMajor*1000+ImplementationMinor)*1000 + ImplementationSub
	ImplementationAbout   = "Intel® Generic Implementation of GraphBLAS* for Go*. Copyright © 2022-2023, Intel Corporation. All rights reserved."
	ImplementationLicense = "BSD 3-Clause License"
)

// SuiteSparse:GraphBLAS extensions
const (
	SpecDate    = C.GxB_SPEC_DATE
	SpecMajor   = C.GxB_SPEC_MAJOR
	SpecMinor   = C.GxB_SPEC_MINOR
	SpecSub     = C.GxB_SPEC_SUB
	SpecVersion = (SpecMajor*1000+SpecMinor)*1000 + SpecSub
	SpecAbout   = C.GxB_SPEC_ABOUT
)

// Spec version
const (
	Version    = C.GRB_VERSION
	Subversion = C.GRB_SUBVERSION
)

func init() {
	if VersionToInt(
		SuiteSparseImplementationMajor,
		SuiteSparseImplementationMinor,
		SuiteSparseImplementationSub,
	) != SuiteSparseImplementation {
		panic("Invalid SuiteSparse version number. This can normally not happen.")
	}
	if VersionToInt(
		ImplementationMajor,
		ImplementationMinor,
		ImplementationSub,
	) != Implementation {
		panic("Invalid forGraphBLASGo version number. This can normally not happen.")
	}
	if VersionToInt(
		SpecMajor,
		SpecMinor,
		SpecSub,
	) != SpecVersion {
		panic("Invalid spec version number. This can normally not happen.")
	}
}

// GetVersion is used to query the major and minor version number of the
// GraphBLAS API specification that the library implements at runtime. The following two
// constants are also defined by the library:
//   - [Version]: 2
//   - [Subversion]: 0
//
// Return Values:
//   - version: Major version number.
//   - subversion: subversion number.
func GetVersion() (version, subversion int) {
	return Version, Subversion
}

type (
	// A Format specifies the external format for [MatrixImport] and [Matrix.Export].
	Format int
	// A Mode specifies the execution mode for the [Init] or [InitWithMalloc] functions.
	Mode int
	// A WaitMode specifies the wait mode for the Wait methods.
	WaitMode int
)

// External formats for [MatrixImport] and [Matrix.Export].
const (
	CsrFormat Format = iota // compressed sparse row matrix format
	CscFormat               // compressed sparse column matrix format
	CooFormat               // sparse coordinate matrix format
)

func (format Format) String() string {
	switch format {
	case CsrFormat:
		return "CSR"
	case CscFormat:
		return "CSC"
	case CooFormat:
		return "COO"
	}
	panic("invalid format")
}

// Execution modes for the [Init] and [InitWithMalloc] functions.
const (
	NonBlocking Mode = iota
	Blocking
)

func (mode Mode) String() string {
	switch mode {
	case NonBlocking:
		return "non-blocking"
	case Blocking:
		return "blocking"
	}
	panic("invalid mode")
}

// Wait modes for the Wait methods.
const (
	Complete WaitMode = iota
	Materialize
)

func (waitMode WaitMode) String() string {
	switch waitMode {
	case Complete:
		return "complete"
	case Materialize:
		return "materialize"
	}
	panic("invalid wait mode")
}

// C types for user-defined allocation functions (replacements for malloc).
type (
	UserMallocFunction  C.user_malloc_function
	UserCallocFunction  C.user_calloc_function
	UserReallocFunction C.user_realloc_function
	UserFreeFunction    C.user_free_function
)

var (
	user_malloc UserMallocFunction = nil
	user_calloc UserCallocFunction = nil
	user_free   UserFreeFunction   = nil
)

/*
func malloc(size int) unsafe.Pointer {
	return C.user_malloc(user_malloc, C.size_t(size))
}
*/

func calloc(size int) unsafe.Pointer {
	if user_calloc == nil {
		return C.memset(C.user_malloc(user_malloc, C.size_t(size)), 0, C.size_t(size))
	}
	if user_malloc == nil {
		panic("inconsistent malloc / calloc")
	}
	return C.user_calloc(user_calloc, 1, C.size_t(size))
}

func free(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}
	C.user_free(user_free, ptr)
}

// Init creates and initializes a GraphBLAS API context. The argument to
// Init defines the mode for the context. The two available modes are:
//
//   - [Blocking]: In this mode, each method in a sequence returns after
//     its computations have completed and output arguments are available to subsequent
//     statements in an application. When executing in Blocking mode, the methods
//     execute in program order.
//
//   - [NonBlocking]: In this mode, methods in a sequence may return after arguments in
//     the method have been tested for dimension and domain compatibility within the method
//     but potentially before their computations complete. Output arguments are available to
//     subsequent GraphBLAS methods in an application. When executing in NonBlocking mode, the
//     methods in a sequence may execute in any order that preserves the mathematical result
//     defined by the sequence.
//
// An application can only create one context per execution instance. An application may only
// call Init once. Calling Init more than once results in undefined behavior.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue]
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func Init(mode Mode) error {
	info := Info(C.GrB_init(C.GrB_Mode(mode)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// InitWithMalloc is identical to [Init], except that it also redefines the memory management
// functions that SuiteSparse:GraphBLAS will use. calloc and realloc are optional, and may be
// nil. The functions passed to InitWithMalloc must be thread-safe.
//
// InitWithMalloc is a SuiteSparse:GraphBLAS extension.
func InitWithMalloc(
	mode Mode,
	malloc UserMallocFunction,
	calloc UserCallocFunction,
	realloc UserReallocFunction,
	free UserFreeFunction,
) error {
	info := Info(C.GxB_init(C.GrB_Mode(mode), malloc, calloc, realloc, free))
	if info == success {
		user_malloc = malloc
		user_calloc = calloc
		user_free = free
		return nil
	}
	return makeError(info)
}

// Finalize terminates and frees any internal resources created to support the
// GraphBLAS API context. Finalize may only be called after a context has been initialized
// by calling [Init] or [InitWithMalloc], or else undefined behavior occurs. After
// Finalize has been called to finalize a GraphBLAS context, calls to any GraphBLAS functions,
// including Finalize, will result in undefined behavior.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func Finalize() error {
	info := Info(C.GrB_finalize())
	if info == success {
		return nil
	}
	return makeError(info)
}
