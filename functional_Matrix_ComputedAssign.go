package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"github.com/intel/forGoParallel/psort"
	"sort"
	"sync"
)

type matrixAssign[D any] struct {
	A                      *matrixReference[D]
	rowIndices, colIndices []int
	index                  func(rj, cj int) (ri, ci int, ok bool)
	rowIndex, colIndex     func(j int) (i int, ok bool)
}

func matrixIndexLookup(rowIndices, colIndices []int) (matrixIndex func(int, int) (int, int, bool), rowIndex, colIndex func(int) (int, bool)) {
	if nRowIndices, allRows := isAll(rowIndices); allRows {
		if nColIndices, allCols := isAll(colIndices); allCols {
			return func(i, j int) (int, int, bool) {
					return i, j, i < nRowIndices && j < nColIndices
				}, func(i int) (int, bool) {
					return i, i < nRowIndices
				}, func(j int) (int, bool) {
					return j, j < nColIndices
				}
		}
		colSearcher := newIntSearcher(colIndices)
		return func(rj, cj int) (int, int, bool) {
				if rj >= nRowIndices {
					return 0, 0, false
				}
				ci, ok := colSearcher.search(cj)
				return rj, ci, ok
			}, func(rj int) (int, bool) {
				return rj, rj < nRowIndices
			}, colSearcher.search
	}
	if nColIndices, allCols := isAll(colIndices); allCols {
		rowSearcher := newIntSearcher(rowIndices)
		return func(rj, cj int) (int, int, bool) {
				if cj >= nColIndices {
					return 0, 0, false
				}
				ri, ok := rowSearcher.search(rj)
				return ri, cj, ok
			}, rowSearcher.search,
			func(cj int) (int, bool) {
				return cj, cj < nColIndices
			}
	}
	var rowSearcher, colSearcher intSearcher
	parallel.Do(func() {
		rowSearcher = newIntSearcher(rowIndices)
	}, func() {
		colSearcher = newIntSearcher(colIndices)
	})
	return func(rj, cj int) (int, int, bool) {
		var ri, ci int
		var rok, cok bool
		parallel.Do(func() {
			ri, rok = rowSearcher.search(rj)
		}, func() {
			ci, cok = colSearcher.search(cj)
		})
		return ri, ci, rok && cok
	}, rowSearcher.search, colSearcher.search
}

func newMatrixAssign[D any](A *matrixReference[D], rowIndices, colIndices []int) computeMatrixT[D] {
	index, rowIndex, colIndex := matrixIndexLookup(rowIndices, colIndices)
	return matrixAssign[D]{
		A:          A,
		rowIndices: rowIndices,
		colIndices: colIndices,
		index:      index,
		rowIndex:   rowIndex,
		colIndex:   colIndex,
	}
}

func (compute matrixAssign[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	var newRowIndices, newColIndices []int
	parallel.Do(func() {
		newRowIndices = resizeAssignIndices(newNRows, compute.rowIndices)
	}, func() {
		newColIndices = resizeAssignIndices(newNCols, compute.colIndices)
	})
	return newMatrixAssign[D](compute.A, newRowIndices, newColIndices)
}

func (compute matrixAssign[D]) assignIndex(row, col int) (int, int, bool) {
	return compute.index(row, col)
}

func (compute matrixAssign[D]) computeElement(row, col int) (result D, ok bool) {
	return compute.A.extractElement(row, col)
}

func addMatrixAssignPipeline[D any](p *pipeline.Pipeline[any], rowIndex, colIndex func(int) (int, bool)) {
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[D])
			slice.filter(func(row, col int, value D) (newRow, newCol int, newValue D, ok bool) {
				if newRow, ok = rowIndex(row); ok {
					newCol, ok = colIndex(col)
					newValue = value
				}
				return
			})
			return slice
		})),
	)
}

func (compute matrixAssign[D]) computePipeline() *pipeline.Pipeline[any] {
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
		}, func(col int) (int, bool) {
			if col < len(compute.colIndices) {
				dstColIndex := compute.colIndices[col]
				return dstColIndex, dstColIndex >= 0
			}
			return -1, false
		})
		if psort.IntsAreSorted(compute.colIndices) {
			return p
		}
	} else if psort.IntsAreSorted(compute.rowIndices) {
		if ncols, allCols := isAll(compute.colIndices); allCols {
			addMatrixAssignPipeline[D](p, func(row int) (int, bool) {
				if row < len(compute.rowIndices) {
					dstRowIndex := compute.rowIndices[row]
					return dstRowIndex, dstRowIndex >= 0
				}
				return -1, false
			}, func(col int) (int, bool) {
				return col, col < ncols
			})
			return p
		}
		addMatrixAssignPipeline[D](p, func(row int) (int, bool) {
			if row < len(compute.rowIndices) {
				dstRowIndex := compute.rowIndices[row]
				return dstRowIndex, dstRowIndex >= 0
			}
			return -1, false
		}, func(col int) (int, bool) {
			if col < len(compute.colIndices) {
				dstColIndex := compute.colIndices[col]
				return dstColIndex, dstColIndex >= 0
			}
			return -1, false
		})
		if psort.IntsAreSorted(compute.colIndices) {
			return p
		}
	} else if ncols, allCols := isAll(compute.colIndices); allCols {
		addMatrixAssignPipeline[D](p, func(row int) (int, bool) {
			if row < len(compute.rowIndices) {
				dstRowIndex := compute.rowIndices[row]
				return dstRowIndex, dstRowIndex >= 0
			}
			return -1, false
		}, func(col int) (int, bool) {
			return col, col < ncols
		})
	} else {
		addMatrixAssignPipeline[D](p, func(row int) (int, bool) {
			if row < len(compute.rowIndices) {
				dstRowIndex := compute.rowIndices[row]
				return dstRowIndex, dstRowIndex >= 0
			}
			return -1, false
		}, func(col int) (int, bool) {
			if col < len(compute.colIndices) {
				dstColIndex := compute.colIndices[col]
				return dstColIndex, dstColIndex >= 0
			}
			return -1, false
		})
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

func (_ matrixAssign[D]) addVectorPipeline(p *pipeline.Pipeline[any], vectorIndices []int) *pipeline.Pipeline[any] {
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
				if index < len(vectorIndices) {
					newIndex = vectorIndices[index]
					newValue = value
					ok = newIndex >= 0
				}
				return
			})
			return slice
		})),
	)
	if psort.IntsAreSorted(vectorIndices) {
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

func (compute matrixAssign[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	if accessRow, ok := compute.rowIndex(row); ok {
		p := compute.A.getRowPipeline(accessRow)
		if p == nil {
			return nil
		}
		return compute.addVectorPipeline(p, compute.colIndices)
	}
	return nil
}

func (compute matrixAssign[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	if accessCol, ok := compute.colIndex(col); ok {
		p := compute.A.getColPipeline(accessCol)
		if p == nil {
			return nil
		}
		return compute.addVectorPipeline(p, compute.rowIndices)
	}
	return nil
}

func (compute matrixAssign[D]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	if n, all := isAll(compute.rowIndices); all {
		target := 0
		for _, p := range ps {
			if p.index < n {
				ps[target].p = compute.addVectorPipeline(p.p, compute.colIndices)
				target++
			}
		}
		ps = ps[:target]
		return ps
	}
	target := 0
	for _, p := range ps {
		if p.index < len(compute.rowIndices) {
			if dstIndex := compute.rowIndices[p.index]; dstIndex >= 0 {
				ps[target].index = dstIndex
				ps[target].p = compute.addVectorPipeline(p.p, compute.colIndices)
				target++
			}
		}
	}
	ps = ps[:target]
	if psort.IntsAreSorted(compute.rowIndices) {
		return ps
	}
	// todo: parallel?
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].index < ps[j].index
	})
	return ps
}

func (compute matrixAssign[D]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	if n, all := isAll(compute.colIndices); all {
		target := 0
		for _, p := range ps {
			if p.index < n {
				ps[target].p = compute.addVectorPipeline(p.p, compute.rowIndices)
				target++
			}
		}
		ps = ps[:target]
		return ps
	}
	target := 0
	for _, p := range ps {
		if p.index < len(compute.colIndices) {
			if dstIndex := compute.colIndices[p.index]; dstIndex >= 0 {
				ps[target].index = dstIndex
				ps[target].p = compute.addVectorPipeline(p.p, compute.rowIndices)
				target++
			}
		}
	}
	ps = ps[:target]
	if psort.IntsAreSorted(compute.colIndices) {
		return ps
	}
	// todo: parallel?
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].index < ps[j].index
	})
	return ps
}

type colAssign[D any] struct {
	u       *vectorReference[D]
	indices []int
	index   func(j int) (i int, ok bool)
	col     int
}

func newColAssign[D any](u *vectorReference[D], rowIndices []int, col int) computeMatrixT[D] {
	return colAssign[D]{
		u:       u,
		indices: rowIndices,
		index:   vectorIndexLookup(rowIndices),
		col:     col,
	}
}

func (compute colAssign[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	if compute.col >= newNCols {
		return newEmptyComputedMatrix[D]()
	}
	return newColAssign(compute.u, resizeAssignIndices(newNRows, compute.indices), compute.col)
}

func (compute colAssign[D]) assignIndex(row, col int) (int, int, bool) {
	if col != compute.col {
		return -1, -1, false
	}
	rowIndex, ok := compute.index(row)
	return rowIndex, col, ok
}

func (compute colAssign[D]) computeElement(row, _ int) (result D, ok bool) {
	return compute.u.extractElement(row)
}

func (compute colAssign[D]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	if n, all := isAll(compute.indices); all {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				var result matrixSlice[D]
				for i, index := range slice.indices {
					if index < n {
						result.rows = append(result.rows, index)
						result.cols = append(result.cols, compute.col)
						result.values = append(result.values, slice.values[i])
					}
				}
				return result
			})),
		)
		return p
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			var result matrixSlice[D]
			for i, index := range slice.indices {
				if index < len(compute.indices) {
					if dstIndex := compute.indices[index]; dstIndex >= 0 {
						result.rows = append(result.rows, dstIndex)
						result.cols = append(result.cols, compute.col)
						result.values = append(result.values, slice.values[i])
					}
				}
			}
			return result
		})),
	)
	if psort.IntsAreSorted(compute.indices) {
		return p
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

func (compute colAssign[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	rowIndex, rowOk := compute.index(row)
	if rowOk {
		if value, valueOk := compute.u.extractElement(rowIndex); valueOk {
			var p pipeline.Pipeline[any]
			p.Source(vectorSource([]int{compute.col}, []D{value}))
			return &p
		}
	}
	return nil
}

func (compute colAssign[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	if col == compute.col {
		p := compute.u.getPipeline()
		if p == nil {
			return nil
		}
		if n, all := isAll(compute.indices); all {
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
					if index < len(compute.indices) {
						newIndex = compute.indices[index]
						newValue = value
						ok = newIndex >= 0
					}
					return
				})
				return slice
			})),
		)
		if psort.IntsAreSorted(compute.indices) {
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
	return nil
}

func (compute colAssign[D]) computeRowPipelines() (ps []matrix1Pipeline) {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	var result vectorSlice[D]
	result.collect(p)
	if n, all := isAll(compute.indices); all {
		for i, index := range result.indices {
			if index < n {
				var p pipeline.Pipeline[any]
				p.Source(vectorSource([]int{compute.col}, []D{result.values[i]}))
				ps = append(ps, matrix1Pipeline{
					index: index,
					p:     &p,
				})
			}
		}
		return
	}
	for i, index := range result.indices {
		if index < len(compute.indices) {
			if dstIndex := compute.indices[index]; dstIndex >= 0 {
				var p pipeline.Pipeline[any]
				p.Source(vectorSource([]int{compute.col}, []D{result.values[i]}))
				ps = append(ps, matrix1Pipeline{
					index: dstIndex,
					p:     &p,
				})
			}
		}
	}
	if psort.IntsAreSorted(compute.indices) {
		return
	}
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].index < ps[j].index
	})
	return
}

func (compute colAssign[D]) computeColPipelines() []matrix1Pipeline {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	var np *pipeline.Pipeline[any]
	if n, all := isAll(compute.indices); all {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					return index, value, index < n
				})
				return slice
			})),
		)
		np = p
	} else {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					if index < len(compute.indices) {
						newIndex = compute.indices[index]
						newValue = value
						ok = newIndex >= 0
					}
					return
				})
				return slice
			})),
		)
		if psort.IntsAreSorted(compute.indices) {
			np = p
		} else {
			var result vectorSlice[D]
			var wg sync.WaitGroup
			wg.Add(1)
			np = new(pipeline.Pipeline[any])
			np.Source(vectorSourceWithWaitGroup(&wg, &result.indices, &result.values))
			np.Notify(func() {
				defer wg.Done()
				result.collect(p)
				vectorSort(result.indices, result.values)
			})
		}
	}
	return []matrix1Pipeline{{
		index: compute.col,
		p:     np,
	}}
}

type rowAssign[D any] struct {
	u       *vectorReference[D]
	indices []int
	index   func(j int) (i int, ok bool)
	row     int
}

func newRowAssign[D any](u *vectorReference[D], row int, colIndices []int) computeMatrixT[D] {
	return rowAssign[D]{
		u:       u,
		indices: colIndices,
		index:   vectorIndexLookup(colIndices),
		row:     row,
	}
}

func (compute rowAssign[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	if compute.row >= newNRows {
		return newEmptyComputedMatrix[D]()
	}
	return newColAssign(compute.u, resizeAssignIndices(newNCols, compute.indices), compute.row)
}

func (compute rowAssign[D]) assignIndex(row, col int) (int, int, bool) {
	if row != compute.row {
		return -1, -1, false
	}
	colIndex, ok := compute.index(col)
	return row, colIndex, ok
}

func (compute rowAssign[D]) computeElement(_, col int) (result D, ok bool) {
	return compute.u.extractElement(col)
}

func (compute rowAssign[D]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	if n, all := isAll(compute.indices); all {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				var result matrixSlice[D]
				for i, index := range slice.indices {
					if index < n {
						result.rows = append(result.rows, compute.row)
						result.cols = append(result.cols, index)
						result.values = append(result.values, slice.values[i])
					}
				}
				return result
			})),
		)
		return p
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			var result matrixSlice[D]
			for i, index := range slice.indices {
				if index < len(compute.indices) {
					if dstIndex := compute.indices[index]; dstIndex >= 0 {
						result.rows = append(result.rows, compute.row)
						result.cols = append(result.cols, dstIndex)
						result.values = append(result.values, slice.values[i])
					}
				}
			}
			return result
		})),
	)
	if psort.IntsAreSorted(compute.indices) {
		return p
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

func (compute rowAssign[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	colIndex, colOk := compute.index(col)
	if colOk {
		if value, valueOk := compute.u.extractElement(colIndex); valueOk {
			var p pipeline.Pipeline[any]
			p.Source(vectorSource([]int{compute.row}, []D{value}))
			return &p
		}
	}
	return nil
}

func (compute rowAssign[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	if row == compute.row {
		p := compute.u.getPipeline()
		if p == nil {
			return nil
		}
		if n, all := isAll(compute.indices); all {
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
					if index < len(compute.indices) {
						newIndex = compute.indices[index]
						newValue = value
						ok = newIndex >= 0
						return
					}
					return
				})
				return slice
			})),
		)
		if psort.IntsAreSorted(compute.indices) {
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
	return nil
}

func (compute rowAssign[D]) computeColPipelines() (ps []matrix1Pipeline) {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	var result vectorSlice[D]
	result.collect(p)
	if n, all := isAll(compute.indices); all {
		for i, index := range result.indices {
			if index < n {
				var p pipeline.Pipeline[any]
				p.Source(vectorSource([]int{compute.row}, []D{result.values[i]}))
				ps = append(ps, matrix1Pipeline{
					index: index,
					p:     &p,
				})
			}
		}
		return
	}
	for i, index := range result.indices {
		if index < len(compute.indices) {
			if dstIndex := compute.indices[index]; dstIndex >= 0 {
				var p pipeline.Pipeline[any]
				p.Source(vectorSource([]int{compute.row}, []D{result.values[i]}))
				ps = append(ps, matrix1Pipeline{
					index: dstIndex,
					p:     &p,
				})
			}
		}
	}
	if psort.IntsAreSorted(compute.indices) {
		return
	}
	sort.SliceStable(ps, func(i, j int) bool {
		return ps[i].index < ps[j].index
	})
	return
}

func (compute rowAssign[D]) computeRowPipelines() []matrix1Pipeline {
	p := compute.u.getPipeline()
	if p == nil {
		return nil
	}
	var np *pipeline.Pipeline[any]
	if n, all := isAll(compute.indices); all {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					return index, value, index < n
				})
				return slice
			})),
		)
		np = p
	} else {
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				slice := data.(vectorSlice[D])
				slice.filter(func(index int, value D) (newIndex int, newValue D, ok bool) {
					if index < len(compute.indices) {
						newIndex = compute.indices[index]
						newValue = value
						ok = newIndex >= 0
					}
					return
				})
				return slice
			})),
		)
		if psort.IntsAreSorted(compute.indices) {
			np = p
		} else {
			var result vectorSlice[D]
			var wg sync.WaitGroup
			wg.Add(1)
			np = new(pipeline.Pipeline[any])
			np.Source(vectorSourceWithWaitGroup(&wg, &result.indices, &result.values))
			np.Notify(func() {
				defer wg.Done()
				result.collect(p)
				vectorSort(result.indices, result.values)
			})
		}
	}
	return []matrix1Pipeline{{
		index: compute.row,
		p:     np,
	}}
}

type matrixAssignConstant[D any] struct {
	value                        D
	rowIndices, colIndices       []int
	rowIndexValid, colIndexValid func(int) bool
}

func newMatrixAssignConstant[D any](value D, rowIndices, colIndices []int) computeMatrixT[D] {
	return matrixAssignConstant[D]{
		value:         value,
		rowIndices:    rowIndices,
		colIndices:    colIndices,
		rowIndexValid: computeIndexValid(rowIndices),
		colIndexValid: computeIndexValid(colIndices),
	}
}

func (compute matrixAssignConstant[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	var newRowIndices, newColIndices []int
	parallel.Do(func() {
		newRowIndices = resizeAssignIndices(newNRows, compute.rowIndices)
	}, func() {
		newColIndices = resizeAssignIndices(newNCols, compute.colIndices)
	})
	return newMatrixAssignConstant[D](compute.value, newRowIndices, newColIndices)
}

func (compute matrixAssignConstant[D]) assignIndex(row, col int) (int, int, bool) {
	return row, col, compute.rowIndexValid(row) && compute.colIndexValid(col)
}

func (compute matrixAssignConstant[D]) computeElement(_, _ int) (result D, ok bool) {
	return compute.value, true
}

func newMatrixAssignConstantPipeline[D any](value D, rowIndices, colIndices []int) *pipeline.Pipeline[any] {
	var values []D
	if nrows, allRows := isAll(rowIndices); allRows {
		if ncols, allCols := isAll(colIndices); allCols {
			var p pipeline.Pipeline[any]
			index := 0
			total := nrows * ncols
			p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
				var result matrixSlice[D]
				if index >= total {
					return result, 0, nil
				}
				if index+size > total {
					size = total - index
				}
				if size < len(values) {
					values = values[:size]
				} else {
					for len(values) < size {
						values = append(values, value)
					}
				}
				result.cow = cowv
				result.rows = make([]int, size)
				result.cols = make([]int, size)
				result.values = values
				for i := index; i < index+size; i++ {
					row, col := indexToCoord(i, nrows, ncols)
					result.rows[i-index] = row
					result.cols[i-index] = col
				}
				index += size
				return result, size, nil
			}))
			return &p
		}
		var p pipeline.Pipeline[any]
		index := 0
		total := nrows * len(colIndices)
		p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
			var result matrixSlice[D]
			if index >= total {
				return result, 0, nil
			}
			if index+size > total {
				size = total - index
			}
			if size < len(values) {
				values = values[:size]
			} else {
				for len(values) < size {
					values = append(values, value)
				}
			}
			result.cow = cowv
			result.rows = make([]int, size)
			result.cols = make([]int, size)
			result.values = values
			for i := index; i < index+size; i++ {
				result.rows[i] = i / len(colIndices)
				result.cols[i] = colIndices[i%len(colIndices)]
			}
			index += size
			return result, size, nil
		}))
		return &p
	}
	if ncols, allCols := isAll(colIndices); allCols {
		var p pipeline.Pipeline[any]
		index := 0
		total := len(rowIndices) * ncols
		p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
			var result matrixSlice[D]
			if index >= total {
				return result, 0, nil
			}
			if index+size > total {
				size = total - index
			}
			if size < len(values) {
				values = values[:size]
			} else {
				for len(values) < size {
					values = append(values, value)
				}
			}
			result.cow = cowv
			result.rows = make([]int, size)
			result.cols = make([]int, size)
			result.values = values
			for i := index; i < index+size; i++ {
				result.rows[i] = rowIndices[i/ncols]
				result.cols[i] = i % ncols
			}
			index += size
			return result, size, nil
		}))
		return &p
	}
	var p pipeline.Pipeline[any]
	index := 0
	total := len(rowIndices) * len(colIndices)
	p.Source(pipeline.NewFunc[any](total, func(size int) (data any, fetched int, err error) {
		var result matrixSlice[D]
		if index >= total {
			return result, 0, nil
		}
		if index+size > total {
			size = total - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, value)
			}
		}
		result.cow = cowv
		result.rows = make([]int, size)
		result.cols = make([]int, size)
		result.values = values
		for i := index; i < index+size; i++ {
			result.rows[i] = rowIndices[i/len(colIndices)]
			result.cols[i] = colIndices[i%len(colIndices)]
		}
		index += size
		return result, size, nil
	}))
	return &p
}

func (compute matrixAssignConstant[D]) computePipeline() *pipeline.Pipeline[any] {
	return newMatrixAssignConstantPipeline(compute.value, compute.rowIndices, compute.colIndices)
}

func newMatrixAssignConstant1DimPipeline[D any](value D, colIndices []int) *pipeline.Pipeline[any] {
	var values []D
	if n, all := isAll(colIndices); all {
		var p pipeline.Pipeline[any]
		index := 0
		p.Source(pipeline.NewFunc[any](n, func(size int) (data any, fetched int, err error) {
			var result vectorSlice[D]
			if index >= n {
				return result, 0, nil
			}
			if index+size > n {
				size = n - index
			}
			if size < len(values) {
				values = values[:size]
			} else {
				for len(values) < size {
					values = append(values, value)
				}
			}
			result.cow = cowv
			result.indices = make([]int, size)
			for i := range result.indices {
				result.indices[i] = index + i
			}
			result.values = values
			index += size
			return result, size, nil
		}))
		return &p
	}
	var p pipeline.Pipeline[any]
	index := 0
	p.Source(pipeline.NewFunc[any](len(colIndices), func(size int) (data any, fetched int, err error) {
		var result vectorSlice[D]
		if index >= len(colIndices) {
			return result, 0, nil
		}
		if index+size > len(colIndices) {
			size = len(colIndices) - index
		}
		if size < len(values) {
			values = values[:size]
		} else {
			for len(values) < size {
				values = append(values, value)
			}
		}
		result.cow = cow0 | cowv
		result.indices = colIndices[index : index+size]
		result.values = values
		index += size
		return result, size, nil
	}))
	return &p
}

func (compute matrixAssignConstant[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	if compute.rowIndexValid(row) {
		return newMatrixAssignConstant1DimPipeline(compute.value, compute.colIndices)
	}
	return nil
}

func newMatrixAssignConstant1DimPipelines[D any](value D, rowIndices, colIndices []int) []matrix1Pipeline {
	if n, all := isAll(rowIndices); all {
		ps := make([]matrix1Pipeline, n)
		for i := 0; i < n; i++ {
			ps[i].index = i
			ps[i].p = newMatrixAssignConstant1DimPipeline(value, colIndices)
		}
		return ps
	}
	ps := make([]matrix1Pipeline, len(rowIndices))
	for i, index := range rowIndices {
		ps[i].index = index
		ps[i].p = newMatrixAssignConstant1DimPipeline(value, colIndices)
	}
	return ps
}

func (compute matrixAssignConstant[D]) computeRowPipelines() []matrix1Pipeline {
	return newMatrixAssignConstant1DimPipelines(compute.value, compute.rowIndices, compute.colIndices)
}

func (compute matrixAssignConstant[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	if compute.colIndexValid(col) {
		return newMatrixAssignConstant1DimPipeline(compute.value, compute.rowIndices)
	}
	return nil
}

func (compute matrixAssignConstant[D]) computeColPipelines() []matrix1Pipeline {
	return newMatrixAssignConstant1DimPipelines(compute.value, compute.colIndices, compute.rowIndices)
}

type deleteMatrix[D any] struct {
	rowIndices, colIndices       []int
	rowIndexValid, colIndexValid func(int) bool
}

func newDeleteMatrix[D any](rowIndices, colIndices []int) computeMatrixT[D] {
	return deleteMatrix[D]{
		rowIndices:    rowIndices,
		colIndices:    colIndices,
		rowIndexValid: computeIndexValid(rowIndices),
		colIndexValid: computeIndexValid(colIndices),
	}
}

func (compute deleteMatrix[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	var newRowIndices, newColIndices []int
	parallel.Do(func() {
		newRowIndices = resizeAssignIndices(newNRows, compute.rowIndices)
	}, func() {
		newColIndices = resizeAssignIndices(newNCols, compute.colIndices)
	})
	return newDeleteMatrix[D](newRowIndices, newColIndices)
}

func (compute deleteMatrix[D]) assignIndex(row, col int) (int, int, bool) {
	return row, col, compute.rowIndexValid(row) && compute.colIndexValid(col)
}

func (_ deleteMatrix[D]) computeElement(_, _ int) (result D, ok bool) {
	return
}

func (_ deleteMatrix[D]) computePipeline() *pipeline.Pipeline[any] {
	return nil
}

func (_ deleteMatrix[D]) computeRowPipeline(_ int) *pipeline.Pipeline[any] {
	return nil
}

func (_ deleteMatrix[D]) computeColPipeline(_ int) *pipeline.Pipeline[any] {
	return nil
}

func (_ deleteMatrix[D]) computeRowPipelines() []matrix1Pipeline {
	return nil
}

func (_ deleteMatrix[D]) computeColPipelines() []matrix1Pipeline {
	return nil
}

type matrixAssignConstantScalar[D any] struct {
	value                        *scalarReference[D]
	rowIndices, colIndices       []int
	rowIndexValid, colIndexValid func(int) bool
}

func newMatrixAssignConstantScalar[D any](value *scalarReference[D], rowIndices, colIndices []int) computeMatrixT[D] {
	v := value.get()
	if v.optimized() {
		if val, ok := v.extractElement(value); ok {
			return newMatrixAssignConstant[D](val, rowIndices, colIndices)
		}
		return newDeleteMatrix[D](rowIndices, colIndices)
	}
	return matrixAssignConstantScalar[D]{
		value:         value,
		rowIndices:    rowIndices,
		colIndices:    colIndices,
		rowIndexValid: computeIndexValid(rowIndices),
		colIndexValid: computeIndexValid(colIndices),
	}
}

func (_ matrixAssignConstantScalar[D]) use() {}

func (compute matrixAssignConstantScalar[D]) resize(newNRows, newNCols int) computeMatrixT[D] {
	var newRowIndices, newColIndices []int
	parallel.Do(func() {
		newRowIndices = resizeAssignIndices(newNRows, compute.rowIndices)
	}, func() {
		newColIndices = resizeAssignIndices(newNCols, compute.colIndices)
	})
	return newMatrixAssignConstantScalar(compute.value, newRowIndices, newColIndices)
}

func (compute matrixAssignConstantScalar[D]) assignIndex(row, col int) (int, int, bool) {
	return row, col, compute.rowIndexValid(row) && compute.colIndexValid(col)
}

func (compute matrixAssignConstantScalar[D]) computeElement(_, _ int) (result D, ok bool) {
	return compute.value.extractElement()
}

func (compute matrixAssignConstantScalar[D]) computePipeline() *pipeline.Pipeline[any] {
	if s, sok := compute.value.extractElement(); sok {
		return newMatrixAssignConstantPipeline(s, compute.rowIndices, compute.colIndices)
	}
	return nil
}

func (compute matrixAssignConstantScalar[D]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	if compute.rowIndexValid(row) {
		if s, sok := compute.value.extractElement(); sok {
			return newMatrixAssignConstant1DimPipeline(s, compute.colIndices)
		}
	}
	return nil
}

func (compute matrixAssignConstantScalar[D]) computeRowPipelines() []matrix1Pipeline {
	if s, sok := compute.value.extractElement(); sok {
		return newMatrixAssignConstant1DimPipelines(s, compute.rowIndices, compute.colIndices)
	}
	return nil
}

func (compute matrixAssignConstantScalar[D]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	if compute.colIndexValid(col) {
		if s, sok := compute.value.extractElement(); sok {
			return newMatrixAssignConstant1DimPipeline(s, compute.rowIndices)
		}
	}
	return nil
}

func (compute matrixAssignConstantScalar[D]) computeColPipelines() []matrix1Pipeline {
	if s, sok := compute.value.extractElement(); sok {
		return newMatrixAssignConstant1DimPipelines(s, compute.colIndices, compute.rowIndices)
	}
	return nil
}
