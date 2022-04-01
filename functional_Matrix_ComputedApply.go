package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type matrixApply[DC, DA any] struct {
	op UnaryOp[DC, DA]
	A  *matrixReference[DA]
}

func newMatrixApply[DC, DA any](op UnaryOp[DC, DA], A *matrixReference[DA]) computeMatrixT[DC] {
	return matrixApply[DC, DA]{op: op, A: A}
}

func (compute matrixApply[DC, DA]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApply[DC, DA]{op: compute.op, A: compute.A.resize(newNRows, newNCols)}
}

func (compute matrixApply[DC, DA]) computeElement(row, col int) (result DC, ok bool) {
	if value, ok := compute.A.extractElement(row, col); ok {
		return compute.op(value), true
	}
	return
}

func (compute matrixApply[DC, DA]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApply[DC, DA]) addVectorPipeline(p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value)
			}
			return result
		})),
	)
}

func (compute matrixApply[DC, DA]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApply[DC, DA]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApply[DC, DA]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

func (compute matrixApply[DC, DA]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

type matrixApplyBinaryOp1st[DC, Ds, DA any] struct {
	op    BinaryOp[DC, Ds, DA]
	value Ds
	A     *matrixReference[DA]
}

func newMatrixApplyBinaryOp1st[DC, Ds, DA any](op BinaryOp[DC, Ds, DA], value Ds, A *matrixReference[DA]) computeMatrixT[DC] {
	return matrixApplyBinaryOp1st[DC, Ds, DA]{op: op, value: value, A: A}
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApplyBinaryOp1st[DC, Ds, DA]{
		op:    compute.op,
		value: compute.value,
		A:     compute.A.resize(newNRows, newNCols),
	}
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) computeElement(row, col int) (result DC, ok bool) {
	if A, Aok := compute.A.extractElement(row, col); Aok {
		return compute.op(compute.value, A), true
	}
	return
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(compute.value, value)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) addVectorPipeline(p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(compute.value, value)
			}
			return result
		})),
	)
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

func (compute matrixApplyBinaryOp1st[DC, Ds, DA]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

type matrixApplyBinaryOp1stScalar[DC, Ds, DA any] struct {
	op    BinaryOp[DC, Ds, DA]
	value *scalarReference[Ds]
	A     *matrixReference[DA]
}

func newMatrixApplyBinaryOp1stScalar[DC, Ds, DA any](op BinaryOp[DC, Ds, DA], value *scalarReference[Ds], A *matrixReference[DA]) computeMatrixT[DC] {
	return matrixApplyBinaryOp1stScalar[DC, Ds, DA]{op: op, value: value, A: A}
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApplyBinaryOp1stScalar[DC, Ds, DA]{
		op:    compute.op,
		value: compute.value,
		A:     compute.A.resize(newNRows, newNCols),
	}
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) computeElement(row, col int) (result DC, ok bool) {
	if a, aok := compute.A.extractElement(row, col); aok {
		if s, sok := compute.value.extractElement(); sok {
			return compute.op(s, a), true
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			computeValue, ok := compute.value.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(computeValue, value)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) addVectorPipeline(p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			computeValue, ok := compute.value.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(computeValue, value)
			}
			return result
		})),
	)
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

func (compute matrixApplyBinaryOp1stScalar[DC, Ds, DA]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

type matrixApplyBinaryOp2nd[DC, DA, Ds any] struct {
	op    BinaryOp[DC, DA, Ds]
	A     *matrixReference[DA]
	value Ds
}

func newMatrixApplyBinaryOp2nd[DC, DA, Ds any](op BinaryOp[DC, DA, Ds], A *matrixReference[DA], value Ds) computeMatrixT[DC] {
	return matrixApplyBinaryOp2nd[DC, DA, Ds]{op: op, A: A, value: value}
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApplyBinaryOp2nd[DC, DA, Ds]{
		op:    compute.op,
		A:     compute.A.resize(newNRows, newNCols),
		value: compute.value,
	}
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) computeElement(row, col int) (result DC, ok bool) {
	if A, Aok := compute.A.extractElement(row, col); Aok {
		return compute.op(A, compute.value), true
	}
	return
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, compute.value)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) addVectorPipeline(p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, compute.value)
			}
			return result
		})),
	)
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

func (compute matrixApplyBinaryOp2nd[DC, DA, Ds]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

type matrixApplyBinaryOp2ndScalar[DC, DA, Ds any] struct {
	op    BinaryOp[DC, DA, Ds]
	A     *matrixReference[DA]
	value *scalarReference[Ds]
}

func newMatrixApplyBinaryOp2ndScalar[DC, DA, Ds any](op BinaryOp[DC, DA, Ds], A *matrixReference[DA], value *scalarReference[Ds]) computeMatrixT[DC] {
	return matrixApplyBinaryOp2ndScalar[DC, DA, Ds]{op: op, A: A, value: value}
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApplyBinaryOp2ndScalar[DC, DA, Ds]{
		op:    compute.op,
		A:     compute.A.resize(newNRows, newNCols),
		value: compute.value,
	}
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) computeElement(row, col int) (result DC, ok bool) {
	if a, aok := compute.A.extractElement(row, col); aok {
		if s, sok := compute.value.extractElement(); sok {
			return compute.op(a, s), true
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			computeValue, ok := compute.value.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, computeValue)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) addVectorPipeline(p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			computeValue, ok := compute.value.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, computeValue)
			}
			return result
		})),
	)
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addVectorPipeline(p)
	return p
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

func (compute matrixApplyBinaryOp2ndScalar[DC, DA, Ds]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addVectorPipeline(p.p)
	}
	return ps
}

type matrixApplyIndexOp[DC, DA, Ds any] struct {
	op IndexUnaryOp[DC, DA, Ds]
	A  *matrixReference[DA]
	s  Ds
}

func newMatrixApplyIndexOp[DC, DA, Ds any](op IndexUnaryOp[DC, DA, Ds], A *matrixReference[DA], s Ds) computeMatrixT[DC] {
	return matrixApplyIndexOp[DC, DA, Ds]{op: op, A: A, s: s}
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApplyIndexOp[DC, DA, Ds]{
		op: compute.op,
		A:  compute.A.resize(newNRows, newNCols),
		s:  compute.s,
	}
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) computeElement(row, col int) (result DC, ok bool) {
	if value, ok := compute.A.extractElement(row, col); ok {
		return compute.op(value, row, col, compute.s), true
	}
	return
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, result.rows[i], result.cols[i], compute.s)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) addRowVectorPipeline(row int, p *pipeline.Pipeline[any]) {
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, row, result.indices[i], compute.s)
			}
			return result
		})),
	)
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) addColVectorPipeline(col int, p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, result.indices[i], col, compute.s)
			}
			return result
		})),
	)
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addRowVectorPipeline(row, p)
	return p
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addColVectorPipeline(col, p)
	return p
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addRowVectorPipeline(p.index, p.p)
	}
	return ps
}

func (compute matrixApplyIndexOp[DC, DA, Ds]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addColVectorPipeline(p.index, p.p)
	}
	return ps
}

type matrixApplyIndexOpScalar[DC, DA, Ds any] struct {
	op IndexUnaryOp[DC, DA, Ds]
	A  *matrixReference[DA]
	s  *scalarReference[Ds]
}

func newMatrixApplyIndexOpScalar[DC, DA, Ds any](op IndexUnaryOp[DC, DA, Ds], A *matrixReference[DA], s *scalarReference[Ds]) computeMatrixT[DC] {
	return matrixApplyIndexOpScalar[DC, DA, Ds]{op: op, A: A, s: s}
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) resize(newNRows, newNCols int) computeMatrixT[DC] {
	return matrixApplyIndexOpScalar[DC, DA, Ds]{
		op: compute.op,
		A:  compute.A.resize(newNRows, newNCols),
		s:  compute.s,
	}
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) computeElement(row, col int) (result DC, ok bool) {
	if a, aok := compute.A.extractElement(row, col); aok {
		if s, sok := compute.s.extractElement(); sok {
			return compute.op(a, row, col, s), true
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[DA])
			result := matrixSlice[DC]{
				cow:    slice.cow &^ cowv,
				rows:   slice.rows,
				cols:   slice.cols,
				values: make([]DC, len(slice.values)),
			}
			computeS, ok := compute.s.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, result.rows[i], result.cols[i], computeS)
			}
			return result
		})),
	)
	return p
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) addRowVectorPipeline(row int, p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			computeS, ok := compute.s.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, row, result.indices[i], computeS)
			}
			return result
		})),
	)
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) addColVectorPipeline(col int, p *pipeline.Pipeline[any]) {
	if p == nil {
		return
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[DA])
			result := vectorSlice[DC]{
				cow:     slice.cow &^ cowv,
				indices: slice.indices,
				values:  make([]DC, len(slice.values)),
			}
			computeS, ok := compute.s.extractElement()
			if !ok {
				panic(EmptyObject)
			}
			for i, value := range slice.values {
				result.values[i] = compute.op(value, result.indices[i], col, computeS)
			}
			return result
		})),
	)
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	compute.addRowVectorPipeline(row, p)
	return p
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	compute.addColVectorPipeline(col, p)
	return p
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addRowVectorPipeline(p.index, p.p)
	}
	return ps
}

func (compute matrixApplyIndexOpScalar[DC, DA, Ds]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addColVectorPipeline(p.index, p.p)
	}
	return ps
}
