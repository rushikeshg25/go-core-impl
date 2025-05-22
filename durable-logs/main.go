package main

import (
	"durablelogs/durablelogs"
)

func main() {
	dl := durablelogs.NewDLServer("./logs", 5)
	for i := 0; i < 10; i++ {
		dl.Log("hello")
	}
	dl.Flush()
}
