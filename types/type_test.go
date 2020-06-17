package types

import (
	"encoding/json"
	"reflect"
	"testing"
)

type User1 struct {
	Id 		int			`db:"id"`
	Name    string		`db:"name"`
	Company				`db:"company"`
}
type Company struct {
	Addr 	string		`db:"addr"`
	Name 	string		`db:"name"`
}

func TestStructUnwrap(t *testing.T) {
	u1 := &User1{}
	structInfo := NewStructInfo()
	typ := reflect.TypeOf(u1)

	structInfo.Unwrap(typ, "db")
	data ,err := json.Marshal(structInfo)
	if err != nil {
		t.Log(err)
	}
	t.Logf("struct info:%v\n",string(data))
}