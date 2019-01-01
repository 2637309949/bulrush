package bulrush

import (
	"fmt"
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/2637309949/bulrush/plugins"
)

// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --httpProxy gin httpProxy, no middles has been used
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type Bulrush struct {
	HTTPProxy 	*gin.Engine
	config 		*WellCfg
	middles 	[]interface{}
	injects 	[]interface{}
}

// New returns a new blank bulrush instance
// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --httpProxy gin httpProxy, no middles has been used
// --router all router that user defined will be hook in a new router
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() *Bulrush {
	HTTPProxy  := gin.New()
	bulrush := &Bulrush {
		config: 	nil,
		HTTPProxy: 	HTTPProxy,
		middles: 	make([]interface{}, 0),
		injects: 	make([]interface{}, 0),
	}
	bulrush.Inject(HTTPProxy)
	// HTTPRouter middles
	HTTPRouter := func(HTTPProxy *gin.Engine, config *WellCfg) *gin.RouterGroup {
		return plugins.HTTPRouter(config.getString("prefix","/api/v1"))(HTTPProxy)
	}
	bulrush.Use(HTTPRouter)
	return bulrush
}

// Default return a new bulrush with some default middles
// --Recovery middle has been register in httpProxy and user router
// --LoggerWithWriter middles has been register in router for print requester
func Default() *Bulrush {
	bulrush := New()
	loggerWithWriter := func (bulrush *Bulrush, LoggerWithWriter func(*Bulrush) gin.HandlerFunc) func(router *gin.RouterGroup) {
		var lowerType interface{}
		lowerType = bulrush
		return plugins.LoggerWithWriter(lowerType, func(c interface{}) gin.HandlerFunc {
			upperType := c.(*Bulrush)
			return LoggerWithWriter(upperType)
		})
	}
	bulrush.Use(plugins.Recovery(), loggerWithWriter(bulrush, LoggerWithWriter))
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *Bulrush) Use(middles ...interface{}) *Bulrush {
	bulrush.middles = append(bulrush.middles, middles...)
	return bulrush
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bulrush *Bulrush) Config(path string) *Bulrush {
	bulrush.config = NewWc(path)
	if mode := bulrush.config.getString("mode",  ""); mode != "" {
		gin.SetMode(mode)
	}
	bulrush.Inject(bulrush.config)
	return bulrush
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush {
	for _, inject := range  injects {
		if typeExists(bulrush.injects, inject) {
			panic(fmt.Errorf("item: %s type %s has exists", inject, reflect.TypeOf(inject)))
		} else {
			bulrush.injects = append(bulrush.injects, inject)
		}
	}
	return bulrush
}

// Run app, something has been done
func (bulrush *Bulrush) Run()  {
	HTTPRun := func(HTTPProxy *gin.Engine, config *WellCfg) {
		port   := config.getString("port",  ":8080")
		if err := HTTPProxy.Run(port);err != nil {
			panic(err)
		}
	}
	bulrush.Use(HTTPRun)
	dynamicMethodsCall(bulrush.middles, &bulrush.injects, func(rs interface{}) {
		bulrush.Inject(rs.([] interface{})...)
	})
}