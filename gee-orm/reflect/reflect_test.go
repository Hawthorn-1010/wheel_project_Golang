package reflect

import (
	"fmt"
	"reflect"
	"testing"
)

type Account struct {
	Name     string
	Password string
}

func TestReflect(t *testing.T) {
	account := Account{
		Name:     "Sally",
		Password: "12345",
	}
	value := reflect.ValueOf(&account)
	//fmt.Println(value.Field(0))
	//fmt.Println(value.Field(1))
	//fmt.Println(value.Type().Name())
	//fmt.Println(reflect.Indirect(value).Type().Field(1).Name)
	//for i := 0; i < reflect.Indirect(value).Type().NumField(); i++ {
	//	fmt.Println(reflect.Indirect(value).Type().Field(i))
	//}

	for i := 0; i < reflect.Indirect(value).Type().NumField(); i++ {
		fmt.Println(reflect.Indirect(value).Field(i).Interface())
	}
	//typ := reflect.Indirect(reflect.ValueOf(&account)).Type()
	//fmt.Println(typ.Name()) // Account
	//
	//for i := 0; i < typ.NumField(); i++ {
	//	field := typ.Field(i)
	//	fmt.Println(field.Name) // Username Password
	//}

	//	// MySQL 数据库连接字符串
	//	connectionString := ""
	//
	//	// 打开数据库连接
	//	db, err := sql.Open("mysql", connectionString)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer db.Close()
	//
	//	// 尝试连接数据库
	//	err = db.Ping()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	fmt.Println("Connected to MySQL database!")
	//
	//	// 在这里可以执行数据库操作
	//	// ...
	//
	//	// 示例：查询数据库中的数据
	//	rows, err := db.Query("SELECT * FROM book")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer rows.Close()
	//
	//	for rows.Next() {
	//		var id int
	//		var name string
	//		err := rows.Scan(&id, &name)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		fmt.Printf("ID: %d, Name: %s\n", id, name)
	//	}
	//
	//	// 检查是否有错误发生
	//	if err := rows.Err(); err != nil {
	//		log.Fatal(err)
	//
	//	}
}

func TestPanic(t *testing.T) {
	//注册捕获panic的函数,必须先注册，若在panic之后则无效
	defer doPanic()
	n := 0
	res := 1 / n
	fmt.Println(res) //panic 之后的代码不会执行
}

// 当捕获到panic时触发此函数
func doPanic() {
	err := recover()
	if err != nil {
		fmt.Println("捕获到panic")
		//panic(err)
	}
}

type MyStruct struct {
	Field1 int    `json:"field1" custom:"tag1"`
	Field2 string `json:"field2" custom:"tag2"`
}

func TestTagStruct(t *testing.T) {
	// Get the type of MyStruct
	myStructType := reflect.TypeOf(MyStruct{})

	// Get the first field of MyStruct
	field := myStructType.Field(0)

	// Get the struct tag for the "json" key
	tag := field.Tag
	jsonTag := tag.Get("json")
	fmt.Println("JSON Tag:", jsonTag)

	// Use Lookup to get the value for a custom key
	customTag, exists := tag.Lookup("custom")
	if exists {
		fmt.Println("Custom Tag:", customTag)
	} else {
		fmt.Println("Custom Tag not found")
	}
}
