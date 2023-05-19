package main

import (
	"fmt"
	"github.com/Li-giegie/go-jeans"
	"log"
	"net"
)

func newClient(addr string)  {

	conn,err := net.Dial("tcp",addr)
	if err != nil {
		log.Fatalln("client 链接服务端失败：",err)
	}
	
	defer conn.Close()

	//buf,err := go_jeans.NewMsgA([]byte("hello ? i'm the client !")).Marshal()

	msgA := go_jeans.NewMsgA([]byte("hello ? i'm the client !"))
	buf,err := msgA.Marshal()
	if err!= nil {
		log.Fatalln("client：",err)
	}
	_,err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatalln("client 发送消息失败：",err)
	}
	_,err = conn.Write(buf.Bytes())
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
	//msgA,err := new(go_jeans.MessageA).Unmarshal(*conn)
	msgA := new(go_jeans.MessageA)
	msgA,err :=msgA.Unmarshal(*conn)
	if err != nil {
		return nil, err
	}
	return msgA,nil
}