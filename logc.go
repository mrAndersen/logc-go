package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"time"
)

type DatabaseColumn struct {
	Title string
	Type  string
}

//ip:
//time:
//date:
//uri: String
//method: String
//protocol: String
//status: UInt16
//bytes: UInt16
//referer: String
//userAgent: String

func HandleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type Table struct {
	Columns []DatabaseColumn
}

type LogMessage struct {
	ip        uint32
	time      string
	date      string
	uri       string
	method    string
	protocol  string
	status    uint16
	bytes     uint16
	referer   string
	userAgent string
}

func configureClickHouseTable() Table {
	table := Table{}

	table.Columns = append(table.Columns, DatabaseColumn{"ip", "UInt32"})
	table.Columns = append(table.Columns, DatabaseColumn{"time", "DateTime"})
	table.Columns = append(table.Columns, DatabaseColumn{"date", "Date default toDate(time)"})
	table.Columns = append(table.Columns, DatabaseColumn{"uri", "String"})
	table.Columns = append(table.Columns, DatabaseColumn{"method", "String"})
	table.Columns = append(table.Columns, DatabaseColumn{"protocol", "String"})
	table.Columns = append(table.Columns, DatabaseColumn{"status", "UInt16"})
	table.Columns = append(table.Columns, DatabaseColumn{"bytes", "UInt16"})
	table.Columns = append(table.Columns, DatabaseColumn{"referer", "String"})
	table.Columns = append(table.Columns, DatabaseColumn{"userAgent", "String"})

	return table
}

func ip2Long(ip string) uint32 {
	var long uint32

	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	HandleError(err)

	return long
}

func main() {
	port := 9222
	address := "0.0.0.0"
	bufSize := 10240

	pc, err := net.ListenPacket("udp", fmt.Sprintf("%s:%d", address, port))
	log.Printf("Started logc-go on %s:%d", address, port)
	HandleError(err)

	nginxRegex, err := regexp.Compile(`<(\d+)>(.*)nginx:\s(.*?)\s\[(.*?)\]\s\"(GET|POST|PUT|HEAD|PATCH|DELETE|UPDATE|OPTIONS|TRACE|PATCH)\s(.*?)\s(.*?)\"\s(\d+)\s(\d+)\s\"(.*?)\"\s\"(.*?)\"$`)
	HandleError(err)

	defer pc.Close()

	for {
		buf := make([]byte, bufSize)

		n, _, err := pc.ReadFrom(buf)

		if err != nil {
			continue
		}

		go serve(nginxRegex, buf[:n])
	}

}

func serve(r *regexp.Regexp, buf []byte) {
	logMessage := LogMessage{}
	match := r.FindSubmatch(buf)

	logMessage.ip = ip2Long(string(match[3]))

	parsed, err := time.Parse("02/Jan/2006:15:04:05 -0700", string(match[4]))
	HandleError(err)

	//Mon Jan 2 15:04:05 MST 2006
	logMessage.time = parsed.Format("2006-01-02 15:04:05")
	logMessage.date = parsed.Format("2006-01-02")

	logMessage.method = string(match[5])
	logMessage.uri = string(match[6])
	logMessage.protocol = string(match[7])
	logMessage.status = binary.BigEndian.Uint16(match[8])
	logMessage.bytes = binary.BigEndian.Uint16(match[9])
	logMessage.referer = string(match[10])
	logMessage.userAgent = string(match[11])

	log.Printf("%s", string(buf))
	os.Exit(1)
}
