package go_jeans

import (
	"bytes"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
)

func NewMsgA_Proto(msg []byte) *MessageA_Proto  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageA_Proto{
		MsgId: count,
		Msg:   msg,
	}
}

func (a *MessageA_Proto) Marshal() (*bytes.Buffer,error) {
	var buf = new(bytes.Buffer)
	pbuf,err := proto.Marshal(a)
	if err != nil {
		return nil, err
	}

	if err = binary.Write(buf,binary.LittleEndian,uint32(len(pbuf))); err != nil {
		return nil, err
	}
	_,err = buf.Write(pbuf)

	return buf,err
}

func (a *MessageA_Proto) Unmarshal(conn io.Reader) (*MessageA_Proto,error) {
	buf,err := _read(conn)
	if err != nil {
		return nil,err
	}
	var tmp = new(MessageA_Proto)

	return tmp,proto.Unmarshal(buf,tmp)
}
