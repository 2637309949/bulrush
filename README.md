# Bulrush Framework

![Bulrush flash](./assets/flash.jpg)
![Bulrush flash](./assets/frame.png)

## Benchmarks
```cmd
Running 3s test @ http://127.0.0.1:3333/api/v1/
8 threads and 50 connections
Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   374.51us  504.86us   8.45ms   89.96%
    Req/Sec    22.71k     4.36k   60.90k    97.94%
549392 requests in 3.10s, 77.02MB read
Requests/sec: 177249.36
Transfer/sec:     24.85MB
```

## Instruction
Quickly build applications and customize special functions through plug-ins, Multiple base plug-ins are provided
Install Bulrush
```shell
$ go get github.com/2637309949/bulrush
```
QuickStart
```go
import (
    "github.com/2637309949/bulrush"
)
// Use attachs a global Recovery middleware to the router
// Use attachs a Recovery middleware to the user router
// Use attachs a LoggerWithWriter middleware to the user router
app := bulrush.Default()
app.Config(CONFIGPATH)
app.Inject("bulrushApp")
app.Use(&models.Model{}, &routes.Route{})
app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
    router.GET("/bulrushApp", func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": 	testInject,
        })
    })
}))
app.Run(func(httpProxy *gin.Engine, config *bulrush.Config) {
    port := config.GetString("port", ":8080")
    port = strings.TrimSpace(port)
    name := config.GetString("name", "")
    if prefix := port[:1]; prefix != ":" {
        port = fmt.Sprintf(":%s", port)
    }
    fmt.Println("\n\n================================")
    fmt.Printf("App: %s\n", name)
    fmt.Printf("Listen on %s\n", port)
    fmt.Println("================================")
    httpProxy.Run(port)
})
```
OR
```go
import (
    "github.com/2637309949/bulrush"
)
// No middlewares has been attached
app := bulrush.New()
app.Config(CONFIGPATH)
app.Run(func(httpProxy *gin.Engine, config *bulrush.Config) {
    httpProxy.Run(config.GetString("port", ":8080"))
})
```
For more details, Please reference to [bulrush-template](https://github.com/2637309949/bulrush-template). 

## API
#### Set app config
```go
app.Config(CONFIGPATH)
```
#### Inject your custom injects
All injects would be provided as plugins params next by next.  
Init injects by Inject function
```go
app.Inject("bulrushApp")
```
Set injects by plugin ret  
```go
// Plugin for role
func (role *Role) Plugin() bulrush.PNRet {
	return func() *Role {
		return role
	}
}
```
#### Import your plugins
```go
app.Use(bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
    router.GET("/bulrushApp", func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": 	testInject,
        })
    })
}))
```
#### Run app
```go
app.Run(func(httpProxy *gin.Engine, config *bulrush.Config) {
    port := config.GetString("port", ":8080")
    port = strings.TrimSpace(port)
    name := config.GetString("name", "")
    if prefix := port[:1]; prefix != ":" {
        port = fmt.Sprintf(":%s", port)
    }
    fmt.Println("\n\n================================")
    fmt.Printf("App: %s\n", name)
    fmt.Printf("Listen on %s\n", port)
    fmt.Println("================================")
    httpProxy.Run(port)
})
```
#### Share state between plug-ins

##### store state
```go
app.Use(bulrush.PNQuick(func(status *bulrush.Status) {
    status.Set("count", 1)
}))
```
##### read state
```go
app.Use(bulrush.PNQuick(func(status *bulrush.Status) {
    status.Get("count")
    status.ALL()
}))
```
#### Plug in communication between plug-ins
```go
app.Use(bulrush.PNQuick(func(events events.EventEmmiter) {
	events.On("hello", func(payload ...interface{}) {
		message := payload[0].(string)
		fmt.Println(message)
	})
}))
```

## Design Philosophy
## Injects
### Built-in Injects
-	EventEmmiter
-	*Status
-	*Validate
-	*Jobrunner
-	*ReverseInject

## Plugins
### Built-in Plugins
- [bulrush-addition](https://github.com/2637309949/bulrush-addition)
- [bulrush-openapi](https://github.com/2637309949/bulrush-openapi)
- [bulrush-captcha](https://github.com/2637309949/bulrush-captcha)
- [bulrush-delivery](https://github.com/2637309949/bulrush-delivery)
- [bulrush-identify](https://github.com/2637309949/bulrush-identify)
- [bulrush-logger](https://github.com/2637309949/bulrush-logger)
- [bulrush-proxy](https://github.com/2637309949/bulrush-proxy)
- [bulrush-role](https://github.com/2637309949/bulrush-role)
- [bulrush-limit](https://github.com/2637309949/bulrush-limit)
- [bulrush-upload](https://github.com/2637309949/bulrush-upload)


### Custom your plugins
If your want to write a user-defined plugins, you should implement PNBase interface or the duck type,
PNRet is a function, and you can get all you want through func parameters, also you can return any type as
`Injects` entity.
```go
PNBase interface{ Plugin() PNRet }
```
EXAMPLE:   
```go
type (
    Override struct { PNBase }
)
func (pn *Override) Plugin() PNRet {
    return func(router *gin.RouterGroup, httpProxy *gin.Engine) {
            return "inject entity"
    }
}
```
OR
```go
bulrush.PNQuick(func(testInject string, router *gin.RouterGroup) {
    router.GET("/test", func (c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": 	testInject,
        })
    })
})

```

### Assemble your plugin` config from bulrush Injects
```go
// Example for my mgo
type conf struct {
	Addrs          []string      `json:"addrs" yaml:"addrs"`
	Timeout        time.Duration `json:"timeout" yaml:"timeout"`
	Database       string        `json:"database" yaml:"database"`
	ReplicaSetName string        `json:"replicaSetName" yaml:"replicaSetName"`
	Source         string        `json:"source" yaml:"source"`
	Service        string        `json:"service" yaml:"service"`
	ServiceHost    string        `json:"serviceHost" yaml:"serviceHost"`
	Mechanism      string        `json:"mechanism" yaml:"mechanism"`
	Username       string        `json:"username" yaml:"username"`
	Password       string        `json:"password" yaml:"password"`
	PoolLimit      int           `json:"poolLimit" yaml:"poolLimit"`
	PoolTimeout    time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
	ReadTimeout    time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	AppName        string        `json:"appName" yaml:"appName"`
	FailFast       bool          `json:"failFast" yaml:"failFast"`
	Direct         bool          `json:"direct" yaml:"direct"`
	MinPoolSize    int           `json:"minPoolSize" yaml:"minPoolSize"`
	MaxIdleTimeMS  int           `json:"maxIdleTimeMS" yaml:"maxIdleTimeMS"`
}
func New(bulCfg *bulrush.Config) *Mongo {
	cf, err := bulCfg.Unmarshal("mongo", conf{})
	if err != nil {
		panic(err)
	}
	conf := cf.(conf)
	session := createSession(&conf)
	mgo := &Mongo{
		m:       make([]map[string]interface{}, 0),
		cfg:     &conf,
		API:     &api{},
		Session: session,
	}
	mgo.API.mgo = mgo
	mgo.AutoHook = autoHook(mgo)
	return mgo
}
``` 
    // Read part and assemble
    func (c *Config) Unmarshal(fieldName string, v interface{}) (interface{}, error)

## Note
    Note go vendor, bulrush needs to reference the same package, otherwise injection fails

## MIT License

Copyright (c) 2018-2020 Double

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.