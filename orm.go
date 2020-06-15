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

func (o *Orm) BuildValues(dest interface{}, cols[]string) ([]interface{}, error) {
	ptr := types.UnpackEFace(dest).Data
	typ := types.UnpackEFace(dest).Typ

	structInfo, err := o.getStructInfo(typ, dest)
	if err != nil {
		return nil, err
	}
	if structInfo == nil {
		return nil, errors.New("nil struct map cache.")
	}
	fields := structInfo.Fields
	return o.BuildValuesByPtr(ptr, fields, cols)
}

//using struct ptr and struct field offset to build values memory addr.
func (o *Orm)BuildValuesByPtr(ptr unsafe.Pointer, fields map[string]types.StructField, cols []string) ([]interface{}, error){
	var values [] interface{}
	for _, col := range cols {
		if v, ok := fields[col]; ok {
			switch v.Typ.Kind() {
			case reflect.Int:
				filedPtr := (*int)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.String:
				filedPtr := (*string)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Int64:
				filedPtr := (*int64)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Float64:
				filedPtr := (*float64)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Uint64:
				filedPtr := (*uint64)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Int32:
				filedPtr := (*int32)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Bool:
				filedPtr := (*bool)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Complex64:
				filedPtr := (*complex64)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Complex128:
				filedPtr := (*complex128)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Int16:
				filedPtr := (*int16)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Int8:
				filedPtr := (*int8)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Uint32:
				filedPtr := (*uint32)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Uint16:
				filedPtr := (*uint16)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Uint8:
				filedPtr := (*uint8)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Float32:
				filedPtr := (*float32)(unsafe.Pointer(uintptr(ptr) + v.Offset))
				values = append(values, filedPtr)
			}
		}
	}
	return values, nil
}

func (o *Orm)getStructInfo(typ unsafe.Pointer, dest interface{}) (*types.StructInfo, error){
	if value, ok := o.structCache.Load(typ); ok {
		return value.(*types.StructInfo), nil
	}
	t := reflect.TypeOf(dest)

	s := types.NewStructInfo()
	err := s.Unwrap(t, o.Tag)
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
	err := s.Unwrap(rTyp, o.Tag)
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
	err := s.Unwrap(t)
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
	typ := types.UnpackEFace(dest).Typ

	structInfo, err := o.getStructInfo(typ, dest)
	if err != nil {
		return  err
	}
	if structInfo == nil {
		return  errors.New("nil struct map cache.")
	}
	name, cols := structInfo.GetNameAndCols()
	if s.tableName == "" {
		s.tableName = name
	}
	if len(s.cols) == 0 {
		s.cols = cols
	}

	return  nil
}


