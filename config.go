/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ghodss/yaml"
)

const confVer = 1

// Config bulrush config struct
type Config struct {
	Version     int    `json:"version" yaml:"version"`
	Name        string `json:"name" yaml:"name"`
	Prefix      string `json:"prefix" yaml:"prefix"`
	Port        string `json:"port" yaml:"port"`
	Mode        string `json:"mode" yaml:"mode"`
	DuckReflect bool   `json:"duckReflect" yaml:"duckReflect"`
	Log         log
	Mongo       mongo
	Redis       redis
}

type log struct {
	Path string `json:"path" yaml:"path"`
}

type mongo struct {
	Addrs          []string      `json:"addrs" yaml:"addrs"`
	Timeout        time.Duration `json:"timeout" yaml:"timeout"`
	Database       string        `json:"database" yaml:"database"`
	ReplicaSetName string        `json:"replicaSetName" yaml:"replicaSetName"`
	Source         string        `json:"source" yaml:"source"`
	Service        string        `json:"service" yaml:"service"`
	ServiceHost    string        `json:"serviceHost" yaml:"serviceHost"`
	Mechanism      string        `json:"mechanism" yaml:"mechanism"`
	Username       string        `json:"username" yaml:"username"`
	Password       string        `json:"password" yaml:"password"`
	PoolLimit      int           `json:"poolLimit" yaml:"poolLimit"`
	PoolTimeout    time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
	ReadTimeout    time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout   time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	AppName        string        `json:"appName" yaml:"appName"`
	FailFast       bool          `json:"failFast" yaml:"failFast"`
	Direct         bool          `json:"direct" yaml:"direct"`
	MinPoolSize    int           `json:"minPoolSize" yaml:"minPoolSize"`
	MaxIdleTimeMS  int           `json:"maxIdleTimeMS" yaml:"maxIdleTimeMS"`
}

type redis struct {
	Network            string        `json:"network" yaml:"network"`
	Addr               string        `json:"addrs" yaml:"addrs"`
	Password           string        `json:"password" yaml:"password"`
	DB                 int           `json:"db" yaml:"db"`
	MaxRetries         int           `json:"maxRetries" yaml:"maxRetries"`
	MinRetryBackoff    time.Duration `json:"minRetryBackoff" yaml:"minRetryBackoff"`
	MaxRetryBackoff    time.Duration `json:"maxRetryBackoff" yaml:"maxRetryBackoff"`
	DialTimeout        time.Duration `json:"dialTimeout" yaml:"dialTimeout"`
	ReadTimeout        time.Duration `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout       time.Duration `json:"writeTimeout" yaml:"writeTimeout"`
	PoolSize           int           `json:"poolSize" yaml:"poolSize"`
	MinIdleConns       int           `json:"minIdleConns" yaml:"minIdleConns"`
	MaxConnAge         time.Duration `json:"maxConnAge" yaml:"maxConnAge"`
	PoolTimeout        time.Duration `json:"poolTimeout" yaml:"poolTimeout"`
	IdleTimeout        time.Duration `json:"idleTimeout" yaml:"idleTimeout"`
	IdleCheckFrequency time.Duration `json:"idleCheckFrequency" yaml:"idleCheckFrequency"`
}

func initConfig() *Config {
	return &Config{
		Version:     1,
		Name:        "bulrush",
		Prefix:      "/api/v1",
		Mode:        "debug",
		DuckReflect: true,
		Mongo: mongo{
			Addrs:    make([]string, 0),
			Timeout:  0,
			Database: "bulrush",
		},
		Redis: redis{
			Password: "",
			DB:       0,
		},
	}
}

// LoadConfig loads the bulrush tool configuration.
// It looks for .yaml or .json in the current path,
// and falls back to default configuration in case not found.
func LoadConfig(path string) *Config {
	conf := initConfig()
	if strings.HasSuffix(path, ".json") {
		err := parseJSON(path, &conf)
		if err != nil {
			panic(fmt.Errorf("Failed to parse JSON file: %s", err))
		}
	} else if strings.HasSuffix(path, ".yaml") {
		err := parseYAML(path, &conf)
		if err != nil {
			panic(fmt.Errorf("Failed to parse YAML file: %s", err))
		}
	}
	// Check format version
	if conf.Version != confVer {
		rushLogger.Warn("Check the latest version of bulrush's configuration file.")
	}
	return conf
}

func parseJSON(path string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	return err
}

func parseYAML(path string, v interface{}) error {
	var (
		data []byte
		err  error
	)
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, v)
	return err
}
