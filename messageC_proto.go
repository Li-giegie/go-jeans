package go_jeans

import (
	"bytes"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"io"
	"strconv"
)

func NewMsgC_Proto(msg []byte,SrcAddr,DestApi,DestAddr string) *MessageC_Proto  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageC_Proto{
		MsgId:         SrcAddr + strconv.Itoa(int(count)),
		Msg:           msg,
		SrcAddr:       SrcAddr,
		DestApi:       DestApi,
		DestAddr:      DestAddr,
	}
}

func (a *MessageC_Proto) Marshal() (*bytes.Buffer,error) {
	var buf = new(bytes.Buffer)
	pbuf,err := proto.Marshal(a)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf,binary.LittleEndian,uint32(len(pbuf)))
	if err != nil {
		return nil, err
	}
	_,err = buf.Write(pbuf)
	if err != nil {
		return nil, err
	}
	return buf,nil
}

func (a *MessageC_Proto) Unmarshal(conn io.Reader) (*MessageC_Proto,error) {
	buf,err := _read(conn)
	if err != nil {
		return nil,err
	}
	var tmp = new(MessageC_Proto)

	return tmp,proto.Unmarshal(buf,tmp)
}