// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
)

type (
	// Callables defined func array
	Callables []PNRet
	// Injects defined bulrush Inject entitys
	Injects  []interface{}
	executor struct {
		callables Callables
		injects   *Injects
	}
)

// concat defined array concat
func (inj *Injects) concat(target *Injects) *Injects {
	injects := append(*inj, *target...)
	return &injects
}

// typeExisted defined inject type is existed or not
func (inj *Injects) typeExisted(item interface{}) bool {
	return typeExists(*inj, item)
}

func (call Callables) toValues() []reflect.Value {
	values := []reflect.Value{}
	for _, ret := range call {
		values = append(values, reflect.ValueOf(ret))
	}
	return values
}

func (exec *executor) execute(inspect func(...interface{})) {
	values := exec.callables.toValues()
	for _, value := range values {
		funcValue := funcValue{value: value}
		funcValue.inputsFrom(*exec.injects)
		ret := funcValue.call().([]interface{})
		inspect(ret...)
	}
}
