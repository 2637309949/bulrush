// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/2637309949/bulrush-utils/sync"
	"github.com/gin-gonic/gin"
	"github.com/kataras/go-events"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc"
)

type (
	// Plugins defined those that can be call by reflect
	// , Plugins passby func or a struct that has `Plugin` func
	Plugins []interface{}
	// HTTPContext defined httpContxt
	HTTPContext struct {
		Chan         chan struct{}
		DeadLineTime time.Time
	}
	// PluginsOption defined plugin option
	PluginsOption interface {
		apply(r *rush) *rush
		check(r *rush) interface{}
	}
)

// PluginsValidOption defined Option of valid
func PluginsValidOption(plugins ...interface{}) PluginsOption {
	return Option(func(r *rush) interface{} {
		funk.ForEach(plugins, func(item interface{}) {
			assert1(isPlugin(item), ErrWith(ErrPlugin, fmt.Sprintf("%v can not be used as plugin", item)))
		})
		return plugins
	})
}

// PrePluginsOption defined Option of PrePlugin
func PrePluginsOption(plugins ...interface{}) PluginsOption {
	return Option(func(r *rush) interface{} {
		r.lock.Acquire("prePlugins", func(async sync.Async) {
			funk.ForEach(plugins, func(item interface{}) {
				r.prePlugins.Put(item)
			})
		})
		return r
	})
}

// PostPluginsOption defined Option of PostPlugin
func PostPluginsOption(plugins ...interface{}) PluginsOption {
	return Option(func(r *rush) interface{} {
		r.lock.Acquire("postPlugins", func(async sync.Async) {
			funk.ForEach(plugins, func(item interface{}) {
				r.postPlugins.Put(item)
			})
		})
		return r
	})
}

// MiddlePluginsOption defined Option of MiddlePlugin
func MiddlePluginsOption(plugins ...interface{}) PluginsOption {
	return Option(func(r *rush) interface{} {
		r.lock.Acquire("plugins", func(async sync.Async) {
			funk.ForEach(plugins, func(item interface{}) {
				r.plugins.Put(item)
			})
		})
		return r
	})
}

func newPlugins(items ...interface{}) *Plugins {
	plugins := make(Plugins, 0)
	plugins = append(plugins, items...)
	return &plugins
}

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

// PutHead defined put ele to head
func (p *Plugins) PutHead(target interface{}) *Plugins {
	*p = append([]interface{}{target}, *p...)
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

// toScopes defined to get `ret` that plugin func return
func (p *Plugins) toScopes(acquire func(reflect.Type) interface{}) *[]scope {
	scopes := funk.Map(*p, func(v interface{}) scope {
		return determineScope(v, acquire)
	}).([]scope)
	return &scopes
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

//	 Recovery         plugin   defined sys recovery
//   HTTPProxy        plugin   defined http proxy
//   HTTPRouter       plugin   defined http router
//   Override         plugin   defined method override
//   Run   plugin   defined httpproxy run
var (
	// Starting defined before all plugin
	Starting = func(event events.EventEmmiter) {
		event.Emit(EventsStarting, EventsStarting)
	}
	// Recovery system rec from panic
	Recovery = func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(recovery())
		router.Use(recovery())
	}
	// HTTPProxy create http proxy
	HTTPProxy = func() *gin.Engine {
		return gin.New()
	}
	// GRPCProxy create grpc proxy
	GRPCProxy = func() *grpc.Server {
		return grpc.NewServer()
	}
	// HTTPRouter create http router
	HTTPRouter = func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
		return httpProxy.Group(config.Prefix)
	}
	// Override http methods
	Override = func(router *gin.RouterGroup, httpProxy *gin.Engine) {
		funk.ForEach([]func(...gin.HandlerFunc) gin.IRoutes{router.Use, httpProxy.Use}, func(Use func(...gin.HandlerFunc) gin.IRoutes) {
			Use(func(c *gin.Context) {
				if c.Request.Method != "POST" {
					c.Next()
				} else {
					method := c.PostForm("_method")
					methods := [3]string{"DELETE", "PUT", "PATCH"}
					if method != "" {
						for _, target := range methods {
							if target == strings.ToUpper(method) {
								c.Request.Method = target
								httpProxy.HandleContext(c)
								break
							}
						}
					}
				}
			})
		})
	}
	// HTTPBooting run http proxy and grpc proxy
	HTTPBooting = func(httpProxy *gin.Engine, gs *grpc.Server, event events.EventEmmiter, config *Config) {
		var err error
		defer func() {
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		addr1 := fixedPortPrefix(strings.TrimSpace(config.Port))
		addr2 := fixedPortPrefix(strings.TrimSpace(config.Port), 1)
		name := config.Name
		grpc, err := net.Listen("tcp", addr2)
		http := &http.Server{Addr: addr1, Handler: httpProxy}
		event.On(EventsShutdown, func(payload ...interface{}) {
			// graceful stop srv
			gs.GracefulStop()
			http.Shutdown(&HTTPContext{
				DeadLineTime: time.Now().Add(3 * time.Second),
				Chan:         make(chan struct{}, 1),
			})
		})
		go func() {
			err = http.ListenAndServe()
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		go func() {
			err = gs.Serve(grpc)
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		rushLogger.Debug("================================")
		rushLogger.Debug("App: %s", name)
		rushLogger.Debug("Env: %s", config.Env)
		rushLogger.Debug("Http Listen on %s", addr1)
		rushLogger.Debug("Grpc Listen on %s", addr2)
		rushLogger.Debug("================================")
	}
	// HTTPTLSBooting run http proxy and grpc proxy
	HTTPTLSBooting = func(httpProxy *gin.Engine, gs *grpc.Server, event events.EventEmmiter, config *Config) {
		var err error
		defer func() {
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		addr1 := fixedPortPrefix(strings.TrimSpace(config.Port))
		addr2 := fixedPortPrefix(strings.TrimSpace(config.Port), 1)
		name := config.Name
		grpc, err := net.Listen("tcp", addr2)
		http := &http.Server{Addr: addr1, Handler: httpProxy}
		event.On(EventsShutdown, func(payload ...interface{}) {
			// graceful stop srv
			gs.GracefulStop()
			http.Shutdown(&HTTPContext{
				DeadLineTime: time.Now().Add(3 * time.Second),
				Chan:         make(chan struct{}, 1),
			})
		})
		go func() {
			err = http.ListenAndServeTLS(config.TLS.CRT, config.TLS.Key)
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		go func() {
			err = gs.Serve(grpc)
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		rushLogger.Debug("================================")
		rushLogger.Debug("App: %s", name)
		rushLogger.Debug("Env: %s", config.Env)
		rushLogger.Debug("Http Listen on %s", addr1)
		rushLogger.Debug("Grpc Listen on %s", addr2)
		rushLogger.Debug("================================")
	}
	// Running defined after all plugin
	Running = func(event events.EventEmmiter) {
		event.Emit(EventsRunning, EventsRunning)
	}
)
