package plugins

import (
	"github.com/gin-gonic/gin"

)

// Recovery -
func Recovery() func(httpProxy *gin.Engine, router *gin.RouterGroup) {
	return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}
}

// LoggerWithWriter -
func LoggerWithWriter(bulrush interface{}, loggerWithWriter func(interface{}) gin.HandlerFunc) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.Use(loggerWithWriter(bulrush))
	}
}
