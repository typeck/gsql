package test

import (
	"github.com/typeck/gsql"
	"testing"
)

func TestQuery(t *testing.T) {
	sql := gsql.SqlInfo{}
	s,values := sql.Table("user").Query("name","company").
		Where("status=? AND company=?",1,"goo").
		And("age>?",10).Raw("email like ?","%com%").Done()
	t.Log(s,values)
}

func TestInsert(t *testing.T) {
	sql := gsql.SqlInfo{}
	s, values := sql.Table("user").Insert("name","company").
		Values("jack","goo").Done()
	t.Log(s, values)
}

func TestUpdate(t *testing.T) {
	sql := gsql.SqlInfo{}
	s, values := sql.Table("user").Update("name", "company").
		Values("jack", "goo").Done()
	t.Log(s, values)
}