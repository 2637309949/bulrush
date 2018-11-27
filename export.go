package bulrush

import (
		  "errors"
		  "github.com/gin-gonic/gin"
		  "github.com/thoas/go-funk"
	ldCfg "github.com/olebedev/config"
)
var bulrushs []*Bulrush

// Master *Bulrush export, just for one instance use
var (
	Config 	*ldCfg.Config
	Mongo 	*MongoGroup
	Redis   *RedisGroup
	Middles []gin.HandlerFunc
	Injects []func(map[string]interface{})
)

// append instance
func appendInstance(b *Bulrush) {
	bulrushs = append(bulrushs, b)
}

// Obtain a application
func Obtain(name string) *Bulrush {
	target := funk.Find(bulrushs, func(item *Bulrush) bool {
		name, _ := item.config.String("name")
		return name == name
	})
	if target == nil {
		panic(errors.New("no such application"))
	}
	return target.(*Bulrush)
}
