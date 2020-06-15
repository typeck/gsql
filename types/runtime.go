package types

import (
	"reflect"
	"unsafe"
)

type EmptyInterface struct {
	Typ  unsafe.Pointer
	Data unsafe.Pointer
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}



func packEFace(rtype unsafe.Pointer, data unsafe.Pointer) interface{} {
	var i interface{}
	e := (*EmptyInterface)(unsafe.Pointer(&i))
	e.Typ = rtype
	e.Data = data
	return i
}

func UnpackEFace(obj interface{}) *EmptyInterface {
	return (*EmptyInterface)(unsafe.Pointer(&obj))
}

func Indirection(p reflect.Type) reflect.Type {
	return p.Elem()
}

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype unsafe.Pointer) unsafe.Pointer