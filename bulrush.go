// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"
	"time"

	"github.com/2637309949/bulrush-utils/sync"
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
		Acquire(reflect.Type) interface{}
		RunImmediately() error
		Run(interface{}) error
		Shutdown() error
	}
	// rush implement Bulrush std
	rush struct {
		events.EventEmmiter
		config      *Config
		injects     *Injects
		prePlugins  *Plugins
		plugins     *Plugins
		postPlugins *Plugins
		lock        *sync.Lock
		httpContext *HTTPContext
		exit        chan struct{}
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
		lock:         sync.NewLock(),
		httpContext:  NewHTTPContext(3 * time.Second),
		exit:         make(chan struct{}, 1),
	})
	bul.
		Clear().
		Inject(builtInInjects(bul)...).
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

// Clear defined empty all exists plugin and inject
// would return a empty bulrush
// should be careful
func (bul *rush) Clear() Bulrush {
	bul.injects = new(Injects)
	bul.prePlugins = new(Plugins)
	bul.plugins = new(Plugins)
	bul.postPlugins = new(Plugins)
	return bul
}

// PreUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PreUse(items ...interface{}) Bulrush {
	if len(items) == 0 {
		return bul
	}
	bul.lock.Acquire("prePlugins", func(async sync.Async) {
		funk.ForEach(items, func(item interface{}) {
			assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypePlugin,
				Err: fmt.Errorf("%v can not be used as plugin", item)}})
			bul.prePlugins.Put(item)
		})
	})
	return bul
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) Use(items ...interface{}) Bulrush {
	if len(items) == 0 {
		return bul
	}
	bul.lock.Acquire("plugins", func(async sync.Async) {
		funk.ForEach(items, func(item interface{}) {
			assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypePlugin,
				Err: fmt.Errorf("%v can not be used as plugin", item)}})
			bul.plugins.Put(item)
		})
	})
	return bul
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PostUse(items ...interface{}) Bulrush {
	if len(items) == 0 {
		return bul
	}
	bul.lock.Acquire("postPlugins", func(async sync.Async) {
		funk.ForEach(items, func(item interface{}) {
			assert1(isPlugin(item), errorMsgs{&Error{Type: ErrorTypePlugin,
				Err: fmt.Errorf("%v can not be used as plugin", item)}})
			bul.postPlugins.Put(item)
		})
	})
	return bul
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bul *rush) Config(path string) Bulrush {
	if len(path) == 0 {
		return bul
	}
	bul.lock.Acquire("config", func(async sync.Async) {
		conf := LoadConfig(path)
		conf.Version = conf.version()
		conf.Name = conf.name()
		conf.Prefix = conf.prefix()
		conf.Mode = conf.mode()
		SetMode(conf.Mode)
		conf.verifyVersion(Version)
		*bul.config = *conf
		bul.Inject(bul.config)
	})
	return bul
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bul *rush) Inject(items ...interface{}) Bulrush {
	if len(items) == 0 {
		return bul
	}
	bul.lock.Acquire("injects", func(async sync.Async) {
		funk.ForEach(items, func(item interface{}) {
			assert1(!bul.injects.Has(item), errorMsgs{&Error{Type: ErrorTypeInject,
				Err: fmt.Errorf("inject %v has existed", reflect.TypeOf(item))}})
			bul.injects.Put(item)
		})
	})
	return bul
}

// Acquire defined acquire inject ele from type
// - match type or match interface{}
// - return nil if no ele match
func (bul *rush) Acquire(ty reflect.Type) interface{} {
	ele := typeMatcher(ty, *bul.injects)
	if ele == nil {
		ele = duckMatcher(ty, *bul.injects)
	}
	return ele
}

// CatchError error which one from outside of recovery pluigns, this rec just for bulrush
// you can CatchError if your error code does not affect the next plug-in
// sometime you should handler all error in plugin
func CatchError(funk interface{}) (err error) {
	defer func() {
		var ok bool
		if ret := recover(); ret != nil {
			err, ok = ret.(error)
			if !ok {
				err = fmt.Errorf("%v", ret)
			}
			if rushLogger != nil {
				rushLogger.Error("%s panic recovered:\n%s\n%s%s",
					timeFormat(time.Now()), err, stack(3), reset)
			}
		}
	}()
	assert1(isFunc(funk), fmt.Errorf("funk %v should be func type", reflect.TypeOf(funk)))
	reflect.ValueOf(funk).Call([]reflect.Value{})
	return
}

// CatchError error which one from outside of recovery pluigns, this rec just for bulrush
// you can CatchError if your error code does not affect the next plug-in
// sometime you should handler all error in plugin
func (bul *rush) CatchError(funk interface{}) error {
	return CatchError(funk)
}

// RunImmediately, excute plugin in orderly
// Quick start application
func (bul *rush) RunImmediately() error {
	return bul.Run(RunImmediately(bul.NewHTTPContext(1 * time.Second)))
}

// NewHTTPContext defined obtain a httpContext for httpProxy
// if you implement run logic, you should obtain a ctx for HttpProxy
// reference RunImmediately plugin
func (bul *rush) NewHTTPContext(duration time.Duration) *HTTPContext {
	bul.httpContext.DeadLineTime = time.Now().Add(duration)
	return bul.httpContext
}

// Shutdown defined bul gracefulExit
// ,, close http or other resources
func (bul *rush) Shutdown() error {
	<-func() chan struct{} {
		bul.httpContext.Exit <- struct{}{}
		return bul.httpContext.Exit
	}()
	rushLogger.Warn("Shutdown: httpProxy Closed")
	<-func() chan struct{} {
		bul.exit <- struct{}{}
		return bul.exit
	}()
	rushLogger.Warn("Shutdown: bulrush Closed")
	return nil
}

// Run application with callback, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens
func (bul *rush) Run(p interface{}) (err error) {
	go func() {
		err = bul.CatchError(func() {
			bul.PostUse(p)
			pcts := bul.
				prePlugins.
				Append(bul.plugins).
				Append(bul.postPlugins).
				toPluginContexts()
			executor := &executor{
				pluginContexts: pcts,
				injects:        bul.injects,
			}
			executor.execute(func(ret ...interface{}) {
				bul.Inject(ret...)
			})
		})
	}()
	bul.exit <- <-bul.exit
	return
}
