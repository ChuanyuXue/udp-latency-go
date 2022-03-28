package main

import (
	"flag"
	"fmt"

	"github.com/ChuanyuXue/udp-latency-go/src"
)

var argType string
var argSavePath string

func init() {
	flag.StringVar(&argType, "type", "c", "c -> client, s -> server")
	flag.StringVar(&argSavePath, "save", "test.csv", "The path you save the log file")
	flag.Parse()
}

func main() {
	if argType == "c" {
		client := src.Client{}
		client.Init("localhost", 12345, "localhost", 54321, "sw", false)
		go client.Listen(1024, argSavePath)
		client.Send(1, 100, 10)
	}

	if argType == "s" {
		fmt.Println("start server")
		server := src.Server{}
		server.Init("localhost", 54321, "localhost", 12345, "sw", false)
		go server.Listen(100)
		server.Send()
	}
}
