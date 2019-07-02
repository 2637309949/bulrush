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
	// DefaultMode default gin mode
	DefaultMode = "debug"
	// defaultBulrush default bulrush
	defaultBulrush = New()
)

type (
	// Middles defined array of PNBase
	Middles []interface{}
	// Injects defined bulrush Inject entitys
	Injects []interface{}
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
	// Bulrush interface defined
	Bulrush interface {
		On(events.EventName, ...events.Listener)
		Once(events.EventName, ...events.Listener)
		Emit(events.EventName, ...interface{})
		PreUse(...interface{}) Bulrush
		Use(...interface{}) Bulrush
		PostUse(...interface{}) Bulrush
		Config(string) Bulrush
		Inject(...interface{}) Bulrush
		RunImmediately()
		Run(interface{})
	}
)

// concat defined array concat
func (inj *Injects) concat(target *Injects) *Injects {
	injects := append(*inj, *target...)
	return &injects
}

// typeExisted defined inject type is existed or not
func (inj *Injects) typeExisted(item interface{}) bool {
	return typeExists(*inj, item)
}

// concat defined array concat
func (mi *Middles) concat(target *Middles) *Middles {
	middles := append(*mi, *target...)
	return &middles
}

// toCallables defined to get `ret` that plugin func return
func (mi *Middles) toCallables() *Callables {
	callables := &Callables{}
	*callables = append(*callables, *mi...)
	return callables
}

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
	bulrush := defaultBulrush
	middles := Middles{
		Recovery,
		Override,
	}
	bulrush.Use(middles...)
	return bulrush
}

// SetMode defined httpProxy mode
func (bulrush *rush) SetMode() Bulrush {
	gin.SetMode(bulrush.config.Mode)
	return bulrush
}

// PreUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) PreUse(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		value := indirectPlugin(item)
		*bulrush.preMiddles = append(*bulrush.preMiddles, value)
	})
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) Use(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		value := indirectPlugin(item)
		*bulrush.middles = append(*bulrush.middles, value)
	})
	return bulrush
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) PostUse(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		value := indirectPlugin(item)
		*bulrush.postMiddles = append(*bulrush.postMiddles, value)
	})
	return bulrush
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bulrush *rush) Config(path string) Bulrush {
	conf := LoadConfig(path)
	conf.Version = conf.version()
	conf.Name = conf.name()
	conf.Prefix = conf.prefix()
	conf.Mode = conf.mode()
	if conf.Version != Version {
		rushLogger.Warn("Please check the latest version of bulrush's configuration file")
	}
	*bulrush.config = *conf
	bulrush.Inject(bulrush.config)
	bulrush.SetMode()
	return bulrush
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bulrush *rush) Inject(items ...interface{}) Bulrush {
	funk.ForEach(items, func(inject interface{}) {
		if bulrush.injects.typeExisted(inject) {
			rushLogger.Error("inject %v has existed", inject)
			panic(fmt.Errorf("inject %v has existed", inject))
		}
	})
	*bulrush.injects = append(*bulrush.injects, items...)
	return bulrush
}

// RunImmediately, excute plugin in orderly
// Quick start application
func (bulrush *rush) RunImmediately() {
	bulrush.Run(RunImmediately)
}

// Run application with callback, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (bulrush *rush) Run(cb interface{}) {
	bulrush.PostUse(cb)
	middles := bulrush.preMiddles.concat(bulrush.middles).concat(bulrush.postMiddles)
	callables := middles.toCallables()
	executor := &executor{
		callables: callables,
		injects:   bulrush.injects,
	}
	executor.execute(func(ret ...interface{}) {
		bulrush.Inject(ret...)
	})
}
