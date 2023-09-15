package GrB

// #include "GraphBLAS.h"
import "C"

// VectorEWiseUnion is like [VectorEWiseAddBinaryOp], except that two scalars are used to
// define how to compute the result when entries are present in one of the two input vectors
// (u and v), but no the other. Each of the two input scalars alpha and beta must contain
// an entry. When computing the result t = u + v, if u(i) is present but v(i) is not, then
// t(i) = u(i) + beta. Likewise, if v(i) is present but u(i) is not, then t(i) = alpha + v(i),
// where + denotes the the binary operator add.
//
// VectorEWiseUnion is a SuiteSparse:GraphBLAS extension.
func VectorEWiseUnion[Dw, Du, Dv any](
	w Vector[Dw],
	mask *Vector[bool],
	accum *BinaryOp[Dw, Dw, Dw],
	op BinaryOp[Dw, Du, Dv],
	u Vector[Du],
	alpha Scalar[Du],
	v Vector[Dv],
	beta Scalar[Dv],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADV(mask, accum, desc)
	info := Info(C.GxB_Vector_eWiseUnion(w.grb, cmask, caccum, op.grb, u.grb, alpha.grb, v.grb, beta.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}

// MatrixEWiseUnion is like [MatrixEWiseAddBinaryOp], except that two scalars are used to
// define how to compute the result when entries are present in one of the two input matrices
// (a and b), but no the other. Each of the two input scalars alpha and beta must contain
// an entry. When computing the result t = a + b, if a(i) is present but b(i) is not, then
// t(i) = a(i) + beta. Likewise, if b(i) is present but a(i) is not, then t(i) = alpha + b(i),
// where + denotes the the binary operator add.
//
// MatrixEWiseUnion is a SuiteSparse:GraphBLAS extension.
func MatrixEWiseUnion[DC, DA, DB any](
	c Matrix[DC],
	mask *Matrix[bool],
	accum *BinaryOp[DC, DC, DC],
	op BinaryOp[DC, DA, DB],
	a Matrix[DA],
	alpha Scalar[DA],
	b Matrix[DB],
	beta Scalar[DB],
	desc *Descriptor,
) error {
	cmask, caccum, cdesc := processMADM(mask, accum, desc)
	info := Info(C.GxB_Matrix_eWiseUnion(c.grb, cmask, caccum, op.grb, a.grb, alpha.grb, b.grb, beta.grb, cdesc))
	if info == success {
		return nil
	}
	return makeError(info)
}
