package main

import (
	"fmt"
)

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

func getLogMessagePrep() string {
	return fmt.Sprintf("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
}

func (s *LogMessage) getSlice() []interface{} {
	result := make([]interface{}, 0)

	result = append(result, s.ip)
	result = append(result, s.time)
	result = append(result, s.date)
	result = append(result, s.uri)
	result = append(result, s.method)
	result = append(result, s.protocol)
	result = append(result, s.status)
	result = append(result, s.bytes)
	result = append(result, s.referer)
	result = append(result, s.userAgent)

	return result
}

func (s *LogMessage) toString() string {
	return fmt.Sprintf(
		"(%d, '%s', '%s', '%s', '%s', '%s', %d, %d, '%s', '%s')",
		s.ip,
		s.time,
		s.date,
		s.uri,
		s.method,
		s.protocol,
		s.status,
		s.bytes,
		s.referer,
		s.userAgent,
	)
}
