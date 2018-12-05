package bulrush

import (
	"github.com/go-redis/redis"
)

// RedisGroup some common function
type RedisGroup struct {
	Client *redis.Client
}

// obtainClient -
func obtainClient(config *WellConfig) *redis.Client{
	addrs, _ := config.String("redis.addrs")
	if addrs != "" {
		options := &redis.Options{}
		options.Addr = addrs
		opts, _  := config.Map("redis.opts")
		if item, ok := opts["password"]; ok {
			options.Password = item.(string)
		}
		if item, ok := opts["db"]; ok {
			options.DB = item.(int)
		}
		client := redis.NewClient(options)
		if _, err := client.Ping().Result(); err != nil {
			panic(err)
		}
		return client
	}
	return nil
}

