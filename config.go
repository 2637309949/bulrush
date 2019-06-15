// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ghodss/yaml"
)

const confVer = 1

type (
	// Config bulrush config struct
	Config struct {
		Version     int    `json:"version" yaml:"version"`
		Name        string `json:"name" yaml:"name"`
		Prefix      string `json:"prefix" yaml:"prefix"`
		Port        string `json:"port" yaml:"port"`
		Mode        string `json:"mode" yaml:"mode"`
		DuckReflect bool   `json:"duckReflect" yaml:"duckReflect"`
		Log         log
		data        []byte
		dataType    string
	}
	log struct {
		Path string `json:"path" yaml:"path"`
	}
)

// Unmarshal defined Unmarshal type
func (c *Config) Unmarshal(v interface{}) error {
	if c.dataType == "json" {
		return json.Unmarshal(c.data, v)
	}
	if c.dataType == "yaml" {
		return yaml.Unmarshal(c.data, v)
	}
	return errors.New("no support")
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

// initConfig defined return a config with default fields
func initConfig() *Config {
	return &Config{
		Version:     1,
		Name:        "bulrush",
		Prefix:      "/api/v1",
		Mode:        "debug",
		DuckReflect: true,
	}
}

// LoadConfig loads the bulrush tool configuration.
// It looks for .yaml or .json in the current path,
// and falls back to default configuration in case not found.
func LoadConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("failed to load file %s", err))
	}
	config := initConfig()
	config.data = data
	if strings.HasSuffix(path, ".json") {
		config.dataType = "json"
		err = json.Unmarshal(data, config)
		if err != nil {
			panic(fmt.Errorf("failed to parse json file: %s", err))
		}
	} else if strings.HasSuffix(path, ".yaml") {
		config.dataType = "yaml"
		err = yaml.Unmarshal(data, config)
		if err != nil {
			panic(fmt.Errorf("failed to parse yaml file: %s", err))
		}
	}
	if config.Version != confVer {
		rushLogger.Warn("please check the latest version of bulrush's configuration file")
	}
	return config
}
