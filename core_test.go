package go_jeans

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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

var btBuf, jsonBuf []byte

func init() {
	var err error
	btBuf, err = Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.bs, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	if err != nil {
		log.Fatalln(err)
	}

	jsonBuf, err = json.Marshal(bt)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("test init")
}

func TestEncode(t *testing.T) {
	buf, err := Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	if err != nil {
		t.Error(err)
		return
	}
	btBuf = buf
	log.Println(buf)
}

func TestDecode(t *testing.T) {
	var s struct{
		A int
		B string
		C bool
	}
	buf, _ := Encode(s.A,s.B,s.C)
	err := Decode(buf,&s.A,&s.B,&s.C)
	if err != nil {
		return
	}
	fmt.Println(s)
}
func TestBytesToBaseType(t *testing.T) {
	TestEncode(t)
	var bt2 baseType

	err := Decode(btBuf, &bt2.bs, &bt2.i, &bt2.i8, &bt2.i16, &bt2.i32, &bt2.i64, &bt2.ui, &bt2.ui8, &bt2.ui16, &bt2.ui32, &bt2.ui64, &bt2.s, &bt2.b, &bt2.f32, &bt2.f64)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(bt2)
}

type Person struct {
	Name    string
	Age     int
	Address Address
}

type Address struct {
	Street  string
	City    string
	Country string
	M       map[string]int
}

// experimental stage
func PrintStructMembers(data interface{}) {
	v := reflect.ValueOf(data)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		fmt.Printf("%s: ", fieldType.Name)

		switch field.Kind() {
		case reflect.Struct:
			fmt.Println()
			PrintStructMembers(field.Interface())
		case reflect.Map:
			fmt.Println()
			keys := field.MapKeys()
			for _, key := range keys {
				value := field.MapIndex(key)
				fmt.Printf("%s: %v\n", key, value)
			}
		default:
			fmt.Println(field.Interface())
		}
	}
}

func TestPrintStructMembers(t *testing.T) {
	person := Person{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
			M: map[string]int{
				"1": 1,
			},
		},
	}

	PrintStructMembers(person)
}

func TestPanic(t *testing.T) {


	var ll = 100000000
	var a = make([]int64,ll)
	for i, _ := range a {
		n := i
		if (n+1) %2 == 0 {
			n=n*-1
		}
		a[i] = int64(n)
	}
	//fmt.Println(a)
	var resA = make([]int64,0,ll)
	buf := make([]byte, binary.MaxVarintLen64)
	for _, x := range a {
		n:= binary.PutVarint(buf, x)
		n1,_ := binary.Varint(buf[:n])
		resA = append(resA, n1)
	}
	//fmt.Println(resA)
	for i, i2 := range a {
		if resA[i] != i2 {
			t.Error("错误：",resA,i2)
			return
		}
	}
	fmt.Println("正确")
}