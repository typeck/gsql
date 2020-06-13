package test

import "testing"

var u = User{
	Id:         0,
	Name:       "test",
	Account:    "test@test.com",
	Company:    "",
	Pwd:        "",
	Phone:      "",
	Email:      "",
	Roles:      0,
	Status:     0,
	CreateTime: "2020-06-12 21:49:48",
	UpdateTime: "2020-06-12 21:49:48",
}

func TestUpdate(t *testing.T) {

	affect,err := db.New().Debug().Where("id=?", 3).Cols("name, account").Update(&u).RowsAffected()
	if err != nil {
		t.Errorf("update failed:%v",err)
		return
	}
	t.Logf("update success; affect rows:%v",affect)
}

func TestInsert(t *testing.T) {
	id,err := db.New().Debug().Create(&u).LastInsertId()
	if err != nil {
		t.Errorf("insert error:%v",err)
		return
	}
	t.Logf("insert success, id:%v", id)
}
