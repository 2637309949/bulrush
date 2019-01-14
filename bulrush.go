package bulrush

import (
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/kataras/go-events"
)

// Bulrush the framework's struct
// --EventEmmiter emit and on
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
type (
	// Event -
	Event events.EventEmmiter
	// Middles -
	Middles []interface{}
	// Injects -
	Injects []interface{}
	// Bulrush interface defined
	Bulrush interface {
		Config(string) Bulrush
		Use(...interface{}) Bulrush
		Inject(...interface{}) Bulrush
		Run(func(error, *Config))
	}
	// rush Implement Bulrush interface
	rush struct {
		Event
		config 		*Config
		middles 	*Middles
		injects 	*Injects
	}
)

const (
	// DebugMode -
	DebugMode = "debug"
	// ReleaseMode -
	ReleaseMode = "release"
	// TestMode -
	TestMode = "test"
)

var (
	// gin httpProxy
	// maybe would be other later
	httpProxy = func() func() *gin.Engine {
		return func() *gin.Engine {
			proxy := gin.New()
			return proxy
		}
	}
	// httpRouter middles
	// gin router
	httpRouter = func() func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
		return func(httpProxy *gin.Engine, config *Config) *gin.RouterGroup {
			httpRouter := httpProxy.Group(config.GetString("prefix","/api/v1"))
			return httpRouter
		}
	}
	// listen proxy
	// call router listen
	runProxy = func(cb func(error, *Config)) func(httpProxy *gin.Engine, config *Config) {
		return func(httpProxy *gin.Engine, config *Config) {
			port := config.GetString("port",  ":8080")
			cb(nil, config)
			err := httpProxy.Run(port)
			cb(err, config)
		}
	}
	// log user req by http
	// save to file and print to console
	loggerWithWriter = func (bulrush *rush, LoggerWithWriter func(*rush) gin.HandlerFunc) func(router *gin.RouterGroup) {
		return func(router *gin.RouterGroup) {
			router.Use(LoggerWithWriter(bulrush))
		}
	}
	// rec system from panic
	recovery = func () func(httpProxy *gin.Engine, router *gin.RouterGroup) {
		return func(httpProxy *gin.Engine, router *gin.RouterGroup) {
			httpProxy.Use(gin.Recovery())
			router.Use(gin.Recovery())
		}
	}
)

// New returns a new blank bulrush instance
// Bulrush is the framework's instance
// --config json or yaml config for bulrush
// --injects struct instance can be reflect by bulrush
// --middles some middles for gin self
func New() Bulrush {
	middles := make(Middles, 0)
	injects := make(Injects, 0)
	emmiter := events.New()
	bulrush := &rush {
		Event: 		  emmiter,
		middles: 	  &middles,
		injects: 	  &injects,
	}
	defaultMiddles := Middles {
		httpProxy(),
		httpRouter(),
	}
	bulrush.Use(defaultMiddles...)
	return bulrush
}

// Default return a new bulrush with some default middles
// --Recovery middle has been register in httpProxy and user router
// --LoggerWithWriter middles has been register in router for print requester
func Default() Bulrush {
	bulrush := New()
	defaultMiddles := Middles {
		recovery(),
		loggerWithWriter(bulrush.(*rush), LoggerWithWriter),
	}
	bulrush.Use(defaultMiddles...)
	return bulrush
}

// Use attachs a global middleware to the router
// just like function in gin, but not been inited util bulrush inited.
// bulrush range these middles in order
func (bulrush *rush) Use(items ...interface{}) Bulrush {
	plugins := funk.Filter(items, func(x interface{}) bool {
		return reflect.Func == reflect.TypeOf(x).Kind()
	}).([] interface{})
	*bulrush.middles = append(*bulrush.middles, plugins...)
	return bulrush
}

// Config load config from string path
// currently, it support loading file that end with .json or .yarm
func (bulrush *rush) Config(path string) Bulrush {
	bulrush.config = NewCfg(path)
	if mode := bulrush.config.GetString("mode",  ""); mode != "" {
		gin.SetMode(mode)
	} else {
		gin.SetMode(DebugMode)
	}
	bulrush.Inject(bulrush.config)
	return bulrush
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bulrush *rush) Inject(items ...interface{}) Bulrush {
	injects := funk.Filter(items, func(x interface{}) bool {
		return !typeExists(*bulrush.injects, x)
	}).([] interface{})
	*bulrush.injects = append(*bulrush.injects, injects...)
	return bulrush
}

// Run app, something has been done
func (bulrush *rush) Run(cb func(error, *Config)) {
	lastMiddles := [] interface{} {
		runProxy(cb),
	}
	bulrush.Use(lastMiddles...)
	dynamicMethodsCall(*bulrush.middles, bulrush.injects, func(rs interface{}) {
		bulrush.Inject(rs.([] interface{})...)
	})
}