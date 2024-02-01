package go_jeans

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestPackUnpack(t *testing.T) {
	// 10 byte
	var data = []byte("hello word!")
	data = Pack(data)
	_, err := Unpack(bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
	}
}

func TestPackUnpackN(t *testing.T) {
	var data = []byte("hello word")
	data, err := PackN(data, PacketHerderLenType_uint16)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = UnpackN(bytes.NewBuffer(data), PacketHerderLenType_uint16)
	if err != nil {
		t.Error(err)
		return
	}
}

type baseType struct {
	b   bool
	i   int
	i8  int8
	i16 int16
	i32 int32
	i64 int64

	ui   uint
	ui8  uint8
	ui16 uint16
	ui32 uint32
	ui64 uint64

	bs []byte

	f32 float32
	f64 float64

	s string
}

var bt = baseType{
	b:    true,
	i:    -1,
	i8:   2,
	i16:  -3,
	i32:  4,
	i64:  -5,
	ui:   6,
	ui8:  7,
	ui16: 8,
	ui32: 9,
	ui64: 10,
	f32:  11.1,
	f64:  12.1234,
	bs:   []byte("test[]byte"),
	s:    "hello word !",
}

var btpb = BaseType_PB{
	B:    true,
	I:    -1,
	I8:   2,
	I16:  -3,
	I32:  4,
	I64:  -5,
	Ui:   6,
	Ui8:  7,
	Ui16: 8,
	Ui32: 9,
	Ui64: 10,
	Bs:   []byte("test[]byte"),
	F32:  11.1,
	F64:  12.1234,
	S:    "hello word !",
}

func TestEncode(t *testing.T) {
	_, err := _Encode()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEncodeWithByte(t *testing.T) {
	buf, itemLen, err := _EncodeWithByte()
	if err != nil {
		t.Error(err)
		return
	}
	buf2, err2 := _Encode()
	if err2 != nil {
		t.Error(err2)
		return
	}
	if !bytes.Equal(buf2, buf) {
		t.Error("TestEncodeWithByte fail -1")
		return
	}
	count := 0
	for _, i2 := range itemLen {
		count += int(i2)
	}
	if count != len(buf) {
		t.Error("TestEncodeWithByte fail -2")
		return
	}
}

func TestDecode(t *testing.T) {
	buf, err := _Encode()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = _Decode(buf)
	if err != nil {
		t.Error(err)
	}
}

func _Encode() ([]byte, error) {
	buf, err := Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	return buf, err
}

func _EncodeWithByte() ([]byte, []int32, error) {
	buf, itemLen, err := EncodeWithLenByteItem(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	return buf, itemLen, err
}
func _Decode(buf []byte) (*baseType, error) {
	b := new(baseType)
	err := Decode(buf, &b.bs, &b.i, &b.i8, &b.i16, &b.i32, &b.i64, &b.ui, &b.ui8, &b.ui16, &b.ui32, &b.ui64, &b.s, &b.b, &b.f32, &b.f64)
	return b, err
}

func TestAA(t *testing.T) {
	fmt.Println(binary.LittleEndian.AppendUint32([]byte{1, 2, 3}, 3))

}
