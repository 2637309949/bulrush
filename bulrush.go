/**
 * @author [double]
 * @email [2637309949@qq.com]
 * @create date 2019-01-15 09:49:33
 * @modify date 2019-01-15 09:49:33
 * @desc [bulrush implement]
 */

package bulrush

import (
	"log"
	"sync"
	"reflect"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/kataras/go-events"
)

const (
	// Version current version number
	Version = "0.0.1"
	// DefaultMode default gin mode
	DefaultMode = "debug"
	// DefaultMaxPlugins is the number of max Plugins
	// `Just for Learning synchronization`
	// a matter of little interest
	DefaultMaxPlugins = 0
	// EnableWarning prints a warning when trying to add an plugin which it's len is equal to the maxPlugins
	// Defaults to false, which means it does not prints a warning
	EnableWarning = false
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
		Emit(events.EventName, ...interface{})
		On(events.EventName, ...events.Listener)
		SetMaxPlugins(int)
		GetMaxPlugins() int
		Config(string) Bulrush
		Use(...interface{}) Bulrush
		Inject(...interface{}) Bulrush
		Run(func(error, *Config))
	}
	// rush implement bulrush interface
	rush struct {
		Event
		config 		*Config
		middles 	*Middles
		injects 	*Injects
		mu          sync.Mutex
		maxPlugins  int
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
		maxPlugins:   DefaultMaxPlugins,
	}
	defaultMiddles := Middles {
		httpProxy(),
		httpRouter(),
	}
	bulrush.Use(defaultMiddles...)
	return bulrush
}

var (
	// silence the compiler
	_   Bulrush = &rush{}
	// defaultRush default rush
	defaultRush = New()
)

// Default return a new bulrush with some default middles
// --Recovery middle has been register in httpProxy and user router
// --LoggerWithWriter middles has been register in router for print requester
func Default() Bulrush {
	bulrush := defaultRush
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
	if len(items) == 0 {
		return nil
	}
	bulrush.mu.Lock()
	defer bulrush.mu.Unlock()
	if bulrush.maxPlugins > 0 && len(*bulrush.middles) == bulrush.maxPlugins {
		if EnableWarning {
			log.Printf(
				`(events) warning: possible Plugins memory 'leak detected. %d Plugin added. '
				 Use app.SetMaxPlugins(n int) to increase limit.
				 `, len(*bulrush.middles))
		}
		return nil
	}
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
	gin.SetMode(bulrush.config.GetString("mode",  DefaultMode))
	bulrush.Inject(bulrush.config)
	return bulrush
}

// Inject `inject` to bulrush
// - inject should be someone that never be pushed in before.
func (bulrush *rush) Inject(items ...interface{}) Bulrush {
	if len(items) == 0 {
		return nil
	}
	injects := funk.Filter(items, func(x interface{}) bool {
		return !typeExists(*bulrush.injects, x)
	}).([] interface{})
	*bulrush.injects = append(*bulrush.injects, injects...)
	return bulrush
}

// SetMaxPlugins obviously this function allows the MaxPlugins
// to be decrease or increase. Set to zero for unlimited
func (bulrush *rush) SetMaxPlugins(n int) {
	if n < 0 {
		if EnableWarning {
			log.Printf("(events) warning: MaxPlugins must be positive number, tried to set: %d", n)
			return
		}
	}
	bulrush.maxPlugins = n
}

// SetMaxPlugins obviously this function allows the MaxPlugins
// to be decrease or increase. Set to zero for unlimited
func SetMaxPlugins(n int) {
	defaultRush.SetMaxPlugins(n)
}

func (bulrush *rush) GetMaxPlugins() int{
	return bulrush.maxPlugins
}

// GetMaxPlugins returns the max Plugins for this bulrush
// see SetMaxPlugins
func GetMaxPlugins() int {
	return defaultRush.GetMaxPlugins()
}

// Run app, something has been done
func (bulrush *rush) Run(cb func(error, *Config)) {
	lastMiddles := [] interface{} {
		runProxy(cb),
	}
	bulrush.Use(lastMiddles...)
	funk.ForEach(*bulrush.middles, func(x interface{}) {
		rs := reflectMethodAndCall(x, *bulrush.injects)
		bulrush.Inject(rs.([] interface{})...)
	})
}
