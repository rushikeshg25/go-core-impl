package main

import (
	"durablelogs/durablelogs"
	"fmt"
)

func main() {
	dl := durablelogs.NewDLServer("./logs", 10)
	for i := 0; i < 1; i++ {
		dl.Log(fmt.Sprintf("Hello World %d", i))
	}
}
