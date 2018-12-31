package bulrush

import (
	"reflect"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --httpProxy gin httpProxy, no middles has been used
// --router all router that user defined will be hook in a new router
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type Bulrush struct {
	httpProxy 		*gin.Engine
	router  	*gin.RouterGroup
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
	httpProxy  := gin.New()
	bulrush := &Bulrush {
		config: 	nil,
		router: 	nil,
		httpProxy: 	httpProxy,
		middles: 	make([]interface{}, 0),
		injects: 	make([]interface{}, 0),
	}
	bulrush.Inject(httpProxy)
	return bulrush
}

// Default return a new bulrush with some default middles
// --Recovery middle has been register in httpProxy and user router
// --LoggerWithWriter middles has been register in router for print requester
func Default() *Bulrush {
	bulrush := New()
	bulrush.middles = append(bulrush.middles, func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}, func(router *gin.RouterGroup) {
		router.Use(LoggerWithWriter(bulrush))
	})
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
	return bulrush
}

// SetMode -
// you should empty mode str in config if you want to set mode in code
// mode will be set in run function again if mode str in config is not empty
func (bulrush *Bulrush) SetMode(mode string) *Bulrush{
	gin.SetMode(mode)
	return bulrush
}

// DebugPrintRouteFunc -
func (bulrush *Bulrush) DebugPrintRouteFunc(handler func(string, string, string, int)) *Bulrush{
	gin.DebugPrintRouteFunc = handler
	return bulrush
}

// Inject -
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush {
	for _, item := range  injects {
		exists := typeExists(bulrush.injects, item)
		if exists {
			panic(fmt.Errorf("item: %s type %s has exists", item, reflect.TypeOf(item)))
		} else {
			bulrush.injects = append(bulrush.injects, item)
		}
	}
	return bulrush
}

// IsDebugging -
func (bulrush *Bulrush) IsDebugging() bool {
	return IsDebugging()
}

// Run app, something has been done
// -- Init a new Router
// -- Register middles in gin
// -- Reflect
// -- List on
func (bulrush *Bulrush) Run()  {
	port   := bulrush.config.getString("port",  ":8080")
	mode   := bulrush.config.getString("mode",  "")
	prefix := bulrush.config.getString("prefix","/api/v1")
	// read configuration first
	if mode != "" {
		bulrush.SetMode(mode)
	}
	// router middle
	bulrush.router = bulrush.httpProxy.Group(prefix)
	bulrush.Inject(bulrush.router, bulrush.config)

	dynamicMethodsCall(bulrush.middles, bulrush.injects)
	err := bulrush.httpProxy.Run(port)
	if err != nil {
		panic(err)
	}
}
