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
	PNBase interface {
		Plugin() PNRet
	}
)

// PNStruct for a quickly Plugin SetUp when you dont want declare PNBase
// PBBase minimize implement
type PNStruct struct{ Quick interface{} }

// Plugin for PNQuick
func (pn *PNStruct) Plugin() PNRet {
	return pn.Quick
}

// PNQuick for PNQuick
func PNQuick(quick interface{}) PNBase {
	pn := &PNStruct{}
	pn.Quick = quick
	return pn
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

// Override http methods
type Override struct{ PNBase }

// Plugin for Override
func (pn *Override) Plugin() PNRet {
	return func(router *gin.RouterGroup, httpProxy *gin.Engine) {
		hc := func(c *gin.Context) {
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
		router.Use(hc)
		httpProxy.Use(hc)
	}
}
