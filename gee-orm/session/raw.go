package session

import (
	"database/sql"
	"geeorm/log"
	"strings"
)

type Session struct {
	db     *sql.DB
	sql    strings.Builder
	sqlVal []interface{}
}

func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// Encapsulation
func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Reset() {
	s.sql.Reset()
	s.sqlVal = nil
}

// 1. row -> 2. exec/query
func (s *Session) Raw(sql string, val ...interface{}) *Session {
	s.sql.WriteString(sql)
	//s.sql.WriteString(" ")
	s.sqlVal = append(s.sqlVal, val...)
	return s
}

/*
1. reset after execution finish
2. log sql info
3. log error if error happen
*/
func (s *Session) Exec() {
	defer s.Reset()
	_, err := s.db.Exec(s.sql.String(), s.sqlVal...)
	log.Info(s.sql.String(), s.sqlVal)
	if err != nil {
		log.Error(err)
	}
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Reset()
	row := s.db.QueryRow(s.sql.String(), s.sqlVal...)
	log.Info(s.sql.String(), s.sqlVal)
	return row
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Reset()
	row, err := s.db.Query(s.sql.String(), s.sqlVal...)
	log.Info(s.sql.String(), s.sqlVal)
	if err != nil {
		log.Error(err)
	}
	return row, err
}
