package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/driver"
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/types"
	"log"
	"os"
	"sync"
	"unsafe"
)

// Wrapper of sql.DB
type gsql struct {
	SqlDb
	driverName string
	driver.Driver
	logger Logger
	orm    *Orm
	pool   *sync.Pool
}

type Db interface {
	//new for every execution
	Pre() *SqlInfo
	//begin transaction
	Begin() (*gsql, error)
	Rollback() error
	Commit() error
	SqlDb
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
	driver.Driver
}

var (
	ErrNoRows = sql.ErrNoRows
)

var defaultLog = log.New(os.Stdout, "[gsql]", log.Lshortfile|log.Ldate|log.Ltime)

func OpenDb(driverName,dataSource string, opts... Option) (Db, error) {
	db, err := sql.Open(driverName,dataSource)
	if err != nil {
		return nil , err
	}
	err = db.Ping()
	if err != nil {
		return nil,err
	}
	gsqlDb := &gsql{
		driverName: driverName,
		SqlDb: 		db,
		orm:		NewOrm(),
		Driver:		driver.MDriver[driverName],
		logger: 	defaultLog,
		pool: 		&sync.Pool{},
	}
	gsqlDb.pool.New = func() interface{} {
		return gsqlDb.NewSqlInfo()
	}
	gsqlDb.WithOptions(opts...)
	return gsqlDb, nil
}

func (db *gsql) WithOptions(opts ...Option) Db {
	c := db.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}


func (db *gsql) Begin() (*gsql, error) {
	beginner, ok := (db.SqlDb).(TxBeginner)
	if !ok{
		return nil, errors.New("begin tx failed.")
	}
	tx, err := beginner.Begin()
	if err != nil {
		return nil, err
	}
	dbClone := db.clone()
	dbClone.SqlDb = tx
	return dbClone, nil
}

func(db *gsql) Rollback() error {
	committer, ok := (db.SqlDb).(TxCommitter)
	if !ok {
		return errors.New("rollback error, wrong caller.")
	}
	return committer.Rollback()
}

func(db *gsql) Commit() error {
	committer, ok := (db.SqlDb).(TxCommitter)
	if !ok {
		return errors.New("rollback error, wrong caller.")
	}
	return committer.Commit()
}


func (db *gsql) clone() *gsql {
	d := &gsql{
		driverName: db.driverName,
		SqlDb:      db.SqlDb,
		logger:     db.logger,
		orm:        db.orm,
		Driver:     db.Driver,
		pool: 		db.pool,
	}
	return d
}

func (db *gsql)Pre() *SqlInfo {
	s := db.pool.Get().(*SqlInfo)
	s.Reset()
	return s
}

func (db *gsql)NewSqlInfo() *SqlInfo{

	return &SqlInfo{
		driverName: db.driverName,
		execer: db,
		sql:	&sqlBuilder{},
	}
}


func (db *gsql) QueryVal(s *SqlInfo, dest... interface{}) Result {
	res := db.query(s)
	scanVal(res, dest...)
	return res
}

func scanVal(scanner Scanner, dest... interface{}) {
	scanner.scanVal(dest...)
}

func(db *gsql)ExecVal(s *SqlInfo, dest... interface{}) Result {
	return db.exec(s)
}

func(db *gsql) Get(s *SqlInfo, dest interface{}) Result {
	orm := db.orm
	res := &result{}
	values, err := orm.BuildValuesCols(s, dest, false)
	if err != nil {
		res.error = err
		return res
	}
	res = db.query(s)
	if res.error != nil {
		return res
	}

	scan(res, values...)
	return res
}

func scan(scanner Scanner,values... interface{}) {
	scanner.scanVal(values...)
}


func (db *gsql) Gets(s *SqlInfo, dest interface{}) Result {
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
	scanAll(res, destPtr, structInfo, sliceInfo)
	return res
}

func invokeCols(s *SqlInfo, structInfo *types.StructInfo) {
	if s.tableName == "" {
		s.tableName = structInfo.GetName()
	}
	if len(s.cols) == 0 {
		s.cols = structInfo.GetCols()
	}
}

func scanAll(scanner Scanner, destPtr unsafe.Pointer,structInfo *types.StructInfo, sliceInfo *types.SliceInfo) {
	scanner.scanAll(destPtr, structInfo, sliceInfo)
}

func (db *gsql) ExecOrm(s *SqlInfo, dest interface{})Result {
	orm := db.orm
	res := &result{}
	values, err := orm.BuildValuesCols(s, dest, true)
	if err != nil {
		res.error = err
		return res
	}
	values = append(values, s.params...)
	s.params = values

	return db.exec(s)
}

func (db *gsql) Debug(format string, v ...interface{}) {
	db.logger.Printf(format, v...)
}

func (db *gsql) query(s *SqlInfo) *result {
	err := s.done()
	if err != nil {
		return &result{
			error: err,
		}
	}
	rows, err := db.Query(s.sql.String(), s.params...)
	//put the sql info into sync pool
	db.pool.Put(s)
	return &result{
		rows: rows,
		error: err,
	}
}

func(db *gsql) exec(s *SqlInfo) *result {
	err := s.done()
	if err != nil {
		return &result{
			error: err,
		}
	}
	res, err := db.Exec(s.sql.String(), s.params...)
	db.pool.Put(s)
	return &result{
		result: res,
		error: err,
	}
}