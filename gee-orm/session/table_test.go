package session

import (
	"testing"
)

type User struct {
	Id   int `geeorm:"PRIMARY KEY"`
	Name string
	Age  int
}

func TestCreateTable(t *testing.T) {
	//engine, _ := geeorm.NewEngine("root:root@tcp(192.168.255.3:3306)/books", "mysql")
	//defer engine.Close()
	//session := engine.NewSession()
	//session.SetTable(&User{})
	session := NewSession().SetTable(&User{})
	session.DropTable()
	session.CreateTable()
	if session.HasTable() {
		t.Log("create Table success!")
	}
}
