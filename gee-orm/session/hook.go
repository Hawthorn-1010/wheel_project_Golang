package session

import (
	"geeorm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) Hook(method string, model interface{}) {
	//f := reflect.Indirect(reflect.ValueOf(s.Table.Model)).MethodByName(method)
	// todo m is a ptr
	m := s.Table.Model
	f := reflect.ValueOf(m).MethodByName(method)
	if model != nil {
		f = reflect.ValueOf(model).MethodByName(method)
	}

	// TODO
	params := []reflect.Value{reflect.ValueOf(s)}
	if f.IsValid() {
		if v := f.Call(params); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
