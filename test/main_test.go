package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/typeck/gsql"
	"time"
)

var db *gsql.DB

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Account    string    `json:"account" example:"type@test.com"`
	Company    string    `json:"company"`
	Pwd        string    `json:"pwd" example:"e10adc3949ba59abbe56e057f20f883e"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Roles      int       `json:"role"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func init() {
	var err error
	db, err = gsql.NewDb("mysql","type:tang@(127.0.0.1)/test")
	if err != nil {
		panic(err)
	}
}

