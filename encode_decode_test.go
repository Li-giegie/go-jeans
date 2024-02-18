package go_jeans

import (
	"encoding/json"
	"fmt"
	faker "github.com/bxcodec/faker/v3"
	"os"
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
	var b = new(Base)
	if err := faker.FakeData(b); err != nil {
		panic(err)
	}
	b.I *= -1
	b.I8 *= -1
	b.I16 *= -1
	b.I32 *= -1
	b.I64 *= -1
	return b
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

type Slice struct {
	Is    []int
	I8s   []int8
	I16s  []int16
	I32s  []int32
	I64s  []int64
	Uis   []uint
	Ui8s  []uint8
	Ui16s []uint16
	Ui32s []uint32
	Ui64s []uint64
	Bos   []bool
	F32s  []float32
	F64s  []float64
	Bs    []byte
	Ss    []string
}

func (b *Slice) FieldNum() int {
	return 15
}

func (b *Slice) FieldsToInterface() []interface{} {
	return []interface{}{
		b.Is, b.I8s, b.I16s, b.I32s, b.I64s,
		b.Uis, b.Ui8s, b.Ui16s, b.Ui32s, b.Ui64s,
		b.Bos, b.Bs, b.Bs,
		b.F32s, b.F64s,
		b.Ss,
	}
}

func (b *Slice) FieldsPointerToInterface() []interface{} {
	return []interface{}{
		&b.Is, &b.I8s, &b.I16s, &b.I32s, &b.I64s,
		&b.Uis, &b.Ui8s, &b.Ui16s, &b.Ui32s, &b.Ui64s,
		&b.Bos, &b.Bs, &b.Bs,
		&b.F32s, &b.F64s,
		&b.Ss,
	}
}

func NewSlice() *Slice {
	return &Slice{
		Is:    []int{1, 2, 3, -1, -2, -3},
		I8s:   []int8{1, 2, 3, -1, -2, -3},
		I16s:  []int16{1, 2, 3, -1, -2, -3},
		I32s:  []int32{1, 2, 3, -1, -2, -3},
		I64s:  []int64{1, 2, 3, -1, -2, -3},
		Uis:   []uint{1, 2, 3},
		Ui8s:  []uint8{1, 2, 3},
		Ui16s: []uint16{1, 2, 3},
		Ui32s: []uint32{1, 2, 3},
		Ui64s: []uint64{1, 2, 3},
		Bos:   []bool{true, false, false},
		F32s:  []float32{1.11, -2.22, 3.333, 4.4444},
		F64s:  []float64{1.11, 2.22, 3.333, 4.4444, 5.55555, -6.666666},
		Bs:    []byte{1, 2, 3, 4, 5, 6},
		Ss:    []string{"abc", "123das", "", "222", ""},
	}
}

func TestEncodeSlice(t *testing.T) {
	for j := 0; j < 1000; j++ {
		s := NewSlice()
		buf, err := EncodeSlice(s.FieldsToInterface()...)
		if err != nil {
			t.Error(err)
			return
		}
		decodeUi32s := new(Slice)
		if err = DecodeSlice(buf, decodeUi32s.FieldsPointerToInterface()...); err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(s, decodeUi32s) {
			buf1, err1 := json.MarshalIndent(s, "", "\t")
			buf2, err2 := json.MarshalIndent(decodeUi32s, "", "\t")
			if err1 != nil || err2 != nil {
				fmt.Println("json err: ", err1, err2)
				return
			}
			err1 = os.WriteFile("./t1.json", buf1, 0666)
			err2 = os.WriteFile("./t2.json", buf2, 0666)
			if err1 != nil || err2 != nil {
				fmt.Println("write file err: ", err1, err2)
				return
			}
			t.Error("DeepEqual fail")
			return
		}
		fmt.Println("success", j)
	}
}
