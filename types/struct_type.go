package types

import (
	"github.com/typeck/gsql/errors"
	"github.com/typeck/gsql/util"
	"reflect"
	"strings"
	"time"
	"unsafe"
)


type StructField struct {
	Typ            reflect.Type
	BaseType 	   reflect.Type
	Kind 		   reflect.Kind
	Offset         uintptr
	TagName        string
	Tags           []string
	OmitEmpty 		bool
	EmbeddedStruct *StructInfo
	EmbeddedSlice  *SliceInfo
}


type StructInfo struct {
	Name 		string
	//type of struct
	Typ 				reflect.Type
	Fields 				map[string] StructField
	//struct field
	FieldsWithStruct 	map[string] StructField
	//slice field
	FieldsWithSlice 	map[string] StructField
}


var (
	TimeType = reflect.TypeOf(time.Time{})
	TimePtrType = reflect.TypeOf(&time.Time{})
)

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
			if len(tagSlice) > 1 && tagSlice[1] == "omitempty"{
				field.OmitEmpty = true
			}
			field.Typ = t.Field(i).Type
			field.Kind = t.Field(i).Type.Kind()
			field.BaseType = GetBaseElem(t.Field(i).Type)
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



//using struct ptr and struct field offset to build values memory addr.
func (s *StructInfo)BuildValuesByPtr(ptr unsafe.Pointer, cols []string) ([]interface{}, error) {
	if ptr == nil {
		return nil, errors.New("nil ptr.")
	}
	fields := s.Fields
	fieldsWithStruct := s.FieldsWithStruct
	var values [] interface{}
	for _, col := range cols {
		if v, ok := fields[col]; ok {
			switch v.Kind {
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
		}else {
			for _, strus := range fieldsWithStruct {
				structInfo := strus.EmbeddedStruct
				if structInfo == nil {
					continue
				}
				var fPtr unsafe.Pointer
				switch strus.Kind {
				case reflect.Ptr:
					fPtr = unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(ptr) + strus.Offset)))
				case reflect.Struct:
					fPtr = unsafe.Pointer(uintptr(ptr) + strus.Offset)
				}
				vs, err := structInfo.BuildValuesByPtr(fPtr, []string{col})
				if err != nil {
					return nil, err
				}
				values = append(values, vs...)
			}
		}
	}
	return values, nil
}

func (s *StructInfo)BuildValuesCols(ptr unsafe.Pointer) ([]interface{}, []string, error){
	if ptr == nil {
		return nil, nil, errors.New("nil ptr.")
	}
	fields := s.Fields
	fieldsWithStruct := s.FieldsWithStruct
	var values [] interface{}
	var cols []string
	for col, field := range fields {
		switch field.Kind {
		case reflect.Int:
			filedPtr := (*int)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.String:
			filedPtr := (*string)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == "" && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Int64:
			filedPtr := (*int64)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Float64:
			filedPtr := (*float64)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Uint64:
			filedPtr := (*uint64)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Int32:
			filedPtr := (*int32)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Bool:
			filedPtr := (*bool)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == false && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Complex64:
			filedPtr := (*complex64)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Complex128:
			filedPtr := (*complex128)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Int16:
			filedPtr := (*int16)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Int8:
			filedPtr := (*int8)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Uint32:
			filedPtr := (*uint32)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Uint16:
			filedPtr := (*uint16)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Uint8:
			filedPtr := (*uint8)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		case reflect.Float32:
			filedPtr := (*float32)(unsafe.Pointer(uintptr(ptr) + field.Offset))
			if *filedPtr == 0 && field.OmitEmpty == true {
				continue
			}
			cols = append(cols, col)
			values = append(values, filedPtr)
		}

		for k, strus := range fieldsWithStruct {
			structInfo := strus.EmbeddedStruct
			if  structInfo == nil {
				continue
			}
			if strus.Typ == TimePtrType || strus.Typ == TimeType {
				cols = append(cols, k)
				if strus.Typ == TimeType {
					values = append(values, (*time.Time)(unsafe.Pointer(uintptr(ptr) + field.Offset)))
				}else {
					values = append(values, (*time.Time)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(ptr) + strus.Offset)))))
				}
				continue
			}
			var fPtr unsafe.Pointer
			switch strus.Kind {
			case reflect.Ptr:
				fPtr = unsafe.Pointer(*(*uintptr)(unsafe.Pointer(uintptr(ptr) + strus.Offset)))
			case reflect.Struct:
				fPtr = unsafe.Pointer(uintptr(ptr) + strus.Offset)
			}
			vs, cs, err := structInfo.BuildValuesCols(fPtr)
			if err != nil {
				return nil, nil, err
			}
			values = append(values, vs...)
			cols = append(cols, cs...)
		}

	}
	return values, cols, nil
}


func(s *StructInfo) GetCols()  []string {
	var cols []string

	for k, _ := range s.Fields {
		cols = append(cols, k)
	}
	for _, v := range s.FieldsWithStruct {
		fCols := v.EmbeddedStruct.GetCols()
		cols = append(cols, fCols...)
	}
	return  cols
}

func (s *StructInfo) GetName() string {
	return s.Name
}

func (s *StructInfo) New() unsafe.Pointer {
	return unsafe_New(UnpackEFace(s.Typ.Elem()).Data)
}

func UnsafeNew(typ reflect.Type)unsafe.Pointer {
	return unsafe_New(UnpackEFace(typ).Data)
}

func NewStructInfo() *StructInfo {
	return &StructInfo{
		Fields: make(map[string]StructField),
		FieldsWithStruct: make(map[string]StructField),
		FieldsWithSlice: make(map[string]StructField),
	}
}

