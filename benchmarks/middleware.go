package main

import (
	"os"
	"fmt"
	"path"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/2637309949/bulrush"
)

// number of middleware
var n, _ = strconv.Atoi(os.Getenv("MW"))
var useAsync, _ = strconv.ParseBool(os.Getenv("USE_ASYNC"))

func main() {
  app := bulrush.Default()
  app.Config(path.Join(".", "cfg.yaml"))
  app.Inject("bulrushApp")

  for ;n > 0; n = n -1 {
    app.Use(bulrush.PNQuick(func (httpProxy *gin.Engine, router *gin.RouterGroup) {
      httpProxy.Use(func (c *gin.Context) {
        c.Next()
      })
		}))
  }
  app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
    router.GET("/", func (c *gin.Context) {
      c.JSON(http.StatusOK, gin.H{
        "data": 	testInject,
        "errcode": 	nil,
        "errmsg": 	nil,
      })
    })
  }))
  app.Run(func(err error, config *bulrush.Config) {
    if err != nil {
      panic(err)
    } else {
      name := config.GetString("name",  "")
      port := config.GetString("port",  "")
      fmt.Println("================================")
      fmt.Printf("App: %s\n", name)
      fmt.Printf("Listen on %s\n", port)
      fmt.Println("================================")
    }
  })
}
