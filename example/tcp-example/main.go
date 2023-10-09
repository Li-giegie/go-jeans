package main

func main(){
	go server("127.0.0.1:9999")
	newClient("127.0.0.1:9999")
}
