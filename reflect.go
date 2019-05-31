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

	"github.com/thoas/go-funk"
)

// DuckReflect indicate inject with duck Type, default is true
var DuckReflect = true

func reflectObjectAndCall(target interface{}, params []interface{}) {
	objType := reflect.TypeOf(target)
	objValue := reflect.ValueOf(target)
	for i := 0; i < objType.NumMethod(); i++ {
		inputs := make([]reflect.Value, 0)
		funcType := objType.Method(i)
		method := objValue.Method(i)
		numIn := funcType.Type.NumIn()
		for index := 1; index < numIn; index++ {
			ptype := funcType.Type.In(index)
			eleValue := reflectTypeMatcher(ptype, params)
			inputs = append(inputs, eleValue.(reflect.Value))
		}
		methodCall(method.Interface(), inputs)
	}
}

func reflectMethodAndCall(target interface{}, params []interface{}) interface{} {
	if reflect.Func == reflect.TypeOf(target).Kind() {
		funcType := reflect.TypeOf(target)
		numIn := funcType.NumIn()
		inputs := make([]reflect.Value, 0)
		for index := 0; index < numIn; index++ {
			ptype := funcType.In(index)
			eleValue := reflectTypeMatcher(ptype, params)
			inputs = append(inputs, eleValue.(reflect.Value))
		}
		return methodCall(target, inputs)
	}
	panic(fmt.Errorf("Invalid plugin type %s", target))
}

func methodCall(method interface{}, inputs []reflect.Value) interface{} {
	funcType := reflect.TypeOf(method)
	funcName := funcType.Name()
	funcValue := reflect.ValueOf(method)
	numIn := funcType.NumIn()
	if funcValue.IsValid() && (numIn == len(inputs)) {
		rs := funcValue.Call(inputs)
		return funk.Map(funk.Filter(rs, func(v reflect.Value) bool {
			return v.IsValid()
		}), func(v reflect.Value) interface{} {
			return v.Interface()
		})
	}
	panic(fmt.Errorf("Invalid method %s", funcName))
}

// duckMatcher match type if from target`type
func typeMatcher(ptype reflect.Type, params []interface{}) interface{} {
	target := funk.Find(params, func(x interface{}) bool {
		return ptype == reflect.TypeOf(x)
	})
	if target != nil {
		return reflect.ValueOf(target)
	}
	return nil
}

// duckMatcher match type if implements target`interface
func duckMatcher(ptype reflect.Type, params []interface{}) interface{} {
	target := funk.Find(params, func(x interface{}) bool {
		if ptype.Kind() == reflect.Interface {
			return reflect.TypeOf(x).Implements(ptype)
		}
		return false
	})
	if target != nil {
		return reflect.ValueOf(target)
	}
	return nil
}

// reflectTypeMatcher match type with type tactics or ducker tactics
func reflectTypeMatcher(ptype reflect.Type, params []interface{}) interface{} {
	eleValue := typeMatcher(ptype, params)
	if eleValue == nil {
		if DuckReflect {
			if eleValue = duckMatcher(ptype, params); eleValue == nil {
				panic(fmt.Errorf("Invalid param in reflectTypeMatcher: %s", ptype))
			}
		} else {
			panic(fmt.Errorf("Invalid param in reflectTypeMatcher: %s", ptype))
		}
	}
	return eleValue
}

func isIteratee(in interface{}) bool {
	arrType := reflect.TypeOf(in)
	kind := arrType.Kind()
	return kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map
}

func typeExists(arr interface{}, target interface{}) bool {
	if !isIteratee(arr) {
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

func createSlice(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	tSlice := reflect.MakeSlice(reflect.SliceOf(tType), 0, 0).Interface()
	return tSlice
}

func createObject(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	tObject := reflect.New(tType).Interface()
	return tObject
}
