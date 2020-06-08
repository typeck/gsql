package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/errors"
	"reflect"
	"unsafe"
)

type result struct {
	error   errors.Error
	result 	sql.Result
	row 	*sql.Row
	rows 	*sql.Rows
}


func (r *result) Err() error {
	return r.error
}

func (r *result)LastInsertId()(int64, error) {
	if r.error != nil {
		return 0, r.error
	}
	if r.result == nil {
		return 0,errors.New("wrong call.")
	}
	return r.result.LastInsertId()
}

func(r *result)RowsAffected()(int64, error) {
	if r.error != nil {
		return 0, r.error
	}
	if r.result == nil {
		return 0,errors.New("wrong call.")
	}

	return r.result.RowsAffected()
}

func (r *result)Rows() (*sql.Rows, error) {
	if r.rows == nil {
		return nil, errors.New("rows is nil.")
	}
	return r.rows, r.error
}

func( r *result)scanValues(dest ...interface{})  {
	if r.error != nil {
		return
	}
	defer r.rows.Close()
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			r.error = errors.New("sql: RawBytes isn't allowed on Row.Scan")
			return
		}
	}
	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			r.error = err
			return
		}
		r.error = sql.ErrNoRows
		return
	}
	err := r.rows.Scan(dest...)
	if err != nil {
		r.error = err
	}
}

func (r *result) scan(orm *Orm, dest interface{}) {
	cols, err := r.rows.Columns()
	if err != nil {
		r.error = err
		return
	}

	values, err := orm.BuildValues(dest, cols)
	if err != nil {
		r.error = err
		return
	}
	r.scanValues(values...)
}

func (r *result) scanAll(orm *Orm, dest interface{}) {
	if r.error != nil {
		return
	}
	defer r.rows.Close()
	cols, err := r.rows.Columns()
	if err != nil{
		r.error = err
		return
	}

	cacheKey := unpackEFace(dest).typ

	sliceInfo,err := orm.getSlice(cacheKey, dest)
	if err != nil {
		r.error = err
		return
	}
	if sliceInfo == nil {
		r.error = errors.New("slice info is nil.")
	}
	structInfo, err := orm.getStructInfoByType(sliceInfo.elemTyp)
	if err != nil {
		r.error = err
		return
	}

	for r.rows.Next() {
		//new struct
		ptr := unsafe_New(unpackEFace(structInfo.typ.Elem()).data)
		values, err := orm.buildValues(ptr, structInfo.fields, cols)
		if err != nil {
			r.error = err
			return
		}
		err = r.rows.Scan(values...)
		if err != nil {
			r.error	= err
			return
		}
		// use reflect type of *struct to build **struct interface v
		pPtr := &ptr
		pTyp := reflect.PtrTo(structInfo.typ)
		v := packEFace(_type(unpackEFace(pTyp).data), unsafe.Pointer(pPtr))

		sliceInfo.typ2.Append(dest, v)
	}
}


