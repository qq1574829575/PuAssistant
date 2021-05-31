package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	fmt.Println(GetCurrentTimeStr())
}

func GetCurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCurrentTimeUnixStr() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
