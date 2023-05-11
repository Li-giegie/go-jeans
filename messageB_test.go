package go_jeans

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)


func Test_MSGB(t *testing.T) {
	msg := NewMsgB([]byte("hello word0"),10,2,3)
	buf,err := msg.Marshal()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(buf.Bytes())
	var _msg = new(MessageB)
	_msg,err  = _msg.Unmarshal(buf)
	if err != nil {
		log.Fatalln(err)
	}

	jbuf,_ := json.Marshal(_msg)
	fmt.Println(string(jbuf))
}

func Test_MSGB_Proto(t *testing.T) {
	msg := NewMsgB_Proto([]byte("hello word~"),10,20,30)
	buf,err := msg.Marshal()
	if err != nil {
		log.Fatalln("-1 ",err)
	}

	var _msg = new(MessageB_Proto)

	_msg,err  = _msg.Unmarshal(buf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(_msg.Msg))
}

func TestB(t *testing.T) {
	msg:=NewMsgB([]byte("hello word B"),1,2,3)
	buf,err := msg.Marshal()
	fmt.Println(err)
	msgB,err := msg.Unmarshal(buf)
	fmt.Println(err)
	fmt.Println(string(msgB.Msg),msgB.MsgId)
}

func BenchmarkB(b *testing.B) {
	for i := 0; i < b.N; i++ {

		buf,_ := NewMsgB([]byte("hello word B"),1,2,3).Marshal()
		//fmt.Println(err)

		(&MessageB{}).Unmarshal(buf)
		//fmt.Println(err)
		//fmt.Println(string(msgB.Msg),msgB.MsgId)
	}
}

