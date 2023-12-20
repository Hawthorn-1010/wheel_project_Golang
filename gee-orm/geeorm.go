package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
	dbName  string
}

func NewEngine(connectionString string, driver string) (e *Engine, err error) {
	// open database connection
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// try to connect database
	err = db.Ping()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		log.Error("Get dialect error!")
		return nil, err
	}
	res := strings.Split(connectionString, "/")
	dbName := res[len(res)-1]
	log.Info("connection setup!")
	return &Engine{db: db, dialect: dialect, dbName: dbName}, nil
}

func (e *Engine) Close() {
	err := e.db.Close()
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("database close successfully!")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect, e.dbName)
}

type TxFunc func(*session.Session) (interface{}, error)

func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}

// difference returns a - b
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.SetTable(value).HasTable() {
			log.Infof("table %s doesn't exist", s.Table.TableName)
			return nil, s.CreateTable()
		}
		table := s.Table
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.TableName)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)

		// add new columns
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.TableName, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		//tmp := "tmp_" + table.TableName
		//fieldStr := strings.Join(table.FieldNames, ", ")
		//s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, fieldStr, table.TableName))
		//s.Raw(fmt.Sprintf("DROP TABLE %s;", table.TableName))
		//s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.TableName))
		//_, err = s.Exec()

		// for mysql to drop columns
		for _, col := range delCols {
			sqlStr := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", table.TableName, col)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}
		return
	})
	return err
}
