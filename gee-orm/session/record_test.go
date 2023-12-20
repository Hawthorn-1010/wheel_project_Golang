package session

import (
	"log"
	"testing"
)

var (
	user1 = &User{1, "Tom", 27}
	user2 = &User{2, "Sam", 12}
	user3 = &User{3, "Jack", 37}
)

func TestInsertRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	s.Insert(user1, user2, user3)
}

func TestFindRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	var user []User
	if err := s.Find(&user); err != nil {
		log.Fatal("Find error!")
	}
	t.Logf("%#v", user)
}
func TestUpdateRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	affected, _ := s.Where("Name = 'Tom'").Update("Age", 32)
	u := &User{}
	_ = s.Where("Name = 'Tom'").First(u)

	if affected != 1 || u.Age != 32 {
		t.Fatal("failed to update")
	}
}
func TestDeleteRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	if rowNum, err := s.Where("ID = 1").Delete(); err != nil {
		log.Fatal("Delete error!")
	} else {
		t.Log(rowNum)
	}
}

func TestCountRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	if rowNum, err := s.Count(); err != nil {
		log.Fatal("Count error!")
	} else {
		t.Log(rowNum)
	}
}

func TestGetFirstRecord(t *testing.T) {
	s := NewSession().SetTable(&User{})
	s.OrderBy("ID DESC")
	var user User
	//user := &User{}
	if err := s.First(&user); err != nil {
		log.Fatal("Get First error!")
	}
	t.Logf("%#v", user)
}
