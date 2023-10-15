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

// Benchmark_Encode-12     10352448               106.9 ns/op
// Benchmark_Encode-12     10415156               105.8 ns/op
// Benchmark_Encode-12     11582350               108.4 ns/op
// Benchmark_Encode-12     10175656               108.1 ns/op
// Benchmark_Encode-12     11331326               108.2 ns/op
func Benchmark_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Encode(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// Benchmark_EncodeWithByte-12      7484763               153.3 ns/op
// Benchmark_EncodeWithByte-12      8154261               149.8 ns/op
// Benchmark_EncodeWithByte-12      7853212               155.0 ns/op
// Benchmark_EncodeWithByte-12      7718113               158.1 ns/op
// Benchmark_EncodeWithByte-12      7619661               156.7 ns/op
func Benchmark_EncodeWithByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf, itemLen, err := EncodeWithLenByteItem(bt.bs, bt.i, bt.i8, bt.i16, bt.i32, bt.i64, bt.ui, bt.ui8, bt.ui16, bt.ui32, bt.ui64, bt.s, bt.b, bt.f32, bt.f64)
		if err != nil {
			b.Error(err, buf, itemLen)
			return
		}
	}
}

// Benchmark_Decode-12     23324839                55.53 ns/op
// Benchmark_Decode-12     24082989                55.65 ns/op
// Benchmark_Decode-12     23618234                52.14 ns/op
// Benchmark_Decode-12     24749820                52.58 ns/op
// Benchmark_Decode-12     24230968                52.99 ns/op
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

// BenchmarkMarshal_JSON-12         7749994               150.6 ns/op
// BenchmarkMarshal_JSON-12         8432164               144.9 ns/op
// BenchmarkMarshal_JSON-12         7165488               145.7 ns/op
// BenchmarkMarshal_JSON-12         8179947               149.2 ns/op
// BenchmarkMarshal_JSON-12         7217180               149.6 ns/op
func BenchmarkMarshal_JSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(bt)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// BenchmarkUnmarshal_JSON-12       6024643               191.4 ns/op
// BenchmarkUnmarshal_JSON-12       6236120               198.2 ns/op
// BenchmarkUnmarshal_JSON-12       6170451               195.0 ns/op
// BenchmarkUnmarshal_JSON-12       6171696               195.9 ns/op
// BenchmarkUnmarshal_JSON-12       5961298               198.6 ns/op
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

// BenchmarkMarshal_MsgPack-12      4822191               247.8 ns/op
// BenchmarkMarshal_MsgPack-12      4894285               250.1 ns/op
// BenchmarkMarshal_MsgPack-12      4561465               259.9 ns/op
// BenchmarkMarshal_MsgPack-12      4750412               255.1 ns/op
// BenchmarkMarshal_MsgPack-12      4505755               252.5 ns/op
func BenchmarkMarshal_MsgPack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := msgpack.Marshal(bt)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// BenchmarkUnmarshal_MsgPack-12            6741754               184.6 ns/op
// BenchmarkUnmarshal_MsgPack-12            6525402               184.6 ns/op
// BenchmarkUnmarshal_MsgPack-12            6755052               181.0 ns/op
// BenchmarkUnmarshal_MsgPack-12            6422390               181.7 ns/op
// BenchmarkUnmarshal_MsgPack-12            6677840               183.8 ns/op
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

// BenchmarkMarshal_ProtoBuf-12             3572792               330.4 ns/op
// BenchmarkMarshal_ProtoBuf-12             3600763               322.9 ns/op
// BenchmarkMarshal_ProtoBuf-12             3762686               316.0 ns/op
// BenchmarkMarshal_ProtoBuf-12             3806902               317.5 ns/op
// BenchmarkMarshal_ProtoBuf-12             3709130               320.5 ns/op
func BenchmarkMarshal_ProtoBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := proto.Marshal(&btpb)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// BenchmarkUnmarshal_ProtoBuf-12           4755826               265.8 ns/op
// BenchmarkUnmarshal_ProtoBuf-12           4398693               262.9 ns/op
// BenchmarkUnmarshal_ProtoBuf-12           4605128               255.1 ns/op
// BenchmarkUnmarshal_ProtoBuf-12           4787779               252.8 ns/op
// BenchmarkUnmarshal_ProtoBuf-12           4703775               252.0 ns/op
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
