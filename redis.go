package bulrush

import (
	"time"
	"encoding/json"
	"github.com/go-redis/redis"
)

// rdsHooks -
type  rdsHooks struct {
	Client  *redis.Client
}

// Rds some common function
type Rds struct {
	Client  *redis.Client
	Hooks   *rdsHooks
	config 	*WellCfg
}

// NewRds -
func NewRds(config *WellCfg) *Rds{
	client  := obClient(config)
	rds     := &Rds {
		Client: client,
		Hooks:	&rdsHooks {
			Client: client,
		},
		config: config,
	}
	return rds
}

// obClient -
func obClient(config *WellCfg) *redis.Client{
	addrs := config.getString("redis.addrs", "")
	if addrs != "" {
		options := &redis.Options{}
		options.Addr 	 = addrs
		options.Password = config.getString("redis.opts.password", "")
		options.DB 	     = config.getInt("redis.opts.db", 0)
		client 		    := redis.NewClient(options)
		if _, err := client.Ping().Result(); err != nil {
			panic(err)
		}
		return client
	}
	return nil
}

// SaveToken -
func (hook *rdsHooks)SaveToken(token map[string]interface{}) {
	accessToken, _  := token["accessToken"]
	refreshToken, _ := token["refreshToken"]
	value, _ := json.Marshal(token)
	hook.Client.Set("TOKEN:" + accessToken.(string),  value, 2*24*time.Hour)
	hook.Client.Set("TOKEN:" + refreshToken.(string), value, 5*24*time.Hour)
}

// RevokeToken -
func (hook *rdsHooks)RevokeToken(accessToken string) bool {
	status, err := hook.Client.Del("TOKEN:" + accessToken).Result()
	if err != nil {
		return false
	} else if status != 1 {
		return false
	}
	return true
}

// FindToken -
func (hook *rdsHooks)FindToken(accessToken string, refreshToken string) map[string]interface{}{
	var imapGet map[string]interface{}
	var token string
	if accessToken != "" {
		token = accessToken
	} else if refreshToken != "" {
		token = refreshToken
	}
	value, err := hook.Client.Get("TOKEN:" + token).Result()
	if err != nil {
		return nil
	}
	json.Unmarshal([]byte(value), &imapGet)
	return imapGet
}


