package go_jeans

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
)

type MessageA_Json struct {
	MsgId         uint32 `json:"MsgId,omitempty"`
	MsgObj interface{}
}

func NewMsgA_Json(msgObj interface{}) *MessageA_Json  {
	lock.Lock()
	count++
	defer lock.Unlock()
	return &MessageA_Json{
		MsgId: count,
		MsgObj: msgObj,
	}
}

func (a *MessageA_Json) Marshal() (*bytes.Buffer,error) {

	jbuf,err := json.Marshal(a.MsgObj)
	if err != nil {
		return nil, err
	}
	var buf = new(bytes.Buffer)
	// 写入总长度
	err = binary.Write(buf,binary.LittleEndian,uint32(4+len(jbuf)))
	if err != nil {
		return nil, err
	}

	if err = binary.Write(buf,binary.LittleEndian,a.MsgId);err != nil {
		return nil, err
	}
	_,err = buf.Write(jbuf)
	return buf,err
}

func (a *MessageA_Json) Unmarshal(conn io.Reader,msgObj interface{}) (*MessageA_Json,error) {
	buf,err := Read(conn)
	if err != nil {
		return nil,err
	}

	err = json.Unmarshal(buf[4:],msgObj)
	if err != nil {
		return nil, err
	}

	a.MsgId = binary.LittleEndian.Uint32(buf[:4])
	return a,nil
}

