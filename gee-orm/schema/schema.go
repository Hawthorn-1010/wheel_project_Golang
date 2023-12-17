package schema

import (
	"geeorm/dialect"
	"reflect"
)

type Column struct {
	Name string
	Type string
	Tag  string
}

// Table represents a database table
type Table struct {
	Model     interface{}
	TableName string
	//columnsSeq    []string
	//columnsMap    map[string][]*Column
	Columns []*Column
	//Indexes       map[string]*Index
	//PrimaryKeys   []string
	//AutoIncrement string
	//Created       map[string]bool
	//Updated       string
	//Deleted       string
	//Version       string
	//StoreEngine   string
	//Charset       string
	//Comment       string
	//Collation     string
}

func Parse(model interface{}, dialect dialect.Dialect) *Table {

	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	var columns []*Column

	for i := 0; i < modelType.NumField(); i++ {
		column := &Column{
			Name: modelType.Field(i).Name,
			Type: dialect.SQLType(modelType.Field(i).Type),
		}
		columns = append(columns, column)
	}

	t := &Table{
		Model:     modelType,
		TableName: modelType.Name(),
		Columns:   columns,
	}
	return t
}
