package gsql

import (
	"github.com/typeck/gsql/util"
	"strings"
)

type SqlInfo struct {
	tableName	string
	sql  		string
	isQuery 	bool
	method 		[]string
	condition	[]string
	values 		[]interface{}
	isDebug 	bool
	driverName 	string
	db 			*DB
}



func (s *SqlInfo)Table(tableName string) *SqlInfo {
	s.tableName = tableName
	return s
}

func (s *SqlInfo)Query(args... string) *SqlInfo {
	rawQuery := util.Join(",", args...)
	s.sql = util.Join(" ","SELECT",  rawQuery, "FROM", s.tableName)
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
	if strings.Contains(strings.ToLower(s.sql),"select") {
		s.isQuery = true
	}
	if s.isDebug {
		str := strings.ReplaceAll(s.sql, "?", "%v")
		s.db.logger.Printf(str, s.values...)
	}
	return s
}

func (s *SqlInfo) Insert(args... string) *SqlInfo {
	if len(args) == 0 {
		return s
	}
	var ph = "?"
	for i := 0; i < len(args)-1; i++ {
		ph = util.Join("," ,ph, "?")
	}
	ph = "VALUES (" + ph + ")"

	columns := util.Join("," , args...)

	s.sql = util.Join(" ", "INSERT INTO ",s.tableName, "(", columns , ")", ph)

	return s
}

func (s *SqlInfo) Update(args... string) *SqlInfo {
	if len(args) == 0 {
		return s
	}
	var param string
	for _,v := range args[:] {
		param = util.Join(",", param, v + "=?")
	}
	s.sql = util.Join(" ","UPDATE", s.tableName, "SET", param[1:])
	return s
}

func (s *SqlInfo)Values(args... interface{}) *SqlInfo {

	s.values = append(s.values,args...)
	return s
}
func (s *SqlInfo)Debug() *SqlInfo{
	s.isDebug = true
	return s
}
//func(s *SqlInfo)Select(sqlInfo *SqlInfo, dest ...interface{}) error {
//	if s.db == nil {
//		return errors.New("must set default db.")
//	}
//
//	return s.db.query(sqlInfo, dest)
//}

func (s *SqlInfo)Done()(string,[]interface{}) {
	s.done()
	return s.sql, s.values
}