package main

import (
	"fmt"
	"github.com/Li-giegie/go-jeans"
	//"fmt"
	"log"
	"net"
)

func server()  {
	lis,err := net.Listen("tcp","127.0.0.1:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer lis.Close()
	for  {
		conn,err := lis.Accept()
		if err != nil {
			log.Println("accept err :",err)
			continue
		}
		
		go process(&conn)
	}
}

func process(conn *net.Conn)  {
	defer (*conn).Close()
	for  {

		//msgA,err := new(go_jeans.MessageA).Unmarshal(*conn)
		msgA,err := new(go_jeans.MessageA).Unmarshal(*conn)
		if err != nil {
			log.Println("read msg err:",err)
			break
		}
		fmt.Println("server receive: ",msgA.MsgId,string(msgA.Msg))
		msgA.Msg = []byte("server reply")
		wbuf,err := msgA.Marshal()
		if err != nil {
			log.Println("打包消息失败",err)
			break
		}
		_,err= (*conn).Write(wbuf.Bytes())
		if err != nil {
			log.Println("server 返回消息失败：",err)
			break
		}
	}

	log.Println("连接断开------")
}