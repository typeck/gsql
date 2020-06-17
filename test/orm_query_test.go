package test

import (
	"github.com/typeck/gsql/util"
	"testing"
)

func TestGet(t *testing.T) {
	var u = &User{
	}
	err := db.New().Debug().Where("id=?",1).Get(u).Err()
	if err != nil {
		t.Errorf("get error:%v",err)
		return
	}
	t.Logf("get success:%s",util.String(u))
}

//func TestGet2(t *testing.T) {
//	var u = &User{}
//	err := db.New().Debug().Table("user").Where("id=?",1).Cols("name, id").Get(u).Err()
//	if err != nil {
//		t.Errorf("get error:%v",err)
//		return
//	}
//	t.Logf("get success:%s",util.String(u))
//}
//
//func TestGets(t *testing.T) {
//	var us []*User
//	err := db.New().Debug().Gets(&us).Err()
//	if err != nil {
//		t.Errorf("gets error:%v",err)
//		return
//	}
//	t.Logf("get success:%s",util.String(us))
//}