package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Writer struct {
	chPort string
	chAddr string

	connection  *sqlx.DB
	flushPeriod time.Duration

	buffer []LogMessage
	table  Table
}

func (s *Writer) watch() {
	for true {
		time.Sleep(s.flushPeriod)
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

	start := time.Now()

	tx, err := s.connection.Begin()
	sql := fmt.Sprintf("insert into %s (%s) values %s", s.table.title, s.table.getColumnsForInsert(), getLogMessagePrep())

	stmt, err := tx.Prepare(sql)
	HandleError(err)

	for _, v := range s.buffer {
		values := v.getSlice()
		_, err := stmt.Exec(values...)

		HandleError(err)
	}

	HandleError(tx.Commit())

	end := time.Now().Sub(start).Seconds() * 1000
	log.Printf("Inserted %d entries in %.2f ms, mem = %.2fMb", len(s.buffer), end, GetMemoryUsage())

	s.buffer = nil
}

func (s *Writer) tryCreateTable() {
	result, err := s.connection.Exec(s.table.getCreateTableSql())
	HandleError(err)

	_ = result
}

func (s *Writer) connect() {
	s.flushPeriod = time.Second * 10

	db, err := sqlx.Open("clickhouse", fmt.Sprintf("tcp://%s:%s", s.chAddr, s.chPort))
	s.connection = db

	HandleError(err)
	HandleError(db.Ping())

	s.table.title = "nginx"
	s.table.engine = "MergeTree(date, (status, time, uri, method, hostname), 8192)"
	s.table.createNginxLayout()
	s.tryCreateTable()

	log.Printf("Connected to clickhouse on %s:%s/%s", s.chAddr, s.chPort, s.table.title)
}
