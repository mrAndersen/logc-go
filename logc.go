package main

import (
	"fmt"
	_ "github.com/kshvakov/clickhouse"
	"log"
	"net"
)


func process(writer *Writer, parser *Parser, buf []byte) {
	message, ok := parser.parse(buf)

	if ok == true {
		writer.add(message)
	}
}

func main() {
	port := 9222
	addr := "0.0.0.0"
	bufSize := 10240

	writer := Writer{}
	writer.connect()

	parser := Parser{}
	parser.init()

	pc, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", addr, port))
	log.Printf("Started UDP proxy on %s:%d", addr, port)
	HandleError(err)

	defer pc.Close()

	for {
		buf := make([]byte, bufSize)

		n, _, err := pc.ReadFrom(buf)

		if err != nil {
			continue
		}

		go writer.watch()
		go process(&writer, &parser, buf[:n])
	}
}
