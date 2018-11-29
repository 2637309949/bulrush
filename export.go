package bulrush

import (
		  "errors"
		  "github.com/gin-gonic/gin"
		  "github.com/thoas/go-funk"
)
var bulrushs []*Bulrush

// master *Bulrush export, just for one instance use
var (
	Config 	*WellConfig
	Mongo 	*MongoGroup
	Redis   *RedisGroup
	Middles []gin.HandlerFunc
	Injects []interface{}
)

// remainInstance instance
func remainInstance(b *Bulrush) {
	bulrushs = append(bulrushs, b)
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
