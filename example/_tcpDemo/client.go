package main

import (
	"fmt"
	"github.com/Li-giegie/go_jeans"
	"log"
	"net"
)

func newClient(addr string)  {

	conn,err := net.Dial("tcp",addr)
	if err != nil {
		log.Fatalln("client 链接服务端失败：",err)
	}
	
	defer conn.Close()

	msgA := go_jeans.NewMsgA_String("hello i'm client !")
	buf,err := msgA.Bytes()
	if err!= nil {
		log.Fatalln("client：",err)
	}
	_,err = conn.Write(buf)
	if err != nil {
		log.Fatalln("client 发送消息失败：",err)
	}
	fmt.Println("client request:",msgA.MsgId,string(msgA.Msg))
	reply,err := client_process(&conn)
	if err != nil {
		log.Fatalln("client 接收消息失败：",err)
	}
	replyA := reply.(*go_jeans.MessageA)
	fmt.Println("client receive:",replyA.MsgId,string(replyA.Msg))

}

func client_process(conn *net.Conn) (interface{},error) {
	msgA,err :=go_jeans.UnpackA(*conn)
	if err != nil {
		return nil, err
	}
	return msgA,nil
}