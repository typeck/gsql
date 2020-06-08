package gsql

import (
	"reflect"
	"unsafe"
)

type emptyInterface struct {
	typ  _type
	data unsafe.Pointer
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}


type _type unsafe.Pointer

func packEFace(rtype _type, data unsafe.Pointer) interface{} {
	var i interface{}
	e := (*emptyInterface)(unsafe.Pointer(&i))
	e.typ = rtype
	e.data = data
	return i
}

func unpackEFace(obj interface{}) *emptyInterface {
	return (*emptyInterface)(unsafe.Pointer(&obj))
}

func indirection(p reflect.Type) reflect.Type {
	return p.Elem()
}

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype unsafe.Pointer) unsafe.Pointer