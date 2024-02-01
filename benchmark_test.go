package go_jeans

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack/v5"
	"testing"
)

//environment：windows11
//ED：go_jeans、JSON、msgpack、protobuf
//goos: windows
//goarch: amd64
//pkg: github.com/Li-giegie/go-jeans
//cpu: AMD Ryzen 5 5600H with Radeon Graphics

// BenchmarkPack-12         6362475               195.1 ns/op
// BenchmarkPack-12         5927304               205.4 ns/op
// BenchmarkPack-12         6263377               196.5 ns/op
// BenchmarkPack-12         6363589               178.6 ns/op
// BenchmarkPack-12         6545034               188.7 ns/op
func BenchmarkPack(b *testing.B) {
	var buf = make([]byte, 1024)
	for i := 0; i < b.N; i++ {
		Pack(buf)
	}
}

// go test -run=none -bench=Benchmark_Encode$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// Benchmark_Encode        13686333                88.77 ns/op          128 B/op          1 allocs/op
// Benchmark_Encode        13080744                85.38 ns/op          128 B/op          1 allocs/op
// Benchmark_Encode        13990644                84.85 ns/op          128 B/op          1 allocs/op
// Benchmark_Encode        14046850                85.49 ns/op          128 B/op          1 allocs/op
// Benchmark_Encode        13442229                85.74 ns/op          128 B/op          1 allocs/op
func Benchmark_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -bench=Benchmark_EncodeV2$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// Benchmark_EncodeV2      26887321                45.09 ns/op            0 B/op          0 allocs/op
// Benchmark_EncodeV2      27020880                44.07 ns/op            0 B/op          0 allocs/op
// Benchmark_EncodeV2      25414467                44.70 ns/op            0 B/op          0 allocs/op
// Benchmark_EncodeV2      27692456                45.06 ns/op            0 B/op          0 allocs/op
// Benchmark_EncodeV2      26919228                44.02 ns/op            0 B/op          0 allocs/op
func Benchmark_EncodeV2(b *testing.B) {
	var buf = make([]byte, 0, 128)
	var err error
	for i := 0; i < b.N; i++ {
		_, err = EncodeV2(buf, bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -bench=Benchmark_EncodeWithByte$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// Benchmark_EncodeWithByte         5000367               241.9 ns/op           272 B/op          9 allocs/op
// Benchmark_EncodeWithByte         4941871               244.5 ns/op           272 B/op          9 allocs/op
// Benchmark_EncodeWithByte         5023605               238.8 ns/op           272 B/op          9 allocs/op
// Benchmark_EncodeWithByte         4944114               244.1 ns/op           272 B/op          9 allocs/op
// Benchmark_EncodeWithByte         4945249               238.4 ns/op           272 B/op          9 allocs/op
func Benchmark_EncodeWithByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf, itemLen, err := EncodeWithLenByteItem(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err, buf, itemLen)
			return
		}
	}
}

// go test -run=none -bench=Benchmark_Decode$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// Benchmark_Decode        26133769                48.00 ns/op            0 B/op          0 allocs/op
// Benchmark_Decode        24945638                51.22 ns/op            0 B/op          0 allocs/op
// Benchmark_Decode        24642883                47.36 ns/op            0 B/op          0 allocs/op
// Benchmark_Decode        24001584                47.83 ns/op            0 B/op          0 allocs/op
// Benchmark_Decode        25729594                48.39 ns/op            0 B/op          0 allocs/o
func Benchmark_Decode(b *testing.B) {
	buf, err := _Encode()
	if err != nil {
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

// go test -run=none -bench=Benchmark_DecodeV2$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// Benchmark_DecodeV2      34651234                34.31 ns/op            0 B/op          0 allocs/op
// Benchmark_DecodeV2      34010134                34.91 ns/op            0 B/op          0 allocs/op
// Benchmark_DecodeV2      35162773                34.61 ns/op            0 B/op          0 allocs/op
// Benchmark_DecodeV2      35196394                34.15 ns/op            0 B/op          0 allocs/op
// Benchmark_DecodeV2      35161228                34.17 ns/op            0 B/op          0 allocs/op
func Benchmark_DecodeV2(b *testing.B) {
	buf, err := _Encode()
	if err != nil {
		b.Error(err)
		return
	}
	obj := new(baseType)
	for i := 0; i < b.N; i++ {
		err = DecodeV2(buf, &obj.bs, &obj.i, &obj.i8, &obj.i16, &obj.i32, &obj.i64, &obj.ui, &obj.ui8, &obj.ui16, &obj.ui32, &obj.ui64, &obj.s, &obj.b, &obj.f32, &obj.f64)
		if err != nil {
			b.Error(err)
		}
	}
}

// go test -run=none -bench=BenchmarkMarshal_JSON$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkMarshal_JSON    9225651               129.7 ns/op           120 B/op          2 allocs/op
// BenchmarkMarshal_JSON    9515396               130.8 ns/op           120 B/op          2 allocs/op
// BenchmarkMarshal_JSON    9454624               128.0 ns/op           120 B/op          2 allocs/op
// BenchmarkMarshal_JSON    9519744               125.2 ns/op           120 B/op          2 allocs/op
// BenchmarkMarshal_JSON    9291312               129.5 ns/op           120 B/op          2 allocs/op
func BenchmarkMarshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bt)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -bench=BenchmarkUnmarshal_JSON$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkUnmarshal_JSON          6519027               189.2 ns/op           152 B/op          2 allocs/op
// BenchmarkUnmarshal_JSON          6481591               183.4 ns/op           152 B/op          2 allocs/op
// BenchmarkUnmarshal_JSON          6295401               187.6 ns/op           152 B/op          2 allocs/op
// BenchmarkUnmarshal_JSON          6403414               183.1 ns/op           152 B/op          2 allocs/op
// BenchmarkUnmarshal_JSON          6658081               182.1 ns/op           152 B/op          2 allocs/op
func BenchmarkUnmarshal_JSON(b *testing.B) {
	var bt2 baseType
	buf, err := json.Marshal(bt)
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

// go test -run=none -bench=BenchmarkMarshal_MsgPack$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkMarshal_MsgPack         5127717               235.2 ns/op           224 B/op          3 allocs/op
// BenchmarkMarshal_MsgPack         5147064               233.7 ns/op           224 B/op          3 allocs/op
// BenchmarkMarshal_MsgPack         5250889               238.5 ns/op           224 B/op          3 allocs/op
// BenchmarkMarshal_MsgPack         5138666               232.9 ns/op           224 B/op          3 allocs/op
// BenchmarkMarshal_MsgPack         5079962               230.3 ns/op           224 B/op          3 allocs/op
func BenchmarkMarshal_MsgPack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := msgpack.Marshal(bt)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -bench=BenchmarkUnmarshal_MsgPack$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkUnmarshal_MsgPack       6983198               176.3 ns/op            48 B/op          1 allocs/op
// BenchmarkUnmarshal_MsgPack       6987830               175.7 ns/op            48 B/op          1 allocs/op
// BenchmarkUnmarshal_MsgPack       7033148               170.8 ns/op            48 B/op          1 allocs/op
// BenchmarkUnmarshal_MsgPack       6829590               174.0 ns/op            48 B/op          1 allocs/op
// BenchmarkUnmarshal_MsgPack       6791074               171.8 ns/op            48 B/op          1 allocs/op
func BenchmarkUnmarshal_MsgPack(b *testing.B) {
	buf, err := msgpack.Marshal(bt)
	if err != nil {
		b.Error(err)
		return
	}
	var bt baseType
	for i := 0; i < b.N; i++ {
		if err = msgpack.Unmarshal(buf, &bt); err != nil {
			b.Error(err)
			return
		}

	}
}

// go test -run=none -bench=BenchmarkMarshal_ProtoBuf$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkMarshal_ProtoBuf        4799300               247.6 ns/op            96 B/op          1 allocs/op
// BenchmarkMarshal_ProtoBuf        4863567               256.9 ns/op            96 B/op          1 allocs/op
// BenchmarkMarshal_ProtoBuf        4702652               254.2 ns/op            96 B/op          1 allocs/op
// BenchmarkMarshal_ProtoBuf        4850224               248.0 ns/op            96 B/op          1 allocs/op
// BenchmarkMarshal_ProtoBuf        4794951               260.1 ns/op            96 B/op          1 allocs/op
func BenchmarkMarshal_ProtoBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(&btpb)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -bench=BenchmarkUnmarshal_ProtoBuf$ -count 5 -cpu 1 -benchmem
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkUnmarshal_ProtoBuf      4988426               236.7 ns/op            32 B/op          2 allocs/op
// BenchmarkUnmarshal_ProtoBuf      5131682               236.8 ns/op            32 B/op          2 allocs/op
// BenchmarkUnmarshal_ProtoBuf      5069283               237.0 ns/op            32 B/op          2 allocs/op
// BenchmarkUnmarshal_ProtoBuf      5141564               243.6 ns/op            32 B/op          2 allocs/op
// BenchmarkUnmarshal_ProtoBuf      5076106               237.3 ns/op            32 B/op          2 allocs/op
func BenchmarkUnmarshal_ProtoBuf(b *testing.B) {
	buf, err := proto.Marshal(&btpb)
	if err != nil {
		b.Error(err)
		return
	}
	var tmp BaseType_PB
	for i := 0; i < b.N; i++ {
		err = proto.Unmarshal(buf, &tmp)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkCountLength(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n, err := countLength(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err, n)
			return
		}
	}
}
