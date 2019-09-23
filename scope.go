// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
)

type (
	// ScopeBase defined scope base type
	ScopeBase struct {
		v       interface{}
		acquire func(reflect.Type) interface{}
	}
	// ConstructorScope defined plugin scope
	ConstructorScope struct {
		ScopeBase
	}
	// ObjectScope defined plugin scope
	ObjectScope struct {
		ScopeBase
	}
	// Arguments defined Scope arguments
	Arguments interface {
		arguments(funk reflect.Type) (args []reflect.Value)
	}
	// ScopeHook defined hook in plugin
	ScopeHook interface {
		Pre() error
		Post() error
		Plugin() []interface{}
	}
	// Scope defined plugin hooks
	Scope interface {
		Arguments
		ScopeHook
		Type() reflect.Type
	}
)

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

// arguments defined obtain arguments before exec plugin
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

func newConstructorScope(src interface{}, acquire func(reflect.Type) interface{}) Scope {
	return &ConstructorScope{
		ScopeBase{
			v:       src,
			acquire: acquire,
		},
	}
}

func newObjectScope(src interface{}, acquire func(reflect.Type) interface{}) Scope {
	return &ObjectScope{
		ScopeBase{
			v:       src,
			acquire: acquire,
		},
	}
}

func newScope(src interface{}, acquire func(reflect.Type) interface{}) (s Scope) {
	vType := reflect.TypeOf(src)
	if vType.Kind() == reflect.Interface {
		vType = vType.Elem()
	}
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	switch true {
	case vType.Kind() == reflect.Func:
		s = newConstructorScope(src, acquire)
	case vType.Kind() == reflect.Struct:
		s = newObjectScope(src, acquire)
	}
	return
}
