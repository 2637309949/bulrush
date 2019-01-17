/**
 * @author [double]
 * @email [2637309949@qq.com]
 * @create date 2019-01-15 09:49:33
 * @modify date 2019-01-15 09:49:33
 * @desc [default plugins for bulrush]
 */

package bulrush

import (
	"strings"
	"github.com/gin-gonic/gin"
)

type (
	// PNRet return a plugin after call Plugin func
	PNRet interface{}
	// PNBase Plugin interface defined
	PNBase interface{
		Plugin() PNRet
	}
	// PNStruct for a quickly Plugin SetUp when you dont want declare PNBase
	PNStruct struct {
		Quick interface{}
	}
	// Override Plugin
	Override struct {
		PNBase
	}
	// Recovery system rec from panic
	Recovery struct {
		PNBase
	}
	// HTTPProxy create http proxy
	HTTPProxy struct {
		PNBase
	}
	// HTTPRouter create http router
	HTTPRouter struct {
		PNBase
	}
	// RUNProxy run proxy
	RUNProxy struct {
		PNBase
	    CallBack func(error, *Config)
	}
	// LoggerWriter log req
	LoggerWriter struct {
		PNBase
		Bulrush 		 *rush
		LoggerWithWriter func(*rush) gin.HandlerFunc
	}
)


// Plugin for PNQuick
func(pnStruct *PNStruct) Plugin() PNRet {
	return pnStruct.Quick
}

// PNQuick for a quickly Plugin SetUp when you dont want declare PNBase
func PNQuick(quick interface{}) PNBase {
	return &PNStruct{
		Quick: quick,
	}
}

// Plugin for Recovery
func(recovery *Recovery) Plugin() PNRet {
	return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}
}

// Plugin for HTTPProxy
func(httpProxy *HTTPProxy) Plugin() PNRet {
	return func() *gin.Engine {
		proxy := gin.New()
		return proxy
	}
}

// Plugin for HTTPRouter
func(httpRouter *HTTPRouter) Plugin() PNRet {
	return func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
		router := httpProxy.Group(config.GetString("prefix","/api/v1"))
		return router
	}
}

// Plugin for RUNProxy
func(runProxy *RUNProxy) Plugin() PNRet {
	return func(httpProxy *gin.Engine, config *Config) {
		port := config.GetString("port",  ":8080")
		runProxy.CallBack(nil, config)
		err := httpProxy.Run(port)
		runProxy.CallBack(err, config)
	}
}

// Plugin for LoggerWriter
func(loggerWriter *LoggerWriter) Plugin() PNRet {
	return func(router *gin.RouterGroup) {
		router.Use(loggerWriter.LoggerWithWriter(loggerWriter.Bulrush))
	}
}

// Plugin for gin
func (override *Override) Plugin() PNRet {
	return func(router *gin.RouterGroup, httpProxy *gin.Engine) {
		httpProxy.Use(func(c *gin.Context) {
			if c.Request.Method != "POST" {
				c.Next()
			} else {
				method := c.PostForm("_method")
				methods := [3]string{"DELETE", "PUT", "PATCH"}
				if method != "" {
					for _, target := range methods {
						if(target == strings.ToUpper(method)) {
							c.Request.Method = target
							httpProxy.HandleContext(c)
							break
						}
					}
				}
			}
		})
		router.Use(func(c *gin.Context) {
			if c.Request.Method != "POST" {
				c.Next()
			} else {
				method := c.PostForm("_method")
				methods := [3]string{"DELETE", "PUT", "PATCH"}
				if method != "" {
					for _, target := range methods {
						if(target == strings.ToUpper(method)) {
							c.Request.Method = target
							httpProxy.HandleContext(c)
							break
						}
					}
				}
			}
		})
	}
}
