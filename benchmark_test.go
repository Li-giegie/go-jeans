package go_jeans

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func init() {
	var err error
	btBuf, err = BaseTypeToBytes(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
	if err != nil {
		log.Fatalln(err)
	}

	jsonBuf, err = json.Marshal(bt)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("benchmark init")
}

// BenchmarkPack-12    	42299990	        26.61 ns/op
func BenchmarkPack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pack([]byte("qwertyuiop1234567890!@#@#$%^&*()_{"))
	}
}

// go test -bench=BenchmarkBaseTypeToBytes$   -benchtime=3s .\ -cpuprofile="BenchmarkBaseTypeToBytes_CPUV1.out"
func Benchmark_BaseTypeToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := BaseTypeToBytes(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func Benchmark_BytesToBaseType(b *testing.B) {
	var bt2 baseType
	for i := 0; i < b.N; i++ {
		err := BytesToBaseType(btBuf, &bt2.bs, &bt2.i, &bt2.i8, &bt2.i16, &bt2.i32, &bt2.i64, &bt2.ui, &bt2.ui8, &bt2.ui16, &bt2.ui32, &bt2.ui64, &bt2.s, &bt2.b, &bt2.f32, &bt2.f64)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkMarshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bt)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkUnmarshal_JSON(b *testing.B) {
	var bt2 baseType
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(jsonBuf, &bt2)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
