package bulrush

import (
		  "errors"
		  "github.com/gin-gonic/gin"
		  "github.com/thoas/go-funk"
)

// all application instance
var bulrushs []*Bulrush
// master *Bulrush export, just for one instance use
var (
	Config 		*WellConfig
	Mongo 		*MongoGroup
	Redis   	*RedisGroup
	Middles 	[]gin.HandlerFunc
	Injects 	[]interface{}
)

// retain instance
func retain(bulrush *Bulrush) {
	bulrushs = append(bulrushs, bulrush)
	if Mongo == nil {
		Mongo = bulrush.mongo
	}
	if Redis == nil {
		Redis = bulrush.redis
	}
	if Middles == nil {
		Middles = bulrush.middles
	}
	if Injects == nil {
		Injects = bulrush.injects
	}
	if Config == nil {
		Config = bulrush.config
	}
}

// Obtain a application
func Obtain(name string) (*Bulrush, error) {
	target := funk.Find(bulrushs, func(item *Bulrush) bool {
		name, _ := item.config.String("name")
		return name == name
	})
	if target == nil {
		return nil, errors.New("no such application")
	}
	return target.(*Bulrush), nil
}
