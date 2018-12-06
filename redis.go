package bulrush

import (
	"time"
	"encoding/json"
	"github.com/go-redis/redis"
)

// redisHooks -
type  redisHooks struct {
	SaveToken    func(token map[string]interface{})
	RevokeToken  func(accessToken string) bool
	FindToken    func(accessToken string, refreshToken string) map[string]interface{}
}

// RedisGroup some common function
type RedisGroup struct {
	Client  *redis.Client
	Hooks   redisHooks
}

// obtainClient -
func obtainClient(config *WellConfig) *redis.Client{
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
func SaveToken(bulrush *Bulrush) func (token map[string]interface{}){
	return func (token map[string]interface{}) {
		accessToken, _  := token["accessToken"]
		refreshToken, _ := token["refreshToken"]
		value, _ := json.Marshal(token)
		bulrush.redis.Client.Set("TOKEN:" + accessToken.(string),  value, 2*24*time.Hour)
		bulrush.redis.Client.Set("TOKEN:" + refreshToken.(string), value, 5*24*time.Hour)
	}
}

// RevokeToken -
func RevokeToken(bulrush *Bulrush) func (accessToken string) bool {
	return func (accessToken string) bool {
		status, err := bulrush.redis.Client.Del("TOKEN:" + accessToken).Result()
		if err != nil {
			return false
		} else if status != 1 {
			return false
		}
		return true
	}
}

// FindToken -
func FindToken(bulrush *Bulrush) func(accessToken string, refreshToken string) map[string]interface{}{
	return func(accessToken string, refreshToken string) map[string]interface{} {
		var imapGet map[string]interface{}
		var token string
		if accessToken != "" {
			token = accessToken
		} else if refreshToken != "" {
			token = refreshToken
		}
		value, err := bulrush.redis.Client.Get("TOKEN:" + token).Result()
		json.Unmarshal([]byte(value), &imapGet)
		if err != nil {
			return nil
		}
		return imapGet
	}
}


