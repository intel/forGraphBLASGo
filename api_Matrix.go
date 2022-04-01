package forGraphBLASGo

import (
	"github.com/intel/forGoParallel/parallel"
	"github.com/intel/forGoParallel/pipeline"
	"sync/atomic"
)

/*
Matrix is the exported representation of matrices. It uses Matrix.ref as an indirection to a matrixReference
representation.

When a Matrix is side-effected, then assignments are made to the .ref pointer. For example:
	func (m *Matrix[T]) SetElement(value T, row, col int) error {
		...
		m.ref = newMatrixReference(m.ref.setElement(value, row, col))
		return nil
	}
In this example, m.ref.setElement is a side-effect-free function that creates a new matrix, which then gets assigned to
m.ref to perform the actual side effect. This is important because the original m.ref might still be used in other
contexts.

The representation underneath Matrix.ref might silently change (for example from listMatrix to csrMatrix), but
these changes do not affect the semantics. (They are side-effect-free in the functional programming sense.) These
changes are properly synchronized.

Assignments to Matrix.ref directly are not synchronized. As per GraphBLAS specification, it is the task of user programs
to take care of synchronizing actual side effects.
*/

type Matrix[T any] struct {
	ref *matrixReference[T]
}

func MatrixNew[T any](nrows, ncols int) (result *Matrix[T], err error) {
	if nrows <= 0 || ncols <= 0 {
		err = InvalidValue
		return
	}
	return &Matrix[T]{newMatrixReference[T](makeEmptyMatrix[T](nrows, ncols), 0)}, nil
}

func (m *Matrix[T]) Dup() (result *Matrix[T], err error) {
	if m == nil || m.ref == nil {
		err = UninitializedObject
		return
	}
	return &Matrix[T]{m.ref}, nil
}

func MatrixDiag[T any](v *Vector[T], k int) (result *Matrix[T], err error) {
	if v == nil || v.ref == nil {
		err = UninitializedObject
		return
	}
	size := v.ref.size() + absInt(k)
	n := atomic.LoadInt64(&v.ref.nvalues)
	return &Matrix[T]{newMatrixReference[T](newDiagonalMatrix[T](size, size, v.ref, k), n)}, nil
}

func (m *Matrix[T]) Resize(nrows, ncols int) error {
	if nrows <= 0 || ncols <= 0 {
		return InvalidValue
	}
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	m.ref = m.ref.resize(nrows, ncols)
	return nil
}

func (m *Matrix[T]) Clear() error {
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	nrows, ncols := m.ref.size()
	m.ref = newMatrixReference[T](makeEmptyMatrix[T](nrows, ncols), 0)
	return nil
}

func (m *Matrix[T]) NRows() (int, error) {
	if m == nil || m.ref == nil {
		return 0, UninitializedObject
	}
	nrows, _ := m.ref.size()
	return nrows, nil
}

func (m *Matrix[T]) NCols() (int, error) {
	if m == nil || m.ref == nil {
		return 0, UninitializedObject
	}
	_, ncols := m.ref.size()
	return ncols, nil
}

func (m *Matrix[T]) Size() (int, int, error) {
	if m == nil || m.ref == nil {
		return 0, 0, UninitializedObject
	}
	nrows, ncols := m.ref.size()
	return nrows, ncols, nil
}

func (m *Matrix[T]) NVals() (int, error) {
	if m == nil || m.ref == nil {
		return 0, UninitializedObject
	}
	return m.ref.nvals(), nil
}

func (m *Matrix[T]) Build(rows, cols []int, values []T, dup BinaryOp[T, T, T]) error {
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	if len(rows) != len(cols) || len(rows) != len(values) {
		return IndexOutOfBounds
	}
	if m.ref.nvals() > 0 {
		return OutputNotEmpty
	}
	nrows, ncols := m.ref.size()
	// todo: use speculative.RangeOr
	if parallel.RangeOr(0, len(rows), func(low, high int) bool {
		for i := low; i < high; i++ {
			if row := rows[i]; row < 0 || row >= nrows {
				return true
			}
			if col := cols[i]; col < 0 || col >= ncols {
				return true
			}
		}
		return false
	}) {
		return IndexOutOfBounds
	}
	rowCopies, colCopies, valueCopies := fpcopy3(rows, cols, values)
	if dup == nil {
		matrixSort(rowCopies, colCopies, valueCopies)
		// todo: use speculative.RangeOr
		if parallel.RangeOr(0, len(rows), func(low, high int) bool {
			for i := low; i < high-1; i++ {
				if rowCopies[i] == rowCopies[i+1] && colCopies[i] == colCopies[i+1] {
					return true
				}
			}
			return high < len(rows) && rowCopies[high-1] == rowCopies[high] && colCopies[high-1] == colCopies[high]
		}) {
			return InvalidValue
		}
		m.ref = newDelayedMatrixReference[T](func() (functionalMatrix[T], int64) {
			newRows, rowSpans := csrRows(rowCopies)
			return newCSRMatrix[T](nrows, ncols, newRows, rowSpans, colCopies, valueCopies), int64(len(valueCopies))
		})
		return nil
	}
	m.ref = newDelayedMatrixReference[T](func() (functionalMatrix[T], int64) {
		matrixSort(rowCopies, colCopies, valueCopies)
		var dups [][2]int
		var p pipeline.Pipeline[any]
		p.Source(newIntervalSource(len(valueCopies)))
		p.Add(
			pipeline.Par(pipeline.Receive(func(_ int, data any) any {
				batch := data.(interval)
				low, high := batch.start, batch.end
				var result [][2]int
				if low > 0 {
					low--
				}
				if high < len(rowCopies) {
					high++
				}
				for i := low; i < high; {
					row := rowCopies[i]
					col := colCopies[i]
					j := i + 1
					for j < high && rowCopies[j] == row && colCopies[j] == col {
						j++
					}
					if j-i > 1 {
						result = append(result, [2]int{i, j})
						i = j
					} else {
						i++
					}
				}
				return result
			})),
			//todo: we can simplify this: since we are already sequential here, maybe we can copy indices and values
			// already to their right destinations (also in VectorBuild)?
			pipeline.Ord(pipeline.Receive(func(_ int, data any) any {
				ndups := data.([][2]int)
				if len(ndups) == 0 {
					return nil
				}
				lx := len(dups)
				if lx == 0 {
					dups = ndups
					return nil
				}
				lx--
				if i, j := dups[lx][0], ndups[0][0]; rowCopies[i] == rowCopies[j] && colCopies[i] == colCopies[j] {
					ndups[0][0] = i
					if lx == 0 {
						dups = ndups
						return nil
					}
					dups = dups[:lx]
				}
				dups = append(dups, ndups...)
				return nil
			})),
		)
		p.Run()
		if err := p.Err(); err != nil {
			panic(err)
		}
		parallel.Range(0, len(dups), func(low, high int) {
			for i := low; i < high; i++ {
				dp := dups[i]
				start, end := dp[0], dp[1]
				for j := start + 1; j < end; j++ {
					valueCopies[start] = dup(valueCopies[start], valueCopies[j])
				}
			}
		})
		dups = append(dups, [2]int{len(rowCopies) - 1, len(rowCopies) - 1})
		delta := 0
		for i := 0; i < len(dups)-1; i++ {
			dstStart := dups[i][0] + 1 - delta
			srcStart := dups[i][1]
			srcEnd := dups[i+1][0] + 1
			copy(rowCopies[dstStart:], rowCopies[srcStart:srcEnd])
			copy(colCopies[dstStart:], colCopies[srcStart:srcEnd])
			copy(valueCopies[dstStart:], valueCopies[srcStart:srcEnd])
			delta += dups[i][1] - dups[i][0] - 1
		}
		rowCopies = rowCopies[:len(rowCopies)-delta]
		colCopies = colCopies[:len(colCopies)-delta]
		valueCopies = valueCopies[:len(valueCopies)-delta]
		newRows, rowSpans := csrRows(rowCopies)
		return newCSRMatrix[T](nrows, ncols, newRows, rowSpans, colCopies, valueCopies), int64(len(valueCopies))
	})
	return nil
}

func (m *Matrix[T]) SetElement(value T, row, col int) error {
	if row < 0 || col < 0 {
		return InvalidIndex
	}
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	nrows, ncols := m.ref.size()
	if row >= nrows || col >= ncols {
		return InvalidIndex
	}
	m.ref = m.ref.setElement(value, row, col)
	return nil
}

func (m *Matrix[T]) RemoveElement(row, col int) error {
	if row < 0 || col < 0 {
		return InvalidIndex
	}
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	nrows, ncols := m.ref.size()
	if row >= nrows || col >= ncols {
		return InvalidIndex
	}
	m.ref = m.ref.removeElement(row, col)
	return nil
}

func (m *Matrix[T]) ExtractElement(row, col int) (result T, err error) {
	if row < 0 || col < 0 {
		err = InvalidIndex
		return
	}
	if m == nil || m.ref == nil {
		err = UninitializedObject
		return
	}
	nrows, ncols := m.ref.size()
	if row >= nrows || col >= ncols {
		err = InvalidIndex
		return
	}
	if value, ok := m.ref.extractElement(row, col); ok {
		return value, nil
	}
	err = NoValue
	return
}

func (m *Matrix[T]) ExtractTuples() (rows, cols []int, values []T, err error) {
	if m == nil || m.ref == nil {
		err = UninitializedObject
		return
	}
	p := m.ref.getPipeline()
	if p == nil {
		atomic.StoreInt64(&m.ref.nvalues, 0)
		return
	}
	var result matrixSlice[T]
	result.collect(p)
	rows = result.rows
	cols = result.cols
	values = result.values
	atomic.StoreInt64(&m.ref.nvalues, int64(len(values)))
	return
}

func (m *Matrix[T]) ExportHint() (Format, error) {
	if m == nil || m.ref == nil {
		return 0, UninitializedObject
	}
	return CSRFormat, nil
}

func (m *Matrix[T]) ExportSize(_ Format) (int, int, int, error) {
	if m == nil || m.ref == nil {
		return 0, 0, 0, UninitializedObject
	}
	panic("todo") // todo
}

func (m *Matrix[T]) Export(_ Format, _, _, _ int) ([]int, []int, []T, error) {
	if m == nil || m.ref == nil {
		return nil, nil, nil, UninitializedObject
	}
	panic("todo") // todo
}

func MatrixImport[T any](_, _ int, _, _ []int, _ []T, _ Format) (result *Matrix[T], err error) {
	panic("todo") // todo
}

/* todo
SerialSize
Serialize
Deserialize
=> ensure compatibility with Go standard library
*/

func (m *Matrix[T]) Wait(mode WaitMode) error {
	if m == nil || m.ref == nil {
		return UninitializedObject
	}
	if mode == Complete {
		return nil
	}
	m.ref.optimize()
	return nil
}

func (m *Matrix[T]) AsMask() *Matrix[bool] {
	if m == nil || m.ref == nil {
		return nil
	}
	n := atomic.LoadInt64(&m.ref.nvalues)
	switch m := any(m).(type) {
	case *Matrix[bool]:
		return m
	case *Matrix[int8]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[int8](m.ref), n)}
	case *Matrix[int16]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[int16](m.ref), n)}
	case *Matrix[int32]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[int32](m.ref), n)}
	case *Matrix[int64]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[int64](m.ref), n)}
	case *Matrix[uint8]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[uint8](m.ref), n)}
	case *Matrix[uint16]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[uint16](m.ref), n)}
	case *Matrix[uint32]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[uint32](m.ref), n)}
	case *Matrix[uint64]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[uint64](m.ref), n)}
	case *Matrix[float32]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[float32](m.ref), n)}
	case *Matrix[float64]:
		return &Matrix[bool]{newMatrixReference[bool](newMatrixAsMask[float64](m.ref), n)}
	}
	return &Matrix[bool]{newMatrixReference[bool](newMatrixAsStructuralMask[T](m.ref), n)}
}
