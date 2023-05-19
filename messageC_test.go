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

func TestNewMsgC_Json(t *testing.T) {
	//msgC := NewMsgC_Json(&A{Str: "hello word"},"a","b","c")
	//buf,err := msgC.Marshal()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//var _a A
	//msgC_json,err := new(MessageC_Json).Unmarshal(buf,&_a)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(_a,msgC_json)
}