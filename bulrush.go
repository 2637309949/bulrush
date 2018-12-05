package bulrush

import (
	"log"
	"github.com/2637309949/bulrush/utils"
	"github.com/gin-gonic/gin"
)

// Bulrush is the framework's instance
type Bulrush struct {
	config 		*WellConfig
	engine 		*gin.Engine
	router  	*gin.RouterGroup
	mongo 		*MongoGroup
	redis   	*RedisGroup
	injects 	[]interface{}
	middles 	[]gin.HandlerFunc
}

// New returns a new blank bulrush instance
func New() *Bulrush {
	var (
		engine  *gin.Engine
		bulrush *Bulrush
	)
	engine  = gin.New()
	bulrush = &Bulrush {
		config: 	nil,
		router: 	nil,
		engine: 	engine,
		injects: 	make([]interface{}, 0),
		middles: 	make([]gin.HandlerFunc, 0),
		mongo: &MongoGroup {
			Session: 	nil,
			Register: 	nil,
			Model: 		nil,
			manifests: 	make([]interface{}, 0),
		},
		redis: &RedisGroup {
			Client:		nil,
		},
	}
	bulrush.mongo.Register   = register(bulrush)
	bulrush.mongo.Model 	 = model(bulrush)

	bulrush.mongo.Hooks.List = list(bulrush)
	bulrush.mongo.Hooks.One  = one(bulrush)
	retain(bulrush)
	return bulrush
}

// Default return a new bulrush with some default middles
func Default() *Bulrush {
	bulrush := New()
	bulrush.engine.Use(gin.Recovery())
	bulrush.middles = append(bulrush.middles, gin.Recovery(), LoggerWithWriter(bulrush))
	return bulrush
}

// Use attachs a global middleware to the router
func (bulrush *Bulrush) Use(middles ...gin.HandlerFunc) *Bulrush{
	bulrush.middles = append(bulrush.middles, middles...)
	return bulrush
}

// Inspect -
func (bulrush *Bulrush) Inspect(target interface{}) interface {} {
	return inspectInvoke(target, bulrush)
}

// LoadConfig load config from string path
func (bulrush *Bulrush) LoadConfig(path string) *Bulrush {
	wc := &WellConfig{ Path: path }
	bulrush.config = utils.LeftSV(wc.LoadFile(path)).(*WellConfig)
	return bulrush
}

// Inject inject params to func
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush{
	bulrush.injects = append(bulrush.injects, injects...)
	return bulrush
}

// Run app
func (bulrush *Bulrush) Run()  {
	port   := bulrush.config.getString("port",  ":8080")
	mode   := bulrush.config.getString("mode",  "debug")
	prefix := bulrush.config.getString("prefix","/api/v1")

	gin.SetMode(mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("%5v %9v\n", httpMethod, absolutePath)
	}

	bulrush.mongo.Session = obtainSession(bulrush.config)
	bulrush.redis.Client  = obtainClient(bulrush.config)
	bulrush.router 		  = bulrush.engine.Group(prefix)

	routeMiddles(bulrush.router, bulrush.middles)
	injectInvoke(bulrush.injects, bulrush)
	err := bulrush.engine.Run(port)
	if err != nil {
		bulrush.mongo.Session.Close()
		panic(err)
	}
}
