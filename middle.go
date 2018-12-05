package bulrush

import (
	"github.com/gin-gonic/gin"
)

// routeMiddles -
func routeMiddles(router *gin.RouterGroup, middles []gin.HandlerFunc) {
	for _, middle := range middles {
		router.Use(middle)
	}
}