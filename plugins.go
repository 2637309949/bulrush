/**
 * @author [double]
 * @email [2637309949@qq.com]
 * @create date 2019-01-15 09:49:33
 * @modify date 2019-01-15 09:49:33
 * @desc [default plugins for bulrush]
 */

package bulrush

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	// PNRet return a plugin after call Plugin func
	// all user custom must implement this interface{}
	PNRet interface{}
	// PNBase Plugin interface defined
	// all user custom must implement this interface{}
	PNBase interface {
		Plugin() PNRet
	}
)

// PNStruct for a quickly Plugin SetUp when you dont want declare PNBase
// PBBase minimize implement
type PNStruct struct{ Quick interface{} }

// Plugin for PNQuick
// if your do not want to implement PNBase interface{}
// use:
// app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
// 		router.GET("/bulrushApp", func (c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{
// 			"message": 	testInject,
//			})
// 		})
// }))
func (pn *PNStruct) Plugin() PNRet {
	return pn.Quick
}

// PNQuick for PNQuick
// if your do not want to implement PNBase interface{}
// use:
// app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
// 		router.GET("/bulrushApp", func (c *gin.Context) {
// 			c.JSON(http.StatusOK, gin.H{
// 			"message": 	testInject,
//			})
// 		})
// }))
func PNQuick(q interface{}) PNBase {
	return &PNStruct{
		Quick: q,
	}
}

// Recovery system rec from panic
type Recovery struct{ PNBase }

// Plugin for Recovery
func (pn *Recovery) Plugin() PNRet {
	return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}
}

// HTTPProxy create http proxy
type HTTPProxy struct{ PNBase }

// Plugin for HTTPProxy
func (pn *HTTPProxy) Plugin() PNRet {
	return func() *gin.Engine {
		return gin.New()
	}
}

// HTTPRouter create http router
type HTTPRouter struct{ PNBase }

// Plugin for create http router
func (pn *HTTPRouter) Plugin() PNRet {
	return func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
		return httpProxy.Group(config.GetString("prefix", "/api/v1"))
	}
}

// RUNProxy run HttpProxy
type RUNProxy struct{ PNBase, CallBack func(error, *Config) }

// Plugin for RUNProxy
func (pn *RUNProxy) Plugin() PNRet {
	return func(httpProxy *gin.Engine, config *Config) {
		var err error
		defer func() { pn.CallBack(err, config) }()
		port := config.GetString("port", ":8080")
		port = strings.TrimSpace(port)
		if prefix := port[:1]; prefix != ":" {
			port = fmt.Sprintf(":%s", port)
		}
		pn.CallBack(nil, config)
		err = httpProxy.Run(port)
	}
}

// Override http methods
type Override struct{ PNBase }

// Plugin for Override
func (pn *Override) Plugin() PNRet {
	return func(router *gin.RouterGroup, httpProxy *gin.Engine) {
		handleContext := func(c *gin.Context) {
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
		}
		httpProxy.Use(handleContext)
		router.Use(handleContext)
	}
}
