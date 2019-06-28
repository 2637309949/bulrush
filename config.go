// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

const (
	// DefaultConfigVersion defined default version
	DefaultConfigVersion = 1.0
	// DefaultConfigName defined default name
	DefaultConfigName = "bulrush"
	// DefaultConfigPrefix defined default prefix
	DefaultConfigPrefix = "/api/v1"
	// DefaultConfigMode defined default mode
	DefaultConfigMode = "debug"
)

type log struct {
	Path string `json:"path" yaml:"path"`
}

// Config bulrush config struct
type Config struct {
	Version     float64 `json:"version" yaml:"version"`
	Name        string  `json:"name" yaml:"name"`
	Prefix      string  `json:"prefix" yaml:"prefix"`
	Port        string  `json:"port" yaml:"port"`
	Mode        string  `json:"mode" yaml:"mode"`
	DuckReflect bool    `json:"duckReflect" yaml:"duckReflect"`
	Log         log
	dataType    string
	data        []byte
}

func (c *Config) version() float64 {
	if c.Version == 0 {
		return DefaultConfigVersion
	}
	return c.Version
}

func (c *Config) name() string {
	if c.Name == "" {
		return DefaultConfigName
	}
	return c.Name
}

func (c *Config) prefix() string {
	if c.Prefix == "" {
		return DefaultConfigPrefix
	}
	return c.Prefix
}

func (c *Config) mode() string {
	if c.Mode == "" {
		return DefaultConfigMode
	}
	return c.Mode
}

func (c *Config) dataTypeByPath(path string) string {
	var dataType string
	if strings.HasSuffix(path, ".json") {
		dataType = "json"
	} else if strings.HasSuffix(path, ".yaml") {
		dataType = "yaml"
	}
	return dataType
}

// LoadConfig loads the bulrush tool configuration.
func LoadConfig(path string) *Config {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(path); err != nil {
		panic(fmt.Errorf("failed to load file %s", err))
	}
	conf := &Config{data: data}
	conf.dataType = conf.dataTypeByPath(path)

	if err := unmarshalByFileType(conf.data, conf, conf.dataType); err != nil {
		panic(fmt.Errorf("failed to parse yaml file type: %v", err))
	}
	return conf
}

// Unmarshal defined Unmarshal type
func (c *Config) Unmarshal(fieldName string, v interface{}) error {
	vType := reflect.TypeOf(v)

	if vType.Kind() == reflect.Ptr || vType.Kind() == reflect.Slice {
		vType = vType.Elem()
	}

	sv := createStruct([]reflect.StructField{
		reflect.StructField{
			Name: strings.Title(fieldName),
			Type: vType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" yaml:"%s"`, fieldName, fieldName)),
		},
	})

	if err := unmarshalByFileType(c.data, sv, c.dataType); err != nil {
		return err
	}

	conf := stealFieldInStruct(strings.Title(fieldName), sv)
	if reflect.TypeOf(v).Kind() == reflect.Ptr || reflect.TypeOf(v).Kind() == reflect.Slice {
		elem := reflect.ValueOf(v).Elem()
		if elem.CanSet() {
			elem.Set(reflect.ValueOf(conf))
			return nil
		}
	}
	return errors.New("can not unmarshal this type")
}
