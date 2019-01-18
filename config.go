/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import (
	"time"
	"errors"
	"strings"
	"github.com/olebedev/config"
	"github.com/2637309949/bulrush/utils"
)

type (
	// Config -
	Config struct {
		*config.Config
		Path string
	}
)

// LoadFile -
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

// NewCfg -
func NewCfg(path string) *Config {
	cfg, _ := LoadFile(path)
	return cfg
}

// GetString -
func (cfg *Config) GetString(key string, init string) string {
	return utils.Some(utils.LeftV(cfg.String(key)), init).(string)
}

// GetInt -
func (cfg *Config) GetInt(key string, init int) int {
	return utils.Some(utils.LeftV(cfg.Int(key)), init).(int)
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
	return utils.Some(utils.LeftV(cfg.Bool(key)), init).(bool)
}

// GetStrList -
func (cfg *Config) GetStrList(key string, init []string) []string {
	value := utils.LeftV(cfg.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToStrArray(value)
}

// GetListInt -
func (cfg *Config) GetListInt(key string, init []int) []int {
	value := utils.LeftV(cfg.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToIntArray(value)
}

