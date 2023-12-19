package session

import (
	"log"
	"testing"
)

var (
	user1 = &User{1, "Tom"}
	user2 = &User{2, "Sam"}
	user3 = &User{3, "Jack"}
)

func TestInsertRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	s.Insert(user2, user1)
}

func TestFindRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	var user []User
	if err := s.Find(&user); err != nil {
		log.Fatal("Find error!")
	}
	t.Logf("%#v", user)
}

func TestDeleteRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	s.Where("ID = 2")
	if rowNum, err := s.Delete(&User{}); err != nil {
		log.Fatal("Delete error!")
	} else {
		t.Log(rowNum)
	}
}

func TestCountRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	if rowNum, err := s.Count(&User{}); err != nil {
		log.Fatal("Count error!")
	} else {
		t.Log(rowNum)
	}
}

func TestGetFirstRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	s.OrderBy("ID DESC")
	//var user User
	user := &User{}
	if err := s.First(user); err != nil {
		log.Fatal("Get First error!")
	}
	t.Logf("%#v", user)
}
