package GrB

// #include "GraphBLAS.h"
import "C"
import "unsafe"

// Layout indicates whether a matrix is laid out by row or by column.
// This corresponds to GxB_FORMAT in SuiteSparse:GraphBLAS.
//
// Layout is a SuiteSparse:GraphBLAS extension.
type Layout int

// SuiteSparse:GraphBLAS extensions
const (
	ByRow Layout = iota
	ByCol
)

func (format Layout) String() string {
	switch format {
	case ByRow:
		return "by row"
	case ByCol:
		return "by col"
	}
	panic("invalid matrix format")
}

// SuiteSparse:GraphBLAS extensions
var (
	LayoutDefault = Layout(C.GxB_FORMAT_DEFAULT) // The default is by row.
	HyperDefault  = float64(C.GxB_HYPER_DEFAULT)
)

// Sparsity indicates the sparsity representation of a matrix.
//
// Sparsity is a SuiteSparse:GraphBLAS extension.
type Sparsity int

// SuiteSparse:GraphBLAS extensions
const (
	Hypersparse Sparsity = 1 << iota
	Sparse
	Bitmap
	Full
	AutoSparsity = Hypersparse + Sparse + Bitmap + Full
)

func (sparsity Sparsity) String() string {
	switch sparsity {
	case Hypersparse:
		return "hypersparse"
	case Sparse:
		return "sparse"
	case Bitmap:
		return "bitmap"
	case Full:
		return "full"
	case AutoSparsity:
		return "auto"
	}
	panic("invalid matrix sparsity")
}

// NBitmapSwitch is the number of float64 values for the
// global bitmap switch setting. See [GlobalSetBitmapSwitch]
// and [GlobalGetBitmapSwitch].
//
// NBitmapSwitch is a SuiteSparse:GraphBLAS extension.
const NBitmapSwitch = 8

// Values for hyper switch settings.
//
// SuiteSparse:GraphBLAS extensions
var (
	AlwaysHyper = float64(C.GxB_ALWAYS_HYPER)
	NeverHyper  = float64(C.GxB_NEVER_HYPER)
)

// JITControl can be used to control the just-in-time compiler
// of SuiteSparse:GraphBLAS.
//
// JITControl is a SuiteSparse:GraphBLAS extension.
type JITControl int

// SuiteSparse:GraphBLAS extensions
const (
	// JITOff means: do not use the JIT and free all JIT kernels if loaded
	JITOff JITControl = iota

	// JITPause means: do not run JIT kernels but keep any loaded
	JITPause

	// JITRun means: run JIT kernels if already loaded; no load/compile
	JITRun

	// JITLoad means: able to load and run JIT kernels; may not compile
	JITLoad

	// JITOn means: full JIT: able to compile, load, and run
	JITOn
)

func (control JITControl) String() string {
	switch control {
	case JITOff:
		return "JIT off"
	case JITPause:
		return "JIT pause"
	case JITRun:
		return "JIT run"
	case JITLoad:
		return "JIT load"
	case JITOn:
		return "JIT on"
	}
	panic("invalid JIT control")
}

// GlobalSetHyperSwitch determines how future matrices are converted between the hypersparse and
// non-hypersparse formats by default.
//
// Parameters:
//
//   - hyperSwitch (IN): A value between 0 and 1. To force matrices to always be non-hypersparse,
//     use [NeverHyper]. To force a matrix to always stay hypersparse, use [AlwaysHyper].
//
// GlobalSetHyperSwitch is a SuiteSparse:GraphBLAS extension.
func GlobalSetHyperSwitch(hyperSwitch float64) error {
	info := Info(C.GxB_Global_Option_set_FP64(C.GxB_HYPER_SWITCH, C.double(hyperSwitch)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetHyperSwitch retrieves the current global switch to hypersparse. See [GlobalSetHyperSwitch].
//
// GlobalGetHyperSwitch is a SuiteSparse:GraphBLAS extension.
func GlobalGetHyperSwitch() (float64, error) {
	var cHyperSwitch C.double
	info := Info(C.GxB_Global_Option_get_FP64(C.GxB_HYPER_SWITCH, &cHyperSwitch))
	if info == success {
		return float64(cHyperSwitch), nil
	}
	return 0, makeError(info)
}

// GlobalSetBitmapSwitch determines how future matrices are converted to the bitmap format by default.
//
// Parameters:
//
//   - bitmapSwitch (IN): A value between 0 and 1.
//
// GlobalSetBitmapSwitch is a SuiteSparse:GraphBLAS extension.
func GlobalSetBitmapSwitch(bitmapSwitch [NBitmapSwitch]float64) error {
	info := Info(C.GxB_Global_Option_set_FP64_ARRAY(C.GxB_BITMAP_SWITCH, (*C.double)(&bitmapSwitch[0])))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetBitmapSwitch retrieves the current global switch to bitmap. See [GlobalSetBitmapSwitch].
//
// GlobalGetBitmapSwitch is a SuiteSparse:GraphBLAS extension.
func GlobalGetBitmapSwitch() (bitmapSwitch [NBitmapSwitch]float64, err error) {
	info := Info(C.GxB_Global_Option_get_FP64(C.GxB_BITMAP_SWITCH, (*C.double)(&bitmapSwitch[0])))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// GlobalSetLayout sets the [Layout] (GxB_FORMAT) of future matrices.
//
// GlobalSetLayout is a SuiteSparse:GraphBLAS extension.
func GlobalSetLayout(format Layout) error {
	info := Info(C.GxB_Global_Option_set_INT32(C.GxB_FORMAT, C.int32_t(format)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetLayout retrieves the current global [Layout] (GxB_FORMAT) of matrices.
//
// GlobalGetLayout a SuiteSparse:GraphBLAS extension.
func GlobalGetLayout() (Layout, error) {
	var format C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_FORMAT, &format))
	if info == success {
		return Layout(format), nil
	}
	return 0, makeError(info)
}

// GlobalSetNThreads controls how many OpenMP threads GraphBLAS operations use
// by default. By default, if set to 0, all available threads are used. If the
// value is more than the available OpenMP threads, then all available threads
// are used.
//
// GlobalSetNThreads is a SuiteSparse:GraphBLAS extension.
func GlobalSetNThreads(nthreads int32) error {
	info := Info(C.GxB_Global_Option_set_INT32(C.GxB_NTHREADS, C.int32_t(nthreads)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetNThreads retrieves the current number of OpenMP threads GraphBLAS
// operations use by default.
//
// GlobalGetNThreads is a SuiteSparse:GraphBLAS extension.
func GlobalGetNThreads() (int32, error) {
	var nthreads C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_NTHREADS, &nthreads))
	if info == success {
		return int32(nthreads), nil
	}
	return 0, makeError(info)
}

// GlobalSetChunk sets a value that controls how many OpenMP threads GraphBLAS
// operations use for small problems. When chunk < 1, a default value
// is used. The default value is usually 65536.
//
// GlobalSetChunk is a SuiteSparse:GraphBLAS extension.
func GlobalSetChunk(chunk float64) error {
	info := Info(C.GxB_Global_Option_set_FP64(C.GxB_CHUNK, C.double(chunk)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetChunk retrieves the current global chunk value. See [GlobalSetChunk].
//
// GlobalGetChunk is a SuiteSparse:GraphBLAS extension.
func GlobalGetChunk() (float64, error) {
	var chunk C.double
	info := Info(C.GxB_Global_Option_get_FP64(C.GxB_CHUNK, &chunk))
	if info == success {
		return float64(chunk), nil
	}
	return 0, makeError(info)
}

// GlobalSetBurble enables or disables the burble.
//
// If enabled, SuiteSparse:GraphBLAS reports which internal kernels it uses,
// and how much times is spent.
//
// GlobalSetBurble is a SuiteSparse:GraphBLAS extension.
func GlobalSetBurble(enableNotDisable bool) error {
	info := Info(C.GxB_Global_Option_set_INT32(C.GxB_BURBLE, gotocbool(enableNotDisable)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetBurble retrieves whether the burble is currently enabled disabled. See [GlobalSetBurble].
//
// GlobalGetBurble is a SuiteSparse:GraphBLAS extension.
func GlobalGetBurble() (enabledNotDisabled bool, err error) {
	var cburble C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_BURBLE, &cburble))
	if info == success {
		return ctogobool(cburble), nil
	}
	return false, makeError(info)
}

// GlobalSetPrint1Based sets whether vector and matrix indices start with 1
// or 0 when printed.
//
// GlobalSetPrint1Based is a SuiteSparse:GraphBLAS extension.
func GlobalSetPrint1Based(print1Based bool) error {
	info := Info(C.GxB_Global_Option_set_INT32(C.GxB_PRINT_1BASED, gotocbool(print1Based)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetPrint1Based retrieves whether vector and matrix indices start with 1
// or 0 when printed.
//
// GlobalGetPrint1Based is a SuiteSparse:GraphBLAS extension.
func GlobalGetPrint1Based() (print1Based bool, err error) {
	var cprint1Based C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_PRINT_1BASED, &cprint1Based))
	if info == success {
		return ctogobool(cprint1Based), nil
	}
	return false, makeError(info)
}

// GlobalGetMode retrieves whether GraphBLAS currently operates in
// blocking or non-blocking mode.
//
// GlobalGetMode is a SuiteSparse:GraphBLAS extension.
func GlobalGetMode() (Mode, error) {
	var mode C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_MODE, &mode))
	if info == success {
		return Mode(mode), nil
	}
	return 0, makeError(info)
}

// GlobalGetOpenMP retrieves whether SuiteSparse:GraphBLAS was compiled
// with OpenMP or not.
//
// GlobalGetOpenMP is a SuiteSparse:GraphBLAS extension.
func GlobalGetOpenMP() (bool, error) {
	var openMP C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_LIBRARY_OPENMP, &openMP))
	if info == success {
		return ctogobool(openMP), nil
	}
	return false, makeError(info)
}

// GlobalSetJITCCompilerName sets the name of the C compiler
// used for JIT kernels.
//
// GlobalSetJITCCompilerName is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCCompilerName(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Global_Option_set_CHAR(C.GxB_JIT_C_COMPILER_NAME, cname))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCCompilerName retrieves the name of the C compiler
// used for JIT kernels.
//
// GlobalGetJITCCompilerName is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCCompilerName() (string, error) {
	var cname *C.char
	info := Info(C.GxB_Global_Option_get_CHAR(C.GxB_JIT_C_COMPILER_NAME, &cname))
	if info == success {
		return C.GoString(cname), nil
	}
	return "", makeError(info)
}

// GlobalSetJITCCompilerFlags sets the flags for the C compiler used
// for JIT kernels.
//
// GlobalSetJITCCompilerFlags is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCCompilerFlags(flags string) error {
	cflags := C.CString(flags)
	defer C.free(unsafe.Pointer(cflags))
	info := Info(C.GxB_Global_Option_set_CHAR(C.GxB_JIT_C_COMPILER_FLAGS, cflags))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCCompilerFlags retrieves the flags for the C compiler used
// for JIT kernels.
//
// GlobalGetJITCCompilerFlags is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCCompilerFlags() (string, error) {
	var cflags *C.char
	info := Info(C.GxB_Global_Option_get_CHAR(C.GxB_JIT_C_COMPILER_FLAGS, &cflags))
	if info == success {
		return C.GoString(cflags), nil
	}
	return "", makeError(info)
}

// GlobalSetJITCLinkerFlags sets the flags for the C linker used
// for JIT kernels.
//
// GlobalSetJITCLinkerFlags is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCLinkerFlags(flags string) error {
	cflags := C.CString(flags)
	defer C.free(unsafe.Pointer(cflags))
	info := Info(C.GxB_Global_Option_set_CHAR(C.GxB_JIT_C_LINKER_FLAGS, cflags))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCLinkerFlags retrieves the flags for the C linker used
// for JIT kernels.
//
// GlobalGetJITCLinkerFlags is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCLinkerFlags() (string, error) {
	var cflags *C.char
	info := Info(C.GxB_Global_Option_get_CHAR(C.GxB_JIT_C_LINKER_FLAGS, &cflags))
	if info == success {
		return C.GoString(cflags), nil
	}
	return "", makeError(info)
}

// GlobalSetJITCLibraries sets the libraries to link against for the C linker used
// for JIT kernels.
//
// GlobalSetJITCLibraries is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCLibraries(libs string) error {
	cname := C.CString(libs)
	defer C.free(unsafe.Pointer(cname))
	info := Info(C.GxB_Global_Option_set_CHAR(C.GxB_JIT_C_LIBRARIES, cname))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCLibraries retrieves the libraries to link against for the C linker used
// for JIT kernels.
//
// GlobalGetJITCLibraries is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCLibraries() (string, error) {
	var clibs *C.char
	info := Info(C.GxB_Global_Option_get_CHAR(C.GxB_JIT_C_LIBRARIES, &clibs))
	if info == success {
		return C.GoString(clibs), nil
	}
	return "", makeError(info)
}

// GlobalSetJITCPreface sets the preface for JIT kernels.
//
// GlobalSetJITCPreface is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCPreface(preface string) error {
	cpreface := C.CString(preface)
	defer C.free(unsafe.Pointer(cpreface))
	info := Info(C.GxB_Global_Option_set_CHAR(C.GxB_JIT_C_PREFACE, cpreface))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCPreface retrieves the preface for JIT kernels.
//
// GlobalGetJITCPreface is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCPreface() (string, error) {
	var cpreface *C.char
	info := Info(C.GxB_Global_Option_get_CHAR(C.GxB_JIT_C_PREFACE, &cpreface))
	if info == success {
		return C.GoString(cpreface), nil
	}
	return "", makeError(info)
}

// GlobalSetJITCControl enables or disables different functionalities of
// the just-in-time compiler. See [JITControl].
//
// GlobalSetJITCControl is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCControl(control JITControl) error {
	info := Info(C.GxB_Global_Option_set_INT32(C.GxB_JIT_C_CONTROL, C.int32_t(control)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCControl retrieves which functionalities of the just-in-time compiler are
// enabled and disabled. See [JITControl].
//
// GlobalGetJITCControl is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCControl() (control JITControl, err error) {
	var ccontrol C.int32_t
	info := Info(C.GxB_Global_Option_get_INT32(C.GxB_JIT_C_CONTROL, &ccontrol))
	if info == success {
		control = JITControl(ccontrol)
		return
	}
	err = makeError(info)
	return
}

// GlobalSetJITCachePath sets the path of the folder where compiled kernels
// are written to.
//
// GlobalSetJITCachePath is a SuiteSparse:GraphBLAS extension.
func GlobalSetJITCachePath(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	info := Info(C.GxB_Global_Option_set_CHAR(C.GxB_JIT_CACHE_PATH, cpath))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GlobalGetJITCachePath retrieves the path of the folder where compiled kernels
// are written to.
//
// GlobalGetJITCachePath is a SuiteSparse:GraphBLAS extension.
func GlobalGetJITCachePath() (string, error) {
	var cpath *C.char
	info := Info(C.GxB_Global_Option_get_CHAR(C.GxB_JIT_CACHE_PATH, &cpath))
	if info == success {
		return C.GoString(cpath), nil
	}
	return "", makeError(info)
}
