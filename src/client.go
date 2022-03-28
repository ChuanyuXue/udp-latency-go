package src

import (
	"encoding/binary"
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"os"
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

	VlanTag bool

	ClobalTime uint64

	Close chan bool
}

func (client *Client) Init(ipLocal string, portLocal int, ipRemote string, portRemote int, devName string, vlan bool) {
	client.IPLocal = ipLocal
	client.PortLocal = portLocal
	client.IPRemote = ipRemote
	client.PortRemote = portRemote
	client.DevName = devName
	client.VlanTag = vlan
	client.Close = make(chan bool, BUFFER_SIZE)

}

func (client *Client) Send(frequency uint16, packetSize uint16, duration uint16) error {
	var index uint32
	var payloadSize uint16
	var startTime uint64
	var currentTime uint64
	var totalPackets uint32
	var durationNano uint64
	var period float64
	conn, err := net.Dial("udp4", client.IPRemote+":"+strconv.Itoa(client.PortRemote))
	if err != nil {
		fmt.Println("[!] Client UDP Error: Failed to dial remote device.")
		return errors.New("UDP error")
	}
	defer conn.Close()

	if packetSize < MIN_PKT_SIZE || packetSize > MAX_PKT_SIZE {
		fmt.Printf("[!] Client Argument Error: Packet size %d is out of range.", packetSize)
		return errors.New("argument error")
	}

	if client.VlanTag {
		payloadSize = packetSize - HEADER_SIZE
	} else {
		payloadSize = packetSize - HEADER_UNTAG_SIZE
	}

	index = 1
	startTime = GetTime(client.DevName)
	totalPackets = uint32(frequency) * uint32(duration)
	durationNano = uint64(duration) * 1e9
	period = (1 / float64(frequency)) * 1e9

	for currentTime < startTime+durationNano || index < totalPackets {
		msg := make([]byte, payloadSize)
		binary.LittleEndian.PutUint32(msg[:4], index)
		currentTime = GetTime(client.DevName)
		binary.LittleEndian.PutUint64(msg[4:12], currentTime) // T1
		_, err = conn.Write(msg)
		if err != nil {
			fmt.Println("[!] Client UDP Error: Unable to send message")
			return err
		}
		index += 1

		currentTime += uint64(period)
		for GetTime(client.DevName) < currentTime {

		}
	}

	for {
		select {
		case <-client.Close:
			return nil
		default:
			msg := make([]byte, 4)
			binary.LittleEndian.PutUint32(msg[:4], 0)
			conn.Write(msg)
			// _, err := conn.Write(msg)
			// if err != nil {
			// 	fmt.Println("[!] Client UDP Error: Unable to send terminating message")
			// 	return err
			// }
		}
	}
}

func (client *Client) Listen(bufferSize uint16, savePath string) error {
	var index uint32
	var currentTime uint64
	var buf []byte

	addr, _ := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(client.PortLocal))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("[!] Client UDP Error: Failed to dial remote device.")
		return err
	}
	defer conn.Close()

	for {
		buf = make([]byte, bufferSize)
		length, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("[!] Client UDP Error: Unable to read incoming message")
			fmt.Println(err)
		}
		currentTime = GetTime(client.DevName)
		index = binary.LittleEndian.Uint32(buf[:4])
		go client.log(index,
			uint16(length),
			binary.LittleEndian.Uint64(buf[4:12]),
			binary.LittleEndian.Uint64(buf[12:20]),
			binary.LittleEndian.Uint64(buf[20:28]),
			currentTime)
		fmt.Println(index,
			uint16(length),
			binary.LittleEndian.Uint64(buf[4:12]),
			binary.LittleEndian.Uint64(buf[12:20]),
			binary.LittleEndian.Uint64(buf[20:28]),
			currentTime)
		if index == 0 {
			break
		}
	}

	client.save(savePath)
	client.Close <- true
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

func (client *Client) save(path string) error {
	csvfile, err := os.Create(path)
	if err != nil {
		fmt.Println("[!] CSV Error: Unable to create log file")
		return err
	}
	wr := csv.NewWriter(csvfile)
	defer wr.Flush()

	for i, _ := range client.LogIndex {
		var row []string
		row = append(row, strconv.FormatUint(uint64(client.LogIndex[i]), 10))
		row = append(row, strconv.FormatUint(uint64(client.LogMsgLen[i]), 10))
		row = append(row, strconv.FormatUint(uint64(client.LogTimestamp1[i]), 10))
		row = append(row, strconv.FormatUint(uint64(client.LogTimestamp2[i]), 10))
		row = append(row, strconv.FormatUint(uint64(client.LogTimestamp3[i]), 10))
		row = append(row, strconv.FormatUint(uint64(client.LogTimestamp4[i]), 10))
		err = wr.Write(row)
		if err != nil {
			fmt.Println("[!] CSV Error: Error in writing record to file")
			return err
		}
	}

	return nil
}
