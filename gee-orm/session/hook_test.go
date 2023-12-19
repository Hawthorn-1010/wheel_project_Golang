package session

import (
	"log"
	"testing"
)

type Account struct {
	ID       int
	Password string
}

func (u *Account) BeforeQuery(s *Session) error {
	log.Printf("________Before Query________")
	return nil
}

func (u *Account) AfterQuery(s *Session) error {
	log.Printf("________After Query________")
	u.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession().SetTable(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "123456"}, &Account{2, "qwerty"})

	u := &Account{}

	err := s.First(u)
	if err != nil || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}
