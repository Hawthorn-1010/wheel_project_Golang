package main

import (
	"fmt"
	"geeorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, _ := geeorm.NewEngine("username:password@tcp(127.0.0.1:3306)/dbname", "mysql")
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

	// 检查是否有错误发生
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

func TestCreateTable(t *testing.T) {
	engine, _ := geeorm.NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	defer engine.Close()
	session := engine.NewSession()
	session.SetTable(&User{})
	session.CreateTable()
	if session.HasTable() {
		t.Log("create table success!")
	}
	session.DropTable()
	if !session.HasTable() {
		t.Log("drop table success!")
	}
}
