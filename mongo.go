package bulrush

import (
		"math"
		"strconv"
		"fmt"
		"time"
		"net/url"
		"net/http"
		"encoding/json"
		"github.com/gin-gonic/gin"
		"github.com/globalsign/mgo"
  ldCfg "github.com/olebedev/config"
		"github.com/2637309949/bulrush/utils"
)

type registerHandler func(map[string]interface{})
type modelHandler 	 func(name string) (*mgo.Collection, map[string]interface {})
type hooksHandler    struct {
	List func(name string, list interface{}) func (c *gin.Context)
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

func model(bulrush *Bulrush) modelHandler{
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
		model 	  := bulrush.mongo.Session.DB(db.(string)).C(name)
		return model, manifest
	}
}

func obtainDialInfo(config *ldCfg.Config) *mgo.DialInfo{
	addrs, _ := config.List("mongo.addrs")
	opts, _  := config.Map("mongo.opts")
	dialInfo := &mgo.DialInfo{
		Addrs: utils.ToStrArray(addrs),
	}
	if item, ok := opts["timeout"]; ok {
		dialInfo.Timeout = time.Duration(item.(int)) * time.Second
	}
	if item, ok := opts["database"]; ok {
		dialInfo.Database = item.(string)
	}
	if item, ok := opts["replicaSetName"]; ok {
		dialInfo.ReplicaSetName = item.(string)
	}
	if item, ok := opts["source"]; ok {
		dialInfo.Source = item.(string)
	}
	if item, ok := opts["service"]; ok {
		dialInfo.Service = item.(string)
	}
	if item, ok := opts["serviceHost"]; ok {
		dialInfo.ServiceHost = item.(string)
	}
	if item, ok := opts["mechanism"]; ok {
		dialInfo.Mechanism = item.(string)
	}
	if item, ok := opts["username"]; ok {
		dialInfo.Username = item.(string)
	}
	if item, ok := opts["password"]; ok {
		dialInfo.Password = item.(string)
	}
	if item, ok := opts["poolLimit"]; ok {
		dialInfo.PoolLimit = item.(int)
	}
	if item, ok := opts["poolTimeout"]; ok {
		dialInfo.PoolTimeout = time.Duration(item.(int)) * time.Second
	}
	if item, ok := opts["readTimeout"]; ok {
		dialInfo.ReadTimeout = time.Duration(item.(int)) * time.Second
	}
	if item, ok := opts["writeTimeout"]; ok {
		dialInfo.WriteTimeout = time.Duration(item.(int)) * time.Second
	}
	if item, ok := opts["appName"]; ok {
		dialInfo.AppName = item.(string)
	}
	if item, ok := opts["failFast"]; ok {
		dialInfo.FailFast = item.(bool)
	}
	if item, ok := opts["direct"]; ok {
		dialInfo.Direct = item.(bool)
	}
	if item, ok := opts["minPoolSize"]; ok {
		dialInfo.MinPoolSize = item.(int)
	}
	if item, ok := opts["maxIdleTimeMS"]; ok {
		dialInfo.MaxIdleTimeMS = item.(int)
	}
	return dialInfo
}

func obtainSession(config *ldCfg.Config) *mgo.Session{
	addrs, _ := config.List("mongo.addrs")
	if addrs != nil && len(addrs) > 0 {
		dialInfo := obtainDialInfo(config)
		session, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
		return session
	}
	return nil
}

func list(bulrush *Bulrush) func(string,interface{}) func (c *gin.Context) {
	return func(name string, list interface{}) func (c *gin.Context) {
		return func (c *gin.Context) {
			Model, _ := bulrush.mongo.Model(name)
			var match map[string]interface{}
			cond  := c.DefaultQuery("cond", "%7B%7D")

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
			totalpages := math.Ceil(float64(totalrecords) / float64(size))
			err = query.All(list)
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
