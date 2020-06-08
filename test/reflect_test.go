package test

import (
	"fmt"
	"github.com/typeck/gsql"
	"testing"
)

type User struct {
	Id 		int 	`db:"id"`
	Name    string	`db:"name"`
	Core    float64 `db:"core"`
}

type Company struct {
	Id 		int 	`db:"id"`
	Name    string	`db:"name"`
	Core    float64 `db:"core"`
}

func Scan( orm *gsql.Orm, v interface{},s *gsql.SqlInfo, cols []string)error {

	values, err := orm.BuildValues(s, v, cols)
	if err != nil {
		return err
	}

	switch fist := values[0].(type) {
	case *int:
		*fist = 10
		fmt.Println(fist)
	}
	switch fist := values[1].(type) {
	case *string:
		*fist = "xiao ming"
	}
	return nil
}

func TestScan(t *testing.T) {
	orm := gsql.NewOrm()
	var u1 = &User{}
	var u2 = &User{}
	var c1 string
	var clos = []string{"id", "name", "score"}

	err := Scan(orm, u1 ,&gsql.SqlInfo{}, clos)
	if err != nil {
		t.Fatal(err)
	}
	err = Scan(orm, u2 ,&gsql.SqlInfo{}, clos)
	if err != nil {
		t.Fatal(err)
	}
	err = Scan(orm, &c1 ,&gsql.SqlInfo{}, clos)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u1,u2)
}