package gsql

import (
	"database/sql"
	"github.com/typeck/gsql/errors"
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
	return r.result.LastInsertId()
}

func(r *result)RowsAffected()(int64, error) {
	if r.error != nil {
		return 0, r.error
	}
	return r.result.RowsAffected()
}

func (r *result)Rows() (*sql.Rows, error) {

	return r.rows, r.error
}

func( r *result)scan(dest ...interface{})  {
	defer r.rows.Close()
	if r.error != nil {
		return
	}
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

