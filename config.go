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
