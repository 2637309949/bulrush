package bulrush

import (
	"log"
	"errors"
	"strings"
	"github.com/2637309949/bulrush/utils"
	"github.com/2637309949/bulrush/middles"
	"github.com/olebedev/config"
	"github.com/gin-gonic/gin"
)

// WellConfig -
type WellConfig struct {
	config.Config
	Path string
}

// LoadFile -
func (wc *WellConfig) LoadFile(path string) (*WellConfig, error) {
	var (
		jsonSuffix = ".json"
		yamlSuffix = ".yaml"
		ErrUNSupported = errors.New("unsupported file type")
		readFile func(filename string) (*config.Config, error)
	)
	if strings.HasSuffix(wc.Path, jsonSuffix) {
		readFile = config.ParseJsonFile
	} else if strings.HasSuffix(wc.Path, yamlSuffix) {
		readFile = config.ParseYamlFile
	} else {
		return nil, ErrUNSupported
	}
	cfg, err := readFile(wc.Path)
	if err != nil {
		return nil, err
	}
	return &WellConfig{ *cfg, wc.Path }, nil
}

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
	bulrush.middles = append(bulrush.middles, LoggerWithWriter(bulrush))
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
	port   := utils.Some(utils.LeftV(bulrush.config.String("port")), 	":8080").(string)
	mode   := utils.Some(utils.LeftV(bulrush.config.String("mode")), 	"debug").(string)
	prefix := utils.Some(utils.LeftV(bulrush.config.String("prefix")),  "/api/v1").(string)

	gin.SetMode(mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("%5v %9v\n", httpMethod, absolutePath)
	}

	bulrush.mongo.Session = obtainSession(bulrush.config)
	bulrush.redis.Client  = obtainClient(bulrush.config)
	bulrush.router 		  = bulrush.engine.Group(prefix)

	middles.RouteMiddles(bulrush.router, bulrush.middles)
	injectInvoke(bulrush.injects, bulrush)
	err := bulrush.engine.Run(port)
	if err != nil {
		bulrush.mongo.Session.Close()
		panic(err)
	}
}
