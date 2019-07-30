// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"

	"github.com/kataras/go-events"
	"github.com/thoas/go-funk"
)

type (

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
		injects     *Injects
		prePlugins  *Plugins
		plugins     *Plugins
		postPlugins *Plugins
	}
)

// New returns a new blank Bulrush instance without any middleware attached.
// By default the configuration is:
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() Bulrush {
	bul := (&rush{
		EventEmmiter: events.New(),
		config:       new(Config),
		injects:      new(Injects),
		prePlugins:   new(Plugins),
		plugins:      new(Plugins),
		postPlugins:  new(Plugins),
	})
	bul.
		Inject(preInjects(bul)...).
		PreUse(Plugins{HTTPProxy, HTTPRouter}...).
		Use(Plugins{}...).
		PostUse(Plugins{}...)
	return bul
}

// Default returns an Bulrush instance with the Override and Recovery middleware already attached.
// --Recovery middle has been register in httpProxy and user router
// --Override middles has been register in router for override req
func Default() Bulrush {
	bul := New()
	bul.Use(Recovery)
	bul.Use(Override)
	return bul
}

// PreUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PreUse(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypePlugin, Err: fmt.Errorf("%v can not be used as plugin", item)}})
		*bul.prePlugins = append(*bul.prePlugins, item)
	})
	return bul
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) Use(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypePlugin, Err: fmt.Errorf("%v can not be used as plugin", item)}})
		*bul.plugins = append(*bul.plugins, item)
	})
	return bul
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PostUse(items ...interface{}) Bulrush {
	funk.ForEach(items, func(item interface{}) {
		assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypePlugin, Err: fmt.Errorf("%v can not be used as plugin", item)}})
		*bul.postPlugins = append(*bul.postPlugins, item)
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
	conf.verifyVersion(Version)
	*bul.config = *conf
	bul.Inject(bul.config)
	return bul
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bul *rush) Inject(items ...interface{}) Bulrush {
	funk.ForEach(items, func(inject interface{}) {
		assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypeInject, Err: fmt.Errorf("inject %v has existed", inject)}})
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
	plugin := bul.prePlugins.Append(bul.plugins).Append(bul.postPlugins)
	pv := plugin.toPluginValues()
	executor := &executor{
		pluginValues: pv,
		injects:      bul.injects,
	}
	executor.execute(func(ret ...interface{}) {
		bul.Inject(ret...)
	})
}
