package main

import (
	"github.com/ChuanyuXue/udp-latency-go/src"
)

func main() {
	client := src.Client{}
	client.Init("localhost", 1234, "localhost", 4321, "sw")
	go client.Listen(1024, " ")
	client.Send(1, 100, 10)

}
