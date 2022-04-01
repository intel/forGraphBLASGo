package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"runtime"
)

type matrixMult[DC, DA, DB any] struct {
	op Semiring[DC, DA, DB]
	A  *matrixReference[DA]
	B  *matrixReference[DB]
}

func newMatrixMult[DC, DA, DB any](
	op Semiring[DC, DA, DB],
	A *matrixReference[DA],
	B *matrixReference[DB],
) computeMatrixT[DC] {
	return matrixMult[DC, DA, DB]{
		op: op,
		A:  A,
		B:  B,
	}
}

func (compute matrixMult[DC, DA, DB]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	var A *matrixReference[DA]
	var B *matrixReference[DB]
	parallel.Do(func() {
		_, ncols := compute.A.size()
		A = compute.A.resize(newNRows, ncols)
	}, func() {
		nrows, _ := compute.B.size()
		B = compute.B.resize(nrows, newNCols)
	})
	return newMatrixMult[DC, DA, DB](compute.op, A, B)
}

func (compute matrixMult[DC, DA, DB]) computeElement(row, col int) (result DC, ok bool) {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	ap := compute.A.getRowPipeline(row)
	if ap == nil {
		return
	}
	bp := compute.B.getColPipeline(col)
	if bp == nil {
		return
	}
	return vectorPipelineReduce(makeVector2SourcePipeline(ap, bp,
		func(index int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
			if aok && bok {
				return mult(aValue, bValue), true
			}
			return
		},
	), add)
}

func (compute matrixMult[DC, DA, DB]) computePipeline() *pipeline.Pipeline[any] {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	ch := make(chan matrixSlice[DC], runtime.GOMAXPROCS(0))
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(_ int) (data any, fetched int, err error) {
		var result matrixSlice[DC]
		done := p.Context().Done()
		for len(result.values) == 0 {
			select {
			case <-done:
				return nil, 0, nil
			case element, ok := <-ch:
				if ok {
					result = element
				} else {
					return result, len(result.values), nil
				}
			}
		}
		return result, len(result.values), nil
	}))
	p.Notify(func() {
		done := p.Context().Done()
		var result matrixSlice[DC]
		for _, row := range compute.A.getRowPipelines() {
			colps := compute.B.getColPipelines()
			if len(colps) == 0 {
				break
			}
			value, ok := vectorPipelineReduce(makeVector2SourcePipeline(row.p, colps[0].p,
				func(_ int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
					if aok && bok {
						return mult(aValue, bValue), true
					}
					return
				},
			), add)
			if ok {
				result.rows = append(result.rows, row.index)
				result.cols = append(result.cols, colps[0].index)
				result.values = append(result.values, value)
				if len(result.values) == 512 {
					select {
					case <-done:
						return
					case ch <- result:
						result.rows = nil
						result.cols = nil
						result.values = nil
					}
				}
			}
			for _, col := range colps[1:] {
				value, ok = vectorPipelineReduce(makeVector2SourcePipeline(compute.A.getRowPipeline(row.index), col.p,
					func(_ int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
						if aok && bok {
							return mult(aValue, bValue), true
						}
						return
					},
				), add)
				if ok {
					result.rows = append(result.rows, row.index)
					result.cols = append(result.cols, col.index)
					result.values = append(result.values, value)
					if len(result.values) == 512 {
						select {
						case <-done:
							return
						case ch <- result:
							result.rows = nil
							result.cols = nil
							result.values = nil
						}
					}
				}
			}
		}
		if len(result.values) > 0 {
			select {
			case <-done:
				return
			case ch <- result:
			}
		}
		close(ch)
	})
	return &p
}

func (compute matrixMult[DC, DA, DB]) makeRowPipeline(row int, rowp *pipeline.Pipeline[any]) *pipeline.Pipeline[any] {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	ch := make(chan vectorSlice[DC], runtime.GOMAXPROCS(0))
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[DC]
		done := p.Context().Done()
		for len(result.values) == 0 {
			select {
			case <-done:
				return nil, 0, nil
			case element, ok := <-ch:
				if ok {
					result = element
				} else {
					return result, len(result.values), nil
				}
			}
		}
		return result, len(result.values), nil
	}))
	p.Notify(func() {
		done := p.Context().Done()
		var result vectorSlice[DC]
		colps := compute.B.getColPipelines()
		if len(colps) == 0 {
			close(ch)
			return
		}
		value, ok := vectorPipelineReduce(makeVector2SourcePipeline(rowp, colps[0].p,
			func(index int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
				if aok && bok {
					return mult(aValue, bValue), true
				}
				return
			}), add)
		if ok {
			result.indices = append(result.indices, colps[0].index)
			result.values = append(result.values, value)
			if len(result.values) == 512 {
				select {
				case <-done:
					return
				case ch <- result:
					result.indices = nil
					result.values = nil
				}
			}
		}
		for _, col := range colps[1:] {
			value, ok = vectorPipelineReduce(makeVector2SourcePipeline(compute.A.getRowPipeline(row), col.p,
				func(index int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
					if aok && bok {
						return mult(aValue, bValue), true
					}
					return
				}), add)
			if ok {
				result.indices = append(result.indices, col.index)
				result.values = append(result.values, value)
				if len(result.values) == 512 {
					select {
					case <-done:
						return
					case ch <- result:
						result.indices = nil
						result.values = nil
					}
				}
			}
		}
		if len(result.values) > 0 {
			select {
			case <-done:
				return
			case ch <- result:
			}
		}
		close(ch)
	})
	return &p
}

func (compute matrixMult[DC, DA, DB]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	rowp := compute.A.getRowPipeline(row)
	if rowp == nil {
		return nil
	}
	return compute.makeRowPipeline(row, rowp)
}

func (compute matrixMult[DC, DA, DB]) makeColPipeline(col int, colp *pipeline.Pipeline[any]) *pipeline.Pipeline[any] {
	add := compute.op.addition().operator()
	mult := compute.op.multiplication()
	ch := make(chan vectorSlice[DC], runtime.GOMAXPROCS(0))
	var p pipeline.Pipeline[any]
	p.Source(pipeline.NewFunc[any](-1, func(size int) (data any, fetched int, err error) {
		var result vectorSlice[DC]
		done := p.Context().Done()
		for len(result.values) == 0 {
			select {
			case <-done:
				return nil, 0, nil
			case element, ok := <-ch:
				if ok {
					result = element
				} else {
					return result, fetched, nil
				}
			}
		}
		return result, fetched, nil
	}))
	p.Notify(func() {
		done := p.Context().Done()
		var result vectorSlice[DC]
		rowps := compute.A.getRowPipelines()
		if len(rowps) == 0 {
			close(ch)
			return
		}
		value, ok := vectorPipelineReduce(makeVector2SourcePipeline(rowps[0].p, colp,
			func(index int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
				if aok && bok {
					return mult(aValue, bValue), true
				}
				return
			}), add)
		if ok {
			result.indices = append(result.indices, rowps[0].index)
			result.values = append(result.values, value)
			if len(result.values) == 512 {
				select {
				case <-done:
					return
				case ch <- result:
					result.indices = nil
					result.values = nil
				}
			}
		}
		for _, row := range rowps[1:] {
			value, ok = vectorPipelineReduce(makeVector2SourcePipeline(row.p, compute.B.getColPipeline(col),
				func(index int, aValue DA, aok bool, bValue DB, bok bool) (result DC, ok bool) {
					if aok && bok {
						return mult(aValue, bValue), true
					}
					return
				}), add)
			if ok {
				result.indices = append(result.indices, row.index)
				result.values = append(result.values, value)
				if len(result.values) == 512 {
					select {
					case <-done:
						return
					case ch <- result:
						result.indices = nil
						result.values = nil
					}
				}
			}
		}
		if len(result.values) > 0 {
			select {
			case <-done:
				return
			case ch <- result:
			}
		}
		close(ch)
	})
	return &p
}

func (compute matrixMult[DC, DA, DB]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	colp := compute.B.getColPipeline(col)
	if colp == nil {
		return nil
	}
	return compute.makeColPipeline(col, colp)
}

func (compute matrixMult[DC, DA, DB]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for i := range ps {
		ps[i].p = compute.makeRowPipeline(ps[i].index, ps[i].p)
	}
	return ps
}

func (compute matrixMult[DC, DA, DB]) computeColPipelines() []matrix1Pipeline {
	ps := compute.B.getColPipelines()
	for i := range ps {
		ps[i].p = compute.makeColPipeline(ps[i].index, ps[i].p)
	}
	return ps
}
