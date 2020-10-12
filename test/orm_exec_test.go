package test

import "testing"



func TestInsert(t *testing.T) {
	var u = User{
		Name:       "test",
		Account:    "test@test.com",
		Company:    "",
		Pwd:        "",
		Phone:      "",
		Email:      "",
		Roles:      0,
		Status:     0,
		Base: &Base{Id: 22},
	}
	id,err := db.Pre().Debug().Create(&u).LastInsertId()
	if err != nil {
		t.Errorf("insert error:%v",err)
		return
	}
	t.Logf("insert success, id:%v", id)
}

func TestUpdate(t *testing.T) {
	var u = User{
		Name:       "test1",
		Account:    "test1@test.com",
		Company:    "cc",
		Pwd:        "12345",
		Phone:      "",
		Email:      "",
		Roles:      0,
		Status:     0,
		Base: &Base{Id: 23},
	}
	affect,err := db.Pre().Debug().Where("id=?", 23).Update(&u).RowsAffected()
	if err != nil {
		t.Errorf("update failed:%v",err)
		return
	}
	t.Logf("update success; affect rows:%v",affect)
}