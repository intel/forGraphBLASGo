package GrB

// #include "GraphBLAS.h"
import "C"

// Context objects control the number of threads used by OpenMP per
// application thread.
//
// Context is a SuiteSparse:GraphBLAS extension.
type Context struct {
	grb C.GxB_Context
}

// ContextWorld is the default [Context] object. It can be modified with
// [Context.SetNThreads] and [Context.SetChunk].
//
// ContextWorld is a SuiteSparse:GraphBLAS extension.
var ContextWorld = Context{
	C.GxB_CONTEXT_WORLD,
}

// ContextNew creates a new [Context] and initializes it with the current
// global settings for [Context.GetNThreads] and [Context.GetChunk].
//
// ContextNew is a SuiteSparse:GraphBLAS extension.
func ContextNew() (context Context, err error) {
	info := Info(C.GxB_Context_new(&context.grb))
	if info == success {
		return
	}
	err = makeError(info)
	return
}

// Free destroys a previously created [Context] and releases any resources associated with
// it. Calling Free on an object that is not [Context.Valid]() is legal. The behavior of a
// program that calls Free on a pre-defined context is undefined.
//
// GraphBLAS execution errors that may cause a panic:
//   - [Panic]
func (context *Context) Free() error {
	info := Info(C.GxB_Context_free(&context.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetNThreads sets the maximum number of threads to use. If nthreads <= 0,
// then one thread is used.
//
// SetNThreads is a SuiteSparse:GraphBLAS extension.
func (context Context) SetNThreads(nthreads int) error {
	info := Info(C.GxB_Context_set_INT32(context.grb, C.GxB_NTHREADS, C.int32_t(nthreads)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// SetChunk sets the chunk size for small problems. If chunk < 1, then
// the default is used.
//
// SetChunk is a SuiteSparse:GraphBLAS extension.
func (context Context) SetChunk(chunk int) error {
	info := Info(C.GxB_Context_set_INT32(context.grb, C.GxB_CHUNK, C.int32_t(chunk)))
	if info == success {
		return nil
	}
	return makeError(info)
}

// GetNThreads returns the maximum number of threads to use.
//
// GetNThreads is a SuiteSparse:GraphBLAS extension.
func (context Context) GetNThreads() (nthreads int, err error) {
	var cnthreads C.int32_t
	info := Info(C.GxB_Context_get_INT32(context.grb, C.GxB_NTHREADS, &cnthreads))
	if info == success {
		return int(cnthreads), nil
	}
	err = makeError(info)
	return
}

// GetChunk returns the chunk size for small problems.
//
// GetChunk is a SuiteSparse:GraphBLAS extension.
func (context Context) GetChunk() (chunk int, err error) {
	var cchunk C.int32_t
	info := Info(C.GxB_Context_get_INT32(context.grb, C.GxB_NTHREADS, &cchunk))
	if info == success {
		return int(cchunk), nil
	}
	err = makeError(info)
	return
}

// Engage sets the provided [Context] object as the [Context] object
// for this user thread. Multiple user threads can share a single [Context].
// Any prior [Context] for this thread is superseded by the new [Context].
// (The prior one is not [Free]d.) Future calls by this user thread will
// use the provided [Context].
//
// Caution: By default, the Go scheduler can freely move goroutines
// between threads. In order for Engage to be meaningful, the current
// goroutine has to be tied to the current thread with [runtime.LockOSThread].
//
// Please also consider that [runtime.LockOSThread] and [runtime.UnlockOSThread]
// nest: A goroutine is only unlocked from the current thread when there have
// been as many calls to [runtime.UnlockOSThread] as there have been to
// [runtime.LockOSThread]. On the other hand, Engage and [Context.Disengage] do
// not nest, but have immediate effects.
//
// GraphBLAS API errors that may be returned:
//   - [NullPointer], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject]
//
// Engage is a SuiteSparse:GraphBLAS extension.
func (context Context) Engage() error {
	info := Info(C.GxB_Context_engage(context.grb))
	if info == success {
		return nil
	}
	return makeError(info)
}

// ContextDisengage disengages the given context from the current thread.
//
// If context is nil, or if it points to [ContextWorld], then any [Context]
// for the current thread is disengaged. If a valid non-nil [Context] is
// provided, and it matches the current [Context] for this user thread,
// it is disengaged. In all of these cases, nil is returned.
//
// If a non-nil [Context] is provided on input that does not match the
// current [Context] for this thread, then [InvalidValue] is returned.
// In that case, the current [Context] for this user thread is unmodified.
//
// Please also read the caution sections in [Context.Engage] carefully.
//
// GraphBLAS API errors that may be returned:
//   - [InvalidValue], [UninitializedObject]
//
// GraphBLAS execution errors that may cause a panic:
//   - [InvalidObject]
//
// ContextDisengage is a SuiteSparse:GraphBLAS extension.
func ContextDisengage(context *Context) error {
	var ctx C.GxB_Context = nil
	if context != nil {
		ctx = context.grb
	}
	info := Info(C.GxB_Context_disengage(ctx))
	if info == success {
		return nil
	}
	return makeError(info)
}
