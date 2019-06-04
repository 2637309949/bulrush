/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import (
	"errors"
	"strings"
	"time"

	"github.com/olebedev/config"
)

type (
	// Config base on olebedev/config
	// Provide some useful tools
	Config struct {
		*config.Config
		Path string
	}
)

// LoadFile reads a YAML or JSON configuration from the given filename.
func LoadFile(path string) (*Config, error) {
	var readFile func(filename string) (*config.Config, error)
	if strings.HasSuffix(path, ".json") {
		readFile = config.ParseJsonFile
	} else if strings.HasSuffix(path, ".yaml") {
		readFile = config.ParseYamlFile
	}
	if readFile != nil {
		if cfg, err := readFile(path); err == nil {
			return &Config{
				Config: cfg,
				Path:   path,
			}, nil
		}
	}
	return nil, errors.New("unsupported file type")
}

// NewCfg create a Config instance
func NewCfg(path string) *Config {
	cfg, err := LoadFile(path)
	if err != nil {
		panic(err)
	}
	return cfg
}

// GetString get string or from default value
func (cfg *Config) GetString(key string, init string) string {
	return Some(LeftV(cfg.String(key)), init).(string)
}

// GetInt get int or from default value
func (cfg *Config) GetInt(key string, init int) int {
	return Some(LeftV(cfg.Int(key)), init).(int)
}

// GetDurationFromSecInt get duration or from default value
func (cfg *Config) GetDurationFromSecInt(key string, init int) time.Duration {
	return time.Duration(cfg.GetInt(key, init)) * time.Second
}

// GetDurationFromMinInt get duration or from default value
func (cfg *Config) GetDurationFromMinInt(key string, init int) time.Duration {
	return time.Duration(cfg.GetInt(key, init)) * time.Minute
}

// GetDurationFromHourInt get duration or from default value
func (cfg *Config) GetDurationFromHourInt(key string, init int) time.Duration {
	return time.Duration(cfg.GetInt(key, init)) * time.Hour
}

// GetBool get bool or from default value
func (cfg *Config) GetBool(key string, init bool) bool {
	return Some(LeftV(cfg.Bool(key)), init).(bool)
}

// GetStrList get list or from default value
func (cfg *Config) GetStrList(key string, init []string) []string {
	value := LeftV(cfg.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return ToStrArray(value)
}

// GetListInt get list or from default value
func (cfg *Config) GetListInt(key string, init []int) []int {
	value := LeftV(cfg.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return ToIntArray(value)
}
