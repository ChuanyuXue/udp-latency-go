package main

import (
	"fmt"

	"github.com/ChuanyuXue/udp-latency-go/src"
)

func main() {
	client := src.Client{}
	client.Init("localhost", 1234, "localhost", 4321, "sw", false)
	go client.Listen(1024, " ")
	err := client.Send(1, 100, 10)
	fmt.Println(err)
}
