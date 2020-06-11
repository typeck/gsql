package gsql

import (
	"github.com/typeck/gsql/errors"
	"reflect"
	"sync"
	"unsafe"
)

type Orm struct {
	tag		 			string
	structCache 		sync.Map//map[_type] StructInfo
	sliceCache 			sync.Map//map[_type] SliceInfo
}

func NewOrm() *Orm {
	return &Orm{
		tag:         "db",
		//StructCache: make(map[_type]StructInfo),
	}
}

func (o *Orm) BuildValues(dest interface{}, cols[]string) ([]interface{}, error) {
	ptr := unpackEFace(dest).data
	typ := unpackEFace(dest).typ

	structInfo, err := o.getStructInfo(typ, dest)
	if err != nil {
		return nil, err
	}
	if structInfo == nil {
		return nil, errors.New("nil struct map cache.")
	}
	fields := structInfo.fields
	return o.buildValues(ptr, fields, cols)
}

func (o *Orm)buildValues(ptr unsafe.Pointer, fields map[string]structField, cols []string) ([]interface{}, error){
	var values [] interface{}
	for _, col := range cols {
		if v, ok := fields[col]; ok {
			switch v.typ.Kind() {
			case reflect.Int:
				filedPtr := (*int)(unsafe.Pointer(uintptr(ptr) + v.offset))
				values = append(values, filedPtr)
			case reflect.String:
				filedPtr := (*string)(unsafe.Pointer(uintptr(ptr) + v.offset))
				values = append(values, filedPtr)
			case reflect.Int64:
				filedPtr := (*int64)(unsafe.Pointer(uintptr(ptr) + v.offset))
				values = append(values, filedPtr)
			case reflect.Float64:
				filedPtr := (*float64)(unsafe.Pointer(uintptr(ptr) + v.offset))
				values = append(values, filedPtr)
			}
		}
	}
	return values, nil
}

func (o *Orm)getStructInfo(typ _type, dest interface{}) (*structInfo, error){
	if value, ok := o.structCache.Load(typ); ok {
		return value.(*structInfo), nil
	}
	t := reflect.TypeOf(dest)

	s := newStructInfo()
	err := s.unwrap(t, o.tag)
	if err != nil {
		return nil, err
	}
	o.structCache.Store(typ, s)

	return s, nil
}

func (o *Orm) getStructInfoByType(rTyp reflect.Type) (*structInfo, error) {
	//get the slice elem type, as a structInfo cacheKey
	structTyp := _type(unpackEFace(rTyp).data)
	return o.getStructInfo(structTyp, rTyp)
}


func(o *Orm)getSlice(typ _type, dest interface{}) (*sliceInfo, error){

	if value, ok := o.structCache.Load(typ); ok {
		return value.(*sliceInfo), nil
	}
	t := reflect.TypeOf(dest)
	s := newSliceInfo()
	err := s.unwrap(t)
	if err != nil {
		return nil, err
	}
	o.structCache.Store(typ, s)

	return s, nil
}

func(o *Orm) invokeCols(s *SqlInfo,dest interface{}) error {
	if s.tableName != "" && len(s.cols) > 0{
		return nil
	}
	typ := unpackEFace(dest).typ

	structInfo, err := o.getStructInfo(typ, dest)
	if err != nil {
		return  err
	}
	if structInfo == nil {
		return  errors.New("nil struct map cache.")
	}
	structInfo.invokeCols(s)

	return  nil
}

