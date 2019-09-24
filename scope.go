// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
)

type (
	// baseScope defined a struct that constructorScope or objectScope must
	// be provided
	bs struct {
		v       interface{}
		acquire func(reflect.Type) interface{}
	}
	// funcScope defined a construct scope like
	// func(r *gin.RouterGroup) {
	// }
	funcScope struct {
		bs
	}
	// objectScope defined scope plugin from object
	// must implement plugin function
	// plugin() []interface{}
	objectScope struct {
		bs
	}
	// args interface defined arguments provided
	args interface {
		arguments(funk reflect.Type) (args []reflect.Value)
	}
	// hook defined plugin hook
	// pre
	// post
	// plugin
	hook interface {
		pre() error
		post() error
		plugin() []interface{}
	}
	// scope defined plugin hooks
	scope interface {
		args
		hook
		Type() reflect.Type
	}
)

// Type defined type info
func (s *funcScope) Type() reflect.Type {
	return reflect.TypeOf(s.v)
}

// pre defined pre hook
func (s *funcScope) pre() error {
	return nil
}

// post defined pre hook
func (s *funcScope) post() error {
	return nil
}

// plugin defined pre hook
func (s *funcScope) plugin() (rets []interface{}) {
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
func (s *funcScope) arguments(funk reflect.Type) (args []reflect.Value) {
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
func (s *objectScope) Type() reflect.Type {
	return reflect.TypeOf(s.v)
}

// pre defined pre hook
func (s *objectScope) pre() error {
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

// post defined post hook
func (s *objectScope) post() error {
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

// plugin defined plugin hook
func (s *objectScope) plugin() (rets []interface{}) {
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

// arguments defined arguments hook
func (s *objectScope) arguments(funk reflect.Type) (args []reflect.Value) {
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

func newFuncScope(src interface{}, acquire func(reflect.Type) interface{}) scope {
	return &funcScope{
		bs{
			v:       src,
			acquire: acquire,
		},
	}
}

func newObjectScope(src interface{}, acquire func(reflect.Type) interface{}) scope {
	return &objectScope{
		bs{
			v:       src,
			acquire: acquire,
		},
	}
}

func determineScope(src interface{}, acquire func(reflect.Type) interface{}) (s scope) {
	vType := reflect.TypeOf(src)
	if vType.Kind() == reflect.Interface {
		vType = vType.Elem()
	}
	if vType.Kind() == reflect.Ptr {
		vType = vType.Elem()
	}
	switch true {
	case vType.Kind() == reflect.Func:
		s = newFuncScope(src, acquire)
	case vType.Kind() == reflect.Struct:
		s = newObjectScope(src, acquire)
	}
	return
}
