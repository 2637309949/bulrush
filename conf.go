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

	utils "github.com/2637309949/bulrush-utils"
	"gopkg.in/yaml.v2"
)

type (
	// Config bulrush config struct
	Config struct {
		Version float64 `json:"version" yaml:"version"`
		Name    string  `json:"name" yaml:"name"`
		Prefix  string  `json:"prefix" yaml:"prefix"`
		Port    string  `json:"port" yaml:"port"`
		Mode    string  `json:"mode" yaml:"mode"`
		suffix  string
		data    []byte
		Log     struct {
			Level string `json:"level" yaml:"level"`
			Path  string `json:"path" yaml:"path"`
		}
	}
)

func (c *Config) verifyVersion(version float64) {
	if c.Version != version {
		rushLogger.Warn("Please check the latest version of bulrush's configuration file")
	}
}

func (c *Config) version() float64 {
	return utils.Some(c.Version, 1.0).(float64)
}

func (c *Config) name() string {
	return utils.Some(c.Name, "bulrush").(string)
}

func (c *Config) prefix() string {
	return utils.Some(c.Prefix, "/api").(string)
}

func (c *Config) mode() string {
	return utils.Some(c.Mode, "debug").(string)
}

func (c *Config) typeBySuffix(path string) string {
	var dataType string
	if strings.HasSuffix(path, ".json") {
		dataType = "json"
	} else if strings.HasSuffix(path, ".yaml") {
		dataType = "yaml"
	}
	return dataType
}

// Unmarshal defined Unmarshal type
func (c *Config) Unmarshal(field string, v interface{}) (err error) {
	vType := reflect.TypeOf(v)
	if vType.Kind() == reflect.Ptr || vType.Kind() == reflect.Slice {
		vType = vType.Elem()
	}
	sv := createStruct([]reflect.StructField{
		reflect.StructField{
			Name: strings.Title(field),
			Type: vType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" yaml:"%s"`, field, field)),
		},
	})
	if err = c.UnmarshalByType(c.data, sv, c.suffix); err != nil {
		return
	}
	conf := stealFieldInStruct(strings.Title(field), sv)
	if reflect.TypeOf(v).Kind() == reflect.Ptr || reflect.TypeOf(v).Kind() == reflect.Slice {
		elem := reflect.ValueOf(v).Elem()
		if elem.CanSet() {
			elem.Set(reflect.ValueOf(conf))
			return
		}
		err = fmt.Errorf("elem %v can not been set", elem)
	}
	return errors.New("can not unmarshal this type")
}

// LoadConfig loads the bulrush tool configuration
func LoadConfig(path string) *Config {
	var (
		err  error
		data []byte
	)
	if data, err = ioutil.ReadFile(path); err == nil {
		c := &Config{data: data}
		c.suffix = c.typeBySuffix(path)
		err = c.UnmarshalByType(c.data, c, c.suffix)
		return c
	}
	panic(fmt.Errorf("failed to load file %s", err))
}

// UnmarshalByType defined unmarshal by diff file type
// suport
// , : json
// , : yaml
func (c *Config) UnmarshalByType(data []byte, v interface{}, dataType string) (err error) {
	switch true {
	case dataType == "json":
		err = json.Unmarshal(data, v)
		if err != nil {
			return err
		}
	case dataType == "yaml":
		err = yaml.Unmarshal(data, v)
		if err != nil {
			return err
		}
	}
	return
}
