package go_jeans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func _TestPackUnpack(t *testing.T) {
	// 10 byte
	var data = []byte("hello word")
	fmt.Println(len(data), data)
	data = Pack(data)
	fmt.Println(len(data), data)

	data2, err := Unpack(bytes.NewBuffer(data))
	fmt.Println(err)
	fmt.Println(len(data2), string(data2), data2)
}

func _TestPackUnpackN(t *testing.T) {
	// 10 byte
	var data = []byte("hello word")
	fmt.Println(len(data), data)
	data, err := PackN(data, PacketHerderLen_16)
	fmt.Println(len(data), data)

	data2, err := UnpackN(bytes.NewBuffer(data), PacketHerderLen_16)
	fmt.Println(err)
	fmt.Println(len(data2), string(data2), data2)
}

type baseType struct {
	b bool
	i int
	i8 int8
	i16 int16
	i32 int32
	i64 int64

	ui uint
	ui8 uint8
	ui16 uint16
	ui32 uint32
	ui64 uint64

	f32 float32
	f64 float64

	s string
}

var bt =baseType{
	b:    true,
	i:    1,
	i8:   2,
	i16:  3,
	i32:  4,
	i64:  5,
	ui:   6,
	ui8:  7,
	ui16: 8,
	ui32: 9,
	ui64: 10,
	f32:  11.1,
	f64:  12.1234,
	s:    "hello word !",
}

var btBuf,jsonBuf []byte

func init(){
	var err error
	btBuf,err = BaseTypeToBytes(bt.i,bt.i8,bt.i16,bt.i32,bt.i64,bt.ui,bt.ui8,bt.ui16,bt.ui32,bt.ui64,bt.s,bt.b,bt.f32,bt.f64)
	if err != nil{
		log.Fatalln(err)
	}

	jsonBuf,err = json.Marshal(bt)
	if err != nil {
		log.Fatalln(err)
		return
	}
	//fmt.Println("init")
}

func TestBaseTypeToBytes(t *testing.T) {
	buf,err := BaseTypeToBytes(bt.i,bt.i8,bt.i16,bt.i32,bt.i64,bt.ui,bt.ui8,bt.ui16,bt.ui32,bt.ui64,bt.s,bt.b,bt.f32,bt.f64)
	if err != nil{
		t.Error(err)
		return
	}
	btBuf = buf
	//log.Println(buf)
}

func TestBytesToBaseType(t *testing.T) {
	TestBaseTypeToBytes(t)
	var bt2 baseType
	err:=BytesToBaseType(btBuf,&bt2.i,&bt2.i8,&bt2.i16,&bt2.i32,&bt2.i64,&bt2.ui,&bt2.ui8,&bt2.ui16,&bt2.ui32,&bt2.ui64,&bt2.s,&bt2.b,&bt2.f32,&bt2.f64)
	if err != nil {
		t.Error(err)
		return
	}
	//log.Println(bt2)
}

//go test -bench=BenchmarkBaseTypeToBytes$   -benchtime=3s .\ -cpuprofile="BenchmarkBaseTypeToBytes_CPUV1.out"
func BenchmarkBaseTypeToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_,err := BaseTypeToBytes(bt.i,bt.i8,bt.i16,bt.i32,bt.i64,bt.ui,bt.ui8,bt.ui16,bt.ui32,bt.ui64,bt.s,bt.b,bt.f32,bt.f64)
		if err != nil{
			b.Error(err)
			return
		}
	}
}

//go test -bench=BenchmarkBaseTypeToBytesV2$   -benchtime=3s .\ -cpuprofile="BenchmarkBaseTypeToBytes_CPUV2.out"
func BenchmarkBaseTypeToBytesV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_,err := BaseTypeToBytesV2(bt.i,bt.i8,bt.i16,bt.i32,bt.i64,bt.ui,bt.ui8,bt.ui16,bt.ui32,bt.ui64,bt.s,bt.b,bt.f32,bt.f64)
		if err != nil{
			b.Error(err)
			return
		}
	}
}

func BenchmarkBytesToBaseType(b *testing.B) {
	var bt2 baseType
	for i := 0; i < b.N; i++ {
		err:=BytesToBaseType(btBuf,&bt2.i,&bt2.i8,&bt2.i16,&bt2.i32,&bt2.i64,&bt2.ui,&bt2.ui8,&bt2.ui16,&bt2.ui32,&bt2.ui64,&bt2.s,&bt2.b,&bt2.f32,&bt2.f64)
		if err != nil {
			b.Error(err)
			return
		}
	}
}


func BenchmarkMarshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_,err := json.Marshal(bt)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkUnmarshal_JSON(b *testing.B) {
	var bt2 baseType
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBuf,&bt2)
		if err != nil {
			b.Error(err)
			return
		}
	}
}