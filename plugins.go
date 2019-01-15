
/**
 * @author [double]
 * @email [2637309949@qq.com]
 * @create date 2019-01-15 09:49:33
 * @modify date 2019-01-15 09:49:33
 * @desc [default plugins for bulrush]
 */

package bulrush

import (
	"github.com/gin-gonic/gin"
)

var (
	// rec system from panic
	recovery = func () func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
			httpProxy.Use(gin.Recovery())
			router.Use(gin.Recovery())
		}
	}
	// gin httpProxy
	// maybe would be other later
	httpProxy = func() func() *gin.Engine {
		return func() *gin.Engine {
			proxy := gin.New()
			return proxy
		}
	}
	// httpRouter middles
	// gin router
	httpRouter = func() func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
		return func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
			httpRouter := httpProxy.Group(config.GetString("prefix","/api/v1"))
			return httpRouter
		}
	}
	// listen proxy
	// call router listen
	runProxy = func(cb func(error, *Config)) func(httpProxy *gin.Engine, config *Config) {
		return func(httpProxy *gin.Engine, config *Config) {
			port := config.GetString("port",  ":8080")
			cb(nil, config)
			err := httpProxy.Run(port)
			cb(err, config)
		}
	}
	// log user req by http
	// save to file and print to console
	loggerWithWriter = func (bulrush *rush, LoggerWithWriter func(*rush) gin.HandlerFunc) func(router *gin.RouterGroup) {
		return func(router *gin.RouterGroup) {
			router.Use(LoggerWithWriter(bulrush))
		}
	}
)