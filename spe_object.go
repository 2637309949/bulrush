// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
)

// ObjectScope defined plugin scope
type ObjectScope struct {
	v       interface{}
	acquire func(reflect.Type) interface{}
}

// Type defined type info
func (s *ObjectScope) Type() reflect.Type {
	return reflect.TypeOf(s.v)
}

// Pre defined pre hook
func (s *ObjectScope) Pre() error {
	funk := reflect.ValueOf(s.v)
	v := funk.MethodByName("Pre")
	if !v.IsValid() && reflect.TypeOf(s.v).Kind() == reflect.Ptr {
		v = funk.Elem().MethodByName("Pre")
	}
	if v.IsValid() {
		args := s.arguments(v.Type())
		v.Call(args)
	}
	return nil
}

// Post defined pre hook
func (s *ObjectScope) Post() error {
	funk := reflect.ValueOf(s.v)
	v := funk.MethodByName("Post")
	if !v.IsValid() && reflect.TypeOf(s.v).Kind() == reflect.Ptr {
		v = funk.Elem().MethodByName("Post")
	}
	if v.IsValid() {
		args := s.arguments(v.Type())
		v.Call(args)
	}
	return nil
}

// Plugin defined pre hook
func (s *ObjectScope) Plugin() (rets []interface{}) {
	funk := reflect.ValueOf(s.v)
	v := funk.MethodByName("Plugin")
	if !v.IsValid() && reflect.TypeOf(s.v).Kind() == reflect.Ptr {
		v = funk.Elem().MethodByName("Plugin")
	}
	if v.IsValid() {
		args := s.arguments(v.Type())
		ret := v.Call(args)
		for _, v := range ret {
			rets = append(rets, v.Interface())
		}
	}
	return
}

func (s *ObjectScope) arguments(funk reflect.Type) (args []reflect.Value) {
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
