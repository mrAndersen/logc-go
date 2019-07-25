package main

import (
	"fmt"
	"strings"
)

type DatabaseColumn struct {
	Title string
	Type  string
}

type Table struct {
	title   string
	engine  string
	Columns []DatabaseColumn
}

func (s *Table) getCreateTableSql() string {
	var columns []string

	for _, column := range s.Columns {
		columns = append(columns, fmt.Sprintf("%s %s", column.Title, column.Type))
	}

	ddl := strings.Join(columns, ",")
	sql := fmt.Sprintf(
		"create table if not exists default.%s (%s) engine = %s",
		s.title,
		ddl,
		s.engine,
	)

	return sql
}

func (s *Table) getColumnsForInsert() string {
	var titles []string

	for _, column := range s.Columns {
		titles = append(titles, column.Title)
	}

	return strings.Join(titles, ",")
}

func (s *Table) createNginxLayout() {
	s.Columns = append(s.Columns, DatabaseColumn{"ip", "UInt32"})
	s.Columns = append(s.Columns, DatabaseColumn{"time", "DateTime"})
	s.Columns = append(s.Columns, DatabaseColumn{"date", "Date default toDate(time)"})
	s.Columns = append(s.Columns, DatabaseColumn{"uri", "String"})
	s.Columns = append(s.Columns, DatabaseColumn{"method", "String"})
	s.Columns = append(s.Columns, DatabaseColumn{"protocol", "String"})
	s.Columns = append(s.Columns, DatabaseColumn{"status", "UInt16"})
	s.Columns = append(s.Columns, DatabaseColumn{"bytes", "UInt16"})
	s.Columns = append(s.Columns, DatabaseColumn{"referer", "String"})
	s.Columns = append(s.Columns, DatabaseColumn{"userAgent", "String"})
	s.Columns = append(s.Columns, DatabaseColumn{"hostname", "String"})
}
