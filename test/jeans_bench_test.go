package test

import (
	"bytes"
	"encoding/json"
	jeans "github.com/Li-giegie/go-jeans"
	"testing"
)

func BenchmarkEncodeBaseSlice(b *testing.B) {
	d := GenBaseSlice()
	fields := d.FieldsInterface()
	for i := 0; i < b.N; i++ {
		data, err := jeans.Encode(fields...)
		if err != nil {
			b.Error(err, len(data))
			return
		}
	}
}

func BenchmarkEncodeFasterBaseSlice(b *testing.B) {
	d := GenBaseSlice()
	fields := d.FieldsInterface()
	buf := make([]byte, 0, 512)
	for i := 0; i < b.N; i++ {
		data, err := jeans.EncodeFaster(buf, fields...)
		if err != nil {
			b.Error(err, len(data))
			return
		}
	}
}

func BenchmarkEncodeJSONBaseSlice(b *testing.B) {
	d := GenBaseSlice()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(d)
		if err != nil {
			b.Error(err, len(data))
			return
		}
	}
}

func BenchmarkEncoderJSONBaseSlice(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := json.NewEncoder(buf)
	d := GenBaseSlice()
	for i := 0; i < b.N; i++ {
		err := encoder.Encode(d)
		if err != nil {
			b.Error(err)
			return
		}
		buf.Reset()
	}
}

func BenchmarkEncodeBufferBaseSlice(b *testing.B) {
	d := GenBaseSlice()
	fields := d.FieldsInterface()
	buf := jeans.NewBuffer(512)
	for i := 0; i < b.N; i++ {
		err := jeans.EncodeBuffer(buf, fields...)
		if err != nil {
			b.Error(err)
			return
		}
		buf.Reset()
	}
}

func BenchmarkDecode(b *testing.B) {
	d := GenBaseSlice()
	fields := d.FieldsInterface()
	data, err := jeans.Encode(fields...)
	if err != nil {
		b.Error(err)
		return
	}
	rd := NewBaseSlice()
	for i := 0; i < b.N; i++ {
		if err = jeans.Decode(data, rd.FieldsPointerToInterface()...); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkDecodeJSON(b *testing.B) {
	d := GenBaseSlice()
	data, err := json.Marshal(d)
	if err != nil {
		b.Error(err)
		return
	}
	rd := NewBaseSlice()
	for i := 0; i < b.N; i++ {
		if err = json.Unmarshal(data, rd); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkUniteEncodeDecode(b *testing.B) {
	d := GenBaseSlice()
	rd := NewBaseSlice()
	rFields := d.FieldsInterface()
	rdFields := rd.FieldsPointerToInterface()
	for i := 0; i < b.N; i++ {
		data, err := jeans.Encode(rFields...)
		if err != nil {
			b.Error(err)
			return
		}
		if err = jeans.Decode(data, rdFields...); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkUniteEncodeFasterDecode(b *testing.B) {
	d := GenBaseSlice()
	rd := NewBaseSlice()
	rFields := d.FieldsInterface()
	rdFields := rd.FieldsPointerToInterface()
	buf := make([]byte, 0, 512)
	for i := 0; i < b.N; i++ {
		data, err := jeans.EncodeFaster(buf, rFields...)
		if err != nil {
			b.Error(err)
			return
		}
		if err = jeans.Decode(data, rdFields...); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkUniteEncodeBufferDecode(b *testing.B) {
	d := GenBaseSlice()
	rd := NewBaseSlice()
	buf := jeans.NewBuffer(512)
	rFields := d.FieldsInterface()
	rdFields := rd.FieldsPointerToInterface()
	for i := 0; i < b.N; i++ {
		err := jeans.EncodeBuffer(buf, rFields...)
		if err != nil {
			b.Error(err)
			return
		}
		if err = jeans.Decode(buf.Data, rdFields...); err != nil {
			b.Error(err)
			return
		}
		buf.Reset()
	}
}

func BenchmarkUniteEncodeDecodeJSON(b *testing.B) {
	d := GenBaseSlice()
	rd := NewBaseSlice()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(d)
		if err != nil {
			b.Error(err)
			return
		}
		if err = json.Unmarshal(data, rd); err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkUniteEncoderDecodeJSON(b *testing.B) {
	d := GenBaseSlice()
	rd := NewBaseSlice()
	buf := bytes.NewBuffer(make([]byte, 0, 1000))
	encode := json.NewEncoder(buf)
	decode := json.NewDecoder(buf)
	for i := 0; i < b.N; i++ {
		err := encode.Encode(d)
		if err != nil {
			b.Error(err)
			return
		}
		if err = decode.Decode(rd); err != nil {
			b.Error(err)
			return
		}
	}
}
