package session

import (
	"fmt"
	"geeorm/log"
	"geeorm/schema"
	"reflect"
	"strings"
)

// 1. initialize Table in session
// 2. create drop

func (s *Session) SetTable(model interface{}) *Session {
	if s.Table == nil || reflect.Indirect(reflect.ValueOf(s.Table.Model)).Type() != reflect.Indirect(reflect.ValueOf(model)).Type() {
		s.Table = schema.Parse(model, s.dialect)
	}
	return s
}

func (s *Session) GetTable() *schema.Table {
	if s.Table == nil {
		log.Error("Table not set!")
	}
	return s.Table
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
	sql, args := s.dialect.TableExist(s.dbName, s.Table.TableName)
	result := s.Raw(sql, args...).QueryRow()

	var tableName string
	err := result.Scan(&tableName)
	if err != nil {
		log.Error(err)
	}
	if tableName == s.Table.TableName {
		return true
	}
	return false
}
