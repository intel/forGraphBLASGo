package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
)

type matrixEWiseMultBinaryOp[DC, DA, DB any] struct {
	op BinaryOp[DC, DA, DB]
	A  *matrixReference[DA]
	B  *matrixReference[DB]
}

func newMatrixEWiseMultBinaryOp[DC, DA, DB any](
	op BinaryOp[DC, DA, DB],
	A *matrixReference[DA],
	B *matrixReference[DB],
) computeMatrixT[DC] {
	return matrixEWiseMultBinaryOp[DC, DA, DB]{
		op: op,
		A:  A,
		B:  B,
	}
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	var a *matrixReference[DA]
	var b *matrixReference[DB]
	parallel.Do(func() {
		a = compute.A.resize(newNRows, newNCols)
	}, func() {
		b = compute.B.resize(newNRows, newNCols)
	})
	return newMatrixEWiseMultBinaryOp[DC, DA, DB](compute.op, a, b)
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) computeElement(row, col int) (result DC, ok bool) {
	if aValue, aok := compute.A.extractElement(row, col); aok {
		if bValue, bok := compute.B.extractElement(row, col); bok {
			return compute.op(aValue, bValue), true
		}
	}
	return
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) computePipeline() *pipeline.Pipeline[any] {
	ap := compute.A.getPipeline()
	if ap == nil {
		return nil
	}
	bp := compute.B.getPipeline()
	if bp == nil {
		return nil
	}
	return makeMatrix2SourcePipeline(ap, bp,
		func(_, _ int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
			if aok && bok {
				return compute.op(aValue, bValue), true
			}
			return
		})
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	ap := compute.A.getRowPipeline(row)
	if ap == nil {
		return nil
	}
	bp := compute.B.getRowPipeline(row)
	if bp == nil {
		return nil
	}
	return makeVector2SourcePipeline(ap, bp,
		func(_ int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
			if aok && bok {
				return compute.op(aValue, bValue), true
			}
			return
		})
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	ap := compute.A.getColPipeline(col)
	if ap == nil {
		return nil
	}
	bp := compute.B.getColPipeline(col)
	if bp == nil {
		return nil
	}
	return makeVector2SourcePipeline(ap, bp,
		func(_ int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
			if aok && bok {
				return compute.op(aValue, bValue), true
			}
			return
		})
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) mergePipelines(aps, bps []matrix1Pipeline) (result []matrix1Pipeline) {
	for {
		if len(aps) == 0 || len(bps) == 0 {
			return
		}
		if aps[0].index == bps[0].index {
			result = append(result, matrix1Pipeline{
				index: aps[0].index,
				p: makeVector2SourcePipeline(aps[0].p, bps[0].p,
					func(_ int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
						if aok && bok {
							return compute.op(aValue, bValue), true
						}
						return
					}),
			})
		} else if aps[0].index < bps[0].index {
			aps = aps[1:]
		} else {
			bps = bps[1:]
		}
	}
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) computeRowPipelines() []matrix1Pipeline {
	return compute.mergePipelines(compute.A.getRowPipelines(), compute.B.getRowPipelines())
}

func (compute matrixEWiseMultBinaryOp[DC, DA, DB]) computeColPipelines() []matrix1Pipeline {
	return compute.mergePipelines(compute.A.getColPipelines(), compute.B.getColPipelines())
}

type matrixEWiseAddBinaryOp[D any] struct {
	op   BinaryOp[D, D, D]
	A, B *matrixReference[D]
}

func newMatrixEWiseAddBinaryOp[D any](
	op BinaryOp[D, D, D],
	A, B *matrixReference[D],
) computeMatrixT[D] {
	return matrixEWiseAddBinaryOp[D]{
		op: op,
		A:  A,
		B:  B,
	}
}

func (compute matrixEWiseAddBinaryOp[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	var a, b *matrixReference[D]
	parallel.Do(func() {
		a = compute.A.resize(newNRows, newNCols)
	}, func() {
		b = compute.B.resize(newNRows, newNCols)
	})
	return newMatrixEWiseAddBinaryOp[D](compute.op, a, b)
}

func (compute matrixEWiseAddBinaryOp[D]) computeElement(row, col int) (result D, ok bool) {
	var aValue, bValue D
	var aok, bok bool
	parallel.Do(func() {
		aValue, aok = compute.A.extractElement(row, col)
	}, func() {
		bValue, bok = compute.B.extractElement(row, col)
	})
	if aok {
		if bok {
			return compute.op(aValue, bValue), true
		}
		return aValue, true
	}
	if bok {
		return bValue, true
	}
	return
}

func (compute matrixEWiseAddBinaryOp[D]) computePipeline() *pipeline.Pipeline[any] {
	ap := compute.A.getPipeline()
	bp := compute.B.getPipeline()
	if ap == nil {
		return bp
	}
	if bp == nil {
		return ap
	}
	return makeMatrix2SourcePipeline(ap, bp,
		func(_, _ int, aValue D, aok bool, bValue D, bok bool) (result D, ok bool) {
			if aok {
				if bok {
					return compute.op(aValue, bValue), true
				}
				return aValue, true
			}
			return bValue, bok
		})
}

func (compute matrixEWiseAddBinaryOp[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	ap := compute.A.getRowPipeline(row)
	bp := compute.B.getRowPipeline(row)
	if ap == nil {
		return bp
	}
	if bp == nil {
		return ap
	}
	return makeVector2SourcePipeline(ap, bp,
		func(_ int, aValue D, aok bool, bValue D, bok bool) (result D, ok bool) {
			if aok {
				if bok {
					return compute.op(aValue, bValue), true
				}
				return aValue, true
			}
			return bValue, bok
		})
}

func (compute matrixEWiseAddBinaryOp[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	ap := compute.A.getColPipeline(col)
	bp := compute.B.getColPipeline(col)
	if ap == nil {
		return bp
	}
	if bp == nil {
		return ap
	}
	return makeVector2SourcePipeline(ap, bp,
		func(_ int, aValue D, aok bool, bValue D, bok bool) (result D, ok bool) {
			if aok {
				if bok {
					return compute.op(aValue, bValue), true
				}
				return aValue, true
			}
			return bValue, bok
		})
}

func (compute matrixEWiseAddBinaryOp[D]) mergePipelines(aps, bps []matrix1Pipeline) (result []matrix1Pipeline) {
	for {
		if len(aps) == 0 {
			return append(result, bps...)
		}
		if len(bps) == 0 {
			return append(result, aps...)
		}
		if aps[0].index == bps[0].index {
			result = append(result, matrix1Pipeline{
				index: aps[0].index,
				p: makeVector2SourcePipeline(aps[0].p, bps[0].p,
					func(_ int, aValue D, aok bool, bValue D, bok bool) (result D, ok bool) {
						if aok {
							if bok {
								return compute.op(aValue, bValue), true
							}
							return aValue, true
						}
						return bValue, bok
					}),
			})
		} else if aps[0].index < bps[0].index {
			result = append(result, aps[0])
			aps = aps[1:]
		} else {
			result = append(result, bps[0])
			bps = bps[1:]
		}
	}
}

func (compute matrixEWiseAddBinaryOp[D]) computeRowPipelines() []matrix1Pipeline {
	return compute.mergePipelines(compute.A.getRowPipelines(), compute.B.getRowPipelines())
}

func (compute matrixEWiseAddBinaryOp[D]) computeColPipelines() []matrix1Pipeline {
	return compute.mergePipelines(compute.A.getColPipelines(), compute.B.getColPipelines())
}
