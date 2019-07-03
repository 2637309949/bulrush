// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"math/rand"
	"reflect"
)

// Some get or a default value
func Some(t interface{}, i interface{}) interface{} {
	if t != nil && t != "" && t != 0 {
		return t
	}
	return i
}

// Find element in arrs with some match condition
func Find(arrs []interface{}, matcher func(interface{}) bool) interface{} {
	var target interface{}
	for _, item := range arrs {
		match := matcher(item)
		if match {
			target = item
			break
		}
	}
	return target
}

// ToStrArray -
func ToStrArray(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}
	return s
}

// ToIntArray -
func ToIntArray(t []interface{}) []int {
	s := make([]int, len(t))
	for i, v := range t {
		s[i] = v.(int)
	}
	return s
}

// RandString gen string
func RandString(n int) string {
	const seeds = "abcdefghijklmnopqrstuvwxyz1234567890"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = seeds[rand.Intn(len(seeds))]
	}
	return string(bytes)
}

// LeftV -
func LeftV(left interface{}, right interface{}) interface{} {
	return left
}

// RightV -
func RightV(left interface{}, right interface{}) interface{} {
	return right
}

// LeftOkV -
func LeftOkV(left interface{}, right ...bool) interface{} {
	var r = true
	if len(right) != 0 {
		r = right[0]
	} else if left == "" || left == nil || left == 0 {
		r = false
	}
	if r {
		return left
	}
	return nil
}

// LeftSV left value or panic
func LeftSV(left interface{}, right error) interface{} {
	if right != nil {
		panic(right)
	}
	return left
}

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

func typeExists(items interface{}, target interface{}) bool {
	if !isIteratee(items) {
		panic("items must be an iteratee")
	}
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

// make slice from reflect type
func createSlice(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	tType = indirectType(tType)
	return reflect.New(reflect.SliceOf(tType)).Interface()
}

// make object from reflect type
func createObject(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	tType = indirectType(tType)
	return reflect.New(tType).Interface()
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
	if value == nil {
		panic(fmt.Errorf("%v can not be used as plugin", item))
	}
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
