package geeorm

import (
	"database/sql"
	"geeorm/log"
	"geeorm/session"
)

type Engine struct {
	db *sql.DB
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
	log.Info("connection setup!")
	return &Engine{db: db}, nil
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
	return session.New(e.db)
}
