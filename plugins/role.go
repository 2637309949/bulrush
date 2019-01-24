/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush role plugin]
 */

// Use:
// app.Use(&plugins.Role{
// 	FailureHandler: func(c *gin.Context, action string) {
// 		c.JSON(http.StatusForbidden, gin.H{
// 			"message": 	"No permission to access this api",
// 		})
// 	},
// 	RoleHandler:	func(c *gin.Context, action string) bool {
// 		fmt.Println("action :", action)
// 		return false
// 	},
// })

// app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup, role *plugins.Role) {
// 	router.GET("/bulrushApp", role.Can("super@testPermts"), func (c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": 	testInject,
// 		})
// 	})
// }))

package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/2637309949/bulrush"
)


type (
	// Role for role
	Role struct {
		bulrush.PNBase
		FailureHandler func(*gin.Context, string)
		RoleHandler    func(*gin.Context, string) bool
	}
)

// Can do what
func (role *Role) Can(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		can := role.RoleHandler(c, action)
		if !can {
			role.FailureHandler(c, action)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

// Plugin for role
func (role *Role) Plugin() bulrush.PNRet {
	return func() *Role {
		return role
	}
}
