package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"sort"
	"sync"
)

type matrixExtract[D any] struct {
	A                        *matrixReference[D]
	rowIndices, colIndices   []int
	index                    func(int, int) (int, int, bool)
	rowIndex, colIndex       func(int) (int, bool)
	rowSearcher, colSearcher intSearcher
}

func newMatrixExtract[D any](A *matrixReference[D], rowIndices, colIndices []int) matrixExtract[D] {
	var index func(i, j int) (int, int, bool)
	var rowIndex, colIndex func(int) (int, bool)
	var rowSearcher, colSearcher intSearcher
	nrows, allRows := isAll(rowIndices)
	if allRows {
		rowIndex = func(row int) (int, bool) {
			return row, row < nrows
		}
	} else {
		rowIndex = func(row int) (int, bool) {
			if row < len(rowIndices) {
				return rowIndices[row], true
			}
			return -1, false
		}
	}
	ncols, allCols := isAll(colIndices)
	if allCols {
		colIndex = func(col int) (int, bool) {
			return col, col < ncols
		}
	} else {
		colIndex = func(col int) (int, bool) {
			if col < len(colIndices) {
				return colIndices[col], true
			}
			return -1, false
		}
	}
	if allRows {
		if allCols {
			index = func(i, j int) (int, int, bool) {
				return i, j, i < nrows && j < ncols
			}
		} else {
			index = func(i, j int) (int, int, bool) {
				if j < len(colIndices) {
					return i, colIndices[j], i < nrows
				}
				return -1, -1, false
			}
			colSearcher = newIntSearcher(colIndices)
		}
	} else {
		parallel.Do(func() {
			if allCols {
				index = func(i, j int) (int, int, bool) {
					if i < len(rowIndices) {
						return rowIndices[i], j, j < ncols
					}
					return -1, -1, false
				}
			} else {
				index = func(i, j int) (int, int, bool) {
					if i < len(rowIndices) && j < len(colIndices) {
						return rowIndices[i], colIndices[j], true
					}
					return -1, -1, false
				}
				colSearcher = newIntSearcher(colIndices)
			}
		}, func() {
			rowSearcher = newIntSearcher(rowIndices)
		})
	}
	return matrixExtract[D]{
		A:           A,
		rowIndices:  rowIndices,
		colIndices:  colIndices,
		index:       index,
		rowIndex:    rowIndex,
		colIndex:    colIndex,
		rowSearcher: rowSearcher,
		colSearcher: colSearcher,
	}
}

func (compute matrixExtract[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	var newRowIndices, newColIndices []int
	if n, all := isAll(compute.rowIndices); newNRows < n {
		if all {
			newRowIndices = All(newNRows)
		} else {
			newRowIndices = compute.rowIndices[:n]
		}
	} else {
		newRowIndices = compute.rowIndices
	}
	if n, all := isAll(compute.colIndices); newNCols < n {
		if all {
			newColIndices = All(newNCols)
		} else {
			newColIndices = compute.colIndices[:n]
		}
	} else {
		newColIndices = compute.colIndices
	}
	return newMatrixExtract[D](compute.A, newRowIndices, newColIndices)
}

func (compute matrixExtract[D]) computeElement(row, col int) (result D, ok bool) {
	if i, j, k := compute.index(row, col); k {
		return compute.A.extractElement(i, j)
	}
	return
}

func (compute matrixExtract[D]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	if nrows, allRows := isAll(compute.rowIndices); allRows {
		if ncols, allCols := isAll(compute.colIndices); allCols {
			addMatrixAssignPipeline[D](p, func(row int) (int, bool) {
				return row, row < nrows
			}, func(col int) (int, bool) {
				return col, col < ncols
			})
			return p
		}
		addMatrixAssignPipeline[D](p, func(row int) (int, bool) {
			return row, row < nrows
		}, compute.colSearcher.search)
		if compute.colSearcher.valuesAreSorted() {
			return p
		}
	} else if compute.rowSearcher.valuesAreSorted() {
		if ncols, allCols := isAll(compute.colIndices); allCols {
			addMatrixAssignPipeline[D](p,
				compute.rowSearcher.search,
				func(col int) (int, bool) {
					return col, col < ncols
				})
			return p
		}
		addMatrixAssignPipeline[D](p,
			compute.rowSearcher.search,
			compute.colSearcher.search,
		)
		if compute.colSearcher.valuesAreSorted() {
			return p
		}
	} else if ncols, allCols := isAll(compute.colIndices); allCols {
		addMatrixAssignPipeline[D](p,
			compute.rowSearcher.search,
			func(col int) (int, bool) {
				return col, col < ncols
			})
	} else {
		addMatrixAssignPipeline[D](p,
			compute.rowSearcher.search,
			compute.colSearcher.search,
		)
	}
	var result matrixSlice[D]
	var wg sync.WaitGroup
	wg.Add(1)
	var np pipeline.Pipeline[any]
	np.Source(matrixSourceWithWaitGroup(&wg, &result.rows, &result.cols, &result.values))
	np.Notify(func() {
		defer wg.Done()
		result.collect(p)
		matrixSort(result.rows, result.cols, result.values)
	})
	return &np
}

func (_ matrixExtract[D]) addVectorPipeline(p *pipeline.Pipeline[any], vectorIndices []int, searcher intSearcher) *pipeline.Pipeline[any] {
	if n, all := isAll(vectorIndices); all {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					return index, value, index < n
				})
				return slice
			})),
		)
		return p
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
				newIndex, ok = searcher.search(index)
				newValue = value
				return
			})
			return slice
		})),
	)
	if searcher.valuesAreSorted() {
		return p
	}
	var result vectorSlice[D]
	var wg sync.WaitGroup
	wg.Add(1)
	var np pipeline.Pipeline[any]
	np.Source(vectorSourceWithWaitGroup(&wg, &result.indices, &result.values))
	np.Notify(func() {
		defer wg.Done()
		result.collect(p)
		vectorSort(result.indices, result.values)
	})
	return &np
}

func (compute matrixExtract[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	if accessRow, ok := compute.rowIndex(row); ok {
		p := compute.A.getRowPipeline(accessRow)
		if p == nil {
			return nil
		}
		return compute.addVectorPipeline(p, compute.colIndices, compute.colSearcher)
	}
	return nil
}

func (compute matrixExtract[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	if accessCol, ok := compute.colIndex(col); ok {
		p := compute.A.getColPipeline(accessCol)
		if p == nil {
			return nil
		}
		return compute.addVectorPipeline(p, compute.rowIndices, compute.rowSearcher)
	}
	return nil
}

func (compute matrixExtract[D]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	if n, all := isAll(compute.rowIndices); all {
		target := 0
		for _, p := range ps {
			if p.index < n {
				ps[target].p = compute.addVectorPipeline(p.p, compute.colIndices, compute.colSearcher)
				target++
			}
		}
		ps = ps[:target]
		return ps
	}
	target := 0
	for _, p := range ps {
		if dstIndex, ok := compute.rowSearcher.search(p.index); ok {
			ps[target].index = dstIndex
			ps[target].p = compute.addVectorPipeline(p.p, compute.colIndices, compute.colSearcher)
			target++
		}
	}
	ps = ps[:target]
	if compute.rowSearcher.valuesAreSorted() {
		return ps
	}
	// todo: parallel?
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].index < ps[j].index
	})
	return ps
}

func (compute matrixExtract[D]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	if n, all := isAll(compute.colIndices); all {
		target := 0
		for _, p := range ps {
			if p.index < n {
				ps[target].p = compute.addVectorPipeline(p.p, compute.rowIndices, compute.rowSearcher)
				target++
			}
		}
		ps = ps[:target]
		return ps
	}
	target := 0
	for _, p := range ps {
		if dstIndex, ok := compute.colSearcher.search(p.index); ok {
			ps[target].index = dstIndex
			ps[target].p = compute.addVectorPipeline(p.p, compute.rowIndices, compute.rowSearcher)
			target++
		}
	}
	ps = ps[:target]
	if compute.colSearcher.valuesAreSorted() {
		return ps
	}
	// todo: parallel?
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].index < ps[j].index
	})
	return ps
}
