package go_jeans

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)


type MessageB_Json struct {
	MsgId         uint32 `json:"MsgId,omitempty"`
	msgObj        interface{}
	SrcAddr       uint32 `json:"SrcAddr,omitempty"`
	DestApi       uint32 `json:"DestApi,omitempty"`
	DestAddr      uint32 `json:"DestAddr,omitempty"`
}

func NewMsgB_Json(msgObj interface{},SrcAddr,DestApi,DestAddr uint32) *MessageB_Json  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageB_Json{
		MsgId:    count,
		msgObj:      msgObj,
		SrcAddr:  SrcAddr,
		DestApi:  DestApi,
		DestAddr: DestAddr,
	}
}

func (a *MessageB_Json) Marshal() (*bytes.Buffer,error) {

	jbuf,err := json.Marshal(a.msgObj)
	if err != nil {
		return nil, err
	}

	var buf = new(bytes.Buffer)

	err = binary.Write(buf,binary.LittleEndian,uint32(16+len(jbuf)))
	if err != nil {
		return nil, err
	}

	if err = binary.Write(buf,binary.LittleEndian,a.MsgId); err != nil {
		return nil, err
	}

	if err = binary.Write(buf,binary.LittleEndian,a.SrcAddr); err != nil {
		return nil, err
	}
	fmt.Println("addr ",buf.Bytes())
	if err = binary.Write(buf,binary.LittleEndian,a.DestApi); err != nil {
		return nil, err
	}
	fmt.Println("dapi",buf.Bytes())
	if err = binary.Write(buf,binary.LittleEndian,a.DestAddr); err != nil {
		return nil, err
	}
	fmt.Println("",buf.Bytes())
	if _,err = buf.Write(jbuf);err != nil {
		return nil, err
	}

	return buf,err
}

func (a *MessageB_Json) Unmarshal(conn io.Reader,msgObj interface{}) (*MessageB_Json,error) {
	buf,err := _read(conn)
	if err != nil {
		return nil,err
	}


	tmpN := 4
	a.MsgId = binary.LittleEndian.Uint32(buf[:tmpN])
	a.SrcAddr = binary.LittleEndian.Uint32(buf[tmpN:tmpN+4])
	tmpN += 4
	a.DestApi = binary.LittleEndian.Uint32(buf[tmpN:tmpN+4])
	tmpN += 4
	a.DestAddr = binary.LittleEndian.Uint32(buf[tmpN:tmpN+4])
	tmpN += 4
	if err = json.Unmarshal(buf[tmpN:],msgObj); err != nil {
		return nil, err
	}

	return a,nil
}