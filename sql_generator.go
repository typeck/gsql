package gsql

import (
	"strings"
)

type SqlInfo struct {
	tableName	string
	sql  		sqlBuilder
	action 		string
	cols 		[]string
	params 		[]interface{}
	method 		[]string
	condition	[]string
	values 		[]interface{}
	rawSql 		string
	isDebug 	bool
	driverName 	string
	db 			*DB
}

type sqlBuilder struct {
	strings.Builder
}

func(s *sqlBuilder)joinWith(with string, sep string, a... string) {
	s.Builder.WriteString(with)
	switch len(a) {
	case 0:
		return
	case 1:
		s.Builder.WriteString(a[0])
		return
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}
	s.Builder.Grow(n)
	s.Builder.WriteString(a[0])
	for _, ss := range a[1:] {
		s.Builder.WriteString(sep)
		s.Builder.WriteString(ss)
	}
}

func (s *sqlBuilder)writeStrings(args... string) {
	for _,v := range args {
		s.Builder.WriteString(v)
	}
}


func (s *SqlInfo)Query(dest... interface{}) *result {
	s.values = dest
	s.action = "SELECT"
	return s.db.queryVal(s, dest...)
}

func (s *SqlInfo) Exec(action string, dest... interface{}) *result {
	s.action = strings.ToUpper(action)
	s.values = dest
	s.values = append(s.values, s.params...)
	return s.db.exec(s)
}

func (s *SqlInfo) Insert(dest... interface{}) *result {
	s.action = "INSERT"
	s.values = dest
	s.values = append(s.values, s.params...)
	return s.db.exec(s)
}

func (s *SqlInfo) Update(dest... interface{}) *result {
	s.action = "UPDATE"
	s.values = dest
	s.values = append(s.values, s.params...)
	return s.db.exec(s)
}

func(s *SqlInfo) Get(dest interface{}) *result {
	s.action = "SELECT"
	return s.db.get(s, dest)
}

func(s *SqlInfo) Gets(dest interface{}) *result {
	s.action = "SELECT"
	return s.db.gets(s, dest)
}

func (s *SqlInfo)Table(tableName string) *SqlInfo {
	s.tableName = tableName
	return s
}

func (s *SqlInfo)Cols(cols... string) *SqlInfo {
	s.cols = cols
	return s
}

func (s *SqlInfo)RawSql(sql string) *SqlInfo {
	s.rawSql = sql
	return s
}

//where(id=?, 2)
//where(id=? AND name=?,2, jack)
func (s *SqlInfo)Where(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"WHERE")
	s.condition = append(s.condition,condition)
	s.params = append(s.params, args...)
	return s
}

func (s *SqlInfo)And(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"AND")
	s.condition = append(s.condition,condition)
	s.params = append(s.params,args...)
	return s
}

func (s *SqlInfo)Or(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"OR")
	s.condition = append(s.condition,condition)
	s.params = append(s.params,args...)
	return s
}

func(s *SqlInfo)Raw(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"")
	s.condition = append(s.condition,condition)
	s.params = append(s.params,args...)
	return s
}

func(s *SqlInfo)buildCondition() {

	for i, method := range s.method {
		s.sql.writeStrings(" ", method, " ", s.condition[i])
	}
}

func(s *SqlInfo)done() {
	switch s.action {
	case "SELECT":
		s.buildSelect()
	case "INSERT":
		s.buildInsert()
	case "UPDATE":
		s.buildUpdate()
	}

	if s.isDebug {
		str := strings.ReplaceAll(s.sql.String(), "?", "%v")
		s.db.logger.Printf(str, s.values...)
	}
}

func(s *SqlInfo)buildSelect() {
	s.sql.Builder.WriteString("SElECT ")
	s.sql.joinWith("", ",", s.cols...)
	s.sql.writeStrings(" FROM ", s.tableName)
	s.buildCondition()
}

func (s *SqlInfo)buildInsert() {
	s.sql.writeStrings("INSERT INTO ", s.tableName)
	var ph [] string
	for _, v := range s.cols {
		s.sql.joinWith(" (", ",", v)
		ph = append(ph, "?")
	}
	s.sql.Builder.WriteString(")")
	s.sql.joinWith(" VALUES (",",", ph...)
	s.sql.Builder.WriteString(")")
	s.buildCondition()
}

func(s *SqlInfo) buildUpdate() {
	s.sql.writeStrings("UPDATE ", s.tableName, "SET ")
	if len(s.cols) > 0 {
		switch len(s.cols) {
		case 1:
			s.sql.writeStrings(s.cols[0], "=?,")
		default:
			s.sql.writeStrings(s.cols[0], "=?,")
			for _, v:= range s.cols[1:] {
				s.sql.writeStrings(",", v, "=?")
			}
		}
	}
	s.buildCondition()
}

func (s *SqlInfo)Debug() *SqlInfo{
	s.isDebug = true
	return s
}


//func (s *SqlInfo)Done()(string,[]interface{}) {
//	s.done()
//	return s.sql, s.values
//}