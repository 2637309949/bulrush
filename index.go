package bulrush

import (
	"github.com/2637309949/bulrush/utils"
	"github.com/gin-gonic/gin"
	ldCfg "github.com/olebedev/config"
)

// Bulrush is the framework's instance
type Bulrush struct {
	config 	*ldCfg.Config
	engine 	*gin.Engine
	router  *gin.RouterGroup
	mongo 	*MongoGroup
	injects []func(map[string]interface{})
	middles []gin.HandlerFunc
}

// New returns a new blank bulrush instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - engine
// 	- RedirectFixedPath:      false
// 	- HandleMethodNotAllowed: false
// 	- ForwardedByClientIP:    true
// 	- UseRawPath:             false
// 	- UnescapePathValues:     true
// - mongo
// 	- Session
// 	- Register
// 	- Model
// 	- manifests
func New() *Bulrush {
	var bulrush *Bulrush
	var engine *gin.Engine
	engine = gin.New()
	bulrush = &Bulrush {
		config: 		nil,
		router: 		nil,
		engine: 		engine,
		injects: 		make([]func(map[string]interface{}), 0),
		middles: 		make([]gin.HandlerFunc, 0),
		mongo: &MongoGroup {
			Session: 		nil,
			Register: 		nil,
			Model: 			nil,
			manifests: 		make([]interface{}, 0),
		},
	}
	bulrush.mongo.Register = register(bulrush)
	bulrush.mongo.Model = model(bulrush)
	Mongo = bulrush.mongo
	return bulrush
}

// Use attachs a global middleware to the router
func (bulrush *Bulrush) Use(middles ...gin.HandlerFunc) *Bulrush{
	bulrush.middles = append(bulrush.middles, middles...)
	return bulrush
}

// LoadConfig load config from string path
// - path
// - m
func (bulrush *Bulrush) LoadConfig(path string, m utils.Mode) *Bulrush {
	cfg, err := utils.LoadConfig(path, m)
	if err != nil {
		panic(err)
	}
	bulrush.config = cfg
	return bulrush
}

// Inject inject params to func
func (bulrush *Bulrush) Inject(injects ...func(map[string]interface{})) *Bulrush{
	bulrush.injects = append(bulrush.injects, injects...)
	return bulrush
}

// Run app
func (bulrush *Bulrush) Run() error {
	port, 	_ 	:= bulrush.config.String("port")
	mode, 	_ 	:= bulrush.config.String("mode")
	prefix, _ 	:= bulrush.config.String("prefix")

	port 	= utils.Some(port, ":8080").(string)
	mode 	= utils.Some(mode, "debug").(string)
	prefix 	= utils.Some(prefix, "/api/v1").(string)

	gin.SetMode(mode)
	bulrush.mongo.Session = obtainSession(bulrush.config)
	bulrush.router = bulrush.engine.Group(prefix)

	// middle
	for _, middle := range bulrush.middles {
		bulrush.router.Use(middle)
	}

	// inject
	for _, callback := range bulrush.injects {
		callback(map[string]interface{} {
			"Engine": bulrush.engine,
			"Router": bulrush.router,
			"Mongo":  bulrush.mongo,
		})
	}
	err := bulrush.engine.Run(port)
	return err
}
