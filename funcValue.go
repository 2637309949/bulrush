// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

const (
	preHookName    = "Pre"
	postHookName   = "Post"
	pluginHookName = "Plugin"
)

type funcValue struct {
	value  reflect.Value
	inputs []reflect.Value
}

func (fv *funcValue) runPre() {
	pre, fromStruct := indirectFunc(fv.value.Interface(), preHookName)
	if pre != nil && fromStruct {
		value := reflect.ValueOf(pre)
		if value.IsValid() {
			value.Call([]reflect.Value{})
		}
	}
}

func (fv *funcValue) runPlugin() interface{} {
	funcItem := indirectPlugin(fv.value.Interface())
	funcValue := reflect.ValueOf(funcItem)
	ret := funcValue.Call(fv.inputs)
	ret = funk.Filter(ret, func(v reflect.Value) bool {
		return v.IsValid()
	}).([]reflect.Value)
	results := funk.Map(ret, func(v reflect.Value) interface{} {
		return v.Interface()
	})
	return results
}

func (fv *funcValue) runPost() {
	pre, fromStruct := indirectFunc(fv.value.Interface(), postHookName)
	if pre != nil && fromStruct {
		value := reflect.ValueOf(pre)
		if value.IsValid() {
			value.Call([]reflect.Value{})
		}
	}
}

func (fv *funcValue) inputsFrom(inputs []interface{}) {
	funcItem := indirectPlugin(fv.value.Interface())
	funcValue := reflect.ValueOf(funcItem)
	if funcValue.Type().Kind() != reflect.Func {
		panic(fmt.Errorf(" %v inputsFrom call with %v error", funcItem, inputs))
	}
	funcType := funcValue.Type()
	numIn := funcType.NumIn()
	for index := 0; index < numIn; index++ {
		ptype := funcType.In(index)
		eleValue := typeMatcher(ptype, inputs)
		if eleValue == nil {
			eleValue = duckMatcher(ptype, inputs)
		}
		if eleValue == nil {
			panic(fmt.Errorf("inputsFrom %v call with %v error", ptype, reflect.TypeOf(inputs[7])))
		}
		fv.inputs = append(fv.inputs, eleValue.(reflect.Value))
	}
}

// duckMatcher match type if from target`type
func typeMatcher(ptype reflect.Type, params []interface{}) interface{} {
	target := retrieveType(ptype, params)
	if target != nil {
		return reflect.ValueOf(target)
	}
	return nil
}

// duckMatcher match type if implements target`interface
func duckMatcher(ptype reflect.Type, params []interface{}) interface{} {
	target := retrieveInterface(ptype, params)
	if target != nil {
		return reflect.ValueOf(target)
	}
	return nil
}

// retrieve type from given types
func retrieveType(ptype reflect.Type, types []interface{}) interface{} {
	target := funk.Find(types, func(x interface{}) bool {
		return ptype == reflect.TypeOf(x)
	})
	return target
}

// retrieve type whether to implement the interface
func retrieveInterface(ptype reflect.Type, types []interface{}) interface{} {
	target := funk.Find(types, func(x interface{}) bool {
		if ptype.Kind() == reflect.Interface {
			return reflect.TypeOf(x).Implements(ptype)
		}
		return false
	})
	return target
}
