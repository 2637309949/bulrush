package bulrush

import (
	"github.com/gin-gonic/gin"
)

// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --engine gin engine, no middles has been used
// --router all router that user defined will be hook in a new router
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type Bulrush struct {
	config 		*WellConfig
	engine 		*gin.Engine
	router  	*gin.RouterGroup
	injects 	[]interface{}
	middles 	[]gin.HandlerFunc
}

// New returns a new blank bulrush instance
// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --engine gin engine, no middles has been used
// --router all router that user defined will be hook in a new router
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() *Bulrush {
	engine  := gin.New()
	bulrush := &Bulrush {
		config: 	nil,
		router: 	nil,
		engine: 	engine,
		injects: 	make([]interface{}, 0),
		middles: 	make([]gin.HandlerFunc, 0),
	}
	return bulrush
}

// Default return a new bulrush with some default middles
// --Recovery middle has been register in engine and user router
// --LoggerWithWriter middles has been register in router for print requester
func Default() *Bulrush {
	bulrush := New()
	bulrush.engine.Use(gin.Recovery())
	bulrush.middles = append(bulrush.middles, gin.Recovery(), LoggerWithWriter(bulrush))
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *Bulrush) Use(middles ...gin.HandlerFunc) *Bulrush {
	bulrush.middles = append(bulrush.middles, middles...)
	return bulrush
}

// Inspect -
// Inspect will be useful if you want to get some params that can not been quote
// by bulrush instance
func (bulrush *Bulrush) Inspect(target interface{}) interface {} {
	return inspectInvoke(target, bulrush)
}

// LoadConfig load config from string path
// currently, it support loading file that end with .json or .yarm
func (bulrush *Bulrush) LoadConfig(path string) *Bulrush {
	bulrush.config = NewWc(path)
	return bulrush
}

// Inject inject params to func
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush {
	bulrush.injects = append(bulrush.injects, injects...)
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

// IsDebugging -
func (bulrush *Bulrush) IsDebugging() bool {
	return IsDebugging()
}

// Run app, something has been done
// -- Init a new mongo session
// -- Init a new Redis Client
// -- Init a new Router
// -- Register middles in gin
// -- Reflect
// -- List on
func (bulrush *Bulrush) Run()  {
	port   := bulrush.config.getString("port",  ":8080")
	mode   := bulrush.config.getString("mode",  "")
	prefix := bulrush.config.getString("prefix","/api/v1")
	if mode != "" {
		bulrush.SetMode(mode)
	}
	bulrush.router 		  = bulrush.engine.Group(prefix)
	routeMiddles(bulrush.router, bulrush.middles)
	injectInvoke(bulrush.injects, bulrush)
	err := bulrush.engine.Run(port)
	if err != nil {
		panic(err)
	}
}
