## ORM

A ORM framework base that supports MySQL dialect. 

### How to use

* New engine

```golang
engine, _ := NewEngine("username:password@tcp(127.0.0.1:3306)/dbname", "mysql")
defer engine.Close()
```

* New session

```go
s := engine.NewSession()
```

* Table

```go
s := NewSession().SetTable(&User{})
// Drop table
s.DropTable() 	
// Create table
s.CreateTable()		
// Determine whether a table exists
if s.HasTable() {
}
```

* Record

```go
s := NewSession().SetTable(&User{})
// Insert record
s.Insert(user1, user2, user3)	
//Delete record
s.Where("ID = 1").Delete()
// Update record
s.Where("Name = 'Tom'").Update("Age", 32)
// Find record
var user []User
s.Find(&user)	
// Count record
rowNum, _ := s.Count()
// Get first record
s.OrderBy("ID DESC")
var user User
s.First(&user)
```

* Support hook, define as a function of the struct

```go
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
```

* Support transaction

```go
_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
   
})
```

* Support migration when altering table structure

```go
engine.Migrate(&User{})
```

### Reference

https://github.com/geektutu/7days-golang

https://github.com/go-xorm/xorm/


