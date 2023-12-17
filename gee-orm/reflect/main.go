package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// MySQL 数据库连接字符串
	connectionString := ""

	// 打开数据库连接
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 尝试连接数据库
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MySQL database!")

	// 在这里可以执行数据库操作
	// ...

	// 示例：查询数据库中的数据
	rows, err := db.Query("SELECT * FROM book")
	if err != nil {
		log.Fatal(err)
	}
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
