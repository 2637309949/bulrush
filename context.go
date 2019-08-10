// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"errors"
	"fmt"
	"reflect"
	"time"

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
	// PluginContext defined plugin value with pre an post
	PluginContext struct {
		Pre    reflect.Value
		Post   reflect.Value
		Plugin reflect.Value
		Inputs []reflect.Value
	}
	// HTTPContext defined httpContxt
	HTTPContext struct {
		Exit         chan struct{}
		Chan         chan struct{}
		DeadLineTime time.Time
	}
)

// Append defined array concat
func (p *Plugins) Append(target *Plugins) *Plugins {
	middles := append(*p, *target...)
	return &middles
}

// Put defined array Put
func (p *Plugins) Put(target interface{}) *Plugins {
	*p = append(*p, target)
	return p
}

// Size defined Plugins Size
func (p *Plugins) Size() int {
	return len(*p)
}

// Get defined index of Plugins
func (p *Plugins) Get(pos int) interface{} {
	return (*p)[pos]
}

// Swap swaps the two values at the specified positions.
func (p *Plugins) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

// toCallables defined to get `ret` that plugin func return
func (p *Plugins) toPluginContexts() *[]PluginContext {
	pluginValus := funk.Map(*p, func(plugin interface{}) PluginContext {
		return *parsePlugin(plugin)
	}).([]PluginContext)
	return &pluginValus
}

// parsePlugin defined parsePlugin
func parsePlugin(src interface{}) *PluginContext {
	pv := PluginContext{}
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
func (pv *PluginContext) runPost() {
	if pv.Post.IsValid() {
		pv.Post.Call([]reflect.Value{})
	}
}

// runPre defined run pre hook in plugin
// , remind that paramters of hook func should be zero
func (pv *PluginContext) runPre() {
	if pv.Pre.IsValid() {
		pv.Pre.Call([]reflect.Value{})
	}
}

// runPlugin defined run plugin hook in plugin
// , and return injects
func (pv *PluginContext) runPlugin() []interface{} {
	ret := pv.Plugin.Call(pv.Inputs)
	return funk.Map(ret, func(v reflect.Value) interface{} {
		return v.Interface()
	}).([]interface{})
}

// inputsFrom defined plugins paramters by type
// , or by interface{} implement
func (pv *PluginContext) inputsFrom(inputs []interface{}) {
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
			if eleValue = duckMatcher(ptype, inputs); eleValue == nil {
				panic(fmt.Errorf("inputsFrom %v call with %v error", ptype, reflect.TypeOf(inputs[7])))
			}
		}
		pv.Inputs = append(pv.Inputs, eleValue.(reflect.Value))
	}
}

// NewHTTPContext defined httpContext
func NewHTTPContext(duration time.Duration) *HTTPContext {
	cxt := &HTTPContext{}
	cxt.DeadLineTime = time.Now().Add(duration)
	cxt.Chan = make(chan struct{}, 1)
	cxt.Exit = make(chan struct{}, 1)
	return cxt
}

// Done defined http done action
func (ctx *HTTPContext) Done() <-chan struct{} {
	if time.Now().After(ctx.DeadLineTime) {
		ctx.Chan <- struct{}{}
	}
	return ctx.Chan
}

// Err defined http action error
func (ctx *HTTPContext) Err() error {
	return errors.New("can't exit before Specified time")
}

// Value nothing
func (ctx *HTTPContext) Value(key interface{}) interface{} {
	return nil
}

// Deadline defined Deadline time
func (ctx *HTTPContext) Deadline() (time.Time, bool) {
	return ctx.DeadLineTime, true
}
