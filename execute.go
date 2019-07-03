// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import "reflect"

type (
	// Callables defined func array
	callables []interface{}
	executor  struct {
		callables *callables
		injects   *injects
	}
)

func (call callables) toValues() []reflect.Value {
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
		funcValue.runPre()
		ret := funcValue.runPlugin().([]interface{})
		funcValue.runPost()
		inspect(ret...)
	}
}
