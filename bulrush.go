// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kataras/go-events"
	"github.com/thoas/go-funk"
)

var (
	// Version current version number
	Version = "0.0.1"
	// DefaultMode default gin mode
	DefaultMode = "debug"
	// Mode bulrush running Mode
	Mode = "debug"
)

// PNRet return a plugin after call Plugin func
type PNRet interface{}

// PNBase defined interface for bulrush Plugin
type PNBase interface {
	Plugin() PNRet
}

// PNStruct for a quickly Plugin SetUp when you dont want declare PNBase
// PBBase minimize implement
type PNStruct struct{ ret PNRet }

// Plugin for PNQuick
func (pns *PNStruct) Plugin() PNRet {
	return pns.ret
}

// Middles defined array of PNBase
type Middles []PNBase

// concat defined array concat
func (mi *Middles) concat(middles *Middles) *Middles {
	newMiddles := append(*mi, *middles...)
	return &newMiddles
}

// toRet defined to get `ret` that plugin func return
func (mi *Middles) toRet() []PNRet {
	return funk.Map(*mi, func(x PNBase) PNRet {
		return x.Plugin()
	}).([]PNRet)
}

// Bulrush the framework's struct
// --EventEmmiter emit and on
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type (
	// Injects -
	Injects []interface{}
	// Bulrush interface defined
	Bulrush interface {
		On(events.EventName, ...events.Listener)
		Once(events.EventName, ...events.Listener)
		Emit(events.EventName, ...interface{})
		PreUse(...PNBase) Bulrush
		Use(...PNBase) Bulrush
		PostUse(...PNBase) Bulrush
		Config(string) Bulrush
		Inject(...interface{}) Bulrush
		RunImmediately()
		Run(interface{})
	}
	// Bulrush is the framework's instance, it contains the muxer,
	// middleware and configuration settings.
	// Create an instance of Bulrush, by using New() or Default()
	rush struct {
		events.EventEmmiter
		config      *Config
		preMiddles  *Middles
		middles     *Middles
		postMiddles *Middles
		injects     *Injects
	}
)

// New returns a new blank Bulrush instance without any middleware attached.
// By default the configuration is:
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() Bulrush {
	preMiddles := make(Middles, 0)
	middles := make(Middles, 0)
	postMiddles := make(Middles, 0)
	injects := make(Injects, 0)
	emmiter := events.New()
	bulrush := &rush{
		EventEmmiter: emmiter,
		preMiddles:   &preMiddles,
		middles:      &middles,
		postMiddles:  &postMiddles,
		injects:      &injects,
		config:       &Config{},
	}
	defaultMiddles := Middles{
		HTTPProxy,
		HTTPRouter,
	}
	defaultInjects := defaultInjects(bulrush)
	bulrush.Use(defaultMiddles...)
	bulrush.Inject(defaultInjects...)
	return bulrush
}

// Default returns an Bulrush instance with the Override and Recovery middleware already attached.
// --Recovery middle has been register in httpProxy and user router
// --Override middles has been register in router for override req
func Default() Bulrush {
	bulrush := defaultApp
	defaultMiddles := Middles{
		Recovery,
		Override,
	}
	bulrush.Use(defaultMiddles...)
	return bulrush
}

// Silence the compiler
var _ = &rush{}

// defaultApp default rush
var defaultApp = New()

// PreUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) PreUse(items ...PNBase) Bulrush {
	*bulrush.preMiddles = append(*bulrush.preMiddles, items...)
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) Use(items ...PNBase) Bulrush {
	*bulrush.middles = append(*bulrush.middles, items...)
	return bulrush
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) PostUse(items ...PNBase) Bulrush {
	*bulrush.postMiddles = append(*bulrush.postMiddles, items...)
	return bulrush
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bulrush *rush) Config(path string) Bulrush {
	*bulrush.config = *LoadConfig(path)
	bulrush.Inject(bulrush.config)
	Mode = bulrush.config.Mode
	gin.SetMode(bulrush.config.Mode)
	reloadRushLogger(bulrush.config.Mode)
	return bulrush
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bulrush *rush) Inject(items ...interface{}) Bulrush {
	injects := funk.Filter(items, func(x interface{}) bool {
		return !typeExists(*bulrush.injects, x)
	}).([]interface{})
	*bulrush.injects = append(*bulrush.injects, injects...)
	return bulrush
}

// Run application with callback, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (bulrush *rush) Run(cb interface{}) {
	bulrush.PostUse(PNQuick(cb))
	middles := bulrush.preMiddles.concat(bulrush.middles).concat(bulrush.postMiddles)
	rets := middles.toRet()
	funk.ForEach(rets, func(ret PNRet) {
		if isFunc(ret) {
			injects := reflectMethodAndCall(ret, *bulrush.injects, struct{ DuckReflect bool }{bulrush.config.DuckReflect})
			bulrush.Inject(injects...)
		} else {
			panic(fmt.Errorf("ret %v is not a func", ret))
		}
	})
}

// RunImmediately, excute plugin in orderly
// Quick start application
func (bulrush *rush) RunImmediately() {
	bulrush.Run(RunImmediately.Plugin())
}
