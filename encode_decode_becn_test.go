package go_jeans

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack/v5"
	"testing"
)

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkEncode$
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkEncode         11916831                87.44 ns/op          112 B/op          1 allocs/op
// BenchmarkEncode         13481463                95.90 ns/op          176 B/op          1 allocs/op
// BenchmarkEncode         13797559                92.99 ns/op          160 B/op          1 allocs/op
// BenchmarkEncode         12452278                94.96 ns/op          176 B/op          1 allocs/op
// BenchmarkEncode         13829329                94.83 ns/op          160 B/op          1 allocs/op
func BenchmarkEncode(b *testing.B) {
	base := NewBase()
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = Encode(base.FieldsToInterface()...); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkDecode$
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkDecode         33155032                36.66 ns/op            0 B/op          0 allocs/op
// BenchmarkDecode         31286827                36.86 ns/op            0 B/op          0 allocs/op
// BenchmarkDecode         34148534                36.32 ns/op            0 B/op          0 allocs/op
// BenchmarkDecode         33340185                36.82 ns/op            0 B/op          0 allocs/op
// BenchmarkDecode         34769002                37.08 ns/op            0 B/op          0 allocs/op
func BenchmarkDecode(b *testing.B) {
	base := NewBase()
	fields := base.FieldsToInterface()
	var buf []byte
	var err error
	if buf, err = Encode(fields...); err != nil {
		b.Error(err)
		return
	}
	var decodeBase = new(Base)
	var decodeFields = decodeBase.FieldsPointerToInterface()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = Decode(buf, decodeFields...); err != nil {
			b.Error(err)
			return
		}
	}

}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkEncodeAndDecode$
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkEncodeAndDecode         8903648               130.6 ns/op           176 B/op          1 allocs/op
// BenchmarkEncodeAndDecode         8226502               124.2 ns/op           128 B/op          1 allocs/op
// BenchmarkEncodeAndDecode         9377050               123.4 ns/op           112 B/op          1 allocs/op
// BenchmarkEncodeAndDecode         9243880               124.1 ns/op           128 B/op          1 allocs/op
// BenchmarkEncodeAndDecode         9512017               124.0 ns/op           112 B/op          1 allocs/op
func BenchmarkEncodeAndDecode(b *testing.B) {
	base := NewBase()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err := Encode(base.FieldsToInterface()...)
		if err != nil {
			b.Error(err)
			return
		}
		var decodeBase = new(Base)
		if err = Decode(buf, decodeBase.FieldsPointerToInterface()...); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkEncodeFaster$
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkEncodeFaster   21978625                54.35 ns/op            0 B/op          0 allocs/op
// BenchmarkEncodeFaster   22538724                53.25 ns/op            0 B/op          0 allocs/op
// BenchmarkEncodeFaster   22194970                53.26 ns/op            0 B/op          0 allocs/op
// BenchmarkEncodeFaster   23554998                54.49 ns/op            0 B/op          0 allocs/op
// BenchmarkEncodeFaster   22003531                53.31 ns/op            0 B/op          0 allocs/op
func BenchmarkEncodeFaster(b *testing.B) {
	base := NewBase()
	var buf = make([]byte, 0, CountLength(base.FieldsToInterface()...))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := EncodeFaster(buf, base.FieldsToInterface()...); err != nil {
			b.Error(err)
			return
		}
	}
}

//---------------------------------------- JSON

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkJsonMarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkJsonMarshal     1490142               686.7 ns/op           224 B/op          1 allocs/op
// BenchmarkJsonMarshal     1471344               828.9 ns/op           448 B/op          2 allocs/op
// BenchmarkJsonMarshal     1522213               725.7 ns/op           256 B/op          1 allocs/op
// BenchmarkJsonMarshal     1692051               723.4 ns/op           256 B/op          1 allocs/op
// BenchmarkJsonMarshal     1669550               681.4 ns/op           208 B/op          1 allocs/op         1 allocs/op
func BenchmarkJsonMarshal(b *testing.B) {
	base := NewBase()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := json.Marshal(base); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkJsonUnmarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkJsonUnmarshal    320841              3919 ns/op             408 B/op         18 allocs/op
// BenchmarkJsonUnmarshal    374895              4095 ns/op             408 B/op         19 allocs/op
// BenchmarkJsonUnmarshal    272271              3763 ns/op             392 B/op         18 allocs/op
// BenchmarkJsonUnmarshal    327273              3675 ns/op             328 B/op         17 allocs/op
// BenchmarkJsonUnmarshal    285213              3797 ns/op             392 B/op         18 allocs/op
func BenchmarkJsonUnmarshal(b *testing.B) {
	base := NewBase()
	buf, err := json.Marshal(base)
	if err != nil {
		b.Error(err)
		return
	}
	var unmarshalBase = new(Base)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(buf, unmarshalBase); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkJsonMarshalAndUnmarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkJsonMarshalAndUnmarshal          265867              5254 ns/op             952 B/op         21 allocs/op
// BenchmarkJsonMarshalAndUnmarshal          228480              5212 ns/op             880 B/op         22 allocs/op
// BenchmarkJsonMarshalAndUnmarshal          256641              5231 ns/op             936 B/op         21 allocs/op
// BenchmarkJsonMarshalAndUnmarshal          221882              5150 ns/op             856 B/op         21 allocs/op
// BenchmarkJsonMarshalAndUnmarshal          247593              5275 ns/op             944 B/op         21 allocs/op
func BenchmarkJsonMarshalAndUnmarshal(b *testing.B) {
	base := NewBase()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err := json.Marshal(base)
		if err != nil {
			b.Error(err)
			return
		}
		unmarshalBase := new(Base)
		if err = json.Unmarshal(buf, unmarshalBase); err != nil {
			b.Error(err)
			return
		}
	}
}

// ---------------------------------------- protoBuf
// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkProtoBufMarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkProtoBufMarshal         5041351               244.2 ns/op           112 B/op          1 allocs/op
// BenchmarkProtoBufMarshal         4920415               239.5 ns/op           128 B/op          1 allocs/op
// BenchmarkProtoBufMarshal         4819473               235.9 ns/op            80 B/op          1 allocs/op
// BenchmarkProtoBufMarshal         4954760               227.7 ns/op            64 B/op          1 allocs/op
// BenchmarkProtoBufMarshal         4888285               241.4 ns/op           112 B/op          1 allocs/op
func BenchmarkProtoBufMarshal(b *testing.B) {
	var base = new(BaseType_PB)
	if err := faker.FakeData(base); err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := proto.Marshal(base); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkProtoBufUnmarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkProtoBufUnmarshal       5480512               239.6 ns/op            80 B/op          2 allocs/op
// BenchmarkProtoBufUnmarshal       4935399               252.3 ns/op           128 B/op          2 allocs/op
// BenchmarkProtoBufUnmarshal       4783920               244.1 ns/op           128 B/op          2 allocs/op
// BenchmarkProtoBufUnmarshal       4892127               249.4 ns/op           128 B/op          2 allocs/op
// BenchmarkProtoBufUnmarshal       4925493               241.5 ns/op           112 B/op          2 allocs/op
func BenchmarkProtoBufUnmarshal(b *testing.B) {
	var base = new(BaseType_PB)
	if err := faker.FakeData(base); err != nil {
		b.Error(err)
		return
	}
	buf, err := proto.Marshal(base)
	if err != nil {
		b.Error(err)
		return
	}
	protoBufBase := new(BaseType_PB)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = proto.Unmarshal(buf, protoBufBase); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkProtoBufMarshalAndUnmarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkProtoBufMarshalAndUnmarshal     2071735               603.5 ns/op           448 B/op          4 allocs/op
// BenchmarkProtoBufMarshalAndUnmarshal     2070318               562.4 ns/op           320 B/op          4 allocs/op
// BenchmarkProtoBufMarshalAndUnmarshal     2127872               586.8 ns/op           320 B/op          4 allocs/op
// BenchmarkProtoBufMarshalAndUnmarshal     1984470               570.0 ns/op           336 B/op          4 allocs/op
// BenchmarkProtoBufMarshalAndUnmarshal     2124021               564.4 ns/op           288 B/op          4 allocs/op
func BenchmarkProtoBufMarshalAndUnmarshal(b *testing.B) {
	var base = new(BaseType_PB)
	if err := faker.FakeData(base); err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err := proto.Marshal(base)
		if err != nil {
			b.Error(err)
			return
		}
		protoBufBase := new(BaseType_PB)
		if err = proto.Unmarshal(buf, protoBufBase); err != nil {
			b.Error(err)
			return
		}
	}
}

// ---------------------------------------messagePack
// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkMsgPackMarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkMsgPackMarshal          1552291               773.1 ns/op           496 B/op          4 allocs/op
// BenchmarkMsgPackMarshal          1531084               780.6 ns/op           496 B/op          4 allocs/op
// BenchmarkMsgPackMarshal          1527271               800.7 ns/op           496 B/op          4 allocs/op
// BenchmarkMsgPackMarshal          1474113               789.5 ns/op           496 B/op          4 allocs/op
// BenchmarkMsgPackMarshal          1533616               788.1 ns/op           496 B/op          4 allocs/op
func BenchmarkMsgPackMarshal(b *testing.B) {
	base := NewBase()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := msgpack.Marshal(base); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkMsgPackUnmarshal
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkMsgPackUnmarshal        1000000              1142 ns/op              80 B/op          2 allocs/op
// BenchmarkMsgPackUnmarshal        1000000              1168 ns/op              80 B/op          2 allocs/op
// BenchmarkMsgPackUnmarshal        1000000              1154 ns/op              80 B/op          2 allocs/op
// BenchmarkMsgPackUnmarshal        1000000              1215 ns/op              80 B/op          2 allocs/op
// BenchmarkMsgPackUnmarshal        1000000              1169 ns/op              80 B/op          2 allocs/op
func BenchmarkMsgPackUnmarshal(b *testing.B) {
	base := NewBase()
	buf, err := msgpack.Marshal(base)
	if err != nil {
		b.Error(err)
		return
	}
	var unmarshalBase = new(Base)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = msgpack.Unmarshal(buf, unmarshalBase); err != nil {
			b.Error(err)
			return
		}
	}
}

// ----------------------------------------gob
// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkGobEncode
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkGobEncode       2749702               438.9 ns/op           244 B/op          0 allocs/op
// BenchmarkGobEncode       2749166               502.2 ns/op           585 B/op          0 allocs/op
// BenchmarkGobEncode       2960156               382.3 ns/op           226 B/op          0 allocs/op
// BenchmarkGobEncode       2627322               412.3 ns/op           255 B/op          0 allocs/op
// BenchmarkGobEncode       2881346               430.9 ns/op           558 B/op          0 allocs/op
func BenchmarkGobEncode(b *testing.B) {
	base := NewBase()
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := enc.Encode(base); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkGobDecode
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans/test
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkGobDecode         53432             22287 ns/op            8672 B/op        253 allocs/op
// BenchmarkGobDecode         52701             21815 ns/op            8672 B/op        253 allocs/op
// BenchmarkGobDecode         53355             21709 ns/op            8672 B/op        253 allocs/op
// BenchmarkGobDecode         53977             21717 ns/op            8672 B/op        253 allocs/op
// BenchmarkGobDecode         54555             21540 ns/op            8672 B/op        253 allocs/op
func BenchmarkGobDecode(b *testing.B) {
	base := NewBase()
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(base)
	if err != nil {
		b.Error("1", err)
		return
	}
	gobBase := new(Base)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dc := gob.NewDecoder(bytes.NewBuffer(buf.Bytes()))
		b.StartTimer()
		if err = dc.Decode(gobBase); err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkEncodeSlice
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkEncodeSlice      551822              2084 ns/op             352 B/op          1 allocs/op
// BenchmarkEncodeSlice      598401              2088 ns/op             352 B/op          1 allocs/op
// BenchmarkEncodeSlice      555241              2103 ns/op             352 B/op          1 allocs/op
// BenchmarkEncodeSlice      583578              2085 ns/op             352 B/op          1 allocs/op
// BenchmarkEncodeSlice      568130              2108 ns/op             352 B/op          1 allocs/op
func BenchmarkEncodeSlice(b *testing.B) {
	s := NewSlice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncodeSlice(s.FieldsToInterface()...)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkDecodeSlice
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkDecodeSlice     2577890               462.0 ns/op           360 B/op         13 allocs/op
// BenchmarkDecodeSlice     2589948               441.3 ns/op           360 B/op         13 allocs/op
// BenchmarkDecodeSlice     2715187               442.8 ns/op           360 B/op         13 allocs/op
// BenchmarkDecodeSlice     2699566               445.8 ns/op           360 B/op         13 allocs/op
// BenchmarkDecodeSlice     2645998               448.0 ns/op           360 B/op         13 allocs/op
func BenchmarkDecodeSlice(b *testing.B) {
	s := NewSlice()
	buf, err := EncodeSlice(s.FieldsToInterface()...)
	if err != nil {
		b.Error(err)
		return
	}
	decodeS := new(Slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = DecodeSlice(buf, decodeS.FieldsPointerToInterface()...)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

// go test -run=none -cpu 1 -count 5 -benchmem -bench=BenchmarkJsonUnmarshalSlice
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go-jeans
// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkJsonUnmarshalSlice       138262              8710 ns/op             352 B/op         35 allocs/op
// BenchmarkJsonUnmarshalSlice       141118              8647 ns/op             352 B/op         35 allocs/op
// BenchmarkJsonUnmarshalSlice       140202              8736 ns/op             352 B/op         35 allocs/op
// BenchmarkJsonUnmarshalSlice       127918              9126 ns/op             352 B/op         35 allocs/op
// BenchmarkJsonUnmarshalSlice       135925              8680 ns/op             352 B/op         35 allocs/op
func BenchmarkJsonUnmarshalSlice(b *testing.B) {
	s := NewSlice()
	buf, err := json.Marshal(s)
	if err != nil {
		b.Error(err)
		return
	}
	decodeS := new(Slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(buf, decodeS); err != nil {
			b.Error(err)
			return
		}
	}
}
