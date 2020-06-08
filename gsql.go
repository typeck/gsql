package gsql

import (
	"database/sql"
	"log"
	"os"
)

// Wrapper of sql.DB
type DB struct {
	driverName string
	*sql.DB
	logger Logger
	orm    *Orm
}

var defaultLog *log.Logger = log.New(os.Stdout, "[gsql]", log.Lshortfile|log.Ldate|log.Ltime)

func NewSql(driverName,dataSource string) (*DB,error) {
	db, err := sql.Open(driverName,dataSource)
	if err != nil {
		return nil , err
	}
	err = db.Ping()
	if err != nil {
		return nil,err
	}
	return &DB{
		driverName: driverName,
		DB: 		db,
		orm:		NewOrm(),
		logger: defaultLog,
	},nil
}

func (db *DB)PrepareSql() *SqlInfo{
	return &SqlInfo{driverName: db.driverName, db: db}
}

func (db *DB)ExecSql(s *SqlInfo, dest ...interface{}) *result {
	s.done()
	if s.isQuery {
		res := db.queryRows(s)
		res.scanValues(dest...)
		return res
	}

	return db.exec(s,dest)
}

func (db *DB) Scan(s *SqlInfo, dest interface{}) *result {
	s.done()
	return db.get(s, dest)
}

func (db *DB) ScanAll(s *SqlInfo, dest interface{}) *result {
	s.done()
	return db.gets(s, dest)
}

func (db *DB) queryRows(s *SqlInfo) *result {
	rows, err := db.Query(s.sql, s.values...)
	return &result{
		rows: rows,
		error: err,
	}

}

func(db *DB) exec(s *SqlInfo, dest ...interface{}) *result {
	if len(dest) != 0 {
		s.values = append(s.values, dest...)
	}
	res, err := db.Exec(s.sql, s.values)
	return &result{
		result: res,
		error: err,
	}
}

func(db *DB) get(s *SqlInfo, dest interface{}) *result {
	orm := db.orm

	res := db.queryRows(s)
	if res.error != nil {
		return res
	}
	res.scan(orm, dest)

	return res
}

func (db *DB) gets(s *SqlInfo, dest interface{}) *result {
	orm := db.orm
	res := db.queryRows(s)
	if res.error != nil {
		return res
	}
	res.scanAll(orm, dest)
	return res
}

