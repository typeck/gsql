package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/types"
	"log"
	"os"
	"unsafe"
)

// Wrapper of sql.DB
type DB struct {
	driverName string
	*sql.DB
	logger Logger
	orm    *Orm
}

type Db interface {
	//new for every execution
	New() *SqlInfo
	//default tag is "db", to customize struct tag.
	SetTag(tagName string)
	//customize debug logger.
	SetLog(log Logger)
}

type Execer interface {
	//query into vals
	QueryVal(s *SqlInfo, dest... interface{}) Result
	//exec vals
	ExecVal(s *SqlInfo, dest... interface{}) Result
	//query into *struct
	Get(s *SqlInfo, dest interface{}) Result
	//query into *[]*struct
	Gets(s *SqlInfo, dest interface{}) Result
	//exec *struct
	ExecOrm(s *SqlInfo, dest interface{})Result
	//use logger to printf debug log
	Debug(format string, v ...interface{})
}

var defaultLog = log.New(os.Stdout, "[gsql]", log.Lshortfile|log.Ldate|log.Ltime)

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
		orm:		NewOrm(),
		logger: defaultLog,
	},nil
}

func (db *DB)New() *SqlInfo{
	return &SqlInfo{driverName: db.driverName, execer: db}
}


func (db *DB)SetTag(tagName string) {
	db.orm.Tag = tagName
}

func(db *DB)SetLog(log Logger) {
	db.logger = log
}


func (db *DB) QueryVal(s *SqlInfo, dest... interface{}) Result {
	res := db.query(s)
	scanVal(res, dest...)
	return res
}

func scanVal(scanner Scanner, dest... interface{}) {
	scanner.scanVal(dest...)
}

func(db *DB)ExecVal(s *SqlInfo, dest... interface{}) Result {
	return db.exec(s)
}

func(db *DB) Get(s *SqlInfo, dest interface{}) Result {
	orm := db.orm
	 res := &result{}
	err := orm.InvokeCols(s, dest)
	if err != nil {
		res.error = err
		return res
	}

	res = db.query(s)
	if res.error != nil {
		return res
	}

	scan(res, orm, dest)
	return res
}

func scan(scanner Scanner,orm *Orm, dest interface{}) {
	scanner.scan(orm, dest)
}


func (db *DB) Gets(s *SqlInfo, dest interface{}) Result {
	orm := db.orm
	res := &result{}
	cacheKey := types.UnpackEFace(dest).Typ
	destPtr := types.UnpackEFace(dest).Data
	sliceInfo,err := orm.GetSliceInfo(cacheKey, dest)
	if err != nil {
		res.error = err
		return res
	}
	if sliceInfo == nil {
		res.error = errors.New("slice info is nil.")
	}
	structInfo, err := orm.GetStructInfoByType(sliceInfo.ElemTyp)
	if err != nil {
		res.error = err
		return res
	}
	invokeCols(s, structInfo)

	res = db.query(s)
	if res.error != nil {
		return res
	}
	scanAll(res,orm, destPtr, structInfo, sliceInfo)
	return res
}

func invokeCols(s *SqlInfo, structInfo *types.StructInfo) {
	name, cols := structInfo.GetNameAndCols()
	if s.tableName == "" {
		s.tableName = name
	}
	if len(s.cols) == 0 {
		s.cols = cols
	}
}

func scanAll(scanner Scanner, orm *Orm, destPtr unsafe.Pointer,structInfo *types.StructInfo, sliceInfo *types.SliceInfo) {
	scanner.scanAll(orm, destPtr, structInfo, sliceInfo)
}

func (db *DB) ExecOrm(s *SqlInfo, dest interface{})Result {
	orm := db.orm
	res := &result{}
	err := orm.InvokeCols(s, dest)
	if err != nil {
		res.error = err
		return res
	}

	values, err := orm.BuildValues(dest, s.cols)
	if err != nil {
		res.error = err
		return res
	}
	s.values = values
	s.values = append(s.values, s.params...)
	s.params = s.values

	return db.exec(s)
}

func (db *DB) Debug(format string, v ...interface{}) {
	db.logger.Printf(format, v...)
}

func (db *DB) query(s *SqlInfo) *result {
	err := s.done()
	if err != nil {
		return &result{
			error: err,
		}
	}
	rows, err := db.Query(s.sql.String(), s.params...)
	return &result{
		rows: rows,
		error: err,
	}
}

func(db *DB) exec(s *SqlInfo) *result {
	err := s.done()
	if err != nil {
		return &result{
			error: err,
		}
	}
	res, err := db.Exec(s.sql.String(), s.params...)
	return &result{
		result: res,
		error: err,
	}
}