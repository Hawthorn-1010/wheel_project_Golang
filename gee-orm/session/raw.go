package session

import (
	"database/sql"
	"geeorm/clause"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/schema"
	"strings"
)

type Session struct {
	dialect dialect.Dialect
	db      *sql.DB
	dbName  string
	clause  clause.Clause
	table   *schema.Table
	sql     strings.Builder
	sqlVal  []interface{}
}

func New(db *sql.DB, dialect dialect.Dialect, dbName string) *Session {
	return &Session{db: db, dialect: dialect, dbName: dbName}
}

// Encapsulation
func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Reset() {
	s.sql.Reset()
	s.sqlVal = nil
	s.clause = clause.Clause{}
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
func (s *Session) Exec() (sql.Result, error) {
	defer s.Reset()
	result, err := s.db.Exec(s.sql.String(), s.sqlVal...)
	log.Info(s.sql.String())
	log.Info(s.sqlVal)
	if err != nil {
		log.Error(err)
	}
	return result, err
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
