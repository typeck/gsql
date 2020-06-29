package test

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/typeck/gsql"
)

var db gsql.Db
//var gormDb *gorm.DB
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
	Id         int    `json:"id" db:"id" gorm:"primary_key" json:"id"`
	Name       string `json:"name" db:"name" gorm:"name"`
	Account    string `json:"account" db:"account" gorm:"account"`
	Company    string `json:"company" db:"company" gorm:"company"`
	Pwd        string `json:"pwd" db:"pwd" gorm:"pwd"`
	Phone      string `json:"phone" db:"phone" gorm:"phone"`
	Email      string `json:"email" db:"email" gorm:"email"`
	Roles      int    `json:"role" db:"roles" gorm:"roles"`
	Status     int    `json:"status" db:"status" gorm:"status"`
	//CreateTime string `json:"create_time" db:"create_time" gorm:"create_time"`
	//UpdateTime string `json:"update_time" db:"update_time" gorm:"update_time"`
	Time      Time  `json:"time"`
}

type Time struct {
	CreateTime string `json:"create_time" db:"create_time" gorm:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time" gorm:"update_time"`
}

func init() {
	var err error
	db, err = gsql.OpenDb("mysql","type:tang@(127.0.0.1)/test")
	if err != nil {
		panic(err)
	}

	//gormDb, err = gorm.Open("mysql", "type:tang@(127.0.0.1)/test")
	//if err != nil {
	//	panic(err)
	//}
}

