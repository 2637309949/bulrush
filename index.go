package bulrush

import (
	"github.com/2637309949/bulrush/utils"
	"github.com/olebedev/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"fmt"
)

// Mode read from json or yaml
type Mode int
const (
	// JSON json mode
	_  Mode = iota + 1
	// JSONMode json mode
	JSONMode
	// YAMLMode yaml mode
	YAMLMode
)

// WellConfig -
type WellConfig struct {
	config.Config
	Path string
	Mode Mode
}

// LoadFile -
func (wc *WellConfig) LoadFile() *WellConfig {
	var (
		cfg *config.Config
		err error
	)
    file, err := ioutil.ReadFile(wc.Path)
    if err != nil {
		panic(err)
    }
	buffer := string(file)
	switch wc.Mode {
		case JSONMode:
			cfg, err = config.ParseJson(buffer)
		case YAMLMode:
			cfg, err = config.ParseYaml(buffer)
		default:
			panic(fmt.Errorf("No support this Mode %d",wc.Mode))
	}
    if err != nil {
		panic(err)
    }
	wellCfg := &WellConfig{ *cfg, wc.Path, wc.Mode }
	return wellCfg
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
func (bulrush *Bulrush) LoadConfig(path string, m Mode) *Bulrush {
	wc := &WellConfig{ Path: path, Mode: m }
	bulrush.config = wc.LoadFile()
	return bulrush
}

// Inject inject params to func
func (bulrush *Bulrush) Inject(injects ...interface{}) *Bulrush{
	bulrush.injects = append(bulrush.injects, injects...)
	return bulrush
}

// Run app
func (bulrush *Bulrush) Run() error {
	port, 	_ 	:= bulrush.config.String("port")
	mode, 	_ 	:= bulrush.config.String("mode")
	prefix, _ 	:= bulrush.config.String("prefix")

	port 	= utils.Some(port, 	 ":8080").(string)
	mode 	= utils.Some(mode, 	 "debug").(string)
	prefix 	= utils.Some(prefix, "/api/v1").(string)

	gin.SetMode(mode)
	bulrush.mongo.Session = obtainSession(bulrush.config)
	bulrush.redis.Client  = obtainClient(bulrush.config)
	bulrush.router = bulrush.engine.Group(prefix)
	// middle
	for _, middle := range bulrush.middles {
		bulrush.router.Use(middle)
	}
	// inject
	for _, target := range bulrush.injects {
		invoke(target, map[string]interface{} {
			"Engine": bulrush.engine,
			"Router": bulrush.router,
			"Mongo":  bulrush.mongo,
			"Config": bulrush.config,
			"Redis":  bulrush.redis,
		})
	}
	err := bulrush.engine.Run(port)
	return err
}
