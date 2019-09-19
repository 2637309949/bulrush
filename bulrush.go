// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"encoding/json"
	"fmt"
	"reflect"
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
		config      *Config
		injects     *Injects
		prePlugins  *Plugins
		plugins     *Plugins
		postPlugins *Plugins
		lock        *sync.Lock
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
		injects:      newInjects(),
		prePlugins:   newPlugins(),
		plugins:      newPlugins(),
		postPlugins:  newPlugins(),
		lock:         sync.NewLock(),
		exit:         make(chan struct{}, 1),
	})
	bul.Empty()
	bul.Inject(builtInInjects(bul)...).
		PreUse(Plugins{Starting, HTTPProxy, GRPCProxy, HTTPRouter}...).
		Use(Plugins{}...).
		PostUse(Plugins{Running}...)
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
	cParams := PluginsValidOption(params...).
		check(bul).([]interface{})
	return PrePluginsOption(cParams...).apply(bul)
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) Use(params ...interface{}) Bulrush {
	cParams := PluginsValidOption(params...).
		check(bul).([]interface{})
	return MiddlePluginsOption(cParams...).apply(bul)
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PostUse(params ...interface{}) Bulrush {
	cParams := PluginsValidOption(params...).
		check(bul).([]interface{})
	return PostPluginsOption(cParams...).apply(bul)
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
	cParams := InjectsValidOption(params...).
		check(bul).([]interface{})
	return InjectsOption(cParams...).apply(bul)
}

// Acquire defined acquire inject ele from type
// - match type or match interface{}
// - return nil if no ele match
func (bul *rush) Acquire(ty reflect.Type) interface{} {
	return bul.injects.Acquire(ty)
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
	bul.Use(func(router *gin.RouterGroup) {
		router.GET(relativePath, handlers...)
	})
	return bul
}

// POST defined HttpProxy handles
// shortcut method
func (bul *rush) POST(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	bul.Use(func(router *gin.RouterGroup) {
		router.POST(relativePath, handlers...)
	})
	return bul
}

// PUT defined HttpProxy handles
// shortcut method
func (bul *rush) PUT(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	bul.Use(func(router *gin.RouterGroup) {
		router.PUT(relativePath, handlers...)
	})
	return bul
}

// DELETE defined HttpProxy handles
// shortcut method
func (bul *rush) DELETE(relativePath string, handlers ...gin.HandlerFunc) Bulrush {
	bul.Use(func(router *gin.RouterGroup) {
		router.DELETE(relativePath, handlers...)
	})
	return bul
}

// CatchError error which one from outside of recovery pluigns, this rec just for bulrush
// you can CatchError if your error code does not affect the next plug-in
// sometime you should handler all error in plugin
func (bul *rush) CatchError(funk interface{}) error {
	return CatchError(funk)
}

// Shutdown defined bul gracefulExit
// ,, close http or other resources
// should call Shutdown after bulrush has running success
func (bul *rush) Shutdown() error {
	defer func() {
		close(bul.exit)
	}()
	// emit shutdown event
	bul.Emit(EventsShutdown)
	// shutdown bulrush
	time.Sleep(time.Second * 5)
	<-func() chan struct{} {
		bul.exit <- struct{}{}
		return bul.exit
	}()
	rushLogger.Warn("Shutdown: bulrush Closed")
	return nil
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
	go func() {
		err = bul.CatchError(func() {
			bul.PostUse(b)
			scopes := bul.
				prePlugins.
				Append(bul.plugins).
				Append(bul.postPlugins).
				toScopes()
			exec := &executor{
				scopes:  scopes,
				injects: bul.injects,
			}
			exec.execute(func(ret ...interface{}) {
				bul.Inject(ret...)
			})
		})
	}()
	bul.exit <- <-bul.exit
	return
}
