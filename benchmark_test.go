package go_jeans

import (
	"encoding/json"
	"testing"
)

// BenchmarkPack-12    	42299990	        26.61 ns/op
func BenchmarkPack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pack([]byte("qwertyuiop1234567890!@#@#$%^&*()_{"))
	}
}

// go test -bench=BenchmarkEncode$   -benchtime=3s .\ -cpuprofile="BenchmarkEncode_CPUV1.out"
func Benchmark_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil{
			b.Error(err)
			return
		}
	}
}

func Benchmark_Decode(b *testing.B) {
	buf,err:=_Encode()
	if err != nil{
		b.Error(err)
		return
	}
	obj := new(baseType)
	for i := 0; i < b.N; i++ {
		err = Decode(buf, &obj.bs, &obj.i, &obj.i8, &obj.i16, &obj.i32, &obj.i64, &obj.ui, &obj.ui8, &obj.ui16, &obj.ui32, &obj.ui64, &obj.s, &obj.b, &obj.f32, &obj.f64)
		if err != nil {
			b.Error(err)
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
	buf,err := json.Marshal(bt)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(buf, &bt2)
		if err != nil {
			b.Error(err)
			return
		}
	}
}
