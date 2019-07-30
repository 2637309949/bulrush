// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

func fixedPortPrefix(port string) string {
	if prefix := port[:1]; prefix != ":" {
		port = fmt.Sprintf(":%s", port)
	}
	return port
}

func isFunc(target interface{}) bool {
	retType := reflect.TypeOf(target)
	return retType.Kind() == reflect.Func
}

// typeExists defined type is exists or not
func typeExists(items interface{}, target interface{}) bool {
	assert1(isIteratee(items), "items must be an iteratee")
	ptype := reflect.ValueOf(target).Type()
	arrValue := reflect.ValueOf(items)
	for i := 0; i < arrValue.Len(); i++ {
		iEle := arrValue.Index(i).Interface()
		iType := reflect.ValueOf(iEle).Type()
		if iType == ptype {
			return true
		}
	}
	return false
}

// retrieve array type
func isIteratee(in interface{}) bool {
	arrType := reflect.TypeOf(in)
	tpKind := arrType.Kind()
	return tpKind == reflect.Array || tpKind == reflect.Slice || tpKind == reflect.Map
}

// make struct from reflect type
func createStruct(sfs []reflect.StructField) interface{} {
	return reflect.New(reflect.StructOf(sfs)).Interface()
}

// get fieldValue by reflect
func stealFieldInStruct(fieldName string, sv interface{}) interface{} {
	svv := indirectValue(reflect.ValueOf(sv))
	return svv.FieldByName(fieldName).Interface()
}

// indirect from ptr
func indirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// indirect from ptr
func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func indirectFunc(item interface{}, funcName string) (interface{}, bool) {
	fromStruct := false
	value := reflect.ValueOf(item)
	if value.Kind() == reflect.Interface && value.Elem().Kind() == reflect.Interface {
		value = value.Elem().Elem()
	}
	if value.Kind() == reflect.Ptr && value.Elem().Kind() == reflect.Struct {
		if value.MethodByName(funcName).IsValid() {
			value = value.MethodByName(funcName)
			fromStruct = true
		} else {
			value = value.Elem()
		}
	}
	if value.Kind() == reflect.Struct {
		value = value.MethodByName(funcName)
		fromStruct = true
	}
	if value.Kind() == reflect.Func && value.IsValid() {
		return value.Interface(), fromStruct
	}
	return nil, fromStruct
}

func indirectPlugin(item interface{}) interface{} {
	value, _ := indirectFunc(item, pluginHookName)
	assert1(value != nil, fmt.Sprintf("%v can not be used as plugin", item))
	return value
}

func indirectPre(item interface{}) interface{} {
	value, fromStruct := indirectFunc(item, preHookName)
	if !fromStruct {
		return value
	}
	return nil
}

func indirectPost(item interface{}) interface{} {
	value, fromStruct := indirectFunc(item, postHookName)
	if !fromStruct {
		return value
	}
	return nil
}

func isPlugin(item interface{}) bool {
	value, _ := indirectFunc(item, pluginHookName)
	return value != nil
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

func assert1(guard bool, err interface{}) {
	if !guard {
		panic(err)
	}
}
