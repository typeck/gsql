package types

import (
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/util"
	"reflect"
	"strings"
)


type StructField struct {
	Typ 		reflect.Type
	Offset 		uintptr
	TagName 	string
	Tags 		[]string
}


type StructInfo struct {
	Name 		string
	//reflect type of *struct
	Typ 		reflect.Type
	Fields 		map[string] StructField
}


// unwrap * struct type into map.
func(s *StructInfo) Unwrap(t reflect.Type, tagName string)  error {
	s.Typ = t
	if t.Kind() != reflect.Ptr {
		return  errors.New("must pass a pointer.")
	}
	t = indirection(t)
	if t.Kind() != reflect.Struct {
		return  errors.New("must pass a struct pointer.")
	}
	s.Name = util.ToSnakeCase(t.Name())

	size := t.NumField()
	var field StructField
	for i := 0; i < size; i++ {
		tags := t.Field(i).Tag.Get(tagName)
		//split "" will also return a slice contains "", so filter nil tags is necessary.
		if tags == "" {
			continue
		}
		tagSlice := strings.Split(tags, ",")
		if len(tagSlice) == 0 {
			continue
		}

		field.Typ = t.Field(i).Type
		field.Offset = t.Field(i).Offset
		field.TagName = tagSlice[0]
		field.Tags = tagSlice

		s.Fields[tagSlice[0]] = field
	}

	return nil
}

func(s *StructInfo) GetNameAndCols() (string, []string){
	var cols []string

	for k, _ := range s.Fields {
		cols = append(cols, k)
	}
	return s.Name, cols
}


func NewStructInfo() *StructInfo {
	return &StructInfo{
		Fields: make(map[string]StructField),
	}
}
