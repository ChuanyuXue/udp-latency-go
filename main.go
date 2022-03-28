package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/ChuanyuXue/udp-latency-go/src"
)

var ipLocal string
var ipRemote string
var portLocal int
var portRemote int

var vlanTag bool

var argType string
var argClock string
var argSavePath string
var argFrequency uint64
var argPktSize uint64
var argDuration uint64
var argOffset string

func init() {

	flag.StringVar(&ipRemote, "ip", "localhost", "Remote ip")
	flag.IntVar(&portLocal, "lp", 12345, "Local port")
	flag.IntVar(&portRemote, "rp", 12345, "Remote port")
	flag.StringVar(&argType, "type", "c", "Agent type: c -> client, s -> server")
	flag.Uint64Var(&argFrequency, "f", 1, "The frequency in Hz")
	flag.Uint64Var(&argPktSize, "n", 100, "The packet size in byte")
	flag.Uint64Var(&argDuration, "t", 10, "The duration in seconds")
	flag.StringVar(&argOffset, "o", "0.0", "The start time of traffic")

	flag.BoolVar(&vlanTag, "vlan", false, "Vlan tag")
	flag.StringVar(&argClock, "clock", "sw", "Clock type: sw -> Linux kernel time, sw0p* -> PTP transparent clock")
	flag.StringVar(&argSavePath, "save", "test.csv", "The path you save the log file")

	flag.Parse()
}

func main() {
	if argType == "c" {
		client := src.Client{}
		client.Init("localhost", portLocal, ipRemote, portRemote, argClock, vlanTag)
		go client.Listen(argSavePath)
		sec, _ := strconv.ParseUint(strings.Split(argOffset, ".")[0], 10, 64)
		nanosec, _ := strconv.ParseUint(strings.Split(argOffset, ".")[1], 10, 64)
		client.Send(uint16(argFrequency), uint16(argPktSize), uint16(argDuration), sec*1e9+nanosec)
		client.Save(argSavePath)
	}

	if argType == "s" {
		fmt.Println("start server")
		server := src.Server{}
		server.Init("localhost", portLocal, ipRemote, portRemote, argClock, vlanTag)
		go server.Listen(uint16(argPktSize))
		server.Send()
	}
}
