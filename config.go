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
	// Config -
	Config struct {
		*config.Config
		Path string
	}
)

// LoadFile reads a YAML or JSON configuration from the given filename.
func LoadFile(path string) (*Config, error) {
	if strings.HasSuffix(path, ".json") {
		if cfg, err := config.ParseJsonFile(path); err == nil {
			return &Config{
				Config: cfg,
				Path:   path,
			}, nil
		}
	} else if strings.HasSuffix(path, ".yaml") {
		if cfg, err := config.ParseYamlFile(path); err == nil {
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
	cfg, _ := LoadFile(path)
	return cfg
}

// GetString -
func (cfg *Config) GetString(key string, init string) string {
	return Some(LeftV(cfg.String(key)), init).(string)
}

// GetInt -
func (cfg *Config) GetInt(key string, init int) int {
	return Some(LeftV(cfg.Int(key)), init).(int)
}

// GetDurationFromSecInt -
func (cfg *Config) GetDurationFromSecInt(key string, init int) time.Duration {
	return time.Duration(cfg.GetInt(key, init)) * time.Second
}

// GetDurationFromMinInt -
func (cfg *Config) GetDurationFromMinInt(key string, init int) time.Duration {
	return time.Duration(cfg.GetInt(key, init)) * time.Minute
}

// GetDurationFromHourInt -
func (cfg *Config) GetDurationFromHourInt(key string, init int) time.Duration {
	return time.Duration(cfg.GetInt(key, init)) * time.Hour
}

// GetBool -
func (cfg *Config) GetBool(key string, init bool) bool {
	return Some(LeftV(cfg.Bool(key)), init).(bool)
}

// GetStrList -
func (cfg *Config) GetStrList(key string, init []string) []string {
	value := LeftV(cfg.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return ToStrArray(value)
}

// GetListInt -
func (cfg *Config) GetListInt(key string, init []int) []int {
	value := LeftV(cfg.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return ToIntArray(value)
}
