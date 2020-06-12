package test

import (
	"github.com/typeck/gsql/util"
	"testing"
)

func TestQuery(t *testing.T) {
	var u = &User{}
	var createTime string
	var updateTime string
	err := db.New().Debug().Table("user").Where("id=1").
		Cols("id", "name", "status", "company", "account", "company", "create_time, update_time ,email, phone, pwd").
		Query(&u.Id, &u.Name, &u.Status, &u.Company, &u.Account, &u.Company, &createTime, &updateTime, &u.Email, &u.Phone, &u.Pwd).Err()
	if err != nil {
		t.Errorf("query error:%v",err)
		return
	}
	t.Logf("query success,user:%s",util.String(u))
	t.Logf("time:%s, %s",createTime,updateTime)
}

func TestExec(t *testing.T) {
	var u = &User{
		Name: "tee",
		Pwd: "13344",
	}
	//id,err := db.New().Table("user").Debug().Cols("name, pwd").Exec("insert",u.Name, u.Pwd).LastInsertId()
	//if err != nil {
	//	t.Errorf("insert error:%v",err)
	//	return
	//}
	//t.Logf("insert id:%v",id)

	affectRows, err := db.New().Table("user").Debug().Where("name=?",u.Name).And("pwd=?",u.Pwd).
		Cols("phone, company").Exec("update", "133","afg").RowsAffected()
	if err != nil {
		t.Errorf("update error:%v",err)
		return
	}
	t.Logf("affect rows:%v",affectRows)
}
