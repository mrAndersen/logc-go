package main

import "fmt"

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

func (s *LogMessage) toString() string {
	return fmt.Sprintf(
		"(%d, '%s', '%s', '%s', '%s', '%s', '%d', '%d', '%s', '%s')",
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
