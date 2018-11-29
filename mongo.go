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

// obtainDialInfo -
func obtainDialInfo(config *WellConfig) *mgo.DialInfo{
	addrs    := utils.LeftV(config.List("mongo.addrs")).([]interface{})
	opts	 := utils.LeftV(config.Map("mongo.opts")).(map[string]interface{})
	dialInfo := &mgo.DialInfo{
		Addrs: utils.ToStrArray(addrs),
	}
	dialInfo.Timeout  		 = time.Duration(utils.Some(utils.LeftOkV(opts["timeout"]), 0).(int)) * time.Second
	dialInfo.Database 		 = utils.Some(utils.LeftOkV(opts["database"]), "").(string)
	dialInfo.ReplicaSetName  = utils.Some(utils.LeftOkV(opts["replicaSetName"]), "").(string)
	dialInfo.Source     	 = utils.Some(utils.LeftOkV(opts["source"]), "").(string)
	dialInfo.Service     	 = utils.Some(utils.LeftOkV(opts["service"]), "").(string)
	dialInfo.ServiceHost     = utils.Some(utils.LeftOkV(opts["serviceHost"]), "").(string)
	dialInfo.Mechanism    	 = utils.Some(utils.LeftOkV(opts["mechanism"]), "").(string)
	dialInfo.Username    	 = utils.Some(utils.LeftOkV(opts["username"]), "").(string)
	dialInfo.Password   	 = utils.Some(utils.LeftOkV(opts["password"]), "").(string)
	dialInfo.PoolLimit 	 	 = utils.Some(utils.LeftOkV(opts["poolLimit"]), 0).(int)
	dialInfo.PoolTimeout 	 = time.Duration(utils.Some(utils.LeftOkV(opts["poolTimeout"]), 0).(int)) * time.Second
	dialInfo.ReadTimeout 	 = time.Duration(utils.Some(utils.LeftOkV(opts["readTimeout"]), 0).(int)) * time.Second
	dialInfo.WriteTimeout 	 = time.Duration(utils.Some(utils.LeftOkV(opts["writeTimeout"]), 0).(int)) * time.Second
	dialInfo.AppName    	 = utils.Some(utils.LeftOkV(opts["appName"]), "").(string)
	dialInfo.FailFast    	 = utils.Some(utils.LeftOkV(opts["failFast"]), false).(bool)
	dialInfo.Direct    		 = utils.Some(utils.LeftOkV(opts["direct"]), false).(bool)
	dialInfo.MinPoolSize 	 = utils.Some(utils.LeftOkV(opts["minPoolSize"]), 0).(int)
	dialInfo.MaxIdleTimeMS 	 = utils.Some(utils.LeftOkV(opts["maxIdleTimeMS"]), 0).(int)
	return dialInfo
}

// obtainSession -
func obtainSession(config *WellConfig) *mgo.Session{
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
			var match map[string]interface{}
			Model, _ := bulrush.mongo.Model(name)
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
			err = query.All(list)
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
