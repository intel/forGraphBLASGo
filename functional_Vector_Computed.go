package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
)

type (
	computeVectorT[D any] interface {
		resize(newSize int) computeVectorT[D]
		computeElement(index int) (D, bool)
		computePipeline() *pipeline.Pipeline[any]
	}

	assignVectorT[D any] interface {
		computeVectorT[D]
		assignIndex(index int) (int, bool)
	}

	computedVector[D any] struct {
		nsize           int
		w               *vectorReference[D]
		mask            *vectorReference[bool]
		accum           BinaryOp[D, D, D]
		compute         computeVectorT[D]
		desc            Descriptor
		computeElement  func(index int) (D, bool)
		computePipeline func() *pipeline.Pipeline[any]
	}
)

func makeVectorComputeElement[D any](
	w *vectorReference[D],
	mask *vectorReference[bool],
	accum BinaryOp[D, D, D],
	compute computeVectorT[D],
	desc Descriptor,
) func(int) (D, bool) {
	isReplace, err := desc.Is(Outp, Replace)
	if err != nil {
		panic(err)
	}
	isStructure, err := desc.Is(Mask, Structure)
	if err != nil {
		panic(err)
	}
	isComp, err := desc.Is(Mask, Comp)
	if err != nil {
		panic(err)
	}
	nothing := func(_ int) (result D, ok bool) { return }
	makeDescDispatch := func(trueMask, falseMask func(int) (D, bool)) func(int) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int) (result D, ok bool) {
					if _, inMask := mask.extractElement(index); !inMask {
						return trueMask(index)
					}
					return falseMask(index)
				}
			}
			return func(index int) (result D, ok bool) {
				if _, inMask := mask.extractElement(index); inMask {
					return trueMask(index)
				}
				return falseMask(index)
			}
		}
		if isComp {
			return func(index int) (result D, ok bool) {
				if inMask, _ := mask.extractElement(index); !inMask {
					return trueMask(index)
				}
				return falseMask(index)
			}
		}
		return func(index int) (result D, ok bool) {
			if inMask, _ := mask.extractElement(index); inMask {
				return trueMask(index)
			}
			return falseMask(index)
		}
	}
	makeMaskDispatch := func(trueMask func(int) (D, bool)) func(int) (D, bool) {
		if mask == nil {
			if isComp {
				if isReplace {
					return nothing
				}
				return w.extractElement
			}
			return trueMask
		}
		if isReplace {
			return makeDescDispatch(trueMask, nothing)
		}
		return makeDescDispatch(trueMask, w.extractElement)
	}
	if accum == nil {
		if compute, ok := compute.(assignVectorT[D]); ok {
			return makeMaskDispatch(func(index int) (D, bool) {
				if i, assignOk := compute.assignIndex(index); assignOk {
					return compute.computeElement(i)
				}
				return w.extractElement(index)
			})
		}
		return makeMaskDispatch(compute.computeElement)
	}
	return makeMaskDispatch(func(index int) (D, bool) {
		var wValue, tValue D
		var wok, tok bool
		parallel.Do(func() {
			wValue, wok = w.extractElement(index)
		}, func() {
			tValue, tok = compute.computeElement(index)
		})
		if wok {
			if tok {
				return accum(wValue, tValue), true
			}
			return wValue, true
		}
		return tValue, tok
	})
}

func makeVectorComputePipeline[D any](
	w *vectorReference[D],
	mask *vectorReference[bool],
	accum BinaryOp[D, D, D],
	compute computeVectorT[D],
	desc Descriptor,
) func() *pipeline.Pipeline[any] {
	isReplace, err := desc.Is(Outp, Replace)
	if err != nil {
		panic(err)
	}
	isStructure, err := desc.Is(Mask, Structure)
	if err != nil {
		panic(err)
	}
	isComp, err := desc.Is(Mask, Comp)
	if err != nil {
		panic(err)
	}
	makeDescDispatch := func(
		trueValue func(index int, wValue D, wok bool, tValue D, tok bool) (D, bool),
		falseValue func(wValue D, wok bool) (D, bool),
	) func(index int, maskValue, maskOk bool, wValue D, wok bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, wValue D, wok bool, tValue D, tok bool) (D, bool) {
					if !maskOk {
						return trueValue(index, wValue, wok, tValue, tok)
					}
					return falseValue(wValue, wok)
				}
			}
			return func(index int, _, maskOk bool, wValue D, wok bool, tValue D, tok bool) (D, bool) {
				if maskOk {
					return trueValue(index, wValue, wok, tValue, tok)
				}
				return falseValue(wValue, wok)
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, wValue D, wok bool, tValue D, tok bool) (D, bool) {
				if !maskValue {
					return trueValue(index, wValue, wok, tValue, tok)
				}
				return falseValue(wValue, wok)
			}
		}
		return func(index int, maskValue, _ bool, wValue D, wok bool, tValue D, tok bool) (D, bool) {
			if maskValue {
				return trueValue(index, wValue, wok, tValue, tok)
			}
			return falseValue(wValue, wok)
		}
	}
	makeMaskDispatch := func(processEntry func(int, D, bool, D, bool) (D, bool)) func() *pipeline.Pipeline[any] {
		if mask == nil {
			if isComp {
				if isReplace {
					return func() *pipeline.Pipeline[any] { return nil }
				}
				return w.getPipeline
			}
			return func() *pipeline.Pipeline[any] {
				return makeVector2SourcePipeline(
					w.getPipeline(),
					compute.computePipeline(),
					processEntry,
				)
			}
		}
		if isReplace {
			return func() *pipeline.Pipeline[any] {
				return makeMaskVector2SourcePipeline(
					mask.getPipeline(),
					w.getPipeline(),
					compute.computePipeline(),
					makeDescDispatch(processEntry, func(_ D, _ bool) (result D, ok bool) {
						return
					}),
				)
			}
		}
		return func() *pipeline.Pipeline[any] {
			return makeMaskVector2SourcePipeline(
				mask.getPipeline(),
				w.getPipeline(),
				compute.computePipeline(),
				makeDescDispatch(processEntry, func(wValue D, wok bool) (D, bool) {
					return wValue, wok
				}),
			)
		}
	}
	makeDesc1Dispatch := func() func(index int, maskValue, maskOk bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, tValue D, tok bool) (result D, ok bool) {
					if !maskOk {
						return tValue, tok
					}
					return
				}
			}
			return func(index int, _, maskOk bool, tValue D, tok bool) (result D, ok bool) {
				if maskOk {
					return tValue, tok
				}
				return
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, tValue D, tok bool) (result D, ok bool) {
				if !maskValue {
					return tValue, tok
				}
				return
			}
		}
		return func(index int, maskValue, _ bool, tValue D, tok bool) (result D, ok bool) {
			if maskValue {
				return tValue, tok
			}
			return
		}
	}
	makeMask1Dispatch := func() func() *pipeline.Pipeline[any] {
		if mask == nil {
			if isComp {
				if isReplace {
					return func() *pipeline.Pipeline[any] { return nil }
				}
				return w.getPipeline
			}
			return compute.computePipeline
		}
		if isReplace {
			return func() *pipeline.Pipeline[any] {
				return makeMaskVector1SourcePipeline(
					mask.getPipeline(),
					compute.computePipeline(),
					makeDesc1Dispatch(),
				)
			}
		}
		return func() *pipeline.Pipeline[any] {
			return makeMaskVector2SourcePipeline(
				mask.getPipeline(),
				w.getPipeline(),
				compute.computePipeline(),
				makeDescDispatch(func(_ int, wValue D, wok bool, tValue D, tok bool) (D, bool) {
					return tValue, tok
				}, func(wValue D, wok bool) (D, bool) {
					return wValue, wok
				}),
			)
		}
	}
	if accum == nil {
		if compute, ok := compute.(assignVectorT[D]); ok {
			return makeMaskDispatch(func(index int, wValue D, wok bool, tValue D, tok bool) (D, bool) {
				if _, assignOk := compute.assignIndex(index); assignOk {
					return tValue, tok
				}
				return wValue, wok
			})
		}
		return makeMask1Dispatch()
	}
	return makeMaskDispatch(func(index int, wValue D, wok bool, tValue D, tok bool) (D, bool) {
		if wok {
			if tok {
				return accum(wValue, tValue), true
			}
			return wValue, true
		}
		return tValue, tok
	})
}

func newComputedVector[D any](
	size int,
	w *vectorReference[D],
	mask *vectorReference[bool],
	accum BinaryOp[D, D, D],
	compute computeVectorT[D],
	desc Descriptor,
) *computedVector[D] {
	return &computedVector[D]{
		nsize:           size,
		w:               w,
		mask:            mask,
		accum:           accum,
		compute:         compute,
		computeElement:  makeVectorComputeElement[D](w, mask, accum, compute, desc),
		computePipeline: makeVectorComputePipeline[D](w, mask, accum, compute, desc),
	}
}

func (v *computedVector[D]) resize(ref *vectorReference[D], newSize int) *vectorReference[D] {
	if newSize == v.nsize {
		return ref
	}
	var w *vectorReference[D]
	var mask *vectorReference[bool]
	var compute computeVectorT[D]
	parallel.Do(func() {
		if v.w != nil {
			w = v.w.resize(newSize)
		}
	}, func() {
		if mask != nil {
			mask = v.mask.resize(newSize)
		}
	}, func() {
		if v.compute != nil {
			compute = compute.resize(newSize)
		}
	})
	return newVectorReference[D](newComputedVector[D](
		newSize,
		w,
		mask,
		v.accum,
		compute,
		v.desc,
	), -1)
}

func (v *computedVector[D]) size() int {
	return v.nsize
}

func (v *computedVector[D]) nvals() (n int) {
	p := v.getPipeline()
	a := newParallelArray[int](16)
	p.Add(
		pipeline.Par(pipeline.Receive(func(seq int, data any) any {
			a.set(seq, len(data.(vectorSlice[D]).values))
			return nil
		})),
	)
	p.Run()
	if err := p.Err(); err != nil {
		panic(err)
	}
	n, _ = reduceParallelArray(a, func(x int) int {
		return x
	}, func(x, y int) int {
		return x + y
	})
	return
}

func (v *computedVector[D]) setElement(ref *vectorReference[D], value D, index int) *vectorReference[D] {
	return newVectorReference[D](newListVector[D](
		v.nsize, ref,
		&vectorValueList[D]{
			col:   index,
			value: value,
		}), -1)
}

func (v *computedVector[D]) removeElement(ref *vectorReference[D], index int) *vectorReference[D] {
	return newVectorReference[D](newListVector[D](
		v.nsize, ref,
		&vectorValueList[D]{
			col: -index,
		}), -1)
}

func (v *computedVector[D]) extractElement(index int) (result D, ok bool) {
	return v.computeElement(index)
}

func (v *computedVector[D]) getPipeline() *pipeline.Pipeline[any] {
	return v.computePipeline()
}

func (v *computedVector[D]) extractTuples() (indices []int, values []D) {
	var result vectorSlice[D]
	result.collect(v.getPipeline())
	indices = result.indices
	values = result.values
	return
}

func (_ *computedVector[D]) optimized() bool {
	return false
}

func (v *computedVector[D]) optimize() functionalVector[D] {
	indices, values := v.extractTuples()
	return newSparseVector[D](v.nsize, indices, values)
}
