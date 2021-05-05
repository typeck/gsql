# gsql
Lightweight orm and extension of raw sql.

## Feathers

* Orm support
* high-performance, without reflect value
* Sql Generator support ORM and raw SQL operation
* Compatible with `database/sql`
* Driver support: [MYSQL](github.com/go-sql-driver/mysql),[postgres](https://github.com/lib/pq)
* omitempty support
* Transaction support

## Usages

* data structure
```go
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

```

* Open database

```go

db, err := gsql.OpenDb("mysql","type:tang@(127.0.0.1)/test?parseTime=true")

```

* query

```go
//Query one 
var u = &User{}
err := db.Pre().Debug().Cols("id", "name").Where("id=?",1).Get(u).Err()

//Query multiple
var us []*User
err := db.Pre().Debug().Where("id>?", 5).Gets(&us).Err()
```

* exec

```go
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

//Insert
id,err := db.Pre().Debug().Create(&u).LastInsertId()

//Update
affect,err := db.Pre().Debug().Where("id=?", 23).Update(&u).RowsAffected()
```