// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/2637309949/bulrush-utils/sync"
	"github.com/kataras/go-events"
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
		Wire(interface{}) error
		Inspect()
		ToJSON() interface{}
		GET(string, ...gin.HandlerFunc) Bulrush
		POST(string, ...gin.HandlerFunc) Bulrush
		DELETE(string, ...gin.HandlerFunc) Bulrush
		PUT(string, ...gin.HandlerFunc) Bulrush
		Run(...interface{}) error
		RunTLS(...interface{}) error
		ExecWithBooting(interface{}) error
		Shutdown() error
	}
	// rush implement Bulrush std
	rush struct {
		events.EventEmmiter
		lifecycle   Lifecycle
		config      *Config
		injects     *Injects
		prePlugins  *Plugins
		plugins     *Plugins
		postPlugins *Plugins
		lock        *sync.Lock
	}
)

// New returns a new blank Bulrush instance without any middleware attached.
// By default the configuration is:
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New(opt ...Option) Bulrush {
	bul := (&rush{
		EventEmmiter: events.New(),
		lifecycle:    &lifecycleWrapper{},
		config:       new(Config),
		injects:      newInjects(),
		prePlugins:   newPlugins(),
		plugins:      newPlugins(),
		postPlugins:  newPlugins(),
		lock:         sync.NewLock(),
	})
	for _, o := range opt {
		o.apply(bul)
	}
	bul.Empty()
	bul.Inject(builtInInjects(bul)...)
	bul.PreUse(Starting).PostUse(Running)
	return bul
}

// Default returns an Bulrush instance with the Override and Recovery middleware already attached.
// --Recovery middle has been register in httpProxy and user router
// --Override middles has been register in router for override req
func Default(opt ...Option) Bulrush {
	bul := New(opt...)
	bul.PreUse(GRPCProxy, HTTPProxy, HTTPRouter, Recovery, Override)
	return bul
}

// Empty defined empty all exists plugin and inject
// would return a empty bulrush
// should be careful
func (bul *rush) Empty() *rush {
	return Empty().apply(bul)
}

// PreUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PreUse(params ...interface{}) Bulrush {
	params = PluginsValidOption(params...).
		check(bul).([]interface{})
	return PrePluginsOption(params...).apply(bul)
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) Use(params ...interface{}) Bulrush {
	params = PluginsValidOption(params...).
		check(bul).([]interface{})
	return MiddlePluginsOption(params...).apply(bul)
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PostUse(params ...interface{}) Bulrush {
	params = PluginsValidOption(params...).
		check(bul).([]interface{})
	return PostPluginsOption(params...).apply(bul)
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bul *rush) Config(path string) Bulrush {
	conf := ConfigValidOption(path).
		check(bul).(*Config)
	return ParseConfigOption(conf).apply(bul)
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bul *rush) Inject(params ...interface{}) Bulrush {
	params = InjectsValidOption(params...).
		check(bul).([]interface{})
	return InjectsOption(params...).apply(bul)
}

// Acquire defined acquire inject ele from type
// - match type or match interface{}
// - return nil if no ele match
func (bul *rush) Acquire(target reflect.Type) interface{} {
	return bul.injects.Acquire(target)
}

// Wire defined wire ele from type
// - match type or match interface{}
// - return err if wire error
func (bul *rush) Wire(target interface{}) (err error) {
	// tv := (*interface{})(unsafe.Pointer(targetValue.Pointer()))
	// va := reflect.ValueOf(&a).Elem()
	// va.Set(reflect.New(va.Type().Elem()))
	return bul.injects.Wire(target)
}

// Return JSON representation.
// We only bother showing settings
func (bul *rush) ToJSON() interface{} {
	return struct {
		Config *Config
		Env    string
	}{
		Config: bul.config,
		Env:    modeName,
	}
}

// Inspect implementation
// We only bother showing settings
func (bul *rush) Inspect() {
	profile := bul.ToJSON()
	jsIndent, _ := json.MarshalIndent(&profile, "", "\t")
	fmt.Println(string(jsIndent))
}

// GET defined HttpProxy handles
// shortcut method
func (bul *rush) GET(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	return bul.Use(func(router *gin.RouterGroup) {
		router.GET(relativePath, handlers...)
	})
}

// POST defined HttpProxy handles
// shortcut method
func (bul *rush) POST(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	return bul.Use(func(router *gin.RouterGroup) {
		router.POST(relativePath, handlers...)
	})
}

// PUT defined HttpProxy handles
// shortcut method
func (bul *rush) PUT(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	return bul.Use(func(router *gin.RouterGroup) {
		router.PUT(relativePath, handlers...)
	})
}

// DELETE defined HttpProxy handles
// shortcut method
func (bul *rush) DELETE(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	return bul.Use(func(router *gin.RouterGroup) {
		router.DELETE(relativePath, handlers...)
	})
}

// CatchError error which one from outside of recovery pluigns, this rec just for bulrush
// you can CatchError if your error code does not affect the next plug-in
// sometime you should handler all error in plugin
func (bul *rush) CatchError(funk interface{}) error {
	return CatchError(funk)
}

// Done returns a channel of signals to block on after starting the
// application
func (bul *rush) Done() <-chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return c
}

// Shutdown defined bul gracefulExit
// ,, close http or other resources
// should call Shutdown after bulrush has running success
func (bul *rush) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return withTimeout(ctx, bul.lifecycle.Stop)
}

// Run application with a booting`plugin, excute plugin in orderly
// Just for HTTPProxy booting
// Note: this method will block the calling goroutine indefinitely unless an error happens
func (bul *rush) Run(b ...interface{}) (err error) {
	var booting interface{} = HTTPBooting
	if len(b) > 0 {
		booting = b[0]
	}
	return bul.ExecWithBooting(booting)
}

// Run application with a booting`plugin, excute plugin in orderly
// Just for HTTPProxy booting
// Note: this method will block the calling goroutine indefinitely unless an error happens
func (bul *rush) RunTLS(b ...interface{}) (err error) {
	var booting interface{} = HTTPTLSBooting
	if len(b) > 0 {
		booting = b[0]
	}
	return bul.ExecWithBooting(booting)
}

// execWithBooting defined exeute plugins with a booting plugins
// Note: this method will block the calling goroutine indefinitely unless an error happens
func (bul *rush) ExecWithBooting(b interface{}) (err error) {
	done := bul.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bul.PostUse(b)
	go func() {
		err = bul.CatchError(func() {
			plugins := bul.prePlugins.Append(bul.plugins).Append(bul.postPlugins)
			scopes := plugins.toScopes(func(t reflect.Type) interface{} {
				return bul.injects.Acquire(t)
			})
			exec := &engine{
				scopes: scopes,
			}
			exec.exec(func(ret ...interface{}) {
				bul.Inject(ret...)
			})
			withTimeout(ctx, bul.lifecycle.Start)
		})
	}()
	<-done
	return
}
