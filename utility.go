package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"
	"runtime"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Ip2Long(ip string) uint32 {
	var long uint32
	_ = binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

func GetMemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return math.Round(float64(m.HeapAlloc) / 1024 / 1024)
}
