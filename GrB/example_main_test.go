package GrB_test

import (
	"github.com/intel/forGraphBLASGo/GrB"
	"testing"
)

func TestMain(m *testing.M) {
	if err := GrB.Init(GrB.NonBlocking); err != nil {
		panic(err)
	}
	defer func() {
		if err := GrB.Finalize(); err != nil {
			panic(err)
		}
	}()
	m.Run()
}
