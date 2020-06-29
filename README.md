# gsql
Lightweight orm and extension of raw sql.

## Feathers

* Orm support
* high-performance, without reflect value
* Sql Generator support ORM and raw SQL operation
* Compatible with `database/sql`
* Driver support: [MYSQL](github.com/go-sql-driver/mysql)
* Transaction support

## Usages

* data structure
```go
type User struct {
	Id         int    `json:"id" db:"id" json:"id"`
	Name       string `json:"name" db:"name"`
	Account    string `json:"account" db:"account"`
	Company    string `json:"company" db:"company"`
	Pwd        string `json:"pwd" db:"pwd"`
	Phone      string `json:"phone" db:"phone"`
	Email      string `json:"email" db:"email"`
	Roles      int    `json:"role" db:"roles"`
	Status     int    `json:"status" db:"status"`
	Time      Time    `json:"time"`
}

type Time struct {
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}
```

* Open database

`db, err = gsql.OpenDb("mysql","type:tang@(127.0.0.1)/test")`

* query

```go
//Query one 
err := db.New().Debug().Table("user").Where("id=1").
	Cols("id", "name", "status", "company", "account", "company", "email, phone, pwd").
	Query(&u.Id, &u.Name, &u.Status, &u.Company, &u.Account, &u.Company,  &u.Email, &u.Phone, &u.Pwd).Err()

err := db.New().Raw("select id,name from user where id=1").And("name=?",name).Query(u.Id, u.Name).Err()

var u = &User{}
err := db.New().Debug().Where("id=?",1).Cols("name, id").Get(u).Err()

//Query multiple
rows, err := db.New().Debug().Table("user").Rows()
for rows.Next(){
    ...
}

var us []*User
err := db.New().Debug().Gets(&us).Err()
```

* exec

```go
id,err := db.New().Table("user").Debug().Cols("name, pwd").Exec("insert",u.Name, u.Pwd).LastInsertId()

affectRows, err := db.New().Table("user").Debug().Where("name=?",&u.Name).And("pwd=?",&u.Pwd).
	Cols("phone, company").Exec("update", "133","afg").RowsAffected()


id,err := db.New().Debug().Create(&u).LastInsertId()

affect,err := db.New().Debug().Where("id=?", 3).Cols("name, account").Update(&u).RowsAffected()

err := db.New().Table("user").Debug().Where("id=?",id).Exec("delete").Err()
```