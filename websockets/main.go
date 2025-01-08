package main

import (
	"fmt"
	"log"
	"net"
)


func main(){
  fmt.Println("Websocket:")
  listner,err:=net.Listen("tcp",":8080");
  if err!=nil{
	log.Fatal("Tcp Server Error");
   }
  for{
    log.Println("Listening on Port 8080")
    conn,err:=listner.Accept()
    if err!=nil{
      log.Fatal("unable to connect with the client");
    }
    go processConn(conn);
  }
}

func processConn(conn net.Conn){
  buf:=make([]byte,1024);
  _,err:=conn.Read(buf);
  if err!=nil{
    log.Fatal("Error Processing the Stream")
  }
  log.Printf("%s",buf);
}