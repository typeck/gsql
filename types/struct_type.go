package types

import (
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/util"
	"reflect"
	"strings"
	"unsafe"
)


type StructField struct {
	Typ            reflect.Type
	Offset         uintptr
	TagName        string
	Tags           []string
	EmbeddedStruct []*StructInfo
}


type StructInfo struct {
	Name 		string
	//reflect type of *struct
	Typ 		reflect.Type
	Fields 		map[string] StructField
}


func (s *StructInfo) typeCheck(t reflect.Type) error {
	if t.Kind() != reflect.Ptr {
		return  errors.New("must pass a pointer.")
	}
	t = Indirection(t)
	if t.Kind() != reflect.Struct {
		return  errors.New("must pass a struct pointer.")
	}
	return nil
}

func(s *StructInfo) SafeUnwrap(t reflect.Type, tagName string)  error {
	err := s.typeCheck(t)
	if err != nil {
		return err
	}
	s.Unwrap(t, tagName)
	return nil
}

// unwrap struct type into map.
func(s *StructInfo) Unwrap(t reflect.Type, tagName string)  {
	s.Typ = t
	t = GetElem(t)
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

}

func(s *StructInfo) GetNameAndCols() (string, []string){
	var cols []string

	for k, _ := range s.Fields {
		cols = append(cols, k)
	}
	return s.Name, cols
}

func (s *StructInfo) New() unsafe.Pointer {
	return unsafe_New(UnpackEFace(s.Typ.Elem()).Data)
}


func NewStructInfo() *StructInfo {
	return &StructInfo{
		Fields: make(map[string]StructField),
	}
}

