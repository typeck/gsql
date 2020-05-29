package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/errors"
)

// Wrapper of sql.DB
type DB struct {
	driverName string
	*sql.DB
	logger Logger
}


var defaultDb map[string] *DB

func SetDefaultDb(name string, db *DB) error {
	if _, ok := defaultDb[name]; ok {
		return errors.New("%s db is exist.",name)
	}
	defaultDb[name] = db
	return nil
}

func NewSql(driverName,dataSource string) (*DB,error) {
	db, err := sql.Open(driverName,dataSource)
	if err != nil {
		return nil , err
	}
	err = db.Ping()
	if err != nil {
		return nil,err
	}
	return &DB{driverName: driverName,DB: db},nil
}

func (db *DB)PrepareSql() *SqlInfo{
	return &SqlInfo{driverName: db.driverName}
}

func (db *DB)Execx(s *SqlInfo, dest ...interface{}) *result {
	var res *result
	if s.isQuery {
		db.query(s).scan(dest...)
		return res
	}

	return db.exec(s,dest)
}

func (db *DB) query(s *SqlInfo) *result {
	s.done()
	rows, err := db.Query(s.sql, s.values...)
	return &result{
		rows: rows,
		error: err,
	}

}

func(db *DB) exec(s *SqlInfo, dest ...interface{}) *result {
	s.done()
	if len(dest) != 0 {
		s.values = append(s.values, dest...)
	}
	res, err := db.Exec(s.sql, s.values)
	return &result{
		result: res,
		error: err,
	}
}

