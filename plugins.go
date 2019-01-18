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
	// all user custom must implement this interface{}
	PNRet      interface{}
	// PNBase Plugin interface defined
	// all user custom must implement this interface{}
	PNBase     interface{ Plugin() PNRet }
	// PNStruct for a quickly Plugin SetUp when you dont want declare PNBase
	// PBBase minimize implement
	PNStruct   struct { Quick interface{} }
	// Override Plugin
	Override   struct { PNBase }
	// Recovery system rec from panic
	Recovery   struct { PNBase }
	// HTTPProxy create http proxy
	HTTPProxy  struct { PNBase }
	// HTTPRouter create http router
	HTTPRouter struct { PNBase }
	// RUNProxy run proxy
	RUNProxy   struct { PNBase, CallBack func(error, *Config) }
)

// Plugin for PNQuick
// if your do not want to implement PNBase interface{}
// use:
// app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
// 	router.GET("/bulrushApp", func (c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"data": 	testInject,
// 			"errcode": 	nil,
// 			"errmsg": 	nil,
// 		})
// 	})
// }))
func(pn *PNStruct) Plugin() PNRet {
	return pn.Quick
}

// PNQuick for PNQuick
// if your do not want to implement PNBase interface{}
// use:
// app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
// 	router.GET("/bulrushApp", func (c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"data": 	testInject,
// 			"errcode": 	nil,
// 			"errmsg": 	nil,
// 		})
// 	})
// }))
func PNQuick(method interface{}) PNBase {
	pn := &PNStruct{
		Quick: method,
	}
	return pn
}

// Plugin for Recovery
// recovery from system panic
// use:
// defaultMiddles := Middles {
// 	&Recovery{},
// 	&Override{},
// }
// bulrush.Use(defaultMiddles...)
func(pn *Recovery) Plugin() PNRet {
	return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}
}

// Plugin for HTTPProxy
func(pn *HTTPProxy) Plugin() PNRet {
	return func() *gin.Engine {
		proxy := gin.New()
		return proxy
	}
}

// Plugin for HTTPRouter
func(pn *HTTPRouter) Plugin() PNRet {
	return func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
		router := httpProxy.Group(config.GetString("prefix","/api/v1"))
		return router
	}
}

// Plugin for RUNProxy
// use:
// lastMiddles := Middles {
// 	&RUNProxy{ CallBack: cb },
// }
// bulrush.Use(lastMiddles...)
func(pn *RUNProxy) Plugin() PNRet {
	return func(httpProxy *gin.Engine, config *Config) {
		port := config.GetString("port",  ":8080")
		pn.CallBack(nil, config)
		err := httpProxy.Run(port)
		pn.CallBack(err, config)
	}
}

// Plugin for Override
// recovery from system panic
// use:
// defaultMiddles := Middles {
// 	&Recovery{},
// 	&Override{},
// }
// bulrush.Use(defaultMiddles...)
func (pn *Override) Plugin() PNRet {
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
