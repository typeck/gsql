package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/typeck/gsql"
	"time"
)

var db gsql.Db

type User struct {
	Name    string `json:"name" db:"name,omitempty"`
	Account string `json:"account" db:"account,omitempty"`
	Company string `json:"company" db:"company,omitempty"`
	Pwd     string `json:"pwd" db:"pwd,omitempty"`
	Phone   string `json:"phone" db:"phone,omitempty"`
	Email   string `json:"email" db:"email,omitempty"`
	Roles   int    `json:"role" db:"roles,omitempty"`
	Status  int    `json:"status" db:"status,omitempty"`
	Base    *Base   `json:"base" db:"base,omitempty"`
}

type Base struct {
	Id         int       `json:"id" db:"id,omitempty"`
	CreateTime time.Time `json:"create_time" db:"create_time,omitempty"`
	UpdateTime time.Time `json:"update_time" db:"update_time,omitempty"`
}

func init() {
	var err error
	db, err = gsql.OpenDb("mysql","type:tang@(127.0.0.1)/test?parseTime=true")
	if err != nil {
		panic(err)
	}
}

