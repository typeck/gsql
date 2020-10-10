package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/typeck/gsql"
	"time"
)

var db gsql.Db

type User struct {
	Name    string `json:"name" db:"name"`
	Account string `json:"account" db:"account"`
	Company string `json:"company" db:"company"`
	Pwd     string `json:"pwd" db:"pwd"`
	Phone   string `json:"phone" db:"phone"`
	Email   string `json:"email" db:"email"`
	Roles   int    `json:"role" db:"roles"`
	Status  int    `json:"status" db:"status"`
	Base    *Base   `json:"base" db:"base"`
}

type Base struct {
	Id         int       `json:"id" db:"id"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

func init() {
	var err error
	db, err = gsql.OpenDb("mysql","type:tang@(127.0.0.1)/test?parseTime=true")
	if err != nil {
		panic(err)
	}
}

