package gsql

import (
	"github.com/typeck/gsql/errors"
	"reflect"
	"strings"
)

type SqlInfo struct {
	tableName	string
	sql  		*sqlBuilder
	action 		string
	cols 		[]string
	params 		[]interface{}
	method 		[]string
	condition	[]string
	values 		[]interface{}
	isDebug 	bool
	driverName 	string
	execer 		Execer
}

func (s *SqlInfo) Reset() {
	s.tableName  = ""
	s.sql = &sqlBuilder{}
	s.action = ""
	s.cols = s.cols[:0]
	s.params = s.params[:0]
	s.method = s.method[:0]
	s.condition = s.condition[:0]
	s.values = s.values[:0]
	s.isDebug = false
}

func (s *SqlInfo)Query(dest... interface{}) Result {
	s.values = dest
	s.action = "SELECT"
	return s.execer.QueryVal(s, dest...)
}

func (s *SqlInfo) Exec(action string, dest... interface{}) Result {
	s.action = strings.ToUpper(action)
	s.values = dest
	s.values = append(s.values, s.params...)

	//exec operation don't need values to return,
	//so we add values into params head.
	s.params = s.values
	return s.execer.ExecVal(s)
}

//var u = User{}; Create(&u)
//need *struct type, but it's safe, the data of dest won't be change
//TODO: support struct type
func (s *SqlInfo) Create(dest interface{}) Result {
	s.action = "INSERT"
	return s.execer.ExecOrm(s, dest)
}

//var u = User{}; Update(&u)
//need *struct type, but it's safe, the data of dest won't be change
func (s *SqlInfo) Update(dest interface{}) Result {
	s.action = "UPDATE"
	return s.execer.ExecOrm(s, dest)
}

//var u = User{}; Get(&u)
func(s *SqlInfo) Get(dest interface{}) Result {
	s.action = "SELECT"
	return s.execer.Get(s, dest)
}

//var us [] *User; Gets(&us)
func(s *SqlInfo) Gets(dest interface{}) Result {
	s.action = "SELECT"
	return s.execer.Gets(s, dest)
}

//table name, Table("user")
func (s *SqlInfo)Table(tableName string) *SqlInfo {
	s.tableName = tableName
	return s
}

//columns, Cols("name, id", "phone")
func (s *SqlInfo)Cols(cols... string) *SqlInfo {
	var tmp [] string
	for i, _ := range cols {
		cols := strings.Split(cols[i], ",")
		for _, c := range cols {

			tmp = append(tmp, strings.TrimSpace(c))
		}
	}
	s.cols = tmp
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

func (s *SqlInfo) Wherem(m map[string]interface{}) *SqlInfo {
	var con = len(s.params)
	var i = con
	var str = &strings.Builder{}
	for k, v := range m {
		if i == con {
			s.method = append(s.method, "WHERE")
		}else {
			s.method = append(s.method, "AND")
		}
		str.WriteString(k)
		str.WriteString("=")
		s.execer.WritePlaceholder(str, i)
		s.condition = append(s.condition, str.String())
		s.params = append(s.params, v)
		str.Reset()
		i ++
	}
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

//raw condition; Condition("limit ?,?", 1,10)
func(s *SqlInfo)Condition(condition string, args... interface{}) *SqlInfo {

	s.method = append(s.method,"")
	s.condition = append(s.condition,condition)
	s.params = append(s.params,args...)
	return s
}

//raw sql; Raw("select * from user where id=?",1).And("name=?","jack")
func(s *SqlInfo)Raw(raw string, dest... interface{})*SqlInfo {
	s.sql.Builder.WriteString(raw)
	s.values = dest
	return s
}

func(s *SqlInfo)buildCondition() {

	for i, method := range s.method {
		s.sql.writeStrings(" ", method, " ", s.condition[i])
	}
}

func(s *SqlInfo)done() error {
	var err error
	if s.sql.Builder.Len() == 0 {
		switch s.action {
		case "SELECT":
			err = s.buildSelect()
		case "INSERT":
			err = s.buildInsert()
		case "UPDATE":
			err = s.buildUpdate()
		case "DELETE":
			err = s.buildDelete()
		}
	}
	if err != nil {
		return err
	}
	s.buildCondition()

	// print this sql.
	if s.isDebug {
		str := strings.ReplaceAll(s.sql.String(), s.execer.GetPlaceholder(0), "%v")
		var args []interface{}
		for i, p := range s.params {
			str = strings.ReplaceAll(str, s.execer.GetPlaceholder(i), "%v")
			t := reflect.TypeOf(p)
			if t.Kind() == reflect.Ptr {
				args = append(args, reflect.ValueOf(p).Elem().Interface())
				continue
			}
			args = append(args, p)
		}
		s.execer.Debug(str, s.params...)
	}
	return nil
}

func(s *SqlInfo)buildSelect() error{
	s.sql.Builder.WriteString("SElECT ")
	s.sql.join(",", s.cols...)
	s.sql.writeStrings(" FROM ", s.tableName)

	return nil
}

func (s *SqlInfo)buildInsert() error{
	if len(s.cols) == 0 {
		return errors.New("missing columns when exec insert.")
	}
	s.sql.writeStrings("INSERT INTO ", s.tableName, " (")
	s.sql.join(",", s.cols...)
	s.sql.Builder.WriteString(") VALUES (")

	s.execer.WritePlaceholder(s.sql, 0)
	for i := 1; i < len(s.cols); i++ {
		s.sql.Builder.WriteString(", ")
		s.execer.WritePlaceholder(s.sql, i)
	}
	s.sql.Builder.WriteString(")")

	return nil
}

func(s *SqlInfo) buildUpdate() error{
	if len (s.cols) == 0 {
		return errors.New("missing columns when exec update.")
	}
	s.sql.writeStrings("UPDATE ", s.tableName, " SET ")

	//switch len(s.cols) {
	//case 1:
	//	s.sql.writeStrings(s.cols[0], "=?")
	//default:
	//	s.sql.writeStrings(s.cols[0], "=?")
	//	for _, v:= range s.cols[1:] {
	//		s.sql.writeStrings(", ", v, "=?")
	//	}
	//}
	s.sql.writeStrings(s.cols[0], "=")
	s.execer.WritePlaceholder(s.sql, 0)
	for i := 1; i < len(s.cols); i++ {
		s.sql.writeStrings(", ", s.cols[i], "=")
		s.execer.WritePlaceholder(s.sql, i)
	}
	return nil
}

func (s *SqlInfo)buildDelete() error {
	s.sql.writeStrings("DELETE FROM ",s.tableName)
	return nil
}

func (s *SqlInfo)Debug() *SqlInfo{
	s.isDebug = true
	return s
}


//func (s *SqlInfo)Done()(string,[]interface{}) {
//	s.done()
//	return s.sql, s.values
//}