package bulrush

import (
	"fmt"
	"time"
	"github.com/globalsign/mgo"
	ldCfg "github.com/olebedev/config"
	"github.com/2637309949/bulrush/utils"
)

// Mode read from json or yaml
type registerHandle func(map[string]interface{})
type modelHandle func(name string) (*mgo.Collection, map[string]interface {})

// MongoGroup some common function
type MongoGroup struct {
	Session *mgo.Session
	Register registerHandle
	Model modelHandle
	manifests []interface{}
}

// register mongo type
func register(bulrush *Bulrush) registerHandle {
	return func(manifest map[string]interface{}) {
		bulrush.mongo.manifests = append(bulrush.mongo.manifests, manifest)
	}
}

// model return register model manifest
func model(bulrush *Bulrush) modelHandle {
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
		Schema, _ := manifest["schema"].(map[string] interface{})
		Model := bulrush.mongo.Session.DB(db.(string)).C(name)
		return Model, Schema
	}
}

// obtainDialInfo obtain dial info
// - config
func obtainDialInfo(config *ldCfg.Config) *mgo.DialInfo{
	addrs, _ := config.List("mongo.addrs")
	opts, _  := config.Map("mongo.opts")
	dialInfo := &mgo.DialInfo{
		Addrs: utils.ToStrArray(addrs),
	}
	utils.SafeMap(opts, "timeout", func(timeout interface{}) {
		dialInfo.Timeout = time.Duration(timeout.(int)) * time.Second
	})
	utils.SafeMap(opts, "database", func(database interface{}) {
		dialInfo.Database = database.(string)
	})
	utils.SafeMap(opts, "replicaSetName", func(replicaSetName interface{}) {
		dialInfo.ReplicaSetName = replicaSetName.(string)
	})
	utils.SafeMap(opts, "source", func(source interface{}) {
		dialInfo.Source = source.(string)
	})
	utils.SafeMap(opts, "service", func(service interface{}) {
		dialInfo.Service = service.(string)
	})
	utils.SafeMap(opts, "serviceHost", func(serviceHost interface{}) {
		dialInfo.ServiceHost = serviceHost.(string)
	})
	utils.SafeMap(opts, "mechanism", func(mechanism interface{}) {
		dialInfo.Mechanism = mechanism.(string)
	})
	utils.SafeMap(opts, "username", func(username interface{}) {
		dialInfo.Username = username.(string)
	})
	utils.SafeMap(opts, "password", func(password interface{}) {
		dialInfo.Password = password.(string)
	})
	utils.SafeMap(opts, "poolLimit", func(poolLimit interface{}) {
		dialInfo.PoolLimit = poolLimit.(int)
	})
	utils.SafeMap(opts, "poolTimeout", func(poolTimeout interface{}) {
		dialInfo.PoolTimeout = time.Duration(poolTimeout.(int)) * time.Second
	})
	utils.SafeMap(opts, "readTimeout", func(readTimeout interface{}) {
		dialInfo.ReadTimeout = time.Duration(readTimeout.(int)) * time.Second
	})
	utils.SafeMap(opts, "writeTimeout", func(writeTimeout interface{}) {
		dialInfo.WriteTimeout = time.Duration(writeTimeout.(int)) * time.Second
	})
	utils.SafeMap(opts, "appName", func(appName interface{}) {
		dialInfo.AppName = appName.(string)
	})
	utils.SafeMap(opts, "failFast", func(failFast interface{}) {
		dialInfo.FailFast = failFast.(bool)
	})
	utils.SafeMap(opts, "direct", func(direct interface{}) {
		dialInfo.Direct = direct.(bool)
	})
	utils.SafeMap(opts, "minPoolSize", func(minPoolSize interface{}) {
		dialInfo.MinPoolSize = minPoolSize.(int)
	})
	utils.SafeMap(opts, "maxIdleTimeMS", func(maxIdleTimeMS interface{}) {
		dialInfo.MaxIdleTimeMS = maxIdleTimeMS.(int)
	})
	return dialInfo
}

// obtainSession
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

// Mongo export
var Mongo *MongoGroup
