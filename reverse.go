// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
)

// ReverseInject defined a inject
// , for reverse inject
type ReverseInject struct {
	injects *Injects
	inspect func(items ...interface{})
}

// Register defiend function for Reverse Injects
// Example:
// func Route(router *gin.RouterGroup, event events.EventEmmiter, ri *bulrush.ReverseInject) {
// 		ri.Register(RegisterMgo)
// 		ri.Register(RegisterCache)
// 		ri.Register(RegisterSeq)
// 		ri.Register(RegisterMq)
// 		ri.Register(RegisterEvent)
// 		ri.Register(RegisterMock)
// 		ri.Register(RegisterGRPC)
// 		event.Emit("hello", "this is my payload to hello router")
// }
func (r *ReverseInject) Register(rFunc interface{}) {
	kind := reflect.TypeOf(rFunc).Kind()
	if kind != reflect.Func {
		panic(fmt.Errorf("rFunc should to be func type"))
	}
	scopes := (&Plugins{rFunc}).toScopes(func(t reflect.Type) interface{} {
		return r.injects.Acquire(t)
	})
	exec := &engine{
		scopes: scopes,
	}
	exec.exec(r.inspect)
}
