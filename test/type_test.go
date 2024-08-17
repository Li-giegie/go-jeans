package test

import (
	"fmt"
	"unsafe"
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
	b := new(Base)
	genRandValue(b)
	return b
}

type Slice struct {
	I     []int
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
	return fmt.Sprintf("Slice {I: %v, I8s: %v, I16s: %v, I32s: %v, I64s: %v, Uis: %v, Ui8s: %v, Ui16s: %v, Ui32s: %v, Ui64s: %v, Bos: %v, F32s: %v, F64s: %v, Bs: %v, Ss: %v}", s.I, s.I8s, s.I16s, s.I32s, s.I64s, s.Uis, s.Ui8s, s.Ui16s, s.Ui32s, s.Ui64s, s.Bos, s.F32s, s.F64s, s.Bs, s.Ss)
}

func (s *Slice) FieldNum() int {
	return 15
}

func (s *Slice) FieldsToInterface() []interface{} {
	return []interface{}{
		s.I, s.I8s, s.I16s, s.I32s, s.I64s,
		s.Uis, s.Ui8s, s.Ui16s, s.Ui32s, s.Ui64s,
		s.Bos, s.Bs, s.Bs,
		s.F32s, s.F64s,
		s.Ss,
	}
}

func (s *Slice) FieldsPointerToInterface() []interface{} {
	return []interface{}{
		&s.I, &s.I8s, &s.I16s, &s.I32s, &s.I64s,
		&s.Uis, &s.Ui8s, &s.Ui16s, &s.Ui32s, &s.Ui64s,
		&s.Bos, &s.Bs, &s.Bs,
		&s.F32s, &s.F64s,
		&s.Ss,
	}
}

func NewSlice() *Slice {
	s := new(Slice)
	genRandValue(s)
	return s
}

type (
	ui   uint
	ui8  uint8
	ui16 uint16
	ui32 uint32
	ui64 uint64
)

type (
	i   int
	i8  int8
	i16 int16
	i32 int32
	i64 int64
)

type (
	f32 float32
	f64 float64
)

type str string
type bl bool
type (
	b8 byte
	bs []byte
)

type CustomType struct {
	Ui   ui
	Ui8  ui8
	Ui16 ui16
	Ui32 ui32
	Ui64 ui64
	I    i
	I8   i8
	I16  i16
	I32  i32
	I64  i64
	F32  f32
	F64  f64
	Str  str
	Bl   bl
	B8   b8
	Bs   bs
}

func (s *CustomType) String() string {
	return fmt.Sprintf("CustomType {I: %v, I8s: %v, I16s: %v, I32s: %v, I64s: %v, Uis: %v, Ui8s: %v, Ui16s: %v, Ui32s: %v, Ui64s: %v, Bos: %v, F32s: %v, F64s: %v, Bs: %v, Ss: %v}", s.I, s.I8, s.I16, s.I32, s.I64, s.Ui, s.Ui8, s.Ui16, s.Ui32, s.Ui64, s.Bl, s.F32, s.F64, s.Bs, s.Str)
}

func (s *CustomType) FieldNum() int {
	return 15
}

func (s *CustomType) FieldsToInterface() []interface{} {
	return []interface{}{
		s.Ui, s.Ui8, s.Ui16, s.Ui32, s.Ui64,
		s.I, s.I8, s.I16, s.I32, s.I64,
		s.Bl, s.Bs,
		s.F32, s.F64,
		s.Str,
	}
}

func (s *CustomType) FieldsPointerToInterface() []interface{} {
	return []interface{}{
		&s.Ui, &s.Ui8, &s.Ui16, &s.Ui32, &s.Ui64,
		&s.I, &s.I8, &s.I16, &s.I32, &s.I64,
		&s.Bl, &s.Bs,
		&s.F32, &s.F64,
		&s.Str,
	}
}

func NewCustomType() *CustomType {
	c := new(CustomType)
	genRandValue(c)
	return c
}

type T struct {
	Base       *Base
	Slice      *Slice
	CustomType *CustomType
}

func NewT() *T {
	t := new(T)
	genRandValue(t)
	return t
}

type BaseSlice struct {
	Base  *Base
	Slice *Slice
}

func GenBaseSlice() *BaseSlice {
	var bs = new(BaseSlice)
	genRandValue(bs)
	return bs
}

func NewBaseSlice() *BaseSlice {
	return &BaseSlice{
		Base:  new(Base),
		Slice: new(Slice),
	}
}

func (s *BaseSlice) FieldsPointerToInterface() []interface{} {
	return append(s.Base.FieldsPointerToInterface(), s.Slice.FieldsPointerToInterface()...)
}
func (s *BaseSlice) FieldsInterface() []interface{} {
	return append(s.Base.FieldsToInterface(), s.Slice.FieldsToInterface()...)
}

type slice struct {
	ptr unsafe.Pointer
	len int
	cap int
}