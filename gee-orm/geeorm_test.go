package geeorm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

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
