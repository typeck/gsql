package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/typeck/gsql"
)

var db *gsql.DB

//type User struct {
//	Id         int    `json:"id"`
//	Name       string `json:"name"`
//	Account    string `json:"account"`
//	Company    string `json:"company"`
//	Pwd        string `json:"pwd"`
//	Phone      string `json:"phone"`
//	Email      string `json:"email"`
//	Roles      int    `json:"role"`
//	Status     int    `json:"status"`
//	CreateTime string `json:"create_time"`
//	UpdateTime string `json:"update_time"`
//}

type User struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	Account    string `json:"account" db:"account" example:"type@test.com"`
	Company    string `json:"company" db:"company"`
	Pwd        string `json:"pwd" db:"pwd" example:"e10adc3949ba59abbe56e057f20f883e"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Roles      int    `json:"role" db:"roles"`
	Status     int    `json:"status" db:"status"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

func init() {
	var err error
	db, err = gsql.NewDb("mysql","type:tang@(127.0.0.1)/test")
	if err != nil {
		panic(err)
	}
}

