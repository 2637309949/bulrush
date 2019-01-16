/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush override plugin]
 */

package plugins

import (
	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
	"strings"
)
type (
	// Override Plugin
	Override struct {
		bulrush.PNBase
	}
)
// Plugin for gin
func (override *Override) Plugin() bulrush.PNRet {
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
