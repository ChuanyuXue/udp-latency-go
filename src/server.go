package src

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	IPLocal    string
	PortLocal  int
	IPRemote   string
	PortRemote int
	DevName    string

	Conn    *net.UDPConn
	VlanTag bool

	ClobalTime uint64

	TimeChan chan []byte
}

func (server *Server) Init(ipLocal string, portLocal int, ipRemote string, portRemote int, devName string) {
	server.IPLocal = ipLocal
	server.PortLocal = portLocal
	server.IPRemote = ipRemote
	server.PortRemote = portRemote
	server.TimeChan = make(chan []byte, BUFFER_SIZE)
}

func (server *Server) Listen(packetSize uint16) error {
	var index uint32
	var currentTime uint64
	var msg []byte

	if packetSize < MIN_PKT_SIZE || packetSize > MAX_PKT_SIZE {
		fmt.Printf("[!] Server Argument Error: Packet size %d is out of range.", packetSize)
		return errors.New("argument error")
	}

	addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(server.PortLocal))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[!] Server UDP Error: Failed to dial remote device.")
		return err
	}
	defer conn.Close()

	for {
		msg = make([]byte, packetSize)
		_, _, err := conn.ReadFromUDP(msg)
		currentTime = getTime(server.DevName)
		if err != nil {
			fmt.Println("[!] Server UDP Error: Unable to read incoming message")
			return err
		}

		index = binary.LittleEndian.Uint32(msg[:4])
		binary.LittleEndian.PutUint64(msg[12:20], currentTime) // T2
		server.TimeChan <- msg

		if index == 0 {
			return nil
		}
	}

}

func (server *Server) Send() error {
	var index uint32
	var currentTime uint64
	var msg []byte

	conn, err := net.Dial("udp4", server.IPRemote+":"+strconv.Itoa(server.PortRemote))
	if err != nil {
		fmt.Println("[!] Server UDP Error: Failed to dial remote device.")
		return errors.New("UDP error")
	}
	defer conn.Close()

	for {
		msg = <-server.TimeChan
		currentTime = getTime(server.DevName)
		binary.LittleEndian.PutUint64(msg[20:28], currentTime) // T3
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("[!] Client UDP Error: Unable to send message")
			return err
		}

		index = binary.LittleEndian.Uint32(msg[:4])

		if index == 0 {
			return nil
		}

	}
}
