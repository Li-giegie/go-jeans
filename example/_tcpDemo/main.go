package main

func main(){
	go server()

	newClient("127.0.0.1:9999")
}
