package go_jeans

import (
	"fmt"
	"reflect"
	"testing"
)

type Base struct {
	I    int
	I8   int8
	I16  int16
	I32  int32
	I64  int64
	Ui   uint
	Ui8  uint8
	Ui16 uint16
	Ui32 uint32
	Ui64 uint64
	Bo   bool
	F32  float32
	F64  float64
	B    byte
	Bs   []byte
	S    string
}

func (b *Base) FieldNum() int {
	return 16
}

func (b *Base) FieldsToInterface() []interface{} {
	return []interface{}{
		b.I, b.I8, b.I16, b.I32, b.I64,
		b.Ui, b.Ui8, b.Ui16, b.Ui32, b.Ui64,
		b.Bo, b.B, b.Bs,
		b.F32, b.F64,
		b.S,
	}
}

func (b *Base) FieldsPointerToInterface() []interface{} {
	return []interface{}{
		&b.I, &b.I8, &b.I16, &b.I32, &b.I64,
		&b.Ui, &b.Ui8, &b.Ui16, &b.Ui32, &b.Ui64,
		&b.Bo, &b.B, &b.Bs,
		&b.F32, &b.F64,
		&b.S,
	}
}

func NewBase() *Base {
	return &Base{
		I:    -64,
		I8:   -8,
		I16:  -16,
		I32:  -32,
		I64:  -64,
		Ui:   64,
		Ui8:  8,
		Ui16: 16,
		Ui32: 32,
		Ui64: 64,
		Bo:   true,
		F32:  3.1415926,
		F64:  1314.520,
		B:    255,
		Bs:   []byte("hello world ~"),
		S:    "say bay ~",
	}
}

func TestEncode(t *testing.T) {
	base := NewBase()
	result, err := Encode(base.FieldsToInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	var decodeBase = new(Base)
	err = Decode(result, decodeBase.FieldsPointerToInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(base, decodeBase) {
		t.Error("DeepEqual fail")
		return
	}
	fmt.Println(base, decodeBase)
}

func TestEncodeFaster(t *testing.T) {
	var buf = make([]byte, 0, 89)
	base := NewBase()
	result, err := EncodeFaster(buf, base.FieldsToInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	var decodeBase = new(Base)
	err = Decode(result, decodeBase.FieldsPointerToInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(base, decodeBase) {
		t.Error("DeepEqual fail")
		return
	}
	fmt.Println(base, decodeBase)
}

func TestEncodeSlice(t *testing.T) {
	var Ui32s = []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	buf, err := EncodeSlice(Ui32s)
	if err != nil {
		t.Error(err)
		return
	}
	var decodeUi32s []uint32
	if err = DecodeSlice(buf, &decodeUi32s); err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(Ui32s, decodeUi32s) {
		t.Error("DeepEqual fail")
		return
	}
	fmt.Println(Ui32s, decodeUi32s)
}
