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
	"strings"
	"reflect"
	"github.com/thoas/go-funk"
)

// dynamicObjectsCall
// - you can call a method in object by this method
// - injects contains injectObject
// - ptrDyn `inject params` that be about to be injected
func dynamicObjectsCall(injects []interface{}, ptrDyn[]interface{}) {
	funk.ForEach(injects, func(x interface{}) {
		dynamicObjectCall(x, ptrDyn)
	})
}

// dynamicObjectCall
// - you can call a method in object by this method
// - injects contains injectObject
// - ptrDyn `inject params` that be about to be injected
func dynamicObjectCall(target interface{}, params[]interface{}) {
	getType  := reflect.TypeOf(target)
	getValue := reflect.ValueOf(target)

	if getValue.Kind() != reflect.Ptr {
		panic("target must be a ptr")
	}
	for i := 0; i < getType.NumMethod(); i++ {
		valid	   := true
		inputs 	   := make([]reflect.Value, 0)
		methodType := getType.Method(i)
		methodName := methodType.Name
		method 	   := getValue.Method(i)
		numIn	   := methodType.Type.NumIn()
		if !strings.HasPrefix(methodName, "Inject") {
			continue
		}
		for index := 1; index < numIn; index ++ {
			ptype := methodType.Type.In(index)
			r := funk.Find(params, func(x interface{}) bool {
				return ptype == reflect.TypeOf(x)
			})
			if r != nil {
				inputs = append(inputs, reflect.ValueOf(r))
			} else {
				valid = false
				break
			}
		}
		if method.IsValid() && valid {
			method.Call(inputs)
		} else {
			panic(fmt.Errorf("Invalid method: %s in inject", methodName))
		}
	}
}

// dynamicMethodCall
// call method by reflect
func dynamicMethodCall(target interface{}, params[]interface{}) interface {} {
	valid 	   := true
	getType    := reflect.TypeOf(target)
	methodName := getType.Name()
	getValue   := reflect.ValueOf(target)
	inputs 	   := make([]reflect.Value, 0)
	numIn	   := getType.NumIn()
	for index := 0; index < numIn; index ++ {
		ptype := getType.In(index)
		r := funk.Find(params, func(x interface{}) bool {
			return ptype == reflect.TypeOf(x)
		})
		if r != nil {
			inputs = append(inputs, reflect.ValueOf(r))
		} else {
			valid = false
			break
		}
	}
	if getValue.IsValid() && valid {
		rs := getValue.Call(inputs)
		return funk.Map(funk.Filter(rs, func(v reflect.Value) bool {
			return v.IsValid()
		}), func(v reflect.Value) interface {}{
			return v.Interface()
		})
	}
	panic(fmt.Errorf("invalid method: %s in inject", methodName))
}

// dynamicMethodsCall
// call method by reflect
func dynamicMethodsCall(plugins []interface{}, params *[]interface{}, cb func(interface{})) {
	funk.ForEach(plugins, func(x interface{}) {
		cb(dynamicMethodCall(x, *params))
	})
}

// typeExists -
func typeExists(injects []interface{}, target interface{}) bool {
	ptype  := reflect.TypeOf(target)
	r := funk.Find(injects, func(x interface{}) bool {
		return ptype == reflect.TypeOf(x)
	})
	if r != nil {
		return true
	}
	return false
}

// createSlice -
func createSlice(target interface{}) interface{} {
	tagetType 	:= reflect.TypeOf(target)
	if tagetType.Kind() == reflect.Ptr {
		tagetType = tagetType.Elem()
	}
	targetSlice := reflect.MakeSlice(reflect.SliceOf(tagetType), 0, 0).Interface()
	return targetSlice
}

// createObject -
func createObject(target interface{}) interface{} {
	tagetType 	 := reflect.TypeOf(target)
	if tagetType.Kind() == reflect.Ptr {
		tagetType = tagetType.Elem()
	}
	targetObject := reflect.New(tagetType).Interface()
	return targetObject
}
