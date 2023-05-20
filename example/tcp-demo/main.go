package main

func main(){
	go server()
	newClient("127.0.0.1:9999")

	//sendBuf := go_jeans.Pack([]byte("ping ~"))
	//
	//replyBuf,err := go_jeans.Unpack(conn)


}
