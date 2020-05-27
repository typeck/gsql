package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/util"
	"strings"
)

// Wrapper of sql.DB
type DB struct {
	DriverName string
	*sql.DB
	logger Logger
	Error  []errors.Error
}

var defaultDb *DB

type SqlInfo struct {
	tableName	string
	sql  		string
	method 		[]string
	condition	[]string
	values 		[]interface{}
}

func SetDb(db *DB) {
	defaultDb = db
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
	return &DB{DriverName: driverName,DB: db},nil
}

func (db *DB)PrepareSql() *SqlInfo{
	return &SqlInfo{}
}

func (db *DB) Select(s *SqlInfo, dest ...interface{}) error {
	s.done()
	rows, err := db.Query(s.sql, s.values...)
	if err != nil {
		return errors.New("err:%v,sql:%s",err,s.sql)
	}
	defer rows.Close()

	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}
	err = rows.Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlInfo)Table(tableName string) *SqlInfo {
	s.tableName = tableName
	return s
}
func (s *SqlInfo)Query(args... string) *SqlInfo {
	rawQuery := strings.Join(args,",")
	s.sql = util.Join(" ","SELECT", "FROM", s.tableName, rawQuery)
	return s
}

func (s *SqlInfo)Where(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"WHERE")
	s.condition = append(s.condition,condition)
	s.values = append(s.values,args...)
	return s
}

func (s *SqlInfo)And(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"AND")
	s.condition = append(s.condition,condition)
	s.values = append(s.values,args...)
	return s
}

func (s *SqlInfo)Or(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"OR")
	s.condition = append(s.condition,condition)
	s.values = append(s.values,args...)
	return s
}

func(s *SqlInfo)Raw(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"")
	s.condition = append(s.condition,condition)
	s.values = append(s.values,args...)
	return s
}

func(s *SqlInfo)done()*SqlInfo {

	for i, method := range s.method {
		s.sql = util.Join(" ",s.sql, method, s.condition[i])
	}
	return s
}

func (s *SqlInfo) Insert(args... string) *SqlInfo {
	var ph string
	for i := 0; i < len(args); i++ {
		ph = util.Join("," ,ph, "?")
	}
	ph = "VALUES (" + ph + ")"

	columns := util.Join("," , args...)

	s.sql = util.Join(" ", "INSERT INTO ",s.tableName, "(", columns , ")", ph)

	return s
}

func (s *SqlInfo) Update(args... string) *SqlInfo {
	var param string
	for _,v := range args {
		param = util.Join(",", param, v + "=?")
	}
	s.sql = util.Join(" ","UPDATE", s.tableName, "SET", param)
	return s
}

func (s *SqlInfo)Values(args... interface{}) *SqlInfo {

	s.values = append(s.values,args...)
	return s
}

func (s *SqlInfo)Done()(string,[]interface{}) {
	s.done()
	return s.sql, s.values
}