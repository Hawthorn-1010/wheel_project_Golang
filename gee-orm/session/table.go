package session

import (
	"fmt"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

// 1. initialize table in session
// 2. create drop

func (s *Session) SetTable(model interface{}) *Session {
	if s.table == nil || reflect.Indirect(reflect.ValueOf(s.table.Model)).Type() != reflect.Indirect(reflect.ValueOf(model)).Type() {
		s.table = schema.Parse(model, s.dialect)
	}
	return s
}

func (s *Session) GetTable() *schema.Table {
	if s.table == nil {
		log.Error("table not set!")
	}
	return s.table
}

func (s *Session) CreateTable() error {
	table := s.GetTable()
	var columns []string
	for _, column := range table.Columns {
		columns = append(columns, fmt.Sprintf("%s %s %s", column.Name, column.Type, column.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.TableName, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	table := s.GetTable()
	_, err := s.Raw(fmt.Sprintf("DROP TABLE %s;", table.TableName)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, args := s.dialect.TableExist(s.dbName, s.table.TableName)
	result := s.Raw(sql, args...).QueryRow()

	var tableName string
	err := result.Scan(&tableName)
	if err != nil {
		log.Error(err)
	}
	if tableName == s.table.TableName {
		return true
	}
	return false
}
