/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush reflect]
 */

package bulrush

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/thoas/go-funk"
)

// reflectObjectAndCall
// - you can call a method in object by this method
// - injects contains injectObject
// - ptrDyn `inject params` that be about to be injected
func reflectObjectAndCall(target interface{}, params []interface{}) {
	objType := reflect.TypeOf(target)
	objValue := reflect.ValueOf(target)

	if objValue.Kind() != reflect.Ptr {
		panic("target must be a ptr")
	}
	for i := 0; i < objType.NumMethod(); i++ {
		inputs := make([]reflect.Value, 0)
		funcType := objType.Method(i)
		funcName := funcType.Name
		method := objValue.Method(i)
		numIn := funcType.Type.NumIn()
		if !strings.HasPrefix(funcName, "Inject") {
			continue
		}
		for index := 1; index < numIn; index++ {
			ptype := funcType.Type.In(index)
			r := funk.Find(params, func(x interface{}) bool {
				return ptype == reflect.TypeOf(x)
			})
			if r != nil {
				inputs = append(inputs, reflect.ValueOf(r))
			}
		}
		if method.IsValid() && (numIn == len(inputs)) {
			method.Call(inputs)
		} else {
			panic(fmt.Errorf("Invalid method in reflectObjectAndCall: %s in inject", funcName))
		}
	}
}

// reflectMethodAndCall call method by reflect
func reflectMethodAndCall(target interface{}, params []interface{}) interface{} {
	funcType := reflect.TypeOf(target)
	funcName := funcType.Name()
	funcValue := reflect.ValueOf(target)
	inputs := make([]reflect.Value, 0)
	numIn := funcType.NumIn()

	for index := 0; index < numIn; index++ {
		ptype := funcType.In(index)
		r := funk.Find(params, func(x interface{}) bool {
			return ptype == reflect.TypeOf(x)
		})
		if r != nil {
			inputs = append(inputs, reflect.ValueOf(r))
		}
	}

	if funcValue.IsValid() && (numIn == len(inputs)) {
		rs := funcValue.Call(inputs)
		return funk.Map(funk.Filter(rs, func(v reflect.Value) bool {
			return v.IsValid()
		}), func(v reflect.Value) interface{} {
			return v.Interface()
		})
	}
	panic(fmt.Errorf("Invalid method in reflectMethodAndCall: %s in inject", funcName))
}

// IsIteratee returns if the argument is an iteratee.
func IsIteratee(in interface{}) bool {
	arrType := reflect.TypeOf(in)
	kind := arrType.Kind()
	return kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}

// Find iterates over elements of collection, returning predicate returns truthy for.
func typeExists(arr interface{}, target interface{}) bool {
	if !IsIteratee(arr) {
		panic("First parameter must be an iteratee")
	}
	ptype := reflect.ValueOf(target).Type()

	arrValue := reflect.ValueOf(arr)

	for i := 0; i < arrValue.Len(); i++ {
		iEle := arrValue.Index(i).Interface()
		iType := reflect.ValueOf(iEle).Type()
		if iType == ptype {
			return true
		}
	}
	return false
}

// createSlice create array from target type
func createSlice(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	tSlice := reflect.MakeSlice(reflect.SliceOf(tType), 0, 0).Interface()
	return tSlice
}

// createObject create object from target type
func createObject(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	tObject := reflect.New(tType).Interface()
	return tObject
}
