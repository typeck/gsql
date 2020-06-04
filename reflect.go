package gsql

import (
	"github.com/typeck/gsql/errors"
	"reflect"
	"strings"
	"unsafe"
)

type emptyInterface struct {
	typ  *struct{}
	data unsafe.Pointer
}

type StructField struct {
	Type 		reflect.Type
	Offset 		uintptr
	TagName 	string
	Tags 		string
}

type StructInfo map[string] StructField

type Orm struct {
	Tag		 			string
	StructCache 		map[string] StructInfo
}

func(o *Orm) Struct2Map(t reflect.Type) StructInfo {
	size := t.NumField()
	m := make(StructInfo, size)

	var field StructField
	for i := 0; i < size; i++ {
		tags := t.Field(i).Tag.Get(o.Tag)
		tagSlice := strings.Split(tags, ",")
		if len(tagSlice) == 0 {
			continue
		}

		field.Type = t.Field(i).Type
		field.Offset = t.Field(i).Offset
		field.TagName = tagSlice[0]
		field.Tags = tags

		m[tagSlice[0]] = field
	}

	return m
}

func (o *Orm) BuildValues(s *SqlInfo, dest interface{}, cols[]string) ([]interface{}, error) {
	var values [] interface{}
	structInfo, err := o.GetStructMap(s, dest)
	if err != nil {
		return nil, err
	}
	if structInfo == nil {
		return nil, errors.New("nil struct cache.")
	}
	ptr := (*emptyInterface)(unsafe.Pointer(&dest)).data
	for _, col := range cols {
		if v, ok := structInfo[col]; ok {
			switch v.Type.Kind() {
			case reflect.Int:
				filedPtr := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + v.Offset))
				values = append(values, filedPtr)
			case reflect.String:
				filedPtr := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Int64:
				filedPtr := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + v.Offset))
				values = append(values, filedPtr)
			case reflect.Float64:
				filedPtr := (*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + v.Offset))
				values = append(values, filedPtr)
			}

		}
	}
	return values, nil
}

func (o *Orm)GetStructMap(s *SqlInfo, dest interface{}) (StructInfo, error){
	var m  StructInfo
	if m, ok := o.StructCache[s.tableName]; !ok {
		t := reflect.TypeOf(dest)
		if t.Kind() != reflect.Ptr {
			return nil, errors.New("must pass a pointer.")
		}
		if t.Elem().Kind() != reflect.Struct {
			return nil, errors.New("must pass a struct pointer.")
		}
		m = o.Struct2Map(t)
		o.StructCache[s.tableName] = m
	}
	return m, nil
}

func NewOrm() *Orm {
	return &Orm{
		Tag:         "db",
		StructCache: make(map[string]StructInfo),
	}
}
