package test

import (
	"fmt"
	"testing"
)

func TestGnRandValue(t *testing.T) {
	bs := GenBaseSlice()
	_ = bs
	fmt.Println(bs.Base.String(), "\n", bs.Slice.String())
}
