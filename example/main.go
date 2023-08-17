package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)


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
func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe Err:", err.Error())
		return
	}
}

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(resp, "Hello World, Are You OK?")
}
