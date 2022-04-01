package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
)

type (
	computeMatrixT[D any] interface {
		resize(newNRows, newNCols int) computeMatrixT[D]
		computeElement(row, col int) (D, bool)
		computePipeline() *pipeline.Pipeline[any]
		computeRowPipeline(row int) *pipeline.Pipeline[any]
		computeColPipeline(col int) *pipeline.Pipeline[any]
		computeRowPipelines() []matrix1Pipeline
		computeColPipelines() []matrix1Pipeline
	}

	assignMatrixT[D any] interface {
		computeMatrixT[D]
		assignIndex(row, col int) (int, int, bool)
	}

	computedMatrix[D any] struct {
		nrows, ncols        int
		C                   *matrixReference[D]
		mask                *matrixReference[bool]
		accum               BinaryOp[D, D, D]
		compute             computeMatrixT[D]
		desc                Descriptor
		computeElement      func(row, col int) (D, bool)
		computePipeline     func() *pipeline.Pipeline[any]
		computeRowPipeline  func(row int) *pipeline.Pipeline[any]
		computeColPipeline  func(col int) *pipeline.Pipeline[any]
		computeRowPipelines func() []matrix1Pipeline
		computeColPipelines func() []matrix1Pipeline
	}
)

func makeMatrixComputeElement[D any](
	C *matrixReference[D],
	mask *matrixReference[bool],
	accum BinaryOp[D, D, D],
	compute computeMatrixT[D],
	desc Descriptor,
) func(row, col int) (D, bool) {
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
	nothing := func(_, _ int) (result D, ok bool) { return }
	makeDescDispatch := func(trueMask, falseMask func(int, int) (D, bool)) func(int, int) (D, bool) {
		if isStructure {
			if isComp {
				return func(row, col int) (result D, ok bool) {
					if _, inMask := mask.extractElement(row, col); !inMask {
						return trueMask(row, col)
					}
					return falseMask(row, col)
				}
			}
			return func(row, col int) (result D, ok bool) {
				if _, inMask := mask.extractElement(row, col); inMask {
					return trueMask(row, col)
				}
				return falseMask(row, col)
			}
		}
		if isComp {
			return func(row, col int) (result D, ok bool) {
				if inMask, _ := mask.extractElement(row, col); !inMask {
					return trueMask(row, col)
				}
				return falseMask(row, col)
			}
		}
		return func(row, col int) (result D, ok bool) {
			if inMask, _ := mask.extractElement(row, col); inMask {
				return trueMask(row, col)
			}
			return falseMask(row, col)
		}
	}
	makeMaskDispatch := func(trueMask func(int, int) (D, bool)) func(int, int) (D, bool) {
		if mask == nil {
			if isComp {
				if isReplace {
					return nothing
				}
				return C.extractElement
			}
			return trueMask
		}
		if isReplace {
			return makeDescDispatch(trueMask, nothing)
		}
		return makeDescDispatch(trueMask, C.extractElement)
	}
	if accum == nil {
		if compute, ok := compute.(assignMatrixT[D]); ok {
			return makeMaskDispatch(func(row, col int) (D, bool) {
				if rowIndex, colIndex, assignOk := compute.assignIndex(row, col); assignOk {
					return compute.computeElement(rowIndex, colIndex)
				}
				return C.extractElement(row, col)
			})
		}
		return makeMaskDispatch(compute.computeElement)
	}
	return makeMaskDispatch(func(row, col int) (D, bool) {
		var cValue, tValue D
		var cok, tok bool
		parallel.Do(func() {
			cValue, cok = C.extractElement(row, col)
		}, func() {
			tValue, tok = compute.computeElement(row, col)
		})
		if cok {
			if tok {
				return accum(cValue, tValue), true
			}
			return cValue, true
		}
		return tValue, tok
	})
}

func makeMatrixComputePipeline[D any](
	C *matrixReference[D],
	mask *matrixReference[bool],
	accum BinaryOp[D, D, D],
	compute computeMatrixT[D],
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
		trueValue func(row, col int, cValue D, cok bool, tValue D, tok bool) (D, bool),
		falseValue func(row, col int, cValue D, cok bool) (D, bool),
	) func(row, col int, maskValue, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(row, col int, _, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
					if !maskOk {
						return trueValue(row, col, cValue, cok, tValue, tok)
					}
					return falseValue(row, col, cValue, cok)
				}
			}
			return func(row, col int, _, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if maskOk {
					return trueValue(row, col, cValue, cok, tValue, tok)
				}
				return falseValue(row, col, cValue, cok)
			}
		}
		if isComp {
			return func(row, col int, maskValue, _ bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if !maskValue {
					return trueValue(row, col, cValue, cok, tValue, tok)
				}
				return falseValue(row, col, cValue, cok)
			}
		}
		return func(row, col int, maskValue, _ bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
			if maskValue {
				return trueValue(row, col, cValue, cok, tValue, tok)
			}
			return falseValue(row, col, cValue, cok)
		}
	}
	makeMaskDispatch := func(processEntry func(int, int, D, bool, D, bool) (D, bool)) func() *pipeline.Pipeline[any] {
		if mask == nil {
			if isComp {
				if isReplace {
					return func() *pipeline.Pipeline[any] { return nil }
				}
				return C.getPipeline
			}
			return func() *pipeline.Pipeline[any] {
				return makeMatrix2SourcePipeline(
					C.getPipeline(),
					compute.computePipeline(),
					processEntry,
				)
			}
		}
		if isReplace {
			return func() *pipeline.Pipeline[any] {
				return makeMaskMatrix2SourcePipeline(
					mask.getPipeline(),
					C.getPipeline(),
					compute.computePipeline(),
					makeDescDispatch(processEntry, func(_, _ int, _ D, _ bool) (result D, ok bool) {
						return
					}),
				)
			}
		}
		return func() *pipeline.Pipeline[any] {
			return makeMaskMatrix2SourcePipeline(
				mask.getPipeline(),
				C.getPipeline(),
				compute.computePipeline(),
				makeDescDispatch(processEntry, func(_, _ int, cValue D, cok bool) (D, bool) {
					return cValue, cok
				}),
			)
		}
	}
	makeDesc1Dispatch := func() func(row, col int, maskValue, maskOk bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(_, _ int, _, maskOk bool, tValue D, tok bool) (result D, ok bool) {
					if !maskOk {
						return tValue, tok
					}
					return
				}
			}
			return func(_, _ int, _, maskOk bool, tValue D, tok bool) (result D, ok bool) {
				if maskOk {
					return tValue, tok
				}
				return
			}
		}
		if isComp {
			return func(_, _ int, maskValue, _ bool, tValue D, tok bool) (result D, ok bool) {
				if !maskValue {
					return tValue, tok
				}
				return
			}
		}
		return func(_, _ int, maskValue, _ bool, tValue D, tok bool) (result D, ok bool) {
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
				return C.getPipeline
			}
			return compute.computePipeline
		}
		if isReplace {
			return func() *pipeline.Pipeline[any] {
				return makeMaskMatrix1SourcePipeline(
					mask.getPipeline(),
					compute.computePipeline(),
					makeDesc1Dispatch(),
				)
			}
		}
		return func() *pipeline.Pipeline[any] {
			return makeMaskMatrix2SourcePipeline(
				mask.getPipeline(),
				C.getPipeline(),
				compute.computePipeline(),
				makeDescDispatch(func(_, _ int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
					return tValue, tok
				}, func(_, _ int, cValue D, cok bool) (D, bool) {
					return cValue, cok
				}),
			)
		}
	}
	if accum == nil {
		if compute, ok := compute.(assignMatrixT[D]); ok {
			return makeMaskDispatch(func(row, col int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if _, _, assignOk := compute.assignIndex(row, col); assignOk {
					return tValue, tok
				}
				return cValue, cok
			})
		}
		return makeMask1Dispatch()
	}
	return makeMaskDispatch(func(_, _ int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
		if cok {
			if tok {
				return accum(cValue, tValue), true
			}
			return cValue, true
		}
		return tValue, tok
	})
}

func make1DimComputePipeline[D any](
	CGetPipeline func(index int) *pipeline.Pipeline[any],
	mask *matrixReference[bool],
	maskGetPipeline func(index int) *pipeline.Pipeline[any],
	accum BinaryOp[D, D, D],
	isAssign bool,
	assignIndex func(int, int) bool,
	computeComputePipeline func(index int) *pipeline.Pipeline[any],
	desc Descriptor,
) func(rowOrCol int) *pipeline.Pipeline[any] {
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
		trueValue func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool),
		falseValue func(index int, cValue D, cok bool) (D, bool),
	) func(index int, maskValue, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
					if !maskOk {
						return trueValue(index, cValue, cok, tValue, tok)
					}
					return falseValue(index, cValue, cok)
				}
			}
			return func(index int, _, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if maskOk {
					return trueValue(index, cValue, cok, tValue, tok)
				}
				return falseValue(index, cValue, cok)
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if !maskValue {
					return trueValue(index, cValue, cok, tValue, tok)
				}
				return falseValue(index, cValue, cok)
			}
		}
		return func(index int, maskValue, _ bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
			if maskValue {
				return trueValue(index, cValue, cok, tValue, tok)
			}
			return falseValue(index, cValue, cok)
		}
	}
	makeMaskDispatch := func(makeProcessEntry func(rowOrCol int) func(int, D, bool, D, bool) (D, bool)) func(rowOrCol int) *pipeline.Pipeline[any] {
		if mask == nil {
			if isComp {
				if isReplace {
					return func(_ int) *pipeline.Pipeline[any] { return nil }
				}
				return CGetPipeline
			}
			return func(rowOrCol int) *pipeline.Pipeline[any] {
				return makeVector2SourcePipeline(
					CGetPipeline(rowOrCol),
					computeComputePipeline(rowOrCol),
					makeProcessEntry(rowOrCol),
				)
			}
		}
		if isReplace {
			return func(rowOrCol int) *pipeline.Pipeline[any] {
				return makeMaskVector2SourcePipeline(
					maskGetPipeline(rowOrCol),
					CGetPipeline(rowOrCol),
					computeComputePipeline(rowOrCol),
					makeDescDispatch(makeProcessEntry(rowOrCol), func(_ int, _ D, _ bool) (result D, ok bool) {
						return
					}),
				)
			}
		}
		return func(rowOrCol int) *pipeline.Pipeline[any] {
			return makeMaskVector2SourcePipeline(
				maskGetPipeline(rowOrCol),
				CGetPipeline(rowOrCol),
				computeComputePipeline(rowOrCol),
				makeDescDispatch(makeProcessEntry(rowOrCol), func(_ int, cValue D, cok bool) (D, bool) {
					return cValue, cok
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
	makeMask1Dispatch := func() func(index int) *pipeline.Pipeline[any] {
		if mask == nil {
			if isComp {
				if isReplace {
					return func(_ int) *pipeline.Pipeline[any] { return nil }
				}
				return CGetPipeline
			}
			return computeComputePipeline
		}
		if isReplace {
			return func(rowOrCol int) *pipeline.Pipeline[any] {
				return makeMaskVector1SourcePipeline(
					maskGetPipeline(rowOrCol),
					computeComputePipeline(rowOrCol),
					makeDesc1Dispatch(),
				)
			}
		}
		return func(rowOrCol int) *pipeline.Pipeline[any] {
			return makeMaskVector2SourcePipeline(
				maskGetPipeline(rowOrCol),
				CGetPipeline(rowOrCol),
				computeComputePipeline(rowOrCol),
				makeDescDispatch(func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
					return tValue, tok
				}, func(index int, cValue D, cok bool) (D, bool) {
					return cValue, cok
				}),
			)
		}
	}
	if accum == nil {
		if isAssign {
			return makeMaskDispatch(func(rowOrCol int) func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				return func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
					if assignIndex(rowOrCol, index) {
						return tValue, tok
					}
					return cValue, cok
				}
			})
		}
		return makeMask1Dispatch()
	}
	return makeMaskDispatch(func(_ int) func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
		return func(_ int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
			if cok {
				if tok {
					return accum(cValue, tValue), true
				}
				return cValue, true
			}
			return tValue, tok
		}
	})
}

func addProcess[D any](p *pipeline.Pipeline[any], processEntry func(index int, value D, vok bool) (D, bool)) {
	if processEntry != nil {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					newIndex = index
					newValue, ok = processEntry(index, value, true)
					return
				})
				return slice
			})),
		)
	}
}

func mxMin2(x, y []matrix1Pipeline) int {
	if len(x) == 0 {
		if len(y) == 0 {
			return -1
		}
		return y[0].index
	}
	if len(y) == 0 {
		return x[0].index
	}
	if x[0].index < y[0].index {
		return x[0].index
	}
	return y[0].index
}

func mxMin3(x, y, z []matrix1Pipeline) int {
	if len(x) == 0 {
		return mxMin2(y, z)
	}
	if len(y) == 0 {
		return mxMin2(x, z)
	}
	if len(z) == 0 {
		return mxMin2(x, y)
	}
	return Min(x[0].index, Min(y[0].index, z[0].index))
}

func make1DimComputePipelines[D any](
	CGetPipelines func() []matrix1Pipeline,
	mask *matrixReference[bool],
	maskGetPipelines func() []matrix1Pipeline,
	accum BinaryOp[D, D, D],
	isAssign bool,
	assignIndex func(int, int) bool,
	computeComputePipelines func() []matrix1Pipeline,
	desc Descriptor,
) func() []matrix1Pipeline {
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
		trueValue func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool),
		falseValue func(cValue D, cok bool) (D, bool),
	) func(index int, maskValue, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
					if !maskOk {
						return trueValue(index, cValue, cok, tValue, tok)
					}
					return falseValue(cValue, cok)
				}
			}
			return func(index int, _, maskOk bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if maskOk {
					return trueValue(index, cValue, cok, tValue, tok)
				}
				return falseValue(cValue, cok)
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
				if !maskValue {
					return trueValue(index, cValue, cok, tValue, tok)
				}
				return falseValue(cValue, cok)
			}
		}
		return func(index int, maskValue, _ bool, cValue D, cok bool, tValue D, tok bool) (D, bool) {
			if maskValue {
				return trueValue(index, cValue, cok, tValue, tok)
			}
			return falseValue(cValue, cok)
		}
	}
	makeCDescDispatch := func(
		trueValue func(index int, cValue D, cok bool) (D, bool),
		falseValue func(cValue D, cok bool) (D, bool),
	) func(index int, maskValue, maskOk bool, cValue D, cok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, cValue D, cok bool) (D, bool) {
					if !maskOk {
						return trueValue(index, cValue, cok)
					}
					return falseValue(cValue, cok)
				}
			}
			return func(index int, _, maskOk bool, cValue D, cok bool) (D, bool) {
				if maskOk {
					return trueValue(index, cValue, cok)
				}
				return falseValue(cValue, cok)
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, cValue D, cok bool) (D, bool) {
				if !maskValue {
					return trueValue(index, cValue, cok)
				}
				return falseValue(cValue, cok)
			}
		}
		return func(index int, maskValue, _ bool, cValue D, cok bool) (D, bool) {
			if maskValue {
				return trueValue(index, cValue, cok)
			}
			return falseValue(cValue, cok)
		}
	}
	makeTDescDispatch := func(
		trueValue func(index int, tValue D, tok bool) (D, bool),
	) func(index int, maskValue, maskOk bool, tValue D, tok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, tValue D, tok bool) (result D, ok bool) {
					if !maskOk {
						return trueValue(index, tValue, tok)
					}
					return
				}
			}
			return func(index int, _, maskOk bool, tValue D, tok bool) (result D, ok bool) {
				if maskOk {
					return trueValue(index, tValue, tok)
				}
				return
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, tValue D, tok bool) (result D, ok bool) {
				if !maskValue {
					return trueValue(index, tValue, tok)
				}
				return
			}
		}
		return func(index int, maskValue, _ bool, tValue D, tok bool) (result D, ok bool) {
			if maskValue {
				return trueValue(index, tValue, tok)
			}
			return
		}
	}
	makeMaskDispatch := func(makeProcessEntry func(rowOrCol int) (func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool), func(index int, cValue D, cok bool) (D, bool), func(index int, tValue D, tok bool) (D, bool))) func() []matrix1Pipeline {
		if mask == nil {
			if isComp {
				if isReplace {
					return func() []matrix1Pipeline { return nil }
				}
				return CGetPipelines
			}
			return func() (result []matrix1Pipeline) {
				CPipelines := CGetPipelines()
				computePipelines := computeComputePipelines()
				for {
					minIndex := mxMin2(CPipelines, computePipelines)
					if minIndex < 0 {
						return
					}
					if len(CPipelines) > 0 && CPipelines[0].index == minIndex {
						if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
							processEntry, _, _ := makeProcessEntry(minIndex)
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p:     makeVector2SourcePipeline(CPipelines[0].p, computePipelines[0].p, processEntry),
							})
							computePipelines = computePipelines[1:]
						} else {
							_, processEntry, _ := makeProcessEntry(minIndex)
							addProcess(CPipelines[0].p, processEntry)
							result = append(result, CPipelines[0])
						}
						CPipelines = CPipelines[1:]
					} else {
						_, _, processEntry := makeProcessEntry(minIndex)
						addProcess(computePipelines[0].p, processEntry)
						result = append(result, computePipelines[0])
						computePipelines = computePipelines[1:]
					}
				}
			}
		}
		return func() (result []matrix1Pipeline) {
			maskPipelines := maskGetPipelines()
			CPipelines := CGetPipelines()
			computePipelines := computeComputePipelines()
			for {
				if len(CPipelines) == 0 && len(computePipelines) == 0 {
					return
				}
			nextm:
				minIndex := mxMin3(maskPipelines, CPipelines, computePipelines)
				if len(maskPipelines) > 0 && maskPipelines[0].index == minIndex {
					if len(CPipelines) > 0 && CPipelines[0].index == minIndex {
						if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
							processEntry, _, _ := makeProcessEntry(minIndex)
							var falseValue func(D, bool) (D, bool)
							if isReplace {
								falseValue = func(_ D, _ bool) (result D, ok bool) {
									return
								}
							} else {
								falseValue = func(cValue D, cok bool) (D, bool) {
									return cValue, cok
								}
							}
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p: makeMaskVector2SourcePipeline(
									maskPipelines[0].p,
									CPipelines[0].p,
									computePipelines[0].p,
									makeDescDispatch(processEntry, falseValue)),
							})
							computePipelines = computePipelines[1:]
						} else {
							var falseValue func(D, bool) (D, bool)
							if isReplace {
								falseValue = func(_ D, _ bool) (result D, ok bool) {
									return
								}
							} else {
								falseValue = func(cValue D, cok bool) (D, bool) {
									return cValue, cok
								}
							}
							_, processEntry, _ := makeProcessEntry(minIndex)
							if processEntry == nil {
								if !isReplace {
									result = append(result, CPipelines[0])
									continue
								}
								processEntry = func(_ int, cValue D, cok bool) (D, bool) {
									return cValue, cok
								}
							}
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p: makeMaskVector1SourcePipeline(
									maskPipelines[0].p,
									CPipelines[0].p,
									makeCDescDispatch(processEntry, falseValue),
								),
							})
						}
						CPipelines = CPipelines[1:]
					} else if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
						_, _, processEntry := makeProcessEntry(minIndex)
						if processEntry == nil {
							processEntry = func(_ int, tValue D, tok bool) (D, bool) {
								return tValue, tok
							}
						}
						result = append(result, matrix1Pipeline{
							index: minIndex,
							p: makeMaskVector1SourcePipeline(
								maskPipelines[0].p,
								computePipelines[0].p,
								makeTDescDispatch(processEntry),
							),
						})
						computePipelines = computePipelines[1:]
					} else {
						maskPipelines = maskPipelines[1:]
						goto nextm
					}
					maskPipelines = maskPipelines[1:]
				} else if len(CPipelines) > 0 && CPipelines[0].index == minIndex {
					if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
						if isComp {
							processEntry, _, _ := makeProcessEntry(minIndex)
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p: makeVector2SourcePipeline(
									CPipelines[0].p,
									computePipelines[0].p,
									processEntry,
								),
							})
						} else if !isReplace {
							result = append(result, CPipelines[0])
						}
						computePipelines = computePipelines[1:]
					} else {
						if isComp {
							_, processEntry, _ := makeProcessEntry(minIndex)
							addProcess(CPipelines[0].p, processEntry)
							result = append(result, CPipelines[0])
						} else if !isReplace {
							result = append(result, CPipelines[0])
						}
					}
					CPipelines = CPipelines[1:]
				} else { // len(computePipelines) > 0 && computePipelines[0].index == minIndex
					if isComp {
						_, _, processEntry := makeProcessEntry(minIndex)
						addProcess(computePipelines[0].p, processEntry)
						result = append(result, computePipelines[0])
					}
					computePipelines = computePipelines[1:]
				}
			}
		}
	}
	desc1Dispatch := func() func(index int, maskValue, maskOk bool, tValue D, tok bool) (D, bool) {
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
	}()
	desc2Dispatch := func() func(index int, maskValue, maskOk bool, cValue D, cok bool) (D, bool) {
		if isStructure {
			if isComp {
				return func(index int, _, maskOk bool, cValue D, cok bool) (result D, ok bool) {
					if !maskOk {
						return
					}
					return cValue, cok
				}
			}
			return func(index int, _, maskOk bool, cValue D, cok bool) (result D, ok bool) {
				if maskOk {
					return
				}
				return cValue, cok
			}
		}
		if isComp {
			return func(index int, maskValue, _ bool, cValue D, cok bool) (result D, ok bool) {
				if !maskValue {
					return
				}
				return cValue, cok
			}
		}
		return func(index int, maskValue, _ bool, cValue D, cok bool) (result D, ok bool) {
			if maskValue {
				return
			}
			return cValue, cok
		}
	}()
	makeMask1Dispatch := func() func() []matrix1Pipeline {
		if mask == nil {
			if isComp {
				if isReplace {
					return func() []matrix1Pipeline { return nil }
				}
				return CGetPipelines
			}
			return computeComputePipelines
		}
		if isReplace {
			return func() (result []matrix1Pipeline) {
				maskPipelines := maskGetPipelines()
				computePipelines := computeComputePipelines()
				for {
					if len(computePipelines) == 0 {
						return
					}
				nextm:
					minIndex := mxMin2(maskPipelines, computePipelines)
					if len(maskPipelines) > 0 && maskPipelines[0].index == minIndex {
						if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p: makeMaskVector1SourcePipeline(
									maskPipelines[0].p,
									computePipelines[0].p,
									desc1Dispatch,
								),
							})
							computePipelines = computePipelines[1:]
						} else {
							maskPipelines = maskPipelines[1:]
							goto nextm
						}
						maskPipelines = maskPipelines[1:]
					} else if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
						if isComp {
							result = append(result, computePipelines[0])
						}
						computePipelines = computePipelines[1:]
					}
				}
			}
		}
		return func() (result []matrix1Pipeline) {
			maskPipelines := maskGetPipelines()
			CPipelines := CGetPipelines()
			computePipelines := computeComputePipelines()
			for {
				if len(CPipelines) == 0 && len(computePipelines) == 0 {
					return
				}
			nextm:
				minIndex := mxMin3(maskPipelines, CPipelines, computePipelines)
				if len(maskPipelines) > 0 && maskPipelines[0].index == minIndex {
					if len(CPipelines) > 0 && CPipelines[0].index == minIndex {
						if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p: makeMaskVector2SourcePipeline(
									maskPipelines[0].p,
									CPipelines[0].p,
									computePipelines[0].p,
									makeDescDispatch(func(_ int, _ D, _ bool, tValue D, tok bool) (D, bool) {
										return tValue, tok
									}, func(cValue D, cok bool) (D, bool) {
										return cValue, cok
									}),
								),
							})
							computePipelines = computePipelines[1:]
						} else {
							result = append(result, matrix1Pipeline{
								index: minIndex,
								p: makeMaskVector1SourcePipeline(
									maskPipelines[0].p,
									CPipelines[0].p,
									desc2Dispatch),
							})
						}
						CPipelines = CPipelines[1:]
					} else if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
						result = append(result, matrix1Pipeline{
							index: minIndex,
							p: makeMaskVector1SourcePipeline(
								maskPipelines[0].p,
								computePipelines[0].p,
								desc1Dispatch),
						})
						computePipelines = computePipelines[1:]
					} else {
						maskPipelines = maskPipelines[1:]
						goto nextm
					}
					maskPipelines = maskPipelines[1:]
				} else if len(CPipelines) > 0 && CPipelines[0].index == minIndex {
					if len(computePipelines) > 0 && computePipelines[0].index == minIndex {
						if isComp {
							result = append(result, computePipelines[0])
						} else {
							result = append(result, CPipelines[0])
						}
						computePipelines = computePipelines[1:]
					} else {
						if !isComp {
							result = append(result, CPipelines[0])
						}
					}
					CPipelines = CPipelines[1:]
				} else {
					if isComp {
						result = append(result, computePipelines[0])
					}
					computePipelines = computePipelines[1:]
				}
			}
		}
	}
	if accum == nil {
		if isAssign {
			return makeMaskDispatch(func(rowOrCol int) (func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool), func(index int, cValue D, cok bool) (D, bool), func(index int, tValue D, tok bool) (D, bool)) {
				return func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
						if assignIndex(rowOrCol, index) {
							return tValue, tok
						}
						return cValue, cok
					}, func(index int, cValue D, cok bool) (result D, ok bool) {
						if assignIndex(rowOrCol, index) {
							return
						}
						return cValue, cok
					}, func(index int, tValue D, tok bool) (result D, ok bool) {
						if assignIndex(rowOrCol, index) {
							return tValue, tok
						}
						return
					}
			})
		}
		return makeMask1Dispatch()
	}
	return makeMaskDispatch(func(_ int) (func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool), func(index int, cValue D, cok bool) (D, bool), func(index int, tValue D, tok bool) (D, bool)) {
		return func(index int, cValue D, cok bool, tValue D, tok bool) (D, bool) {
			if cok {
				if tok {
					return accum(cValue, tValue), true
				}
				return cValue, true
			}
			return tValue, tok
		}, nil, nil
	})
}

func newComputedMatrix[D any](
	nrows, ncols int,
	C *matrixReference[D],
	mask *matrixReference[bool],
	accum BinaryOp[D, D, D],
	compute computeMatrixT[D],
	desc Descriptor,
) *computedMatrix[D] {
	assign, assignOk := compute.(assignMatrixT[D])
	return &computedMatrix[D]{
		nrows:           nrows,
		ncols:           ncols,
		C:               C,
		mask:            mask,
		accum:           accum,
		compute:         compute,
		computeElement:  makeMatrixComputeElement[D](C, mask, accum, compute, desc),
		computePipeline: makeMatrixComputePipeline[D](C, mask, accum, compute, desc),
		computeRowPipeline: make1DimComputePipeline(
			C.getRowPipeline,
			mask, mask.getRowPipeline,
			accum,
			assignOk, func(row, col int) bool {
				_, _, ok := assign.assignIndex(row, col)
				return ok
			},
			compute.computeRowPipeline,
			desc,
		),
		computeColPipeline: make1DimComputePipeline(
			C.getColPipeline,
			mask, mask.getColPipeline,
			accum,
			assignOk, func(col, row int) bool {
				_, _, ok := assign.assignIndex(row, col)
				return ok
			},
			compute.computeColPipeline,
			desc,
		),
		computeRowPipelines: make1DimComputePipelines(
			C.getRowPipelines,
			mask, mask.getRowPipelines,
			accum,
			assignOk, func(row, col int) bool {
				_, _, ok := assign.assignIndex(row, col)
				return ok
			},
			compute.computeRowPipelines,
			desc,
		),
		computeColPipelines: make1DimComputePipelines(
			C.getColPipelines,
			mask, mask.getColPipelines,
			accum,
			assignOk, func(col, row int) bool {
				_, _, ok := assign.assignIndex(row, col)
				return ok
			},
			compute.computeColPipelines,
			desc,
		),
	}
}

func (m *computedMatrix[D]) resize(ref *matrixReference[D], newNRows, newNCols int) *matrixReference[D] {
	if newNRows == m.nrows && newNCols == m.ncols {
		return ref
	}
	var C *matrixReference[D]
	var mask *matrixReference[bool]
	var compute computeMatrixT[D]
	parallel.Do(func() {
		if m.C != nil {
			C = m.C.resize(newNRows, newNCols)
		}
	}, func() {
		if m.mask != nil {
			mask = m.mask.resize(newNRows, newNCols)
		}
	}, func() {
		if m.compute != nil {
			compute = m.compute.resize(newNRows, newNCols)
		}
	})
	return newMatrixReference[D](newComputedMatrix[D](
		newNRows, newNCols,
		C,
		mask,
		m.accum,
		compute,
		m.desc,
	), -1)
}

func (m *computedMatrix[D]) size() (int, int) {
	return m.nrows, m.ncols
}

func (m *computedMatrix[D]) nvals() (n int) {
	p := m.getPipeline()
	a := newParallelArray[int](16)
	p.Add(
		pipeline.Par(pipeline.Receive(func(seq int, data any) any {
			a.set(seq, len(data.(matrixSlice[D]).values))
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

func (m *computedMatrix[D]) setElement(ref *matrixReference[D], value D, row, col int) *matrixReference[D] {
	return newMatrixReference[D](newListMatrix[D](
		m.nrows, m.ncols, ref,
		&matrixValueList[D]{
			row:   row,
			col:   col,
			value: value,
		}), -1)
}

func (m *computedMatrix[D]) removeElement(ref *matrixReference[D], row, col int) *matrixReference[D] {
	return newMatrixReference[D](newListMatrix[D](
		m.nrows, m.ncols, ref,
		&matrixValueList[D]{
			row: -row,
			col: -col,
		}), -1)
}

func (m *computedMatrix[D]) extractElement(row, col int) (result D, ok bool) {
	return m.computeElement(row, col)
}

func (m *computedMatrix[D]) getPipeline() *pipeline.Pipeline[any] {
	return m.computePipeline()
}

func (m *computedMatrix[D]) getRowPipeline(row int) *pipeline.Pipeline[any] {
	return m.computeRowPipeline(row)
}

func (m *computedMatrix[D]) getColPipeline(col int) *pipeline.Pipeline[any] {
	return m.computeColPipeline(col)
}

func (m *computedMatrix[D]) getRowPipelines() []matrix1Pipeline {
	return m.computeRowPipelines()
}

func (m *computedMatrix[D]) getColPipelines() []matrix1Pipeline {
	return m.computeColPipelines()
}

func (m *computedMatrix[D]) extractTuples() (rows, cols []int, values []D) {
	var result matrixSlice[D]
	result.collect(m.getPipeline())
	rows = result.rows
	cols = result.cols
	values = result.values
	return
}

func (_ *computedMatrix[D]) optimized() bool {
	return false
}

func (m *computedMatrix[D]) optimize() functionalMatrix[D] {
	// todo: the next two lines may be more efficient once extractTuples is improved
	rows, cols, values := m.extractTuples()
	newRows, rowSpans := csrRows(rows)
	return newCSRMatrix[D](m.nrows, m.ncols, newRows, rowSpans, cols, values)
}
