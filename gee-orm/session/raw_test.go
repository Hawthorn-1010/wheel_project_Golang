package session

import (
	"database/sql"
	"geeorm/dialect"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"testing"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
	dbName      string
)

func TestMain(m *testing.M) {
	connectionString := "root:root@tcp(192.168.255.3:3306)/books"
	res := strings.Split(connectionString, "/")
	TestDB, _ = sql.Open("mysql", connectionString)
	dbName = res[len(res)-1]
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial, dbName)
}
