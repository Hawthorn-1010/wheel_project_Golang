package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct {
}

func init() {
	RegisterDialect("mysql", &mysql{})
}

func (db *mysql) SQLType(t reflect.Type) (s string) {
	switch k := t.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "Int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "Int Unsigned"
	case reflect.Int64:
		return "BigInt"
	case reflect.Uint64:
		return "BigInt Unsigned"
	case reflect.Float32:
		return "Float"
	case reflect.Float64:
		return "Double"
	case reflect.Array, reflect.Slice:
		return "Blob"
	case reflect.Bool:
		return "Bool"
	case reflect.String:
		return "Text"
	case reflect.Struct:
		if t.ConvertibleTo(reflect.TypeOf((*time.Time)(nil)).Elem()) {
			return "Datetime"
		}
		//default:
		//	return "Text"
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", t.Name(), t.Kind()))
	return
}

func (db *mysql) TableExist(tableSchema string, tableName string) (string, []interface{}) {
	args := []interface{}{tableSchema, tableName}
	sql := "SELECT table_name FROM information_schema.tables WHERE table_schema = ? AND table_name = ?"
	return sql, args
}
