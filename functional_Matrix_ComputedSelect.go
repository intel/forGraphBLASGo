package forGraphBLASGo

import "github.com/intel/forGoParallel/pipeline"

type matrixSelect[D, Ds any] struct {
	op    IndexUnaryOp[bool, D, Ds]
	A     *matrixReference[D]
	value Ds
}

func newMatrixSelect[D, Ds any](op IndexUnaryOp[bool, D, Ds], A *matrixReference[D], value Ds) computeMatrixT[D] {
	return matrixSelect[D, Ds]{op: op, A: A, value: value}
}

func (compute matrixSelect[D, Ds]) resize(newNRows, newNCols int) computeMatrixT[D] {
	return newMatrixSelect[D, Ds](compute.op, compute.A.resize(newNRows, newNCols), compute.value)
}

func (compute matrixSelect[D, Ds]) computeElement(row, col int) (result D, ok bool) {
	if A, Aok := compute.A.extractElement(row, col); Aok {
		if compute.op(A, row, col, compute.value) {
			return A, true
		}
	}
	return
}

func (compute matrixSelect[D, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[D])
			slice.filter(func(row, col int, value D) (newRow, newCol int, newValue D, ok bool) {
				return row, col, value, compute.op(value, row, col, compute.value)
			})
			return slice
		})),
	)
	return p
}

func (compute matrixSelect[D, Ds]) addRowPipeline(row int, p *pipeline.Pipeline[any]) {
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(col int, value D) (newCol int, newValue D, ok bool) {
				return col, value, compute.op(value, row, col, compute.value)
			})
			return slice
		})),
	)
}

func (compute matrixSelect[D, Ds]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	if p == nil {
		return nil
	}
	compute.addRowPipeline(row, p)
	return p
}

func (compute matrixSelect[D, Ds]) addColPipeline(col int, p *pipeline.Pipeline[any]) {
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(row int, value D) (newRow int, newValue D, ok bool) {
				return row, value, compute.op(value, row, col, compute.value)
			})
			return slice
		})),
	)
}

func (compute matrixSelect[D, Ds]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	if p == nil {
		return nil
	}
	compute.addColPipeline(col, p)
	return p
}

func (compute matrixSelect[D, Ds]) computeRowPipelines() []matrix1Pipeline {
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addRowPipeline(p.index, p.p)
	}
	return ps
}

func (compute matrixSelect[D, Ds]) computeColPipelines() []matrix1Pipeline {
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addColPipeline(p.index, p.p)
	}
	return ps
}

type matrixSelectScalar[D, Ds any] struct {
	op    IndexUnaryOp[bool, D, Ds]
	A     *matrixReference[D]
	value *scalarReference[Ds]
}

func newMatrixSelectScalar[D, Ds any](op IndexUnaryOp[bool, D, Ds], A *matrixReference[D], value *scalarReference[Ds]) computeMatrixT[D] {
	return matrixSelectScalar[D, Ds]{op: op, A: A, value: value}
}

func (compute matrixSelectScalar[D, Ds]) resize(newNRows, newNCols int) computeMatrixT[D] {
	return newMatrixSelectScalar[D, Ds](compute.op, compute.A.resize(newNRows, newNCols), compute.value)
}

func (compute matrixSelectScalar[D, Ds]) computeElement(row, col int) (result D, ok bool) {
	if a, aok := compute.A.extractElement(row, col); aok {
		if value, vok := compute.value.extractElement(); vok {
			if compute.op(a, row, col, value) {
				return a, true
			}
		} else {
			panic(EmptyObject)
		}
	}
	return
}

func (compute matrixSelectScalar[D, Ds]) computePipeline() *pipeline.Pipeline[any] {
	p := compute.A.getPipeline()
	if p == nil {
		return nil
	}
	v, vok := compute.value.extractElement()
	if !vok {
		panic(EmptyObject)
	}
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(matrixSlice[D])
			slice.filter(func(row, col int, value D) (newRow, newCol int, newValue D, ok bool) {
				return row, col, value, compute.op(value, row, col, v)
			})
			return slice
		})),
	)
	return p
}

func (compute matrixSelectScalar[D, Ds]) addRowPipeline(row int, p *pipeline.Pipeline[any], v Ds) {
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(col int, value D) (newCol int, newValue D, ok bool) {
				return col, value, compute.op(value, row, col, v)
			})
			return slice
		})),
	)
}

func (compute matrixSelectScalar[D, Ds]) computeRowPipeline(row int) *pipeline.Pipeline[any] {
	p := compute.A.getRowPipeline(row)
	if p == nil {
		return nil
	}
	v, vok := compute.value.extractElement()
	if !vok {
		panic(EmptyObject)
	}
	compute.addRowPipeline(row, p, v)
	return p
}

func (compute matrixSelectScalar[D, Ds]) addColPipeline(col int, p *pipeline.Pipeline[any], v Ds) {
	p.Add(
		pipeline.Par(pipeline.Receive(func(_ int, data any) any {
			slice := data.(vectorSlice[D])
			slice.filter(func(row int, value D) (newRow int, newValue D, ok bool) {
				return row, value, compute.op(value, row, col, v)
			})
			return slice
		})),
	)
}

func (compute matrixSelectScalar[D, Ds]) computeColPipeline(col int) *pipeline.Pipeline[any] {
	p := compute.A.getColPipeline(col)
	if p == nil {
		return nil
	}
	v, vok := compute.value.extractElement()
	if !vok {
		panic(EmptyObject)
	}
	compute.addColPipeline(col, p, v)
	return p
}

func (compute matrixSelectScalar[D, Ds]) computeRowPipelines() []matrix1Pipeline {
	v, vok := compute.value.extractElement()
	if !vok {
		panic(EmptyObject)
	}
	ps := compute.A.getRowPipelines()
	for _, p := range ps {
		compute.addRowPipeline(p.index, p.p, v)
	}
	return ps
}

func (compute matrixSelectScalar[D, Ds]) computeColPipelines() []matrix1Pipeline {
	v, vok := compute.value.extractElement()
	if !vok {
		panic(EmptyObject)
	}
	ps := compute.A.getColPipelines()
	for _, p := range ps {
		compute.addColPipeline(p.index, p.p, v)
	}
	return ps
}
