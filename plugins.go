// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kataras/go-events"
	"github.com/thoas/go-funk"
)

//	 Recovery         plugin   defined sys recovery
//   HTTPProxy        plugin   defined http proxy
//   HTTPRouter       plugin   defined http router
//   Override         plugin   defined method override
//   RunImmediately   plugin   defined httpproxy run
var (
	// Recovery system rec from panic
	Recovery = func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(recovery())
		router.Use(recovery())
	}
	// HTTPProxy create http proxy
	HTTPProxy = func() *gin.Engine {
		return gin.New()
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
	// RunImmediately run app
	RunImmediately = func(httpProxy *gin.Engine, event events.EventEmmiter, config *Config) {
		var err error
		defer func() {
			if err != nil {
				rushLogger.Error(fmt.Sprintf("%v", err))
			}
		}()
		addr := fixedPortPrefix(strings.TrimSpace(config.Port))
		name := config.Name
		rushLogger.Debug("================================")
		rushLogger.Debug("App: %s", name)
		rushLogger.Debug("Listen on %s", addr)
		rushLogger.Debug("================================")
		event.Emit(EventSysBulrushPluginRunImmediately, EventSysBulrushPluginRunImmediately)
		server := &http.Server{Addr: addr, Handler: httpProxy}
		event.On("shutdown", func(payload ...interface{}) {
			server.Shutdown(NewHTTPContext(3 * time.Second))
		})
		err = server.ListenAndServe()
	}
)
