package src

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	IPLocal    string
	PortLocal  int
	IPRemote   string
	PortRemote int
	DevName    string

	LogIndex      []uint32
	LogMsgLen     []uint16
	LogTimestamp1 []uint64
	LogTimestamp2 []uint64
	LogTimestamp3 []uint64
	LogTimestamp4 []uint64

	Conn    *net.UDPConn
	VlanTag bool

	ClobalTime uint64

	Close chan bool
}

func (client *Client) send(frequency uint16, packetSize uint16, duration uint16, q chan bool) error {
	var index uint32
	var payloadSize uint16
	var startTime uint64
	var currentTime uint64
	var totalPackets uint32
	var durationNano uint64
	var period float64
	conn, err := net.Dial("udp4", client.IPRemote+":"+strconv.Itoa(client.PortRemote))
	defer conn.Close()
	if err != nil {
		fmt.Println("[!] UDP Error: Failed to dial remote device.")
		return errors.New("UDP error")
	}

	if packetSize < MIN_PKT_SIZE || packetSize > MAX_PKT_SIZE {
		fmt.Printf("[!] Argument Error: Packet size %d is out of range.", packetSize)
		return errors.New("argument error")
	}

	if client.VlanTag {
		payloadSize = packetSize - HEADER_SIZE
	} else {
		packetSize = packetSize - HEADER_UNTAG_SIZE
	}

	index = 1
	startTime = getTime(client.DevName)
	totalPackets = uint32(frequency) * uint32(duration)
	durationNano = uint64(duration) * 1e9
	period = 1 / float64(frequency)

	var msg []byte
	for currentTime < startTime+durationNano || index < totalPackets {
		msg = make([]byte, payloadSize)
		binary.LittleEndian.PutUint32(msg[:4], index)
		currentTime = getTime(client.DevName)
		binary.LittleEndian.PutUint64(msg[4:12], currentTime) // T1
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("[!] UDP Error: Unable to send message")
			return errors.New("UDP error")
		}
		index += 1

		currentTime += uint64(period)
		for getTime(client.DevName) < currentTime {

		}
	}

	for len(q) == 0 {
		msg = make([]byte, 4)
		binary.LittleEndian.PutUint32(msg[:4], 0)
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("[!] UDP Error: Unable to send terminating message")
			return errors.New("UDP error")
		}
	}

	return nil
}

func (client *Client) listen(bufferSize uint16, savePath string, q chan bool) error {
	var index uint32
	var currentTime uint64
	var buf []byte

	addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(client.PortLocal))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[!] UDP Error: Failed to dial remote device.")
		return errors.New("UDP error")
	}
	defer conn.Close()

	for {
		buf = make([]byte, bufferSize)
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		currentTime = getTime(client.DevName)
		index = binary.LittleEndian.Uint32(buf[:4])
		go client.log(index,
			uint16(length),
			binary.LittleEndian.Uint64(buf[4:12]),
			binary.LittleEndian.Uint64(buf[12:20]),
			binary.LittleEndian.Uint64(buf[20:28]),
			currentTime)
		if index == 0 {
			break
		}
	}

	q <- true
	return nil
}

func (client *Client) log(index uint32, length uint16, t1 uint64, t2 uint64, t3 uint64, t4 uint64) error {
	client.LogIndex = append(client.LogIndex, index)
	client.LogMsgLen = append(client.LogMsgLen, length)
	client.LogTimestamp1 = append(client.LogTimestamp1, t1)
	client.LogTimestamp2 = append(client.LogTimestamp2, t2)
	client.LogTimestamp3 = append(client.LogTimestamp3, t3)
	client.LogTimestamp4 = append(client.LogTimestamp4, t4)
	return nil
}

func (client *Client) evaluate() error {
	return nil
}
