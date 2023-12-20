package geeorm

import (
	"errors"
	"fmt"
	"geeorm/session"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"testing"
)

type Staff struct {
	ID   int `geeorm:"PRIMARY KEY"`
	Name string
	Age  int
}

type User struct {
	ID   int
	Name string
}

func TestEngine(t *testing.T) {
	engine, _ := NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	defer engine.Close()
	session := engine.NewSession()
	session.Raw("DROP TABLE IF EXISTS User;").Exec()
	session.Raw("CREATE TABLE User(Name text);").Exec()
	session.Raw("CREATE TABLE User(Name text);").Exec()
	session.Raw("SELECT * FROM book")
	rows, _ := session.QueryRows()
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.SetTable(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.SetTable(&User{}).CreateTable()
		_, err = s.Insert(&User{3, "Tom"})
		return nil, errors.New("Error")
	})
	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.SetTable(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.SetTable(&User{}).CreateTable()
		_, err = s.Insert(&User{1, "Tom"})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS Staff;").Exec()
	_, _ = s.Raw("CREATE TABLE Staff(ID INT PRIMARY KEY, Name text, XXX INT);").Exec()
	_, _ = s.Raw("INSERT INTO Staff(`ID`, `Name`) values (?, ?), (?, ?)", 1, "Tom", 2, "Sam").Exec()
	engine.Migrate(&Staff{})

	rows, _ := s.Raw("SELECT * FROM Staff").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"ID", "Name", "Age"}) {
		t.Fatal("Failed to migrate table Staff, got columns", columns)
	}
}
