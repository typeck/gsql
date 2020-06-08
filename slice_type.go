package gsql

import (
	"github.com/modern-go/reflect2"
	"github.com/typeck/gsql/errors"
	"reflect"
)

type sliceInfo struct {
	typ2 		reflect2.SliceType
	elemTyp 	reflect.Type// *struct
}

func(s *sliceInfo) unwrap(t reflect.Type)error {
	if t.Kind() != reflect.Ptr {
		return errors.New("slice must be *[]*struct.")
	}
	t1 := indirection(t)
	if t1.Kind() != reflect.Slice {
		return errors.New("slice must be *[]*struct.")
	}
	if t1 = indirection(t1); t1.Kind() != reflect.Ptr {
		return errors.New("slice must be *[]*struct.")
	}
	s.elemTyp = t1
	if t1 = indirection(t1); t1.Kind() != reflect.Struct {
		return  errors.New("slice must be *[]*struct.")
	}

	//  reflect2 has cache of slice type.
	//we got a []*struct slice type, not a *[]*struct type.
	//so type2.(reflect2.SliceType)  is safe.
	type2 := reflect2.Type2(t.Elem())
	s.typ2 = type2.(reflect2.SliceType)
	return nil
}

func newSliceInfo() *sliceInfo {
	return &sliceInfo{
	}
}