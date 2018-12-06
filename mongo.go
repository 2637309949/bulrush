package bulrush

import (
	"math"
	"fmt"
	"time"
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

type registerHandler func(map[string]interface{})
type modelHandler 	 func(name string) (*mgo.Collection, map[string]interface {})
type mgoHooks    struct {
	One  func(name string) func (c *gin.Context)
	List func(name string) func (c *gin.Context)
}

// MongoGroup -
type MongoGroup struct {
	Session 	*mgo.Session
	Register 	registerHandler
	Model 		modelHandler
	Hooks 		mgoHooks
	manifests 	[]interface{}
}

// register -
func register(bulrush *Bulrush) registerHandler {
	return func(manifest map[string]interface{}) {
		var ok = true
		_, ok = manifest["name"]
		_, ok = manifest["reflector"]
		if !ok {
			panic(errors.New("name and reflector params must be provided"))
		}
		bulrush.mongo.manifests = append(bulrush.mongo.manifests, manifest)
	}
}

func model(bulrush *Bulrush) modelHandler {
	return func(name string) (*mgo.Collection, map[string]interface {}) {
		var db string
		var collect string
		manifest := utils.Find(bulrush.mongo.manifests, func (item interface{}) bool {
			flag := item.(map [string] interface{})["name"].(string) == name
			return flag
		}).(map[string]interface{})
		if manifest == nil {
			panic(fmt.Errorf("manifest %s not found", name))
		}

		if dbName, ok := manifest["db"]; ok && dbName.(string) != "" {
			db = dbName.(string)
		} else {
			db = bulrush.config.getString("mongo.opts.database", "bulrush")
		}
		
		if ctName, ok := manifest["collection"]; ok && ctName.(string) != "" {
			collect = ctName.(string)
		} else {
			collect = name
		}
		model := bulrush.mongo.Session.DB(db).C(collect)
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
	dial.PoolTimeout 	 = config.getDurationFromSecInt("mongo.opts.poolTimeout", 0)
	dial.ReadTimeout 	 = config.getDurationFromSecInt("mongo.opts.readTimeout", 0)
	dial.WriteTimeout 	 = config.getDurationFromSecInt("mongo.opts.writeTimeout", 0)
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

// List -
func List(bulrush *Bulrush) func(string) func (c *gin.Context) {
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

// One -
func One(bulrush *Bulrush) func(string) func (c *gin.Context) {
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