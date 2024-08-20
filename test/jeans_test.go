package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jeans "github.com/Li-giegie/go-jeans"
	"os"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	d := GenBaseSlice()
	data, err := jeans.Encode(d.FieldsInterface()...)
	if err != nil {
		t.Error(err)
		return
	}
	rd := NewBaseSlice()
	if err = jeans.Decode(data, rd.FieldsPointerToInterface()...); err != nil {
		t.Error(err)
		return
	}
	if err = Equal(d, rd); err != nil {
		t.Error(err)
	}
}

func TestPackUnpack(t *testing.T) {
	data := make([]byte, 0xFF)
	data[0] = 55
	p1, err := jeans.PackN(data, jeans.PacketType8)
	if err != nil {
		t.Error("1", err)
		return
	}
	p2, err := jeans.UnpackN(bytes.NewBuffer(p1), jeans.PacketType8)
	if err != nil {
		t.Error("2", err)
		return
	}
	if !reflect.DeepEqual(p2, data) {
		t.Error("packN unpackN error")
		return
	}
}

var comErrTxt = []byte("[Error tagging]----------")

func Equal(a, b interface{}) error {
	if !reflect.DeepEqual(a, b) {
		data1, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			return err
		}
		data2, err := json.MarshalIndent(b, "", "  ")
		if err != nil {
			return err
		}
		data1 = append(data1, 10)
		for i := 0; i < len(data1); i++ {
			if len(data2) <= i {
				data1 = append(data1[:i], append(comErrTxt, data1[i:]...)...)
				data2 = append(data2[:i], append(comErrTxt, data2[i:]...)...)
				break
			}
			if data1[i] != data2[i] {
				data1 = append(data1[:i], append(comErrTxt, data1[i:]...)...)
				data2 = append(data2[:i], append(comErrTxt, data2[i:]...)...)
				fmt.Println("not equal char: \"", string(data1[i-10:i+10])+"\"")
				fmt.Println("not equal char: \"", string(data2[i-10:i+10])+"\"")
				break
			}
		}
		os.WriteFile("src.json", data1, 0666)
		os.WriteFile("dst.json", data2, 0666)
		return errors.New("not equal")
	}
	return nil
}

func TestGenerateJeansFuncs(t *testing.T) {
	f, err := os.OpenFile("templateFuncs.txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()
	b := new(Base)
	genRandValue(b)
	if err = jeans.GenerateJeansFuncs(f, b, jeans.ModeType_All); err != nil {
		t.Error(err)
		return
	}
	data, err := b.Encode()
	if err != nil {
		t.Error(err)
		return
	}
	rb := new(Base)
	if err = rb.Decode(data); err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(b, rb) {
		t.Error("Decode invalid")
		return
	}

	bs := GenBaseSlice()
	if err = jeans.GenerateJeansFuncs(f, bs, jeans.ModeType_All); err != nil {
		t.Error(err)
		return
	}
	data, err = bs.Encode()
	if err != nil {
		t.Error(err)
		return
	}
	rbs := NewBaseSlice()
	if err = rbs.Decode(data); err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(bs, rbs) {
		t.Error("decode BaseSlice invalid")
	}
	fmt.Println(bs.String())
	fmt.Println(rbs.String())
}

type A struct {
	B
	//b2 B
	a []string
}

type B struct {
	b1 string `jeans:"enable"`
	*B
	b2 *int
	C
}

type C struct {
	i int
}

type u8 uint8
