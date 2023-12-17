package main

import (
	"fmt"
	"reflect"
)

type Account struct {
	Name     string
	Password string
}

func main() {
	account := Account{
		Name:     "Sally",
		Password: "12345",
	}
	value := reflect.ValueOf(&account)
	//fmt.Println(value.Field(0))
	//fmt.Println(value.Field(1))
	//fmt.Println(value.Type().Name())
	fmt.Println(reflect.Indirect(value).Type().Field(1))

	typ := reflect.Indirect(reflect.ValueOf(&account)).Type()
	fmt.Println(typ.Name()) // Account

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Println(field.Name) // Username Password
	}

}
