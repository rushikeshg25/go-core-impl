package main

import (
	"log"
	"net"
	"time"
)

func processStream(conn net.Conn){
	buf:=make([]byte,1024)
	_,err:=conn.Read(buf)
	if err!=nil{
		log.Fatal("Error reading from TCP stream")
	}
	time.Sleep(5*time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello\r\n"))  
	conn.Close()
}

func main(){
	listner,err:=net.Listen("tcp",":3333")
	if err!=nil{
		log.Fatal("Error setting up TCP server")
	}
	for{
		log.Println("Listening on port 3333")
		conn,err:=listner.Accept()
		if err!=nil{
			log.Fatal("Error setting up TCP server")
		}
		go processStream(conn);

	}
}