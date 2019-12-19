package main

import (
	"fmt"
	_ "github.com/kshvakov/clickhouse"
	"log"
	"net"
	"os"
	"strings"
)

func process(writer *Writer, parser *Parser, buf []byte) {
	message, ok := parser.parse(buf)

	if ok == true {
		writer.add(message)
	}
}

func GetExecParameters() map[string]string {
	parameters := map[string]string{
		"UDP_PORT":    "9222",
		"UDP_ADDR":    "0.0.0.0",
		"CH_PORT":     "9000",
		"CH_USER":     "default",
		"CH_PASSWORD": "",
		"CH_ADDR":     "localhost",
	}

	for k, _ := range parameters {
		if os.Getenv(k) != "" {
			parameters[k] = os.Getenv(k)
		}

	}

	return parameters
}

func main() {
	p := GetExecParameters()

	debug := false
	bufSize := 10240

	writer := Writer{}
	writer.chPort = p["CH_PORT"]
	writer.chAddr = p["CH_ADDR"]
	writer.chUser = p["CH_USER"]
	writer.chPassword = p["CH_PASSWORD"]
	writer.connect()

	parser := Parser{}
	parser.init()

	connection, err := net.ListenPacket("udp", fmt.Sprintf("%s:%s", p["UDP_ADDR"], p["UDP_PORT"]))
	log.Printf("Started UDP proxy on %s:%s", p["UDP_ADDR"], p["UDP_PORT"])

	HandleError(err)

	defer connection.Close()
	go writer.watch()

	for {
		buf := make([]byte, bufSize)
		n, _, err := connection.ReadFrom(buf)

		if err != nil {
			continue
		}

		clean := buf[:n]

		if debug {
			message := strings.TrimSpace(string(clean))
			log.Printf("-> %s\n", message)
		}

		go process(&writer, &parser, clean)
	}
}
