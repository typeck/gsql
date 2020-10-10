package test

import (
	"github.com/typeck/gsql/util"
	"testing"
)

//func TestGet(t *testing.T) {
//	var u = &User{
//	}
//	err := db.Pre().Debug().Where("id=?",1).Get(u).Err()
//	if err != nil {
//		t.Errorf("get error:%v",err)
//		return
//	}
//	t.Logf("get success:%s",util.String(u))
//}
//
//func TestGetWithCols(t *testing.T) {
//	var u = &User{
//	}
//	err := db.Pre().Debug().Cols("id", "name").Where("id=?",1).Get(u).Err()
//	if err != nil {
//		t.Errorf("get error:%v",err)
//		return
//	}
//	t.Logf("get success:%s",util.String(u))
//}

func TestGets(t *testing.T) {
	var us []*User
	err := db.Pre().Debug().Where("id>?", 5).Gets(&us).Err()
	if err != nil {
		t.Errorf("gets error:%v",err)
		return
	}
	t.Logf("get success:%s",util.String(us))
}