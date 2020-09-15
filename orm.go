package gsql

import (
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/types"
	"reflect"
	"sync"
	"unsafe"
)

type Orm struct {
	Tag		 			string
	structCache 		sync.Map//map[_type] StructInfo
	sliceCache 			sync.Map//map[_type] SliceInfo
}

func NewOrm() *Orm {
	return &Orm{
		Tag:         "db",
	}
}

func (o *Orm) BuildValuesCols(s *SqlInfo, dest interface{}) ([]interface{}, error) {
	ptr := types.UnpackEFace(dest).Data
	var values []interface{}
	structInfo, err := o.GetStructInfo(dest)
	if err != nil {
		return nil, err
	}
	if structInfo == nil {
		return nil, errors.New("nil struct map cache.")
	}
	if len(s.cols) != 0 {
		values, err = structInfo.BuildValuesByPtr(ptr, s.cols)
		if err != nil {
			return nil, err
		}
	}else {
		values, s.cols, err = structInfo.BuildValuesCols(ptr)
		if err != nil {
			return nil, err
		}
	}
	if s.tableName == "" {
		s.tableName = structInfo.Name
	}
	return values, nil
}

func (o *Orm)GetStructInfo(dest interface{}) (*types.StructInfo, error){
	typ := types.UnpackEFace(dest).Typ
	if value, ok := o.structCache.Load(typ); ok {
		return value.(*types.StructInfo), nil
	}
	t := reflect.TypeOf(dest)

	s := types.NewStructInfo()
	err := s.SafeUnwrap(t, o.Tag)
	if err != nil {
		return nil, err
	}
	o.structCache.Store(typ, s)

	return s, nil
}

func (o *Orm) GetStructInfoByType(rTyp reflect.Type) (*types.StructInfo, error) {
	//get the slice elem type, as a structInfo cacheKey
	structTyp := types.UnpackEFace(rTyp).Data

	if value, ok := o.structCache.Load(structTyp); ok {
		return value.(*types.StructInfo), nil
	}

	s := types.NewStructInfo()
	err := s.SafeUnwrap(rTyp, o.Tag)
	if err != nil {
		return nil, err
	}
	o.structCache.Store(structTyp, s)

	return s, nil

}


func(o *Orm)GetSliceInfo(typ unsafe.Pointer, dest interface{}) (*types.SliceInfo, error){

	if value, ok := o.structCache.Load(typ); ok {
		return value.(*types.SliceInfo), nil
	}
	t := reflect.TypeOf(dest)
	s := types.NewSliceInfo()
	err := s.SafeUnwrap(t)
	if err != nil {
		return nil, err
	}
	o.structCache.Store(typ, s)

	return s, nil
}
// if tableName or columns is nil, we use struct name and tags to fill it.
func(o *Orm) InvokeCols(s *SqlInfo,dest interface{}) error {
	if s.tableName != "" && len(s.cols) > 0{
		return nil
	}
	structInfo, err := o.GetStructInfo(dest)
	if err != nil {
		return  err
	}
	if structInfo == nil {
		return  errors.New("nil struct map cache.")
	}
	if s.tableName == "" {
		s.tableName = structInfo.GetName()
	}
	if len(s.cols) == 0 {
		s.cols = structInfo.GetCols()
	}

	return  nil
}



