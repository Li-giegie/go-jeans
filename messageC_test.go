package go_jeans

import (
	"fmt"
	"testing"
)

func Test_MSGC(t *testing.T) {
	msg := NewMsgC([]byte("hello i'm msg ~"),"srcaddr","destapi","destaddr")
	fmt.Println(msg)
	buf,err := msg.Marshal()
	fmt.Println(buf.Bytes(),buf.String(),err)
	var _msg *MessageC
	fmt.Println(_msg.Unmarshal(buf))
}

func Test_MSGC_proto(t *testing.T) {
	msg := NewMsgC_Proto([]byte("hello i'm msg ~"),"srcaddr","destapi","destaddr")
	//fmt.Println(msg)
	buf,err := msg.Marshal()
	fmt.Println(err)
	var msgC MessageC_Proto
	fmt.Println(msgC.Unmarshal(buf))
}
