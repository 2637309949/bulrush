package bulrush

import (
	"fmt"
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --httpProxy gin httpProxy, no middles has been used
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type Bulrush struct {
	config 		*WellCfg
	middles 	[]interface{}
	injects 	[]interface{}
}

// httpProxy middles, use gin as proxy
var httpProxy = func() func() *gin.Engine {
	return func() *gin.Engine {
		proxy := gin.New()
		return proxy
	}
}

// httpRouter middles, use gin as proxy
var httpRouter = func() func(httpProxy *gin.Engine, config *WellCfg) *gin.RouterGroup {
	return func(httpProxy *gin.Engine, config *WellCfg) *gin.RouterGroup {
		httpRouter := httpProxy.Group(config.getString("prefix","/api/v1"))
		return httpRouter
	}
}

// listen proxy
var runProxy = func() func(httpProxy *gin.Engine, config *WellCfg) {
	return func(httpProxy *gin.Engine, config *WellCfg) {
		port := config.getString("port",  ":8080")
		if err := httpProxy.Run(port); err != nil {
			panic(err)
		}
	}
}

// user req log
var loggerWithWriter = func (bulrush *Bulrush, LoggerWithWriter func(*Bulrush) gin.HandlerFunc) func(router *gin.RouterGroup) {
	return func(router *gin.RouterGroup) {
		router.Use(LoggerWithWriter(bulrush))
	}
}

// rec system from panic
var recovery = func () func(httpProxy *gin.Engine, router *gin.RouterGroup) {
	return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		httpProxy.Use(gin.Recovery())
		router.Use(gin.Recovery())
	}
}

// New returns a new blank bulrush instance
// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() *Bulrush {
	bulrush := &Bulrush {
		middles: 	make([]interface{}, 0),
		injects: 	make([]interface{}, 0),
	}
	defaultMiddles := []interface{} {
		httpProxy(),
		httpRouter(),
	}
	bulrush.Use(defaultMiddles...)
	return bulrush
}

// Default return a new bulrush with some default middles
// --Recovery middle has been register in httpProxy and user router
// --LoggerWithWriter middles has been register in router for print requester
func Default() *Bulrush {
	bulrush := New()
	bulrush.Use(recovery())
	bulrush.Use(loggerWithWriter(bulrush, LoggerWithWriter))
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *Bulrush) Use(items ...interface{}) *Bulrush {
	plugins := funk.Filter(items, func(x interface{}) bool {
		return reflect.Func == reflect.TypeOf(x).Kind()
	}).([] interface{})
	bulrush.middles = append(bulrush.middles, plugins...)
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
func (bulrush *Bulrush) Inject(items ...interface{}) *Bulrush {
	funk.ForEach(items, func(x interface{}) {
		if exist := typeExists(bulrush.injects, x);exist {
			panic(fmt.Errorf("Item: %s type %s already exist", x, reflect.TypeOf(x)))
		}
		bulrush.injects = append(bulrush.injects, x)
	})
	return bulrush
}

// Run app, something has been done
func (bulrush *Bulrush) Run() {
	lastMiddles := [] interface{} {runProxy()}
	bulrush.Use(lastMiddles...)
	dynamicMethodsCall(bulrush.middles, &bulrush.injects, func(rs interface{}) {
		bulrush.Inject(rs.([] interface{})...)
	})
}