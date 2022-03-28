package src

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	HEADER_UNTAG_SIZE = 42
	HEADER_SIZE       = 46
	AUG_SIZE          = 4 + 8 + 8 + 8
	MAX_PKT_SIZE      = 1500
	MIN_PKT_SIZE      = AUG_SIZE + HEADER_SIZE
	BUFFER_SIZE       = 1024
)

func GetTime(devName string) uint64 {
	if devName != "sw" {
		path := fmt.Sprintf("/sys/class/net/%s/ieee8021ST/CurrentTime", devName)
		timeCur, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("[!] No current time.")
			curr := uint64(time.Now().UnixNano())
			return curr
		}
		timeStr := strings.Split(string(timeCur), ".")

		secondTime, _ := strconv.ParseUint(timeStr[0][:10], 10, 64)
		nanoTime, _ := strconv.ParseUint(timeStr[1][:9], 10, 64)
		return secondTime*1e9 + nanoTime
	} else {
		return uint64(time.Now().UnixNano())
	}
}


func ArrayToString(a []int, delim string) string {
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
