package main

import (
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/2637309949/bulrush"
	"github.com/gin-gonic/gin"
)

// number of middleware
var n, _ = strconv.Atoi(os.Getenv("MW"))

func main() {
	app := bulrush.Default()
	app.Config(path.Join(".", "cfg.yaml"))
	app.Inject("bulrushApp")

	for ; n > 0; n = n - 1 {
		app.Use(bulrush.PNQuick(func(httpProxy *gin.Engine, router *gin.RouterGroup) {
			httpProxy.Use(func(c *gin.Context) {
				c.Next()
			})
		}))
	}
	app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": testInject,
			})
		})
	}))
	app.RunImmediately()
}
