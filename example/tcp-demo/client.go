package main

import (
	"fmt"
	go_jeans "github.com/Li-giegie/go-jeans"
	"log"
	"net"
	"strconv"
)

func newClient(addr string)  {

	conn,err := net.Dial("tcp",addr)
	if err != nil {
		log.Fatalln("client 链接服务端失败：",err)
	}
	defer conn.Close()
	var i int
	for {
		i++
		err = go_jeans.Write(conn,[]byte("ping ------ from client "+ strconv.Itoa(i)))
		if err != nil {
			log.Fatalln(err)
		}
		if i % 10 == 0 {
			go_jeans.Write(conn,[]byte("exit"))
			fmt.Println("bye ~")
			break
		}
		buf,err := go_jeans.Read(conn)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(buf))

	}
}
