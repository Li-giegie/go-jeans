package go_jeans

import (
	"bytes"
	"encoding/binary"
	"io"
)

type MessageA struct {
	MsgId         uint32 `json:"MsgId,omitempty"`
	Msg           []byte `json:"Msg,omitempty"`
}

func NewMsgA(msg []byte) *MessageA  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageA{
		MsgId: count,
		Msg:   msg,
	}
}

func (a *MessageA) Marshal() (*bytes.Buffer,error) {
	var buf = new(bytes.Buffer)
	// 写入总长度
	err := binary.Write(buf,binary.LittleEndian,uint32(4+len(a.Msg)))
	if err != nil {
		return nil, err
	}

	if err = binary.Write(buf,binary.LittleEndian,a.MsgId);err != nil {
		return nil, err
	}
	_,err = buf.Write(a.Msg)
	return buf,err
}

func (a *MessageA) Unmarshal(conn io.Reader) (*MessageA,error) {
	buf,err := _read(conn)
	if err != nil {
		return nil,err
	}
	var tmp =  &MessageA{
		MsgId: binary.LittleEndian.Uint32(buf[:4]),
		Msg:   buf[4:],
	}

	return tmp,nil
}

