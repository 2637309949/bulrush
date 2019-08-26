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
		exit:         make(chan struct{}, 1),
	})
	bul.Empty()
	bul.Inject(builtInInjects(bul)...).
		PreUse(Plugins{Starting, HTTPProxy, HTTPRouter}...).
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
func (bul *rush) PreUse(items ...interface{}) Bulrush {
	return PrePluginsOption(items...).apply(bul)
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) Use(items ...interface{}) Bulrush {
	return MiddlePluginsOption(items...).apply(bul)
}

// PostUse attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bul *rush) PostUse(items ...interface{}) Bulrush {
	return PostPluginsOption(items...).apply(bul)
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bul *rush) Config(path string) Bulrush {
	return ParseConfigOption(path).apply(bul)
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bul *rush) Inject(items ...interface{}) Bulrush {
	return InjectsOption(items...).apply(bul)
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

// CatchError error which one from outside of recovery pluigns, this rec just for bulrush
// you can CatchError if your error code does not affect the next plug-in
// sometime you should handler all error in plugin
func CatchError(funk interface{}) (err error) {
	defer func() {
		if ret := recover(); ret != nil {
			ok, bulError := false, &Error{Code: ErrNu.Code, Err: err}
			if err, ok = ret.(error); !ok {
				err = fmt.Errorf("%v", ret)
			}
			if bulError, ok = ErrOut(err); !ok {
				bulError.Err = err
			}
			if rushLogger != nil {
				rushLogger.Error("%s panic recovered:\n%s\n%s%s", timeFormat(time.Now()), bulError.Err, stack(3), reset)
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

// RunImmediately, excute plugin in orderly
// Quick start application
func (bul *rush) RunImmediately() error {
	return bul.Run(HTTPBooting)
}

// Run application with callback, excute plugin in orderly
// Note: this method will block the calling goroutine indefinitely unless an error happens
func (bul *rush) Run(p interface{}) (err error) {
	go func() {
		err = bul.CatchError(func() {
			bul.PostUse(p)
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
