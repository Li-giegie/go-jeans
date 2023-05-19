package main

import (
	"fmt"
	go_jeans "github.com/Li-giegie/go-jeans"
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
		buf,err := go_jeans.ReadN(*conn,go_jeans.PacketHerderLen_16)
		if err != nil {
			log.Fatalln(err)
		}
		if string(buf) == "exit" {
			log.Fatalln("bye ~")
		}

		fmt.Println("server receive:",string(buf))
		err = go_jeans.WriteN(*conn,[]byte("pong pong pong ~"),go_jeans.PacketHerderLen_16)
		if err != nil {
			log.Fatalln(err)
		}
	}

}