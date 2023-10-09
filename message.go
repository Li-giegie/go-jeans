package go_jeans

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync/atomic"
)

var _count uint32

type Message struct {
	Id   uint32
	Data []byte
}

func NewMsg(data []byte) *Message {
	return &Message{
		Id:   atomic.AddUint32(&_count, 1),
		Data: data,
	}
}

func (a *Message) Marshal() ([]byte, error) {
	var buf = new(bytes.Buffer)
	var err error
	if err = binary.Write(buf, binary.LittleEndian, a.Id); err != nil {
		return nil, err
	}
	_, err = buf.Write(a.Data)
	return buf.Bytes(), err
}

func (a *Message) Unmarshal(buf []byte) *Message {
	a.Id = binary.LittleEndian.Uint32(buf[:4])
	a.Data = buf[4:]
	return a
}

func (a *Message) String() string {
	return fmt.Sprintf("id:%v data:%s", a.Id, a.Data)
}

func (a *Message) debug() {
	fmt.Println(a.String())
}
