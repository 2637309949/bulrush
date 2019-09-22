// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/thoas/go-funk"
)

func fixedPortPrefix(port string, plus ...int) string {
	port = strings.ReplaceAll(port, ":", "")
	number, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	if len(plus) > 0 {
		number = plus[0] + number
	}
	return fmt.Sprintf(":%v", number)
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

func isPlugin(src interface{}) (is bool) {
	vType := reflect.TypeOf(src)
	if vType.Kind() == reflect.Ptr && vType.Elem().Kind() == reflect.Func {
		is = true
	} else if vType.Kind() == reflect.Ptr && vType.Elem().Kind() == reflect.Struct {
		_, e := vType.MethodByName("Plugin")
		if e {
			is = true
		} else {
			_, e = vType.Elem().MethodByName("Plugin")
			if e {
				is = true
			}
		}
	} else if vType.Kind() == reflect.Func {
		is = true
	} else if vType.Kind() == reflect.Struct {
		_, e := vType.MethodByName("Plugin")
		if e {
			is = true
		}
	}
	return
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

// resolveAddress defined ipaddress
func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}
