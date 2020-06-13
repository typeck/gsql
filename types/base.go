package types

import "reflect"

func indirection(p reflect.Type) reflect.Type {
	return p.Elem()
}