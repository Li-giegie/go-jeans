package go_jeans

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncodeV2(t *testing.T) {
	var buf = make([]byte, 0, 89)
	result, err := EncodeV2(buf, bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("result ", result, len(result), cap(result))
	var bt1 baseType
	err = DecodeV2(result, &bt1.bs, &bt1.i, &bt1.i8, &bt1.i16, &bt1.i32, &bt1.i64, &bt1.ui, &bt1.ui8, &bt1.ui16, &bt1.ui32, &bt1.ui64, &bt1.s, &bt1.b, &bt1.f32, &bt1.f64)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(bt, bt1, reflect.DeepEqual(bt, bt1))
}
