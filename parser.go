package main

import (
	"regexp"
	"strconv"
	"time"
)

type Parser struct {
	regex *regexp.Regexp
}

func (s *Parser) init() {
	regex, err := regexp.Compile(`^<(\d+)>(.*)nginx:\s(.*?)\s(.*?)\s-\s(.*?)\[(.*?)\]\s\"(GET|POST|PUT|HEAD|PATCH|DELETE|UPDATE|OPTIONS|TRACE|PATCH)\s(.*?)\s(.*?)\"\s(\d+)\s(\d+)\s\"(.*?)\"\s\"(.*?)\"$`)
	HandleError(err)

	s.regex = regex
}

func (s *Parser) parse(bytes []byte) (LogMessage, bool) {
	logMessage := LogMessage{}
	match := s.regex.FindSubmatch(bytes)

	if len(match) != 14 {
		return logMessage, false
	}

	if string(match[1]) != "190" {
		return logMessage, false
	}

	logMessage.hostname = string(match[3])
	logMessage.ip = Ip2Long(string(match[4]))


	parsed, err := time.Parse("02/Jan/2006:15:04:05 -0700", string(match[6]))
	HandleError(err)

	loc, _ := time.LoadLocation("UTC")
	parsed = parsed.In(loc)

	//Mon Jan 2 15:04:05 MST 2006
	logMessage.time = parsed.Format("2006-01-02 15:04:05")
	logMessage.date = parsed.Format("2006-01-02")

	logMessage.method = string(match[7])
	logMessage.uri = string(match[8])
	logMessage.protocol = string(match[9])

	value, err := strconv.ParseInt(string(match[10]), 10, 64)
	HandleError(err)
	logMessage.status = value

	value, err = strconv.ParseInt(string(match[11]), 10, 64)
	HandleError(err)
	logMessage.bytes = value

	logMessage.referer = string(match[12])
	logMessage.userAgent = string(match[13])

	return logMessage, true
}
