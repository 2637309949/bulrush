package bulrush

import (
	"fmt"
	"strings"
	"reflect"
	"github.com/thoas/go-funk"
)

// dynamicObjectsCall
// call method by reflect
func dynamicObjectsCall(injects []interface{}, ptrDyn[]interface{}) {
	for _, target := range injects {
		dynamicObjectCall(target, ptrDyn)
	}
}

// dynamicObjectsCall
// call method by reflect
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
		if len(rs) > 0 && rs[0].IsValid() {
			return rs[0].Interface()
		}
	} else {
		panic(fmt.Errorf("invalid method: %s in inject", methodName))
	}
	return nil
}

// dynamicMethodsCall
// call method by reflect
func dynamicMethodsCall(injects []interface{}, params[]interface{}) {
	for _, target := range injects {
		dynamicMethodCall(target, params)
	}
}

// typeExists
// check type if exists or not
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

// createSlice
// create slice by reflect
func createSlice(target interface{}) interface{} {
	tagetType 	:= reflect.TypeOf(target)
	if tagetType.Kind() == reflect.Ptr {
		tagetType = tagetType.Elem()
	}
	targetSlice := reflect.MakeSlice(reflect.SliceOf(tagetType), 0, 0).Interface()
	return targetSlice
}

// createObject
// create object by reflect
func createObject(target interface{}) interface{} {
	tagetType 	 := reflect.TypeOf(target)
	if tagetType.Kind() == reflect.Ptr {
		tagetType = tagetType.Elem()
	}
	targetObject := reflect.New(tagetType).Interface()
	return targetObject
}
