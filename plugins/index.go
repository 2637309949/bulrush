package plugins

import (
	"github.com/gin-gonic/gin"

)

// Recovery -
// -recovery from panic
func Recovery() func(httpProxy *gin.Engine, router *gin.RouterGroup) {
	return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}
}

// LoggerWithWriter -
// log user req
func LoggerWithWriter(bulrush interface{}, loggerWithWriter func(interface{}) gin.HandlerFunc) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.Use(loggerWithWriter(bulrush))
	}
}

// HTTPRouter -
// return a router
func HTTPRouter(prefix string) func(HTTPProxy *gin.Engine) *gin.RouterGroup {
	return func(HTTPProxy *gin.Engine) *gin.RouterGroup {
		return HTTPProxy.Group(prefix)
	}
}