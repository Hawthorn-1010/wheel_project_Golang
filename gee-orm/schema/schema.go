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
	Model      interface{}
	TableName  string
	FieldNames []string
	ColumnsMap map[string]*Column
	Columns    []*Column
}

func (table *Table) GetField(name string) *Column {
	return table.ColumnsMap[name]
}

func Parse(model interface{}, dialect dialect.Dialect) *Table {

	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	var columns []*Column
	var fieldNames []string
	columnsMap := make(map[string]*Column)

	for i := 0; i < modelType.NumField(); i++ {
		column := &Column{
			Name: modelType.Field(i).Name,
			Type: dialect.SQLType(modelType.Field(i).Type),
		}
		if v, ok := modelType.Field(i).Tag.Lookup("geeorm"); ok {
			column.Tag = v
		}
		columns = append(columns, column)
		fieldNames = append(fieldNames, modelType.Field(i).Name)
		columnsMap[modelType.Field(i).Name] = column
	}

	t := &Table{
		// todo
		//Model:     modelType,
		Model:      model,
		TableName:  modelType.Name(),
		Columns:    columns,
		FieldNames: fieldNames,
		ColumnsMap: columnsMap,
	}
	return t
}
