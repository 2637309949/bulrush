// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
)

// ReverseInject Inject
type ReverseInject struct {
	config  *Config
	injects *Injects
	inspect func(items ...interface{})
}

// Register function for Reverse Injects
// If the function you're injecting is a black box,
// then you can try this
// Example: github.com/2637309949/bulrush-template/models.go
func (r *ReverseInject) Register(rFunc interface{}) {
	kind := reflect.TypeOf(rFunc).Kind()
	if kind != reflect.Func {
		panic(fmt.Errorf("rFunc should to be func type"))
	}
	pv := NewPluginValue(rFunc)
	pv.inputsFrom(*r.injects)
	pv.runPre()
	r.inspect(pv.runPlugin()...)
	pv.runPost()
}
