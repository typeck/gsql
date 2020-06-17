package types

import (
	"github.com/modern-go/reflect2"
	"github.com/typeck/gsql/errors"
	"reflect"
	"unsafe"
)

type SliceInfo struct {
	Typ2 		reflect2.SliceType
	ElemTyp 	reflect.Type// *struct
}

func (s *SliceInfo) TypeCheck(t reflect.Type) error {
	if t.Kind() != reflect.Ptr {
		return errors.New("slice must be *[]*struct.")
	}
	t1 := Indirection(t)
	if t1.Kind() != reflect.Slice {
		return errors.New("slice must be *[]*struct.")
	}
	if t1 = Indirection(t1); t1.Kind() != reflect.Ptr {
		return errors.New("slice must be *[]*struct.")
	}
	s.ElemTyp = t1
	if t1 = Indirection(t1); t1.Kind() != reflect.Struct {
		return  errors.New("slice must be *[]*struct.")
	}
	return nil
}

func(s *SliceInfo) SafeUnwrap(t reflect.Type) error {
	if err := s.TypeCheck(t); err != nil {
		return err
	}

	s.Unwrap(t.Elem())
	return nil
}

// t must be []struct type.
func(s *SliceInfo) Unwrap(t reflect.Type)  {
	type2 := reflect2.Type2(t)
	s.Typ2 = type2.(reflect2.SliceType)
	s.ElemTyp = s.GetSliceElemType(t)

}

func (s *SliceInfo) GetSliceElemType(t reflect.Type) reflect.Type {
	for {
		if t.Kind() != reflect.Ptr {
			return t.Elem()
		}
		t = t.Elem()
	}
}

func (s *SliceInfo) Append(destPtr unsafe.Pointer, pPtr unsafe.Pointer) {
	s.Typ2.UnsafeAppend(destPtr, unsafe.Pointer(pPtr))
}

func NewSliceInfo() *SliceInfo {
	return &SliceInfo{
	}
}