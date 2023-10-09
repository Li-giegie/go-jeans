package go_jeans

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
)

func TestPackUnpack(t *testing.T) {
	// 10 byte
	var data = []byte("hello word!")
	data = Pack(data)
	fmt.Println("Pack:", len(data), data)

	data2, err := Unpack(bytes.NewBuffer(data))
	fmt.Println("Unpack:", err, len(data2), string(data2), data2)
}

func TestPackUnpackN(t *testing.T) {
	// 10 byte
	var data = []byte("hello word")
	fmt.Println(len(data), data)
	data, err := PackN(data, PacketHerderLen_16)
	fmt.Println("Pack:", len(data), data)

	data2, err := UnpackN(bytes.NewBuffer(data), PacketHerderLen_16)
	fmt.Println("Unpack:", err, len(data2), string(data2), data2)
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

func TestEncode(t *testing.T) {
	_Encode()
}

func TestDecode(t *testing.T) {
	buf,err := _Encode()
	if err != nil {
		t.Error(err)
		return
	}
	b,err := _Decode(buf)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(b)
}


func _Encode() ([]byte,error){
	buf, err := Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	return buf,err
}

func _Decode( buf []byte) (*baseType,error){
	b := new(baseType)
	err := Decode(buf, &b.bs, &b.i, &b.i8, &b.i16, &b.i32, &b.i64, &b.ui, &b.ui8, &b.ui16, &b.ui32, &b.ui64, &b.s, &b.b, &b.f32, &b.f64)
	return b,err
}

var str = `commit cd4011e37162afabf4908f64b4ad9c5af16fc2f1
Author: Li-giegie <1261930106@qq.com>
Date:   Thu May 11 13:16:47 2023 +0800
Please make sure you have the correct access rights
and the repository exists.
PS D:\_project\GO Project\go-jeans\example> git push origin master
`
var strList = strings.Split(str," ")

//BenchmarkEncodeString-12    	  732082	      1550 ns/op
func BenchmarkEncodeString(b *testing.B) {
	var buf = make([]byte,0,1024)
	for i := 0; i < b.N; i++ {
		for _, s := range strList {
			_b,err := Encode(s)
			if err != nil {
				b.Error(err)
				return
			}
			buf = append(buf, _b...)
		}
	}
}

var sp = "\\SpLiT"
func StringListPack(s []string) []byte{
	strPack := strings.Join(s,sp)
	buf := Pack([]byte(strPack))
	return buf
}

func StringListUnpack(buf []byte) []string{
	_= binary.LittleEndian.Uint32(buf[:4])
	return strings.Split(string(buf[4:]),sp)
}

//BenchmarkStringList-12    	 2632795	       428.1 ns/op
func BenchmarkStringList(b *testing.B) {
	buf := StringListPack(strList)
	for i := 0; i < b.N; i++ {
		StringListUnpack(buf)
	}
}

func TestStringList(t *testing.T) {
	buf := StringListPack(strList)
	sl := StringListUnpack(buf)
	for i, s := range strList {
		if sl[i] != s {
			t.Error("解包失败")
		}
	}
}