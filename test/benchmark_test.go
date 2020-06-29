package test

//import (
//	"fmt"
//	"testing"
//)
//
//func gsqlGet(u *User){
//	err := db.New().Where("id=?",1).Get(u).Err()
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//func gormGet(u *User) {
//	err := gormDb.Table("user").Where("id=?",1).First(u).Error
//	if err != nil {
//		fmt.Println(err)
//	}
//}
//
//func rawSqlGet(u *User) {
//	rows,err := db.Query(`SElECT name,account,company,phone,email,update_time,id,pwd,roles,status,create_time FROM user WHERE id=? `,1)
//	if err != nil {
//		fmt.Println(err)
//	}
//	for rows.Next() {
//		err := rows.Scan(&u.Name, &u.Account, &u.Company, &u.Phone, &u.Email, &u.UpdateTime, &u.Id, &u.Pwd, &u.Roles, &u.Status, &u.CreateTime)
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}
//
//func BenchmarkGsql(b *testing.B) {
//	var u = &User{}
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		gsqlGet(u)
//	}
//}
//
//func BenchmarkGorm(b *testing.B) {
//	var u = &User{}
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		gormGet(u)
//	}
//}
//
//func BenchmarkRawSql(b *testing.B) {
//	var u = &User{}
//	b.ReportAllocs()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		rawSqlGet(u)
//	}
//}
