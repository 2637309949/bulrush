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

type (
	// middles defined those that can be call by reflect
	// , middles passby func or a struct that has `Plugin` func
	middles []interface{}
	// injects defined some entitys that can be inject to middle
	// , inject would panic if repetition
	// , inject can be go base tyle or struct or ptr or interface{}
	injects []interface{}
	// Bulrush interface{} defined all framework should be
	// , also sys provide a default Bulrush - `rush`
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
	// rush implement Bulrush std
	rush struct {
		events.EventEmmiter
		config      *Config
		injects     *injects
		preMiddles  *middles
		middles     *middles
		postMiddles *middles
	}
)

// concat defined array concat
func (src *injects) concat(target *injects) *injects {
	injects := append(*src, *target...)
	return &injects
}

// typeExisted defined inject type is existed or not
func (src *injects) typeExisted(item interface{}) bool {
	return typeExists(*src, item)
}

// concat defined array concat
func (src *middles) concat(target *middles) *middles {
	middles := append(*src, *target...)
	return &middles
}

// toCallables defined to get `ret` that plugin func return
func (src *middles) toCallables() *callables {
	cbs := &callables{}
	*cbs = append(*cbs, *src...)
	return cbs
}

// New returns a new blank Bulrush instance without any middleware attached.
// By default the configuration is:
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() Bulrush {
	preMid := make(middles, 0)
	poMid := make(middles, 0)
	mid := make(middles, 0)
	injects := make(injects, 0)
	emmiter := events.New()
	bulrush := &rush{
		EventEmmiter: emmiter,
		config:       &Config{},
		injects:      &injects,
		preMiddles:   &preMid,
		middles:      &mid,
		postMiddles:  &poMid,
	}
	defaultMiddles := middles{
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
	bul := New()
	bul.Use(middles{
		Recovery,
		Override,
	}...)
	return bul
}

// SetMode defined httpProxy mode
func (bul *rush) setMode() Bulrush {
	gin.SetMode(bul.config.Mode)
	return bul
}

// PreUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PreUse(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		value := indirectPlugin(item)
		*bul.preMiddles = append(*bul.preMiddles, value)
	})
	return bul
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) Use(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		value := indirectPlugin(item)
		*bul.middles = append(*bul.middles, value)
	})
	return bul
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PostUse(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		value := indirectPlugin(item)
		*bul.postMiddles = append(*bul.postMiddles, value)
	})
	return bul
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bul *rush) Config(path string) Bulrush {
	conf := LoadConfig(path)
	conf.Version = conf.version()
	conf.Name = conf.name()
	conf.Prefix = conf.prefix()
	conf.Mode = conf.mode()
	if conf.Version != Version {
		rushLogger.Warn("Please check the latest version of bulrush's configuration file")
	}
	*bul.config = *conf
	bul.Inject(bul.config)
	bul.setMode()
	return bul
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bul *rush) Inject(items ...interface{}) Bulrush {
	funk.ForEach(items, func(inject interface{}) {
		if bul.injects.typeExisted(inject) {
			rushLogger.Error("inject %v has existed", inject)
			panic(fmt.Errorf("inject %v has existed", inject))
		}
	})
	*bul.injects = append(*bul.injects, items...)
	return bul
}

// RunImmediately, excute plugin in orderly
// Quick start application
func (bul *rush) RunImmediately() {
	bul.Run(RunImmediately)
}

// Run application with callback, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (bul *rush) Run(cb interface{}) {
	bul.PostUse(cb)
	middles := bul.preMiddles.concat(bul.middles).concat(bul.postMiddles)
	callables := middles.toCallables()
	executor := &executor{
		callables: callables,
		injects:   bul.injects,
	}
	executor.execute(func(ret ...interface{}) {
		bul.Inject(ret...)
	})
}
