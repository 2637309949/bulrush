package bulrush

import (
	"math"
	"fmt"
	"time"
	"net/url"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/2637309949/bulrush/utils"
	"github.com/globalsign/mgo/bson"
)

type registerHandler func(map[string]interface{})
type modelHandler 	 func(name string) (*mgo.Collection, map[string]interface {})
type hooksHandler    struct {
	List func(name string) func (c *gin.Context)
	One func(name string) func (c *gin.Context)
}

// MongoGroup some common function
type MongoGroup struct {
	Session 	*mgo.Session
	Register 	registerHandler
	Model 		modelHandler
	Hooks 		hooksHandler
	manifests 	[]interface{}
}

func register(bulrush *Bulrush) registerHandler {
	return func(manifest map[string]interface{}) {
		bulrush.mongo.manifests = append(bulrush.mongo.manifests, manifest)
	}
}

func model(bulrush *Bulrush) modelHandler {
	return func(name string) (*mgo.Collection, map[string]interface {}) {
		manifest := utils.Find(bulrush.mongo.manifests, func (item interface{}) bool {
			flag := item.(map [string] interface{})["name"].(string) == name
			return flag
		}).(map[string]interface{})
		if manifest == nil {
			panic(fmt.Errorf("manifest %s not found", name))
		}
		db, ok := manifest["db"]
		if !ok || db == "" {
			db, _ = bulrush.config.String("mongo.opts.database")
		}
		model := bulrush.mongo.Session.DB(db.(string)).C(name)
		return model, manifest
	}
}

// obtainDialInfo -
func obtainDialInfo(config *WellConfig) *mgo.DialInfo {
	addrs    := config.getStrList("mongo.addrs", nil)
	dial := &mgo.DialInfo {}

	dial.Addrs 			 = addrs
	dial.Timeout  		 = time.Duration(config.getInt("mongo.opts.timeout", 0)) * time.Second
	dial.Database 		 = config.getString("mongo.opts.database", "")

	dial.ReplicaSetName  = config.getString("mongo.opts.replicaSetName", "")
	dial.Source     	 = config.getString("mongo.opts.source", "")
	dial.Service     	 = config.getString("mongo.opts.service", "")
	dial.ServiceHost     = config.getString("mongo.opts.serviceHost", "")
	dial.Mechanism    	 = config.getString("mongo.opts.mechanism", "")
	dial.Username    	 = config.getString("mongo.opts.username", "")
	dial.Password   	 = config.getString("mongo.opts.password", "")
	dial.PoolLimit 	 	 = config.getInt("mongo.opts.poolLimit", 0)
	dial.PoolTimeout 	 = time.Duration(config.getInt("mongo.opts.poolTimeout", 0)) * time.Second
	dial.ReadTimeout 	 = time.Duration(config.getInt("mongo.opts.readTimeout", 0)) * time.Second
	dial.WriteTimeout 	 = time.Duration(config.getInt("mongo.opts.writeTimeout", 0)) * time.Second
	dial.AppName    	 = config.getString("mongo.opts.appName", "")
	dial.FailFast    	 = config.getBool("mongo.opts.failFast", false)
	dial.Direct    		 = config.getBool("mongo.opts.direct", false)
	dial.MinPoolSize 	 = config.getInt("mongo.opts.minPoolSize", 0)
	dial.MaxIdleTimeMS 	 = config.getInt("mongo.opts.maxIdleTimeMS", 0)
	return dial
}

// obtainSession -
func obtainSession(config *WellConfig) *mgo.Session {
	addrs, _ := config.List("mongo.addrs")
	if addrs != nil && len(addrs) > 0 {
		dial := obtainDialInfo(config)
		session := utils.LeftSV(mgo.DialWithInfo(dial)).(*mgo.Session)
		return session
	}
	return nil
}

// list -
func list(bulrush *Bulrush) func(string) func (c *gin.Context) {
	return func(name string) func (c *gin.Context) {
		return func (c *gin.Context) {
			var match map[string]interface{}
			Model, manifest := bulrush.mongo.Model(name)
			target := utils.LeftOkV(manifest["reflector"])
			list := createSlice(target)

			cond  	 := c.DefaultQuery("cond", "%7B%7D")
			page, _  := strconv.Atoi(c.DefaultQuery("page", "1"))
			size, _  := strconv.Atoi(c.DefaultQuery("size", "20"))
			_range 	 := c.DefaultQuery("range", "PAGE")
			unescapeCond, err := url.QueryUnescape(cond)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	err.Error(),
				})
				return
			}
			err = json.Unmarshal([]byte(unescapeCond), &match)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	err.Error(),
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
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"data": map[string]interface{}{
					"range": _range,
					"page": page,
					"totalpages": totalpages,
					"size":  size,
					"totalrecords": totalrecords,
					"cond": match,
					"list": list,
				},
				"errcode": 	nil,
				"errmsg": 	nil,
			})
		}
	}
}

// one -
func one(bulrush *Bulrush) func(string) func (c *gin.Context) {
	return func(name string) func (c *gin.Context) {
		return func (c *gin.Context) {
			id := c.Param("id")
			Model, manifest := bulrush.mongo.Model(name)
			target := utils.LeftOkV(manifest["reflector"])
			one := createObject(target)
			isOj := bson.IsObjectIdHex(id)
			if !isOj {
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	"not a valid id",
				})
				return
			}
			err := Model.FindId(bson.ObjectIdHex(id)).One(one)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"data": 	nil,
					"errcode": 	500,
					"errmsg": 	err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"data": one,
				"errcode": 	nil,
				"errmsg": 	nil,
			})
		}
	}
}