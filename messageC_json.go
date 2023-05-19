package go_jeans

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"strconv"
)


type MessageC_Json struct {
	MsgId         string `json:"MsgId,omitempty"`
	MsgObj        interface{}
	SrcAddr       string `json:"SrcAddr,omitempty"`
	DestApi       string `json:"DestApi,omitempty"`
	DestAddr      string `json:"DestAddr,omitempty"`

}

func NewMsgC_Json(msgObj interface{},SrcAddr,DestApi,DestAddr string) *MessageC_Json  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageC_Json{
		MsgId:    SrcAddr + strconv.Itoa(int(count)),
		MsgObj:      msgObj,
		SrcAddr:  SrcAddr,
		DestApi:  DestApi,
		DestAddr: DestAddr,
	}
}

func (a *MessageC_Json) Marshal() (*bytes.Buffer,error) {

	jbuf,err := json.Marshal(a.MsgObj)
	if err != nil {
		return nil, err
	}
	var buf = new(bytes.Buffer)
	// 消息包 长度
	err = binary.Write(buf,binary.LittleEndian,uint32(16+len(jbuf)+len(a.MsgId)+len(a.DestAddr)+len(a.DestApi)+len(a.SrcAddr)))
	if err != nil {
		return nil, err
	}
	// 消息ID 长度
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(a.MsgId))); err != nil {
		return nil, err
	}
	// 消息 长度
	if err = binary.Write(buf,binary.LittleEndian, uint32(len(jbuf))); err != nil {
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
	if _,err = buf.Write(jbuf); err != nil {
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

func (a *MessageC_Json) Unmarshal(conn io.Reader,msgObj interface{}) (*MessageC_Json,error) {
	buf,err := Read(conn)
	if err != nil {
		return nil,err
	}
	id_len := binary.LittleEndian.Uint32(buf[:4])
	msg_len := binary.LittleEndian.Uint32(buf[4:8])
	srcaddr_len := binary.LittleEndian.Uint32(buf[8:12])
	DestApi_len := binary.LittleEndian.Uint32(buf[12:16])
	tmpN := 16+id_len

	a.MsgId = string(buf[16:tmpN])
	if err = json.Unmarshal(buf[tmpN:tmpN+msg_len],&msgObj); err != nil {
		return nil, err
	}

	tmpN = tmpN+msg_len
	a.SrcAddr = string(buf[tmpN:tmpN+srcaddr_len])
	tmpN = tmpN+srcaddr_len
	a.DestApi = string(buf[tmpN:tmpN+DestApi_len])
	tmpN = tmpN+DestApi_len
	a.DestAddr = string(buf[tmpN:])

	return a,nil
}