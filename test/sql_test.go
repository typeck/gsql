package test

import (
	"github.com/typeck/gsql/util"
	"testing"
)

func TestQuery(t *testing.T) {
	var u = &User{}
	err := db.Pre().Debug().Table("user").Where("id=1").
		Cols("id", "name", "status", "company", "account", "company", "create_time, update_time ,email, phone, pwd").
		Query(&u.Base.Id, &u.Name, &u.Status, &u.Company, &u.Account, &u.Company, &u.Base.CreateTime, &u.Base.UpdateTime, &u.Email, &u.Phone, &u.Pwd).Err()
	if err != nil {
		t.Errorf("query error:%v",err)
		return
	}
	t.Logf("query success,user:%s",util.String(u))
}

func TestExecInsert(t *testing.T) {
	var u = &User{
		Name: "tee",
		Pwd: "13344",
	}
	id,err := db.Pre().Table("user").Debug().Cols("name, pwd").Exec("insert",u.Name, u.Pwd).LastInsertId()
	if err != nil {
		t.Errorf("insert error:%v",err)
		return
	}
	t.Logf("insert id:%v",id)
}

func TestExecUpdate(t *testing.T) {
	var u = &User{
		Base: Base{Id: 3},
	}
	affectRows, err := db.Pre().Table("user").Debug().Where("name=?",&u.Name).And("pwd=?",&u.Pwd).
		Cols("phone, company").Exec("update", "133","afg").RowsAffected()
	if err != nil {
		t.Errorf("update error:%v",err)
		return
	}
	t.Logf("affect rows:%v",affectRows)
}

func TestRaw(t *testing.T) {
	var u = &User{}
	err := db.Pre().Debug().Raw("select id,name from user where id=?",1).And("name=?","type").Query(&u.Base.Id, &u.Name).Err()
	if err != nil {
		t.Errorf("query error:%v", err)
		return
	}
	t.Logf("query success,user:%s",util.String(u))
}