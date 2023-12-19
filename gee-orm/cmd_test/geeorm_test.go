package main

import (
	"fmt"
	"geeorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, _ := geeorm.NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
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

func TestInterface(t *testing.T) {
	args := []interface{}{"tableSchema", "tableName"}
	t.Log(args)
}

type User struct {
	Id   int
	Name string
}

var (
	user1 = &User{1, "Tom"}
	user2 = &User{2, "Sam"}
	user3 = &User{3, "Jack"}
)

func TestCreateTable(t *testing.T) {
	engine, _ := geeorm.NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	defer engine.Close()
	session := engine.NewSession()
	session.SetTable(&User{})
	session.CreateTable()
	if session.HasTable() {
		t.Log("create table success!")
	}
	session.Insert(user1)
}

func TestInsertRecord(t *testing.T) {
	engine, _ := geeorm.NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	defer engine.Close()
	s := engine.NewSession()
	s.Insert(user2, user3)
}

func TestFindRecord(t *testing.T) {
	engine, _ := geeorm.NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	defer engine.Close()
	s := engine.NewSession()
	var user []User
	if err := s.Find(&user); err != nil {
		log.Fatal("Find error!")
	}
	t.Logf("%#v", user)
}
