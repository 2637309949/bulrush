package bulrush

import (
	"fmt"
		  "github.com/go-redis/redis"
	ldCfg "github.com/olebedev/config"
)

// RedisGroup some common function
type RedisGroup struct {
	Client *redis.Client
}

func obtainClient(config *ldCfg.Config) *redis.Client{
	addrs, _ := config.String("redis.addrs")
	if addrs != "" {
		options := &redis.Options{
			Addr:     addrs,
			Password: "",
			DB:       0,
		}
		options.Addr = addrs
		opts, _  := config.Map("redis.opts")
		if item, ok := opts["password"]; ok {
			options.Password = item.(string)
		}
		if item, ok := opts["db"]; ok {
			options.DB = item.(int)
		}
		client := redis.NewClient(options)
		resu, err := client.Ping().Result()
		fmt.Print(resu)
		fmt.Println(err)
		if err != nil {
			panic(err)
		}
		return client
	}
	return nil
}


