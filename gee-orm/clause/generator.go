package clause

import (
	"fmt"
	"strings"
)

type generator func(val ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	// TODO
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderBy
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(val ...interface{}) (string, []interface{}) {
	tableName := val[0]
	// TODO
	vals := val[1].([]string)
	fields := strings.Join(vals, ", ")
	sql := fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields)
	return sql, []interface{}{}
}

func _values(vals ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	vars := []interface{}{}
	sql.WriteString("VALUES ")

	for i, val := range vals {
		v := val.([]interface{})
		sqlVars := genBindVars(len(v))
		sql.WriteString(fmt.Sprintf("(%v)", sqlVars))
		if i+1 != len(vals) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}

	return sql.String(), vars
}

func _select(val ...interface{}) (string, []interface{}) {
	tableName := val[0]
	fields := strings.Join(val[1].([]string), ", ")
	sql := fmt.Sprintf("SELECT %s FROM %s", fields, tableName)
	return sql, []interface{}{}
}

func _where(val ...interface{}) (string, []interface{}) {
	sql := fmt.Sprintf("WHERE %s", val[0].(string))
	return sql, []interface{}{}
}

func _limit(val ...interface{}) (string, []interface{}) {
	return "LIMIT ?", val
}

func _orderBy(val ...interface{}) (string, []interface{}) {
	sort := val[0]
	// order 先跟一个字段
	fields := val[1]
	sql := fmt.Sprintf("ORDER BY %s %s", fields, sort)
	return sql, []interface{}{}
}
