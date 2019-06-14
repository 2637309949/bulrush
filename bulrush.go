// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"sync"

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
	// DefaultMaxPlugins is the number of max Plugins
	// `Just for Learning synchronization`
	// a matter of little interest
	DefaultMaxPlugins = 0
	// DuckReflect indicate inject with duck Type, default is true
	DuckReflect = true
)

// Bulrush the framework's struct
// --EventEmmiter emit and on
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type (
	// Middles -
	Middles []PNBase
	// Injects -
	Injects []interface{}
	// Bulrush interface defined
	Bulrush interface {
		On(events.EventName, ...events.Listener)
		Once(events.EventName, ...events.Listener)
		Emit(events.EventName, ...interface{})
		SetMaxPlugins(int)
		GetMaxPlugins() int
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
		maxPlugins  int
		mu          sync.Mutex
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
	status := statusStorage(emmiter)
	bulrush := &rush{
		EventEmmiter: emmiter,
		preMiddles:   &preMiddles,
		middles:      &middles,
		postMiddles:  &postMiddles,
		injects:      &injects,
		maxPlugins:   DefaultMaxPlugins,
	}
	defaultMiddles := Middles{
		HTTPProxy,
		HTTPRouter,
	}
	defaultInjects := Injects{
		emmiter,
		status,
		&ReverseInject{
			injects: bulrush.injects,
		},
	}
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
	if len(items) == 0 {
		return bulrush
	}
	bulrush.mu.Lock()
	defer bulrush.mu.Unlock()
	if bulrush.maxPlugins > 0 && len(*bulrush.preMiddles) == bulrush.maxPlugins {
		rushLogger.Warn(`warning: possible plugins memory 'leak detected. %d plugin added.
			'Use app.SetMaxPlugins(n int) to increase limit.`, len(*bulrush.preMiddles))
		return bulrush
	}
	*bulrush.preMiddles = append(*bulrush.preMiddles, items...)
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) Use(items ...PNBase) Bulrush {
	if len(items) == 0 {
		return bulrush
	}
	bulrush.mu.Lock()
	defer bulrush.mu.Unlock()
	if bulrush.maxPlugins > 0 && len(*bulrush.middles) == bulrush.maxPlugins {
		rushLogger.Warn(`warning: possible plugins memory 'leak detected. %d plugin added.
			'Use app.SetMaxPlugins(n int) to increase limit.`, len(*bulrush.middles))
		return bulrush
	}
	*bulrush.middles = append(*bulrush.middles, items...)
	return bulrush
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) PostUse(items ...PNBase) Bulrush {
	if len(items) == 0 {
		return bulrush
	}
	bulrush.mu.Lock()
	defer bulrush.mu.Unlock()
	if bulrush.maxPlugins > 0 && len(*bulrush.postMiddles) == bulrush.maxPlugins {
		rushLogger.Warn(`warning: possible plugins memory 'leak detected. %d plugin added.
			'Use app.SetMaxPlugins(n int) to increase limit.`, len(*bulrush.postMiddles))
		return bulrush
	}
	*bulrush.postMiddles = append(*bulrush.postMiddles, items...)
	return bulrush
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bulrush *rush) Config(path string) Bulrush {
	bulrush.config = LoadConfig(path)
	bulrush.Inject(bulrush.config)

	DuckReflect = bulrush.config.DuckReflect
	Mode = bulrush.config.Mode

	gin.SetMode(bulrush.config.Mode)
	reloadRushLogger(bulrush.config.Mode)
	return bulrush
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bulrush *rush) Inject(items ...interface{}) Bulrush {
	if len(items) == 0 {
		return bulrush
	}
	injects := funk.Filter(items, func(x interface{}) bool {
		return !typeExists(*bulrush.injects, x)
	}).([]interface{})
	*bulrush.injects = append(*bulrush.injects, injects...)
	return bulrush
}

// SetMaxPlugins obviously this function allows the MaxPlugins
// to be decrease or increase. Set to zero for unlimited
func (bulrush *rush) SetMaxPlugins(n int) {
	if n < 0 {
		rushLogger.Warn("(events) warning: MaxPlugins must be positive number, tried to set: %d", n)
		return
	}
	bulrush.maxPlugins = n
}

// SetMaxPlugins obviously this function allows the MaxPlugins
// to be decrease or increase. Set to zero for unlimited
func SetMaxPlugins(n int) {
	defaultApp.SetMaxPlugins(n)
}

func (bulrush *rush) GetMaxPlugins() int {
	return bulrush.maxPlugins
}

// GetMaxPlugins returns the max Plugins for this bulrush
// see SetMaxPlugins
func GetMaxPlugins() int {
	return defaultApp.GetMaxPlugins()
}

// return middles contain middles, preMiddles and postMiddles
func (bulrush *rush) allMiddles() *Middles {
	middles := append(append(*bulrush.preMiddles, *bulrush.middles...), *bulrush.postMiddles...)
	return &middles
}

// return plugins that contained in every middles
func (bulrush *rush) middle2Plugins(middles *Middles) interface{} {
	plugins := funk.Map(*middles, func(x PNBase) PNRet {
		return x.Plugin()
	})
	return plugins
}

func (bulrush *rush) execPlugins(plugins interface{}) {
	funk.ForEach(plugins, func(x interface{}) {
		rs := reflectMethodAndCall(x, *bulrush.injects)
		bulrush.Inject(rs.([]interface{})...)
	})
}

// Exec middles, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (bulrush *rush) execMiddles() Bulrush {
	middles := bulrush.allMiddles()
	plugins := bulrush.middle2Plugins(middles)
	bulrush.execPlugins(plugins)
	return bulrush
}

// Run application with callback, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (bulrush *rush) Run(cb interface{}) {
	bulrush.PostUse(PNQuick(cb))
	bulrush.execMiddles()
}

// RunImmediately, excute plugin in orderly
// Quick start application
func (bulrush *rush) RunImmediately() {
	bulrush.PostUse(RunImmediately)
	bulrush.execMiddles()
}
