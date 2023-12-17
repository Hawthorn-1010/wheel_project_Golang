package geeorm

import (
	"database/sql"
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
