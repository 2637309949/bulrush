package middles

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// Override method
func Override(r *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.Next()
		} else {
			method := c.PostForm("_method")
			methods := [3]string{"DELETE", "PUT", "PATCH"}
			if method != "" {
				for _, target := range methods {
					if(target == strings.ToUpper(method)) {
						c.Request.Method = target
						r.HandleContext(c)
						break
					}
				}
			}
		}
	}
}