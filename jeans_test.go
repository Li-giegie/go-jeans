package go_jeans

import (
	"fmt"
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

func (b *Base) String() string {
	return fmt.Sprintf("Base {I: %v, I8: %v, I16: %v, I32: %v, I64: %v, Ui: %v, Ui8: %v, Ui16: %v, Ui32: %v, Ui64: %v, Bo: %v, F32: %v, F64: %v, B: %v, Bs: %s, S: %v}", b.I, b.I8, b.I16, b.I32, b.I64, b.Ui, b.Ui8, b.Ui16, b.Ui32, b.Ui64, b.Bo, b.F32, b.F64, b.B, b.Bs, b.S)
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

func (s *Slice) String() string {
	return fmt.Sprintf("Slice {Is: %v, I8s: %v, I16s: %v, I32s: %v, I64s: %v, Uis: %v, Ui8s: %v, Ui16s: %v, Ui32s: %v, Ui64s: %v, Bos: %v, F32s: %v, F64s: %v, Bs: %v, Ss: %v}", s.Is, s.I8s, s.I16s, s.I32s, s.I64s, s.Uis, s.Ui8s, s.Ui16s, s.Ui32s, s.Ui64s, s.Bos, s.F32s, s.F64s, s.Bs, s.Ss)
}

func (s *Slice) FieldNum() int {
	return 15
}

func (s *Slice) FieldsToInterface() []interface{} {
	return []interface{}{
		s.Is, s.I8s, s.I16s, s.I32s, s.I64s,
		s.Uis, s.Ui8s, s.Ui16s, s.Ui32s, s.Ui64s,
		s.Bos, s.Bs, s.Bs,
		s.F32s, s.F64s,
		s.Ss,
	}
}

func (s *Slice) FieldsPointerToInterface() []interface{} {
	return []interface{}{
		&s.Is, &s.I8s, &s.I16s, &s.I32s, &s.I64s,
		&s.Uis, &s.Ui8s, &s.Ui16s, &s.Ui32s, &s.Ui64s,
		&s.Bos, &s.Bs, &s.Bs,
		&s.F32s, &s.F64s,
		&s.Ss,
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

func TestStringToBase(t *testing.T) {
	s := "hello jeans"
	b := stringToBytes(&s)
	if s != string(*b) || len(s) != len(*b) || len(s) != cap(*b) {
		t.Error("stringsToBytes error")
	}
	println("pass", string(*b), len(*b), cap(*b))
}

// 684890460
func Benchmark_Ref(b *testing.B) {
	var o interface{}
	for i := 0; i < b.N; i++ {
		o = i
		//reflect.TypeOf(o).Kind()

		switch v := o.(type) {
		case int:
			_ = v
		}
	}
}
