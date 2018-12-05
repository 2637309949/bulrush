package middles

import (
	"github.com/gin-gonic/gin"
)

// RouteMiddles -
func RouteMiddles(router *gin.RouterGroup, middles []gin.HandlerFunc) {
	for _, middle := range middles {
		router.Use(middle)
	}
}