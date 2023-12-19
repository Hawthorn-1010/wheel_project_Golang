package clause

import (
	"geeorm/log"
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql    map[Type]string
	sqlVal map[Type][]interface{}
}

func (c *Clause) Set(name Type, val ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVal = make(map[Type][]interface{})
	}
	t := generators
	log.Info(t)
	sql, sqlVal := generators[name](val...)
	c.sql[name] = sql
	c.sqlVal[name] = sqlVal
}

func (c *Clause) Build(names ...Type) (string, []interface{}) {
	var sqls []string
	var sqlsVals []interface{}
	for _, name := range names {
		if sql, ok := c.sql[name]; ok {
			sqls = append(sqls, sql)
			// TODO
			sqlsVals = append(sqlsVals, c.sqlVal[name]...)
		}
	}
	return strings.Join(sqls, " "), sqlsVals
}
