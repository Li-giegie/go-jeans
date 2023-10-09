package main

import (
	"fmt"
	go_jeans "github.com/Li-giegie/go-jeans"
	"log"
	"net"
)

func server(addr string)  {
	lis,err := net.Listen("tcp",addr)
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
		go process(conn)
	}
}

func process(conn net.Conn)  {
	defer (conn).Close()
	for  {
		buf,err := go_jeans.Read(conn)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("server receive:",string(buf))
		err = go_jeans.Write(conn,[]byte("pong pong pong ~"))
		if err != nil {
			log.Fatalln(err)
		}
	}

}