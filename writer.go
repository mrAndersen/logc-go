package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
	"time"
)

type Writer struct {
	connection *sqlx.DB
	period     time.Duration

	buffer []LogMessage
	table  Table
}

func (s *Writer) watch() {
	for true {
		time.Sleep(s.period)
		s.write()
	}
}

func (s *Writer) add(message LogMessage) {
	s.buffer = append(s.buffer, message)
}

func (s *Writer) write() {
	if len(s.buffer) == 0 {
		return
	}

	var messages []string

	for _, v := range s.buffer {
		messages = append(messages, v.toString())
	}

	transaction, err := s.connection.Begin()
	HandleError(err)

	sql := fmt.Sprintf("insert into %s (%s) values %s", s.table.title, s.table.getColumnsForInsert(), strings.Join(messages, ","))
	s.buffer = nil

	result, err := transaction.Exec(sql)
	HandleError(err)

	err = transaction.Commit()
	HandleError(err)

	_ = result
}

func (s *Writer) tryCreateTable() {
	result, err := s.connection.Exec(s.table.getCreateTableSql())
	HandleError(err)

	_ = result
}

func (s *Writer) connect() {
	s.period = time.Second * 1

	clickhousePort := 10000
	clickhouseAddr := "localhost"

	db, err := sqlx.Open("clickhouse", fmt.Sprintf("tcp://%s:%d", clickhouseAddr, clickhousePort))

	HandleError(err)
	HandleError(db.Ping())

	s.table.title = "nginx"
	s.table.engine = "MergeTree(date, (status, time, uri, method), 8192)"
	s.table.createNginxLayout()
	s.tryCreateTable()

	s.connection = db
	log.Printf("Connected to clickhouse on %s:%d/%s", clickhouseAddr, clickhousePort, s.table.title)
}
