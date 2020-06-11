package gsql

import (
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/util"
	"reflect"
	"strings"
)


type structField struct {
	typ 		reflect.Type
	offset 		uintptr
	tagName 	string
	tags 		[]string
}


type structInfo struct {
	name 		string
	//reflect type of *struct
	typ 		reflect.Type
	fields 		map[string] structField
}



func(s *structInfo) unwrap(t reflect.Type, tagName string)  error {
	s.typ = t
	if t.Kind() != reflect.Ptr {
		return  errors.New("must pass a pointer.")
	}
	t = indirection(t)
	if t.Kind() != reflect.Struct {
		return  errors.New("must pass a struct pointer.")
	}
	s.name = util.ToSnakeCase(t.Name())

	size := t.NumField()
	var field structField
	for i := 0; i < size; i++ {
		tags := t.Field(i).Tag.Get(tagName)
		tagSlice := strings.Split(tags, ",")
		if len(tagSlice) == 0 {
			continue
		}

		field.typ = t.Field(i).Type
		field.offset = t.Field(i).Offset
		field.tagName = tagSlice[0]
		field.tags = tagSlice

		s.fields[tagSlice[0]] = field
	}

	return nil
}

func(s *structInfo) invokeCols(sl *SqlInfo) {
	if sl.tableName == "" {
		sl.tableName = s.name
	}
	if len(sl.cols) == 0 {
		for k, _ := range s.fields {
			sl.cols = append(sl.cols, k)
		}
	}
}


func newStructInfo() *structInfo {
	return &structInfo{
		fields: make(map[string]structField),
	}
}

