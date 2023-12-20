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
	dbName  string
	db      *sql.DB
	tx      *sql.Tx
	clause  clause.Clause
	Table   *schema.Table
	sql     strings.Builder
	sqlVal  []interface{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect, dbName string) *Session {
	return &Session{db: db, dialect: dialect, dbName: dbName}
}

// Encapsulation
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
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
