// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"

	"github.com/thoas/go-funk"
)

const (
	preHookName    = "Pre"
	postHookName   = "Post"
	pluginHookName = "Plugin"
)

type (
	// Plugins defined those that can be call by reflect
	// , Plugins passby func or a struct that has `Plugin` func
	Plugins []interface{}
	// PluginValue defined plugin value with pre an post
	PluginValue struct {
		Pre    reflect.Value
		Post   reflect.Value
		Plugin reflect.Value
		Inputs []reflect.Value
	}
)

// Append defined array concat
func (p *Plugins) Append(target *Plugins) *Plugins {
	middles := append(*p, *target...)
	return &middles
}

// toCallables defined to get `ret` that plugin func return
func (p *Plugins) toPluginValues() *[]*PluginValue {
	pluginValus := funk.Map(*p, func(plugin interface{}) *PluginValue {
		return NewPluginValue(plugin)
	}).([]*PluginValue)
	return &pluginValus
}

// NewPluginValue defined pluginValue
func NewPluginValue(src interface{}) *PluginValue {
	pv := PluginValue{}
	// Pre hook
	if pre, fromStruct := indirectFunc(src, preHookName); pre != nil && fromStruct {
		value := reflect.ValueOf(pre)
		if value.IsValid() && value.Type().NumIn() == 0 {
			pv.Pre = value
		}
	}
	// plugin hook
	if plugin := indirectPlugin(src); plugin != nil {
		value := reflect.ValueOf(plugin)
		if value.IsValid() {
			pv.Plugin = value
		}
	}
	// post hook
	if post, fromStruct := indirectFunc(src, postHookName); post != nil && fromStruct {
		value := reflect.ValueOf(post)
		if value.IsValid() && value.Type().NumIn() == 0 {
			pv.Post = value
		}
	}
	return &pv
}

// runPost defined run post hook in plugin
// , remind that paramters of hook func should be zero
func (pv *PluginValue) runPost() {
	if pv.Post.IsValid() {
		pv.Post.Call([]reflect.Value{})
	}
}

// runPre defined run pre hook in plugin
// , remind that paramters of hook func should be zero
func (pv *PluginValue) runPre() {
	if pv.Pre.IsValid() {
		pv.Pre.Call([]reflect.Value{})
	}
}

// runPlugin defined run plugin hook in plugin
// , and return injects
func (pv *PluginValue) runPlugin() []interface{} {
	ret := pv.Plugin.Call(pv.Inputs)
	return funk.Map(ret, func(v reflect.Value) interface{} {
		return v.Interface()
	}).([]interface{})
}

// inputsFrom defined plugins paramters by type
// , or by interface{} implement
func (pv *PluginValue) inputsFrom(inputs []interface{}) {
	funcItem := indirectPlugin(pv.Plugin.Interface())
	funcValue := reflect.ValueOf(funcItem)
	if funcValue.Type().Kind() != reflect.Func {
		panic(fmt.Errorf(" %v inputsFrom call with %v error", funcItem, inputs))
	}
	funcType := funcValue.Type()
	numIn := funcType.NumIn()
	for index := 0; index < numIn; index++ {
		ptype := funcType.In(index)
		eleValue := typeMatcher(ptype, inputs)
		if eleValue == nil {
			eleValue = duckMatcher(ptype, inputs)
		}
		if eleValue == nil {
			panic(fmt.Errorf("inputsFrom %v call with %v error", ptype, reflect.TypeOf(inputs[7])))
		}
		pv.Inputs = append(pv.Inputs, eleValue.(reflect.Value))
	}
}
