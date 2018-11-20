package utils

import (
	"math/rand"
	"fmt"
	"io/ioutil"
	ldCfg "github.com/olebedev/config"
)

// Mode read from json or yaml
type Mode int
const (
	// JSON json mode
	_  Mode = iota + 1
	// JSONMode json mode
	JSONMode
	// YAMLMode yaml mode
	YAMLMode
)


// GetOrElse get or a default value
func GetOrElse(target map[string]interface{}, key string, initValue interface{}) interface{} {
	value, ok := target[key]
	if !ok {
		return initValue
	}
	return value
}

// Some get or a default value
func Some(target interface{}, initValue interface{}) interface{}{
	if target != nil && target != "" {
		return target
	}
	return initValue
}

// LoadConfig load config from string path
// - config
// - error
func LoadConfig(path string, m Mode) (*ldCfg.Config, error) {
	var (
		cfg *ldCfg.Config
		err error
	)
    file, err := ioutil.ReadFile(path)
    if err != nil {
		panic(err)
    }
	buffer := string(file)
	switch m {
		case JSONMode:
			cfg, err = ldCfg.ParseJson(buffer)
		case YAMLMode:
			cfg, err = ldCfg.ParseYaml(buffer)
		default:
			panic(fmt.Errorf("No support this Mode %d", m))
	}
	return cfg, err
}

// Find a element in a array
func Find(arrs []interface{}, matcher func(interface{}) bool) interface{} {
	var target interface{}
	for _, item := range arrs {
		match := matcher(item)
		if match {
			target = item
			break
		}
	}
	return target
}


// ToStrArray convert []interface{} to []string
func ToStrArray(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}
	return s
}

// SafeMap field safely map Map
func SafeMap(source map [string]interface{}, key string, mapFunc func(interface{})) interface{}{
	if source != nil {
		item, ok := source[key]
		if ok  {
			mapFunc(item)
			return item
		}
	}
	return nil
}

// RandString string
func RandString(n int) string {
	const seeds = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = seeds[rand.Intn(len(seeds))]
	}
	return string(bytes)
}
