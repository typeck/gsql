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
	EmbeddedStruct *StructInfo
	EmbeddedSlice  *SliceInfo
}


type StructInfo struct {
	Name 		string
	//type of struct
	Typ 				reflect.Type
	Fields 				map[string] StructField
	FieldsWithStruct 	map[string] StructField
	FieldsWithSlice 	map[string] StructField
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

func(s *StructInfo) Unwrap(t reflect.Type, tagName string)  {
	s.Typ = t
	t = GetBaseElem(t)
		s.Name = util.ToSnakeCase(t.Name())
		size := t.NumField()
		var field StructField
		for i := 0; i < size; i++ {
			tags := t.Field(i).Tag.Get(tagName)
			if tags =="-" {
				continue
			}
			//if tags is nil, we use field name as tag name.
			if tags == "" {
				tags = util.ToSnakeCase(t.Field(i).Name)
			}
			tagSlice := strings.Split(tags, ",")
			if len(tagSlice) == 0 {
				continue
			}
			field.Typ = t.Field(i).Type
			field.Offset = t.Field(i).Offset
			field.TagName = tagSlice[0]
			field.Tags = tagSlice

			tf := t.Field(i).Type
			te := GetBaseElem(t.Field(i).Type)
			if tf.Kind() == reflect.Ptr || tf.Kind() == reflect.Slice || tf.Kind() == reflect.Struct {
				if te.Kind() == reflect.Struct {
					structInfo := NewStructInfo()
					structInfo.Unwrap(t.Field(i).Type, tagName)
					field.EmbeddedStruct = structInfo
					s.FieldsWithStruct[tagSlice[0]] = field
				}
				if te.Kind() == reflect.Slice {
					sliceInfo := NewSliceInfo()
					sliceInfo.Unwrap(t.Field(i).Type)
					field.EmbeddedSlice = sliceInfo
					s.FieldsWithSlice[tagSlice[0]] = field
				}
				continue
			}
			s.Fields[tagSlice[0]] = field
		}
}



func(s *StructInfo) GetNameAndCols() (string, []string){
	var cols []string

	for k, _ := range s.Fields {
		cols = append(cols, k)
	}
	for _, v := range s.FieldsWithStruct {
		_, fCols := v.EmbeddedStruct.GetNameAndCols()
		cols = append(cols, fCols...)
	}
	return s.Name, cols
}

func (s *StructInfo) New() unsafe.Pointer {
	return unsafe_New(UnpackEFace(s.Typ.Elem()).Data)
}


func NewStructInfo() *StructInfo {
	return &StructInfo{
		Fields: make(map[string]StructField),
		FieldsWithStruct: make(map[string]StructField),
		FieldsWithSlice: make(map[string]StructField),
	}
}

