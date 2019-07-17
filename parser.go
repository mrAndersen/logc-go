package main

import (
	"encoding/binary"
	"regexp"
	"time"
)

type Parser struct {
	regex *regexp.Regexp
}

func (s *Parser) init() {
	regex, err := regexp.Compile(`<(\d+)>(.*)nginx:\s(.*?)\s\[(.*?)\]\s\"(GET|POST|PUT|HEAD|PATCH|DELETE|UPDATE|OPTIONS|TRACE|PATCH)\s(.*?)\s(.*?)\"\s(\d+)\s(\d+)\s\"(.*?)\"\s\"(.*?)\"$`)
	HandleError(err)

	s.regex = regex
}

func (s *Parser) parse(bytes []byte) (LogMessage, bool) {
	logMessage := LogMessage{}
	match := s.regex.FindSubmatch(bytes)

	if len(match) != 12 {
		return logMessage, false
	}

	if string(match[1]) != "190" {
		return logMessage, false
	}

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

	return logMessage, true
}