// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
)

// ConstructorScope defined plugin scope
type ConstructorScope struct {
	v       interface{}
	acquire func(reflect.Type) interface{}
}

// Type defined type info
func (s *ConstructorScope) Type() reflect.Type {
	return reflect.TypeOf(s.v)
}

// Pre defined pre hook
func (s *ConstructorScope) Pre() error {
	return nil
}

// Post defined pre hook
func (s *ConstructorScope) Post() error {
	return nil
}

// Plugin defined pre hook
func (s *ConstructorScope) Plugin() (rets []interface{}) {
	funk := reflect.ValueOf(s.v)
	if funk.IsValid() {
		args := s.arguments(funk.Type())
		ret := funk.Call(args)
		for _, v := range ret {
			rets = append(rets, v.Interface())
		}
	}
	return
}

func (s *ConstructorScope) arguments(funk reflect.Type) (args []reflect.Value) {
	numIn := funk.NumIn()
	for index := 0; index < numIn; index++ {
		ptype := funk.In(index)
		v := s.acquire(ptype)
		if v == nil {
			panic(fmt.Errorf("invalid scope arguments type %v", ptype))
		}
		args = append(args, reflect.ValueOf(v))
	}
	return
}
