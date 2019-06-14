// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	// PNRet return a plugin after call Plugin func
	PNRet interface{}
	// PNBase Plugin interface defined
	PNBase interface {
		Plugin() PNRet
	}
)

// PNStruct for a quickly Plugin SetUp when you dont want declare PNBase
// PBBase minimize implement
type PNStruct struct{ pn interface{} }

// Plugin for PNQuick
func (p *PNStruct) Plugin() PNRet {
	return p.pn
}

// PNQuick for PNQuick
var PNQuick = func(pn interface{}) PNBase {
	return &PNStruct{
		pn: pn,
	}
}

// Recovery system rec from panic
var Recovery = PNQuick(func(httpProxy *gin.Engine, router *gin.RouterGroup) {
	httpProxy.Use(gin.Recovery())
	router.Use(gin.Recovery())
})

// HTTPProxy create http proxy
var HTTPProxy = PNQuick(func() *gin.Engine {
	return gin.New()
})

// HTTPRouter create http router
var HTTPRouter = PNQuick(func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
	return httpProxy.Group(config.Prefix)
})

// Override http methods
var Override = PNQuick(func(router *gin.RouterGroup, httpProxy *gin.Engine) {
	router.Use(func(c *gin.Context) {
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

// RunImmediately run app
var RunImmediately = PNQuick(func(httpProxy *gin.Engine, config *Config) {
	port := fixedPortPrefix(strings.TrimSpace(config.Port))
	name := config.Name
	rushLogger.Debug("================================")
	rushLogger.Debug("App: %s", name)
	rushLogger.Debug("Listen on %s", port)
	rushLogger.Debug("================================")
	httpProxy.Run(port)
})
