package bulrush

import (
	"math"
	"fmt"
	"errors"
	"net/url"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/2637309949/bulrush/utils"
	"github.com/globalsign/mgo/bson"
)

// mgoHooks -
type mgoHooks struct {
	One  func(name string) func (c *gin.Context)
	List func(name string) func (c *gin.Context)
}

// Mgo -
type Mgo struct {
	Session 	*mgo.Session
	Hooks 		*mgoHooks
	config 		*Config
	manifests 	[]interface{}
}

// NewMgo -
func NewMgo(config *Config) *Mgo{
	session := obSession(config)
	mgo		:= &Mgo{
		Hooks: 	   &mgoHooks{},
		Session:   session,
		manifests: make([]interface{}, 0),
		config:    config,
	}
	mgo.Hooks.List = list(mgo)
	mgo.Hooks.One = one(mgo)
	return mgo
}

// Register -
func (mgo *Mgo)Register(manifest map[string]interface{}) {
	var ok = true
	_, ok = manifest["name"]
	_, ok = manifest["reflector"]
	if !ok {
		panic(errors.New("name and reflector params must be provided"))
	}
	mgo.manifests = append(mgo.manifests, manifest)
}

// Model -
func (mgo *Mgo)Model(name string) (*mgo.Collection, map[string]interface {}) {
	var db string
	var collect string
	manifest := utils.Find(mgo.manifests, func (item interface{}) bool {
		flag := item.(map [string] interface{})["name"].(string) == name
		return flag
	}).(map[string]interface{})
	if manifest == nil {
		panic(fmt.Errorf("manifest %s not found", name))
	}

	if dbName, ok := manifest["db"]; ok && dbName.(string) != "" {
		db = dbName.(string)
	} else {
		db = mgo.config.GetString("mongo.opts.database", "bulrush")
	}
	
	if ctName, ok := manifest["collection"]; ok && ctName.(string) != "" {
		collect = ctName.(string)
	} else {
		collect = name
	}
	model := mgo.Session.DB(db).C(collect)
	return model, manifest
}

// obDialInfo -
func obDialInfo(config *Config) *mgo.DialInfo {
	addrs    := config.GetStrList("mongo.addrs", nil)
	dial := &mgo.DialInfo {}
	dial.Addrs 			 = addrs
	dial.Timeout  		 = config.GetDurationFromSecInt("mongo.opts.timeout", 0)
	dial.Database 		 = config.GetString("mongo.opts.database", "")
	dial.ReplicaSetName  = config.GetString("mongo.opts.replicaSetName", "")
	dial.Source     	 = config.GetString("mongo.opts.source", "")
	dial.Service     	 = config.GetString("mongo.opts.service", "")
	dial.ServiceHost     = config.GetString("mongo.opts.serviceHost", "")
	dial.Mechanism    	 = config.GetString("mongo.opts.mechanism", "")
	dial.Username    	 = config.GetString("mongo.opts.username", "")
	dial.Password   	 = config.GetString("mongo.opts.password", "")
	dial.PoolLimit 	 	 = config.GetInt("mongo.opts.poolLimit", 0)
	dial.PoolTimeout 	 = config.GetDurationFromSecInt("mongo.opts.poolTimeout", 0)
	dial.ReadTimeout 	 = config.GetDurationFromSecInt("mongo.opts.readTimeout", 0)
	dial.WriteTimeout 	 = config.GetDurationFromSecInt("mongo.opts.writeTimeout", 0)
	dial.AppName    	 = config.GetString("mongo.opts.appName", "")
	dial.FailFast    	 = config.GetBool("mongo.opts.failFast", false)
	dial.Direct    		 = config.GetBool("mongo.opts.direct", false)
	dial.MinPoolSize 	 = config.GetInt("mongo.opts.minPoolSize", 0)
	dial.MaxIdleTimeMS 	 = config.GetInt("mongo.opts.maxIdleTimeMS", 0)
	return dial
}

// obSession -
func obSession(config *Config) *mgo.Session {
	addrs, _ := config.List("mongo.addrs")
	if addrs != nil && len(addrs) > 0 {
		dial := obDialInfo(config)
		session := utils.LeftSV(mgo.DialWithInfo(dial)).(*mgo.Session)
		return session
	}
	return nil
}

// list -
func list(mgo *Mgo) func(string) func (c *gin.Context) {
	return func(name string) func (c *gin.Context) {
		return func (c *gin.Context) {
			var match map[string]interface{}
			Model, manifest := mgo.Model(name)
			target := utils.LeftOkV(manifest["reflector"])
			list := createSlice(target)

			cond  	 := c.DefaultQuery("cond", "%7B%7D")
			page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
			size, _  := strconv.Atoi(c.DefaultQuery("size", "20"))
			_range 	 := c.DefaultQuery("range", "PAGE")
			unescapeCond, err := url.QueryUnescape(cond)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": 	err.Error(),
				})
				return
			}
			err = json.Unmarshal([]byte(unescapeCond), &match)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": 	err.Error(),
				})
				return
			}
			// return mapping body， not json in db
			query := Model.Find(match)
			totalrecords, _ := query.Count()
			if _range != "ALL" {
				query = query.Skip((page - 1) * size).Limit(size)
			}
			err = query.All(&list)
			totalpages := math.Ceil(float64(totalrecords) / float64(size))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": 	err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"range": _range,
				"page": page,
				"totalpages": totalpages,
				"size":  size,
				"totalrecords": totalrecords,
				"cond": match,
				"list": list,
			})
		}
	}
}

// one -
func one(mgo *Mgo) func(string) func (c *gin.Context) {
	return func(name string) func (c *gin.Context) {
		return func (c *gin.Context) {
			id := c.Param("id")
			Model, manifest := mgo.Model(name)
			target := utils.LeftOkV(manifest["reflector"])
			one := createObject(target)
			isOj := bson.IsObjectIdHex(id)
			if !isOj {
				c.JSON(http.StatusNotAcceptable, gin.H{
					"message": 	"not a valid id",
				})
				return
			}
			err := Model.FindId(bson.ObjectIdHex(id)).One(one)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": 	err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, one)
		}
	}
}

// create
func create(mgo *Mgo) func(string) func (c *gin.Context) {
	return func(name string) func (c *gin.Context) {
		return func (c *gin.Context) {

		}
	}
}