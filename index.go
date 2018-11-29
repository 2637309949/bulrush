package bulrush

import (
	"log"
	"errors"
	"strings"
	"github.com/2637309949/bulrush/utils"
	"github.com/olebedev/config"
	"github.com/gin-gonic/gin"
)

// WellConfig -
type WellConfig struct {
	config.Config
	Path string
}

// LoadFile -
func (wc *WellConfig) LoadFile() (*WellConfig, error) {
	if strings.HasSuffix(wc.Path, ".json") {
		cfg, err := config.ParseJsonFile(wc.Path)
		if err != nil {
			return nil, err
		}
		return &WellConfig{ *cfg, wc.Path }, nil
	} else if strings.HasSuffix(wc.Path, ".yaml") {
		cfg, err := config.ParseYamlFile(wc.Path)
		if err != nil {
			return nil, err
		}
		return &WellConfig{ *cfg, wc.Path }, nil
	} else {
		return nil, errors.New("unsupported file type")
	}
}

// Bulrush is the framework's instance
type Bulrush struct {
	config 	*WellConfig
	engine 	*gin.Engine
	router  *gin.RouterGroup
	mongo 	*MongoGroup
	redis   *RedisGroup
	injects []interface{}
	middles []gin.HandlerFunc
}

// New returns a new blank bulrush instance
func New() *Bulrush {
	var bulrush *Bulrush
	var engine *gin.Engine
	engine = gin.New()
	bulrush = &Bulrush {
		config: 		nil,
		router: 		nil,
		engine: 		engine,
		injects: 		make([]interface{}, 0),
		middles: 		make([]gin.HandlerFunc, 0),
		mongo: &MongoGroup {
			Session: 		nil,
			Register: 		nil,
			Model: 			nil,
			manifests: 		make([]interface{}, 0),
		},
		redis: &RedisGroup {
			Client:			nil,
		},
	}
	bulrush.mongo.Register   = register(bulrush)
	bulrush.mongo.Model 	 = model(bulrush)
	bulrush.mongo.Hooks.List = list(bulrush)

	Mongo 	= bulrush.mongo
	Redis	= bulrush.redis
	Middles =  bulrush.middles
	Injects = bulrush.injects
	Config 	= bulrush.config
	remainInstance(bulrush)
	return bulrush
}

// Use attachs a global middleware to the router
func (bulrush *Bulrush) Use(middles ...gin.HandlerFunc) *Bulrush{
	bulrush.middles = append(bulrush.middles, middles...)
	return bulrush
}

// LoadConfig load config from string path
func (bulrush *Bulrush) LoadConfig(path string) *Bulrush {
	wc := &WellConfig{ Path: path }
	bulrush.config = utils.LeftSV(wc.LoadFile()).(*WellConfig)
	return bulrush
}

// Inject inject params to func
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush{
	bulrush.injects = append(bulrush.injects, injects...)
	return bulrush
}

// Run app
func (bulrush *Bulrush) Run() error {
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
	for _, middle := range bulrush.middles {
		bulrush.router.Use(middle)
	}
	injectInvoke(bulrush.injects, bulrush)
	err := bulrush.engine.Run(port)
	return err
}
