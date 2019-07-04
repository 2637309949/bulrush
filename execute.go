// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"

	"github.com/thoas/go-funk"
)

type (
	callables []interface{}
	executor  struct {
		callables *callables
		injects   *injects
	}
)

func (call callables) toValues() []reflect.Value {
	values := []reflect.Value{}
	funk.ForEach(call, func(item interface{}) {
		values = append(values, reflect.ValueOf(item))
	})
	return values
}

func (exec *executor) execute(inspect func(...interface{})) {
	values := exec.callables.toValues()
	for _, value := range values {
		fv := parseValue(value)
		fv.inputsFrom(*exec.injects)
		fv.runPre()
		inspect(fv.runPlugin()...)
		fv.runPost()
	}
}
