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

// reflect func in struct and call it with params
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
		reflectCall(method.Interface(), inputs)
	}
}

// reflect func and call it with params
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
		return reflectCall(target, inputs)
	}
	panic(fmt.Errorf("Invalid plugin type %s", target))
}

// relect func type and call with input args
func reflectCall(method interface{}, inputs []reflect.Value) interface{} {
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

// reflectTypeMatcher match type with type tactics or ducker tactics
func reflectTypeMatcher(ptype reflect.Type, params []interface{}) interface{} {
	eleValue := typeMatcher(ptype, params)
	if eleValue == nil && DuckReflect {
		eleValue = duckMatcher(ptype, params)
	}
	if eleValue == nil {
		panic(fmt.Errorf("Invalid param in reflectTypeMatcher: %s", ptype))
	}
	return eleValue
}

// retrieve array type
func isIteratee(in interface{}) bool {
	arrType := reflect.TypeOf(in)
	tpKind := arrType.Kind()
	return tpKind == reflect.Array || tpKind == reflect.Slice || tpKind == reflect.Map
}

// check type exists or not
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

// make slice from reflect type
func createSlice(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	return reflect.MakeSlice(reflect.SliceOf(tType), 0, 0).Interface()
}

// make object from reflect type
func createObject(target interface{}) interface{} {
	tType := reflect.ValueOf(target).Type()
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	return reflect.New(tType).Interface()
}
