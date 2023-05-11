package go_jeans

import (
	"bytes"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

func NewMsgB_Proto(msg []byte,SrcAddr,DestApi,DestAddr uint32) *MessageB_Proto  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageB_Proto{
		MsgId:         count,
		Msg:           msg,
		SrcAddr:       SrcAddr,
		DestApi:       DestApi,
		DestAddr:      DestAddr,
	}
}

func (a *MessageB_Proto) Marshal() (*bytes.Buffer,error) {
	var buf = new(bytes.Buffer)

	pbuf,err := proto.Marshal(a)
	if err != nil {
		return nil, err
	}
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(pbuf))); err != nil {
		return nil, err
	}
	_,err = buf.Write(pbuf)
	return buf,err
}

func (a *MessageB_Proto) Unmarshal(conn io.Reader) (*MessageB_Proto,error) {
	buf,err := _read(conn)

	if err != nil {
		return nil,err
	}

	var tmp = new(MessageB_Proto)

	return tmp,proto.Unmarshal(buf,tmp)
}