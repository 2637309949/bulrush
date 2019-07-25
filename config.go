// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config bulrush config struct
type Config struct {
	Version     float64 `json:"version" yaml:"version"`
	Name        string  `json:"name" yaml:"name"`
	Prefix      string  `json:"prefix" yaml:"prefix"`
	Port        string  `json:"port" yaml:"port"`
	Mode        string  `json:"mode" yaml:"mode"`
	DuckReflect bool    `json:"duckReflect" yaml:"duckReflect"`
	Log         struct {
		Level string `json:"level" yaml:"level"`
		Path  string `json:"path" yaml:"path"`
	}
	dataType string
	data     []byte
}

// LoadConfig loads the bulrush tool configuration
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

func (c *Config) version() float64 {
	if c.Version == 0 {
		return 1.0
	}
	return c.Version
}

func (c *Config) name() string {
	if c.Name == "" {
		return "bulrush"
	}
	return c.Name
}

func (c *Config) prefix() string {
	if c.Prefix == "" {
		return "/api/v1"
	}
	return c.Prefix
}

func (c *Config) mode() string {
	if c.Mode == "" {
		return "debug"
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

// unmarshalByFileType defined unmarshal by diff file type
func unmarshalByFileType(data []byte, v interface{}, dataType string) error {
	switch true {
	case dataType == "json":
		err := json.Unmarshal(data, v)
		if err != nil {
			return err
		}
	case dataType == "yaml":
		err := yaml.Unmarshal(data, v)
		if err != nil {
			return err
		}
	}
	return nil
}
