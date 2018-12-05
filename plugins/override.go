package plugins

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// Override -
type Override struct {
}

// Inject for gin
func (over *Override) Inject(injects map[string]interface{}) {
	engine, _ := injects["Engine"].(*gin.Engine)
	router, _ := injects["Router"].(*gin.RouterGroup)
	engine.Use(func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.Next()
		} else {
			method := c.PostForm("_method")
			methods := [3]string{"DELETE", "PUT", "PATCH"}
			if method != "" {
				for _, target := range methods {
					if(target == strings.ToUpper(method)) {
						c.Request.Method = target
						engine.HandleContext(c)
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
						engine.HandleContext(c)
						break
					}
				}
			}
		}
	})
}
