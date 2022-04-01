package forGraphBLASGo

import (
	"math"
	"reflect"
)

type (
	Signed interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64
	}

	Unsigned interface {
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}

	Integer interface {
		Signed | Unsigned
	}

	Float interface {
		~float32 | ~float64
	}

	Number interface {
		Integer | Float
	}

	Ordered interface {
		Number | ~string
	}
)

func Identity[T any](x T) T {
	return x
}

func Abs[T Number](x T) T {
	return T(math.Abs(float64(x)))
}

func AInv[T Number](x T) T {
	return -x
}

func MInv[T Float](x T) T {
	return 1 / x
}

func LNot(x bool) bool {
	return !x
}

func BNot[T Integer](x T) T {
	return ^x
}

func LOr(x, y bool) bool {
	return x || y
}

func LAnd(x, y bool) bool {
	return x && y
}

func LXor(x, y bool) bool {
	return x != y
}

func LXNor(x, y bool) bool {
	return x == y
}

func BOr[T Integer](x, y T) T {
	return x | y
}

func BAnd[T Integer](x, y T) T {
	return x & y
}

func BXor[T Integer](x, y T) T {
	return x ^ y
}

func BXNor[T Integer](x, y T) T {
	return ^(x ^ y)
}

func Eq[T comparable](x, y T) bool {
	return x == y
}

func Ne[T comparable](x, y T) bool {
	return x != y
}

func Gt[T Ordered](x, y T) bool {
	return x > y
}

func Lt[T Ordered](x, y T) bool {
	return x < y
}

func Ge[T Ordered](x, y T) bool {
	return x >= y
}

func Le[T Ordered](x, y T) bool {
	return x <= y
}

func Oneb[Dout Number, Din1, Din2 any](_ Din1, _ Din2) Dout {
	return 1
}

func Trueb[Din1, Din2 any](_ Din1, _ Din2) bool {
	return true
}

func First[Din1, Din2 any](x Din1, _ Din2) Din1 {
	return x
}

func Second[Din1, Din2 any](_ Din1, y Din2) Din2 {
	return y
}

func Min[T Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Plus[T Number](x, y T) T {
	return x + y
}

func Minus[T Number](x, y T) T {
	return x - y
}

func Times[T Number](x, y T) T {
	return x * y
}

func Div[T Number](x, y T) T {
	return x / y
}

func RowIndex[D Number, Din any](_ Din, i, _, s int) D {
	return D(i + s)
}

func ColIndex[D Number, Din any](_ Din, _, j, s int) D {
	return D(j + s)
}

func DiagIndex[D Number, Din any](_ Din, i, j, s int) D {
	return D(j - i + s)
}

func TriL[Din any](_ Din, i, j, s int) bool {
	return j <= i+s
}

func TriU[Din any](_ Din, i, j, s int) bool {
	return j >= i+s
}

func Diag[Din any](_ Din, i, j, s int) bool {
	return j == i+s
}

func OffDiag[Din any](_ Din, i, j, s int) bool {
	return j != i+s
}

func ColLE[Din any](_ Din, _, j, s int) bool {
	return j <= s
}

func ColGT[Din any](_ Din, _, j, s int) bool {
	return j > s
}

func RowLE[Din any](_ Din, i, _, s int) bool {
	return i <= s
}

func RowGT[Din any](_ Din, i, _, s int) bool {
	return i > s
}

func ValueEq[Din comparable](value Din, _, _ int, s Din) bool {
	return value == s
}

func ValueNE[Din comparable](value Din, _, _ int, s Din) bool {
	return value != s
}

func ValueLT[Din Ordered](value Din, _, _ int, s Din) bool {
	return value < s
}

func ValueLE[Din Ordered](value Din, _, _ int, s Din) bool {
	return value <= s
}

func ValueGT[Din Ordered](value Din, _, _ int, s Din) bool {
	return value > s
}

func ValueGE[Din Ordered](value Din, _, _ int, s Din) bool {
	return value >= s
}

func PlusMonoid[T Number]() (addition BinaryOp[T, T, T], identity T) {
	return Plus[T], 0
}

func TimesMonoid[T Number]() (multiplication BinaryOp[T, T, T], identity T) {
	return Times[T], 1
}

func max[T Number]() T {
	switch reflect.ValueOf(T(0)).Kind() {
	case reflect.Int:
		x := int(math.MaxInt)
		return T(x)
	case reflect.Int8:
		x := int8(math.MaxInt8)
		return T(x)
	case reflect.Int16:
		x := int16(math.MaxInt16)
		return T(x)
	case reflect.Int32:
		x := int32(math.MaxInt32)
		return T(x)
	case reflect.Int64:
		x := int64(math.MaxInt64)
		return T(x)
	case reflect.Uint:
		x := uint(math.MaxUint)
		return T(x)
	case reflect.Uint8:
		x := uint8(math.MaxUint8)
		return T(x)
	case reflect.Uint16:
		x := uint16(math.MaxUint16)
		return T(x)
	case reflect.Uint32:
		x := uint32(math.MaxUint32)
		return T(x)
	case reflect.Uint64:
		x := uint64(math.MaxUint64)
		return T(x)
	case reflect.Uintptr:
		x := ^uintptr(0)
		return T(x)
	case reflect.Float32:
		x := float32(math.Inf(1))
		return T(x)
	case reflect.Float64:
		x := math.Inf(1)
		return T(x)
	default:
		panic("invalid type")
	}
}

func MinMonoid[T Number]() (minimum BinaryOp[T, T, T], identity T) {
	return Min[T], max[T]()
}

func min[T Number]() T {
	switch reflect.ValueOf(T(0)).Kind() {
	case reflect.Int:
		x := int(math.MinInt)
		return T(x)
	case reflect.Int8:
		x := int8(math.MinInt8)
		return T(x)
	case reflect.Int16:
		x := int16(math.MinInt16)
		return T(x)
	case reflect.Int32:
		x := int32(math.MinInt32)
		return T(x)
	case reflect.Int64:
		x := int64(math.MinInt64)
		return T(x)
	case reflect.Uint:
		x := 0
		return T(x)
	case reflect.Uint8:
		x := uint8(0)
		return T(x)
	case reflect.Uint16:
		x := uint16(0)
		return T(x)
	case reflect.Uint32:
		x := uint32(0)
		return T(x)
	case reflect.Uint64:
		x := uint64(0)
		return T(x)
	case reflect.Uintptr:
		x := uintptr(0)
		return T(x)
	case reflect.Float32:
		x := float32(math.Inf(-1))
		return T(x)
	case reflect.Float64:
		x := math.Inf(-1)
		return T(x)
	default:
		panic("invalid type")
	}
}

func MaxMonoid[T Number]() (maximum BinaryOp[T, T, T], identity T) {
	return Max[T], min[T]()
}

func LOrMonoid() (or BinaryOp[bool, bool, bool], identity bool) {
	return LOr, false
}

func LAndMonoid() (and BinaryOp[bool, bool, bool], identity bool) {
	return LAnd, true
}

func LXorMonoid() (xor BinaryOp[bool, bool, bool], identity bool) {
	return LXor, false
}

func LXNorMonoid() (xnor BinaryOp[bool, bool, bool], identity bool) {
	return LXNor, true
}

func PlusTimesSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return PlusMonoid[T], Times[T], 0
}

func MinPlusSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return MinMonoid[T], Plus[T], max[T]()
}

func MaxPlusSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return MaxMonoid[T], Plus[T], min[T]()
}

func MinTimesSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return MinMonoid[T], Times[T], max[T]()
}

func MinMaxSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return MinMonoid[T], Max[T], max[T]()
}

func MaxMinSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return MaxMonoid[T], Min[T], min[T]()
}

func MaxTimesSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return MaxMonoid[T], Times[T], min[T]()
}

func PlusMinSemiring[T Number]() (addition Monoid[T], multiplication BinaryOp[T, T, T], identity T) {
	return PlusMonoid[T], Min[T], 0
}

func LOrLAndSemiring() (addition Monoid[bool], multiplication BinaryOp[bool, bool, bool], identity bool) {
	return LOrMonoid, LAnd, false
}

func LAndLOrSemiring() (addition Monoid[bool], multiplication BinaryOp[bool, bool, bool], identity bool) {
	return LAndMonoid, LOr, true
}

func LXorLAndSemiring() (addition Monoid[bool], multiplication BinaryOp[bool, bool, bool], identity bool) {
	return LXorMonoid, LAnd, false
}

func LXNorLOrSemiring() (addition Monoid[bool], multiplication BinaryOp[bool, bool, bool], identity bool) {
	return LXNorMonoid, LOr, true
}

func MinFirstSemiring[T Number, Din2 any]() (addition Monoid[T], multiplication BinaryOp[T, T, Din2], identity T) {
	return MinMonoid[T], First[T, Din2], 0
}

func MinSecondSemiring[T Number, Din1 any]() (addition Monoid[T], multiplication BinaryOp[T, Din1, T], identity T) {
	return MinMonoid[T], Second[Din1, T], 0
}

func MaxFirstSemiring[T Number, Din2 any]() (addition Monoid[T], multiplication BinaryOp[T, T, Din2], identity T) {
	return MaxMonoid[T], First[T, Din2], 0
}

func MaxSecondSemiring[T Number, Din1 any]() (addition Monoid[T], multiplication BinaryOp[T, Din1, T], identity T) {
	return MaxMonoid[T], Second[Din1, T], 0
}
