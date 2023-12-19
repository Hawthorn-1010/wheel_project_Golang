package session

import (
	"errors"
	"geeorm/clause"
	"geeorm/log"
	"reflect"
)

// [Tom 18]  model:poniter
func (s *Session) GenInsertVals(model interface{}) []interface{} {
	var value []interface{}
	// dereference
	modelValue := reflect.Indirect(reflect.ValueOf(model))
	modelType := reflect.Indirect(modelValue).Type()
	for i := 0; i < modelType.NumField(); i++ {
		value = append(value, modelValue.Field(i).Interface())
	}
	return value
}

func (s *Session) GenInsertColumns() []string {
	var columns []string
	for _, column := range s.table.Columns {
		columns = append(columns, column.Name)
	}
	return columns
}

// [[Tome 18][]] Insert(user1, user2) default: insert the same type
func (s *Session) Insert(model ...interface{}) (int64, error) {
	//test := model[0]
	s.SetTable(model[0])
	s.clause.Set(clause.INSERT, s.table.TableName, s.GenInsertColumns())
	// [[][]]
	vals := []interface{}{}

	for i := 0; i < len(model); i++ {
		vals = append(vals, s.GenInsertVals(model[i]))
	}
	s.clause.Set(clause.VALUES, vals...)
	sql, sqlVal := s.clause.Build(clause.INSERT, clause.VALUES)
	// 一开始没加... 传入变成数组了
	res, err := s.Raw(sql, sqlVal...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Find(models interface{}) error {
	s.Hook(BeforeQuery, nil)
	modelSlice := reflect.Indirect(reflect.ValueOf(models))
	modelType := modelSlice.Type().Elem()
	table := s.SetTable(reflect.New(modelType).Elem().Interface()).table

	s.clause.Set(clause.SELECT, table.TableName, s.GenInsertColumns())

	sql, sqlVal := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

	rows, err := s.Raw(sql, sqlVal...).QueryRows()
	if err != nil {
		log.Error(err)
	}

	for rows.Next() {
		dest := reflect.New(modelType).Elem()
		log.Info(dest.Type())
		var values []interface{}

		for _, column := range table.Columns {
			values = append(values, dest.FieldByName(column.Name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.Hook(AfterQuery, dest.Addr().Interface())
		modelSlice.Set(reflect.Append(modelSlice, dest))
	}
	return rows.Close()
}

func (s *Session) Update(model ...interface{}) (int64, error) {
	//test := model[0]
	s.SetTable(model[0])
	s.clause.Set(clause.INSERT, s.table.TableName, s.GenInsertColumns())
	// [[][]]
	vals := []interface{}{}

	for i := 0; i < len(model); i++ {
		vals = append(vals, s.GenInsertVals(model[i]))
	}
	s.clause.Set(clause.VALUES, vals...)
	sql, sqlVal := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(sql, sqlVal...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Delete(model ...interface{}) (int64, error) {
	//test := model[0]
	s.SetTable(model[0])
	s.clause.Set(clause.DELETE, s.table.TableName)
	sql, sqlVal := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(sql, sqlVal...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Count(model ...interface{}) (int64, error) {
	s.SetTable(model[0])
	s.clause.Set(clause.COUNT, s.table.TableName)

	sql, sqlVal := s.clause.Build(clause.COUNT)
	row := s.Raw(sql, sqlVal...).QueryRow()
	var rowNum int64
	if err := row.Scan(&rowNum); err != nil {
		return 0, err
	}
	return rowNum, nil
}

// return first record
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	//modelType := modelSlice.Type().Elem()
	//s.clause.Set(clause.SELECT, s.table.TableName, s.GenInsertColumns())
	//if err := s.Limit(1).Find(&destSlice); err != nil {
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("DATA NOT FOUND!")
	}

	// todo 区别于 dest = destSlice.Index(0)
	dest.Set(destSlice.Index(0))
	return nil
}

// method chaining
func (s *Session) Where(desc string) *Session {
	s.clause.Set(clause.WHERE, desc)
	return s
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}
