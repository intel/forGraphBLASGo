package forGraphBLASGo

import (
	"errors"
	"math"
	"reflect"
)

const (
	Version    = 2
	Subversion = 0
)

func GetVersion() (version, subversion int) {
	return Version, Subversion
}

type (
	Index    = int // remark: not uint64!
	Mode     int
	WaitMode int
	Format   int
)

const IndexMax = math.MaxInt

func All(size int) []int {
	if size == 0 {
		return nil
	}
	return []int{-size}
}

func isAll(indices []int) (size int, ok bool) {
	switch len(indices) {
	case 0:
		return 0, true
	case 1:
		if size := indices[0]; size < 0 {
			return -size, true
		}
	}
	return len(indices), false
}

const (
	NonBlocking Mode = iota
	Blocking
)

const (
	Complete WaitMode = iota
	Materialize
)

const (
	CSRFormat Format = iota
	CSCFormat
	COOFormat
	homFormat
)

type Type = reflect.Kind

var (
	Bool   = reflect.Bool
	Int8   = reflect.Int8
	Uint8  = reflect.Uint8
	Int16  = reflect.Int16
	Uint16 = reflect.Uint16
	Int32  = reflect.Int32
	Uint32 = reflect.Uint32
	Int64  = reflect.Int64
	Uint64 = reflect.Uint64
	FP32   = reflect.Float32
	FP64   = reflect.Float64
)

type (
	UnaryOp[Dout, Din any] func(in Din) (out Dout)

	BinaryOp[Dout, Din1, Din2 any] func(in1 Din1, in2 Din2) (out Dout)

	IndexUnaryOp[Dout, Din1, Din2 any] func(in1 Din1, row, col int, in2 Din2) (out Dout)

	/*
		Monoid and Semiring are defined as higher-order functions rather than structs, because this allows us
		to have generic predefined (and user-defined) concrete monoids and semirings. See the file api_operators.go
		for examples.

		Go doesn't have generic variables or constants, so it's not possible to express this with structs. Consider, for example:

		var PlusMonoid[T Number] struct {
			operator: Plus[T],
			identity: T(0),
		}
	*/

	Monoid[D any] func() (operator BinaryOp[D, D, D], identity D)

	Semiring[Dout, Din1, Din2 any] func() (addition Monoid[Dout], multiplication BinaryOp[Dout, Din1, Din2], identity Dout)
)

func (m Monoid[D]) operator() BinaryOp[D, D, D] {
	op, _ := m()
	return op
}

func (m Monoid[D]) identity() D {
	_, id := m()
	return id
}

func (s Semiring[Dout, Din1, Din2]) addition() Monoid[Dout] {
	add, _, _ := s()
	return add
}

func (s Semiring[Dout, Din1, Din2]) multiplication() BinaryOp[Dout, Din1, Din2] {
	_, mult, _ := s()
	return mult
}

func (s Semiring[Dout, Din1, Din2]) identity() Dout {
	_, _, id := s()
	return id
}

func MonoidNew[D any](operator BinaryOp[D, D, D], identity D) Monoid[D] {
	return func() (BinaryOp[D, D, D], D) {
		return operator, identity
	}
}

func SemiringNew[Dout, Din1, Din2 any](addition Monoid[Dout], multiplication BinaryOp[Dout, Din1, Din2], identity Dout) Semiring[Dout, Din1, Din2] {
	return func() (Monoid[Dout], BinaryOp[Dout, Din1, Din2], Dout) {
		return addition, multiplication, identity
	}
}

type Info = error

// todo: we need error values with more information about specific error causes
var (
	// error return values
	Success = error(nil)
	NoValue = errors.New("no value")

	UninitializedObject = errors.New("uninitialized object")
	NullPointer         = errors.New("null pointer")
	InvalidValue        = errors.New("invalid value")
	InvalidIndex        = errors.New("invalid index")
	DomainMismatch      = errors.New("domain mismatch")
	DimensionMismatch   = errors.New("dimension mismatch")
	OutputNotEmpty      = errors.New("output not empty")
	NotImplemented      = errors.New("not implemented")

	Panic             = errors.New("panic")
	OutOfMemory       = errors.New("out of memory")
	InsufficientSpace = errors.New("insufficient space")
	InvalidObject     = errors.New("invalid object")
	IndexOutOfBounds  = errors.New("index out of bounds")
	EmptyObject       = errors.New("empty object")
)

func Init(mode Mode) error {
	if mode != NonBlocking {
		panic("not implemented yet") // todo: add a call to materialize after each non-blocking method
	}
	return nil
}

func Finalize() error {
	return nil
}

/*
TypeNew, UnaryOpNew, BinaryOpNew not needed due to generics
*/

// todo: try to define a polymorphic interface, by using any as parameters, and then dispatch on actual types, if at all possible.
