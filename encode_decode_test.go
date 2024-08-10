package go_jeans

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestEncode(t *testing.T) {
	base := NewBase()
	slice := NewSlice()
	fields := base.FieldsToInterface()
	fields = append(fields, slice.FieldsToInterface()...)
	decodeBase := new(Base)
	decodeSlice := new(Slice)
	decodeFields := decodeBase.FieldsPointerToInterface()
	decodeFields = append(decodeFields, decodeSlice.FieldsPointerToInterface()...)
	t1 := time.Now()
	buf := make([]byte, 0, 472)
	for i := 0; i < 1000000; i++ {
		result, err := EncodeFaster(buf, fields...)
		if err != nil {
			t.Error(err)
			return
		}
		err = Decode(result, decodeFields...)
		if err != nil {
			t.Error(err)
			return
		}
	}
	println(time.Since(t1).String())
	if !reflect.DeepEqual(base, decodeBase) || !reflect.DeepEqual(slice, decodeSlice) {
		t.Error("DeepEqual fail")
		return
	}
	println("TestEncode pass")
}

func BenchmarkEncode(b *testing.B) {
	base := NewBase()
	slice := NewSlice()
	fields := base.FieldsToInterface()
	fields = append(fields, slice.FieldsToInterface()...)
	decodeBase := new(Base)
	decodeSlice := new(Slice)
	decodeFields := decodeBase.FieldsPointerToInterface()
	decodeFields = append(decodeFields, decodeSlice.FieldsPointerToInterface()...)
	buf := make([]byte, 0, 472)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := EncodeFaster(buf, fields...)
		if err != nil {
			b.Error(err)
			return
		}
		err = Decode(data, decodeFields...)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

type bs struct {
	B *Base
	S *Slice
}

func TestJson(t *testing.T) {
	bs1 := new(bs)
	bs1.S = NewSlice()
	bs1.B = NewBase()
	bs2 := new(bs)
	t1 := time.Now()
	for i := 0; i < 1000000; i++ {
		data, err := json.Marshal(bs1)
		if err != nil {
			t.Error(err)
			return
		}
		err = json.Unmarshal(data, bs2)
		if err != nil {
			t.Error(err)
			return
		}
	}
	println(time.Since(t1).String())
	if !reflect.DeepEqual(bs1, bs2) {
		t.Error("TestJson err")
	}
}

func BenchmarkJson(b *testing.B) {
	bs1 := new(bs)
	bs1.S = NewSlice()
	bs1.B = NewBase()
	bs2 := new(bs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(bs1)
		if err != nil {
			b.Error(err)
			return
		}
		err = json.Unmarshal(data, bs2)
		if err != nil {
			b.Error(err)
			return
		}
	}
}

func BenchmarkEncodeBase(b *testing.B) {
	base := NewBase()
	fields := base.FieldsToInterface()
	resB := new(Base)
	resFields := resB.FieldsPointerToInterface()
	for i := 0; i < b.N; i++ {
		result, err := Encode(fields...)
		if err != nil {
			b.Error(err)
			return
		}
		if err = Decode(result, resFields...); err != nil {
			b.Error(err)
			return
		}
	}
}

func TestEncodeFaster(t *testing.T) {
	buf := make([]byte, 0, 89)
	base := NewBase()
	result, err := EncodeFaster(buf, base.FieldsToInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	decodeBase := new(Base)
	err = Decode(result, decodeBase.FieldsPointerToInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(base, decodeBase) {
		t.Error("DeepEqual fail")
		return
	}
	fmt.Println(base, decodeBase)
}

func TestEncodeSlice(t *testing.T) {
	for i := 0; i < 1000; i++ {
		s := NewSlice()
		buf, err := EncodeSlice(s.FieldsToInterface()...)
		if err != nil {
			t.Error(err)
			return
		}
		decodeUi32s := new(Slice)
		if err = DecodeSlice(buf, decodeUi32s.FieldsPointerToInterface()...); err != nil {
			t.Error(err)
			return
		}
		if !reflect.DeepEqual(s, decodeUi32s) {
			buf1, err1 := json.MarshalIndent(s, "", "\t")
			buf2, err2 := json.MarshalIndent(decodeUi32s, "", "\t")
			if err1 != nil || err2 != nil {
				fmt.Println("json err: ", err1, err2)
				return
			}
			err1 = os.WriteFile("./t1.json", buf1, 0666)
			err2 = os.WriteFile("./t2.json", buf2, 0666)
			if err1 != nil || err2 != nil {
				fmt.Println("write file err: ", err1, err2)
				return
			}
			t.Error("DeepEqual fail")
			return
		}
	}
}

func TestEncodeBaseAndSlice(t *testing.T) {
	for k := 0; k < 100; k++ {
		b := NewBase()
		s := NewSlice()
		args := make([]interface{}, 0, b.FieldNum()+s.FieldNum())
		for _, i := range b.FieldsToInterface() {
			args = append(args, i)
		}
		for _, i := range s.FieldsToInterface() {
			args = append(args, i)
		}
		buf, err := Encode(args...)
		if err != nil {
			t.Error(err)
			return
		}
		rb := new(Base)
		rs := new(Slice)
		rargs := make([]interface{}, 0, rb.FieldNum()+rs.FieldNum())
		for _, i := range rb.FieldsPointerToInterface() {
			rargs = append(rargs, i)
		}
		for _, i := range rs.FieldsPointerToInterface() {
			rargs = append(rargs, i)
		}
		if err = Decode(buf, rargs...); err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(b, rb) {
			t.Error("base DeepEqual fail")
		}
		if !reflect.DeepEqual(s, rs) {
			t.Error("slice DeepEqual fail")
		}
	}
}
