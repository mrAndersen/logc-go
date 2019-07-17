package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

func HandleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ip2Long(ip string) uint32 {
	var long uint32

	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	HandleError(err)

	return long
}
