package main

import (
	"fmt"
	"net"
)

func main() {
	c, err := net.Dial("udp4", "192.168.10.13:19981")
	if err != nil {
		fmt.Println("[!] UDP Error", err)
	}
	bs := make([]byte, 0)
	_, err = c.Write(bs)
}
