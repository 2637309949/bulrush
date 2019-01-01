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
// --router all router that user defined will be hook in a new router
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type Bulrush struct {
	HTTPProxy 	*gin.Engine
	config 		*WellConfig
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
	bulrush.middles = append(bulrush.middles, func(HTTPProxy *gin.Engine, config *WellConfig) *gin.RouterGroup {
		prefix := config.getString("prefix","/api/v1")
		return HTTPProxy.Group(prefix)
	})
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
	bulrush.middles = append(bulrush.middles, plugins.Recovery(), loggerWithWriter(bulrush, LoggerWithWriter))
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
	bulrush.Inject(bulrush.config)
	return bulrush
}

// Inject `inject` to bulrush
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush {
	for _, inject := range  injects {
		exists := typeExists(bulrush.injects, inject)
		if exists {
			panic(fmt.Errorf("item: %s type %s has exists", inject, reflect.TypeOf(inject)))
		} else {
			bulrush.injects = append(bulrush.injects, inject)
		}
	}
	return bulrush
}

// Run app, something has been done
// -- Init a new Router
// -- Register middles in gin
// -- Reflect
// -- List on
func (bulrush *Bulrush) Run()  {
	port   := bulrush.config.getString("port",  ":8080")
	mode   := bulrush.config.getString("mode",  "")
	// read configuration first
	if mode != "" {
		SetMode(mode)
	}
	dynamicMethodsCall(bulrush.middles, &bulrush.injects, func(rs interface{}) {
		bulrush.Inject(rs.([] interface{})...)
	})
	err := bulrush.HTTPProxy.Run(port)
	if err != nil {
		panic(err)
	}
}

// DebugPrintRouteFunc function in gin
func DebugPrintRouteFunc(handler func(string, string, string, int)) {
	gin.DebugPrintRouteFunc = handler
}

// SetMode function in gin
// you should empty mode str in config if you want to set mode in code
// mode will be set in run function again if mode str in config is not empty
func SetMode(mode string) {
	gin.SetMode(mode)
}
