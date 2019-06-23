// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

type funcValue struct {
	value  reflect.Value
	inputs []reflect.Value
}

func (fv *funcValue) call() interface{} {
	fvType := fv.value.Type()
	numFvIn := fvType.NumIn()
	numPutIn := len(fv.inputs)
	if fv.value.IsValid() && (numFvIn == numPutIn) {
		ret := fv.value.Call(fv.inputs)
		ret = funk.Filter(ret, func(v reflect.Value) bool {
			return v.IsValid()
		}).([]reflect.Value)
		results := funk.Map(ret, func(v reflect.Value) interface{} {
			return v.Interface()
		})
		return results
	}
	panic(fmt.Sprintf("funcValue %v call with %v error", fv.value, fv.inputs))
}

func (fv *funcValue) inputsFrom(inputs []interface{}) {
	if fv.value.Type().Kind() != reflect.Func {
		panic(fmt.Errorf("inputsFrom %v call with %v error", fv.value.Type().Kind() == reflect.Func, inputs))
	}
	funcType := fv.value.Type()
	numIn := funcType.NumIn()
	for index := 0; index < numIn; index++ {
		ptype := funcType.In(index)
		eleValue := typeMatcher(ptype, inputs)
		if eleValue == nil {
			eleValue = duckMatcher(ptype, inputs)
		}
		if eleValue == nil {
			panic(fmt.Errorf("inputsFrom %v call with %v error", fv.value, inputs))
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
