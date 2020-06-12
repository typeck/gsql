package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/errors"
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

func NewDb(driverName,dataSource string) (*DB,error) {
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
		orm:		newOrm(),
		logger: defaultLog,
	},nil
}

func (db *DB)New() *SqlInfo{
	return &SqlInfo{driverName: db.driverName, db: db}
}

func (db *DB) queryVal(s *SqlInfo, dest... interface{}) *result {
	res := db.query(s)
	res.scanValues(dest...)
	return res
}

func (db *DB) query(s *SqlInfo) *result {
	s.done()
	rows, err := db.Query(s.sql.String(), s.params...)
	return &result{
		rows: rows,
		error: err,
	}

}

func(db *DB) exec(s *SqlInfo) *result {
	s.done()
	res, err := db.Exec(s.sql.String(), s.values...)
	return &result{
		result: res,
		error: err,
	}
}

func(db *DB) get(s *SqlInfo, dest interface{}) *result {
	orm := db.orm
	 res := &result{}
	err := orm.invokeCols(s, dest)
	if err != nil {
		res.error = err
		return res
	}

	res = db.query(s)
	if res.error != nil {
		return res
	}
	res.scan(orm, dest)

	return res
}


func (db *DB) gets(s *SqlInfo, dest interface{}) *result {
	orm := db.orm
	res := &result{}
	cacheKey := unpackEFace(dest).typ
	destPtr := unpackEFace(dest).data
	sliceInfo,err := orm.getSlice(cacheKey, dest)
	if err != nil {
		res.error = err
		return res
	}
	if sliceInfo == nil {
		res.error = errors.New("slice info is nil.")
	}
	structInfo, err := orm.getStructInfoByType(sliceInfo.elemTyp)
	if err != nil {
		res.error = err
		return res
	}
	structInfo.invokeCols(s)

	res = db.query(s)
	if res.error != nil {
		return res
	}
	res.scanAll(orm, destPtr, structInfo, sliceInfo)
	return res
}

