package main

import (
	"fmt"
	go_jeans "github.com/Li-giegie/go-jeans"
	"log"
	"net"
	"strconv"
	"time"
)

func newClient(addr string)  {
	conn,err := net.Dial("tcp",addr)
	if err != nil {
		log.Fatalln("client 链接服务端失败：",err)
	}
	defer conn.Close()
	for i:=0;i<10;i++{
		err = go_jeans.Write(conn,[]byte("ping ------ from client "+ strconv.Itoa(i)))
		if err != nil {
			log.Fatalln(err)
		}
		buf,err := go_jeans.Read(conn)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("client receive:",string(buf))
		time.Sleep(time.Second*1)
	}
}
