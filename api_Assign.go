package forGraphBLASGo

import "github.com/intel/forGoParallel/parallel"

func isAnyIndexOutOfBounds(indices []int, size int) bool {
	return parallel.RangeOr(0, len(indices), func(low, high int) bool {
		for i := low; i < high; i++ {
			if index := indices[i]; index < 0 || index >= size {
				return true
			}
		}
		return false
	})
}

func checkIndices(indices []int, size int, checkIndexSize func(int) error) (nindices int, all bool, err error) {
	nindices, all = isAll(indices)
	if all {
		if err = checkIndexSize(nindices); err != nil {
			return
		}
		if size < nindices {
			err = IndexOutOfBounds
		}
		return
	}
	if err = checkIndexSize(nindices); err != nil {
		return
	}
	if isAnyIndexOutOfBounds(indices, size) {
		err = IndexOutOfBounds
	}
	return
}

func vectorAssignBody[D any](
	w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], indices []int, desc Descriptor,
	checkIndexSize func(int) error,
	simpleAssign func(size int) *vectorReference[D],
	complexAssign func() computeVectorT[D],
) error {
	size, err := w.Size()
	if err != nil {
		return err
	}
	nindices, _, err := checkIndices(indices, size, checkIndexSize)
	if err != nil {
		return err
	}
	isComp, err := desc.Is(Mask, Comp)
	if err != nil {
		panic(err)
	}
	if size == nindices && mask == nil && !isComp && accum == nil {
		w.ref = simpleAssign(size)
		return nil
	}
	maskAsStructure, err := vectorMask(mask, size)
	if err != nil {
		return err
	}
	w.ref = newVectorReference[D](newComputedVector[D](
		size, w.ref,
		maskAsStructure, accum,
		complexAssign(),
		desc,
	), -1)
	return nil
}

func VectorAssign[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], u *Vector[D], indices []int, desc Descriptor) error {
	return vectorAssignBody[D](
		w, mask, accum, indices, desc,
		u.expectSize,
		func(_ int) *vectorReference[D] {
			return u.ref
		}, func() computeVectorT[D] {
			return newVectorAssign[D](u.ref, fpcopy(indices))
		})
}

func VectorAssignConstant[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], value D, indices []int, desc Descriptor) error {
	return vectorAssignBody[D](
		w, mask, accum, indices, desc,
		func(_ int) error {
			return nil
		}, func(size int) *vectorReference[D] {
			return newVectorReference[D](newHomVectorConstant[D](size, value), int64(size))
		}, func() computeVectorT[D] {
			return newVectorAssignConstant(value, fpcopy(indices))
		})
}

func VectorAssignConstantScalar[D any](w *Vector[D], mask *Vector[bool], accum BinaryOp[D, D, D], scalar *Scalar[D], indices []int, desc Descriptor) error {
	if scalar == nil || scalar.ref == nil {
		return UninitializedObject
	}
	return vectorAssignBody[D](
		w, mask, accum, indices, desc,
		func(_ int) error {
			return nil
		}, func(size int) *vectorReference[D] {
			return newVectorReference[D](newHomVectorScalar[D](size, scalar.ref), -1)
		}, func() computeVectorT[D] {
			return newVectorAssignConstantScalar[D](scalar.ref, fpcopy(indices))
		})
}

func matrixAssignBody[D any](
	C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], rowIndices, colIndices []int, desc Descriptor,
	checkRowIndexSize, checkColIndexSize func(nindices int) error,
	simpleAssign func(int, int) *matrixReference[D],
	complexAssign func() computeMatrixT[D],
) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return nil
	}
	nRowIndices, allRows, rowErr := checkIndices(rowIndices, nrows, checkRowIndexSize)
	if rowErr != nil {
		return rowErr
	}
	nColIndices, allCols, colErr := checkIndices(colIndices, ncols, checkColIndexSize)
	if colErr != nil {
		return colErr
	}
	isComp, err := desc.Is(Mask, Comp)
	if err != nil {
		panic(err)
	}
	if allRows && nRowIndices == nrows &&
		allCols && nColIndices == ncols &&
		mask == nil && !isComp && accum == nil {
		C.ref = simpleAssign(nrows, ncols)
		return nil
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref,
		maskAsStructure, accum,
		complexAssign(),
		desc,
	), -1)
	return nil
}

func MatrixAssign[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], A *Matrix[D], rowIndices, colIndices []int, desc Descriptor) error {
	ANRows, ANCols, err := A.Size()
	if err != nil {
		return err
	}
	AIsTran, err := desc.Is(Inp0, Tran)
	if err != nil {
		return err
	}
	if AIsTran {
		ANRows, ANCols = ANCols, ANRows
	}
	return matrixAssignBody(
		C, mask, accum, rowIndices, colIndices, desc,
		func(nindices int) error {
			if nindices != ANRows {
				return DimensionMismatch
			}
			return nil
		}, func(nindices int) error {
			if nindices != ANCols {
				return DimensionMismatch
			}
			return nil
		}, func(_, _ int) *matrixReference[D] {
			return maybeTran(A.ref, AIsTran)
		}, func() computeMatrixT[D] {
			rowIndicesCopy, colIndicesCopy := fpcopy2(rowIndices, colIndices)
			return newMatrixAssign[D](maybeTran(A.ref, AIsTran), rowIndicesCopy, colIndicesCopy)
		})
}

func ColAssign[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], u *Vector[D], rowIndices []int, col int, desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	if _, _, err = checkIndices(rowIndices, nrows, u.expectSize); err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newColAssign[D](u.ref, fpcopy(rowIndices), col),
		desc,
	), -1)
	return nil
}

func RowAssign[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], u *Vector[D], row int, colIndices []int, desc Descriptor) error {
	nrows, ncols, err := C.Size()
	if err != nil {
		return err
	}
	if _, _, err = checkIndices(colIndices, ncols, u.expectSize); err != nil {
		return err
	}
	maskAsStructure, err := matrixMask(mask, nrows, ncols)
	if err != nil {
		return err
	}
	C.ref = newMatrixReference[D](newComputedMatrix[D](
		nrows, ncols, C.ref, maskAsStructure, accum,
		newRowAssign[D](u.ref, row, fpcopy(colIndices)),
		desc,
	), -1)
	return nil
}

func MatrixAssignConstant[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], value D, rowIndices, colIndices []int, desc Descriptor) error {
	return matrixAssignBody(
		C, mask, accum, rowIndices, colIndices, desc,
		func(nindices int) error {
			return nil
		}, func(nindices int) error {
			return nil
		}, func(nrows, ncols int) *matrixReference[D] {
			return newMatrixReference[D](newHomMatrixConstant[D](nrows, ncols, value), int64(nrows*ncols))
		}, func() computeMatrixT[D] {
			rowIndicesCopy, colIndicesCopy := fpcopy2(rowIndices, colIndices)
			return newMatrixAssignConstant(value, rowIndicesCopy, colIndicesCopy)
		})
}

func MatrixAssignConstantScalar[D any](C *Matrix[D], mask *Matrix[bool], accum BinaryOp[D, D, D], scalar *Scalar[D], rowIndices, colIndices []int, desc Descriptor) error {
	if scalar == nil || scalar.ref == nil {
		return UninitializedObject
	}
	return matrixAssignBody(
		C, mask, accum, rowIndices, colIndices, desc,
		func(nindices int) error {
			return nil
		}, func(nindices int) error {
			return nil
		}, func(nrows, ncols int) *matrixReference[D] {
			return newMatrixReference[D](newHomMatrixScalar[D](nrows, ncols, scalar.ref), -1)
		}, func() computeMatrixT[D] {
			rowIndicesCopy, colIndicesCopy := fpcopy2(rowIndices, colIndices)
			return newMatrixAssignConstantScalar(scalar.ref, rowIndicesCopy, colIndicesCopy)
		})
}
