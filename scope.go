// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

type (
	// Scope defined plugin scope
	Scope struct {
		Value  interface{}
		Inputs []reflect.Value
	}
)

const (
	preHookName    = "Pre"
	postHookName   = "Post"
	pluginHookName = "Plugin"
)

func (scope *Scope) inValue(t reflect.Type, inputs []interface{}) interface{} {
	return typeMatcher(t, inputs).(reflect.Value).Interface()
}

func (scope *Scope) methodCall(m reflect.Value, inputs []interface{}) {
	if m.IsValid() {
		switch method := m.Interface().(type) {
		case func():
			method()
		case func(*Config):
			method(scope.inValue(reflect.TypeOf(&Config{}), inputs).(*Config))
		default:
		}
	}
}

func (scope *Scope) reflectCall(m reflect.Value, ins []reflect.Value) []interface{} {
	return funk.Map(m.Call(ins), func(v reflect.Value) interface{} {
		return v.Interface()
	}).([]interface{})
}

func (scope *Scope) indirectFunc(name string) reflect.Value {
	if funk, fromStruct := indirectFunc(scope.Value, name); funk != nil && fromStruct {
		value := reflect.ValueOf(funk)
		if value.IsValid() {
			return value
		}
	}
	return reflect.Value{}
}

func (scope *Scope) indirectPlugin() reflect.Value {
	if funk := indirectPlugin(scope.Value, pluginHookName); funk != nil {
		value := reflect.ValueOf(funk)
		if value.IsValid() {
			return value
		}
	}
	return reflect.Value{}
}

func (scope *Scope) inFrom(inputs *Injects) {
	funk := scope.indirectPlugin()
	if funk.Type().Kind() != reflect.Func {
		panic(fmt.Errorf(" %v inputsFrom call with %v error", funk, inputs))
	}
	funcType := funk.Type()
	numIn := funcType.NumIn()
	for index := 0; index < numIn; index++ {
		ptype := funcType.In(index)
		v := inputs.Acquire(ptype)
		if v == nil {
			panic(fmt.Errorf("inputsFrom %v call with %v error", ptype, reflect.TypeOf(inputs)))
		}
		scope.Inputs = append(scope.Inputs, reflect.ValueOf(v))
	}
}

func newScope(src interface{}) *Scope {
	return &Scope{
		Value: src,
	}
}
