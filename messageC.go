package go_jeans

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
)


type MessageC struct {
	MsgId         string `json:"MsgId,omitempty"`
	Msg           []byte `json:"Msg,omitempty"`
	SrcAddr       string `json:"SrcAddr,omitempty"`
	DestApi       string `json:"DestApi,omitempty"`
	DestAddr      string `json:"DestAddr,omitempty"`

}

func NewMsgC(msg []byte,SrcAddr,DestApi,DestAddr string) *MessageC  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageC{
		MsgId:    SrcAddr + strconv.Itoa(int(count)),
		Msg:      msg,
		SrcAddr:  SrcAddr,
		DestApi:  DestApi,
		DestAddr: DestAddr,
	}
}

func (a *MessageC) Marshal() (*bytes.Buffer,error) {
	var buf = new(bytes.Buffer)
	// 消息包 长度
	err := binary.Write(buf,binary.LittleEndian,uint32(16+len(a.Msg)+len(a.MsgId)+len(a.DestAddr)+len(a.DestApi)+len(a.SrcAddr)))
	if err != nil {
		return nil, err
	}
	// 消息ID 长度
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(a.MsgId))); err != nil {
		return nil, err
	}
	// 消息 长度
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(a.Msg))); err != nil {
		return nil, err
	}
	// 消息源地址 长度
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(a.SrcAddr))); err != nil {
		return nil, err
	}
	// 消息目的api 长度
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(a.DestApi))); err != nil {
		return nil, err
	}
	//// 消息目的地址 长度
	//if err = binary.Write(buf,binary.LittleEndian, uint32(len(a.DestAddr))); err != nil {
	//	return nil, err
	//}

	if _,err = buf.WriteString(a.MsgId); err != nil {
		return nil, err
	}
	if _,err = buf.Write(a.Msg); err != nil {
		return nil, err
	}
	if _,err = buf.WriteString(a.SrcAddr); err != nil {
		return nil, err
	}
	if _,err = buf.WriteString(a.DestApi); err != nil {
		return nil, err
	}
	_,err = buf.WriteString(a.DestAddr)

	return buf,err
}

func (a *MessageC) Unmarshal(conn io.Reader) (*MessageC,error) {
	buf,err := _read(conn)
	if err != nil {
		return nil,err
	}
	var tmp = new(MessageC)

	id_len := binary.LittleEndian.Uint32(buf[:4])
	msg_len := binary.LittleEndian.Uint32(buf[4:8])
	srcaddr_len := binary.LittleEndian.Uint32(buf[8:12])
	DestApi_len := binary.LittleEndian.Uint32(buf[12:16])
	tmpN := 16+id_len
	tmp.MsgId = string(buf[16:tmpN])
	tmp.Msg = buf[tmpN:tmpN+msg_len]
	tmpN = tmpN+msg_len
	tmp.SrcAddr = string(buf[tmpN:tmpN+srcaddr_len])
	tmpN = tmpN+srcaddr_len
	tmp.DestApi = string(buf[tmpN:tmpN+DestApi_len])
	tmpN = tmpN+DestApi_len
	tmp.DestAddr = string(buf[tmpN:])

	return tmp,nil
}