package go_jeans

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPackUnpack(t *testing.T) {
	// 10 byte
	var data = []byte("hello word")
	fmt.Println(len(data),data)
	data = Pack(data)
	fmt.Println(len(data),data)

	data2,err := Unpack(bytes.NewBuffer(data))
	fmt.Println(err)
	fmt.Println(len(data2),string(data2),data2)
}

func TestPackUnpackN(t *testing.T) {
	// 10 byte
	var data = []byte("hello word")
	fmt.Println(len(data),data)
	data,err := PackN(data,PacketHerderLen_16)
	fmt.Println(len(data),data)

	data2,err := UnpackN(bytes.NewBuffer(data),PacketHerderLen_16)
	fmt.Println(err)
	fmt.Println(len(data2),string(data2),data2)
}
