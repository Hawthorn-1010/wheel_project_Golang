package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

type Dialect interface {
	SQLType(t reflect.Type) (s string)
	TableExist(tableSchema string, tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialectMap[name]
	return dialect, ok
}
