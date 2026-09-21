package main

import (
	"durablelogs/durablelogs"
	"log"
)

func main() {
	d, e := durablelogs.Open("./logs", 5)
	if e != nil {
		log.Fatal(e)
	}
	for i := 0; i < 10; i++ {
		if e = d.Log("hello"); e != nil {
			d.Close()
			log.Fatal(e)
		}
	}
	if e = d.Close(); e != nil {
		log.Fatal(e)
	}
}
