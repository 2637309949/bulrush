// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
)

type (
	// Scope defined plugin hooks
	Scope interface {
		Pre() error
		Post() error
		Plugin() []interface{}
		Type() reflect.Type
	}
)

func newConstructorScope(src interface{}, acquire func(reflect.Type) interface{}) Scope {
	return &ConstructorScope{
		v:       src,
		acquire: acquire,
	}
}

func newObjectScope(src interface{}, acquire func(reflect.Type) interface{}) Scope {
	return &ObjectScope{
		v:       src,
		acquire: acquire,
	}
}

func newScope(src interface{}, acquire func(reflect.Type) interface{}) (s Scope) {
	vType := reflect.TypeOf(src)
	if vType.Kind() == reflect.Interface {
		vType = vType.Elem()
	}
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	switch true {
	case vType.Kind() == reflect.Func:
		s = newConstructorScope(src, acquire)
	case vType.Kind() == reflect.Struct:
		s = newObjectScope(src, acquire)
	}
	return
}
