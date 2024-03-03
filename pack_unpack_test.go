package go_jeans

import (
	"bytes"
	"testing"
)

func TestPackUnpack(t *testing.T) {
	// 10 byte
	data := []byte("hello word!")
	data = Pack(data)
	_, err := Unpack(bytes.NewBuffer(data))
	if err != nil {
		t.Error(err)
	}
}

func TestPackUnpackN(t *testing.T) {
	data := []byte("hello word")
	data, err := PackN(data, PacketHerderLenType_uint16)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = UnpackN(bytes.NewBuffer(data), PacketHerderLenType_uint16)
	if err != nil {
		t.Error(err)
		return
	}
}
