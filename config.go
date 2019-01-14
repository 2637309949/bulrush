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

// Config -
type (
	Config struct {
		config.Config
		Path string
	}
)

// NewCfg -
func NewCfg(path string) *Config{
	wellCfg := &Config{ Path: path }
	return utils.LeftSV(wellCfg.loadFile(path)).(*Config)
}

// loadFile -
func (wc *Config) loadFile(path string) (*Config, error) {
	if strings.HasSuffix(wc.Path, ".json") {
		if cfg, err := config.ParseJsonFile(wc.Path); err == nil {
			return &Config{ *cfg, wc.Path }, nil
		}
	} else if strings.HasSuffix(wc.Path, ".yaml") {
		if cfg, err := config.ParseYamlFile(wc.Path); err == nil {
			return &Config{ *cfg, wc.Path }, nil
		}
	}
	return nil, errors.New("unsupported file type")
}

// GetString -
func (wc *Config) GetString(key string, init string) string {
	return utils.Some(utils.LeftV(wc.String(key)), init).(string)
}

// GetInt -
func (wc *Config) GetInt(key string, init int) int {
	return utils.Some(utils.LeftV(wc.Int(key)), init).(int)
}

// GetDurationFromSecInt -
func (wc *Config) GetDurationFromSecInt(key string, init int) time.Duration {
	return time.Duration(wc.GetInt(key, init)) * time.Second
}

// GetDurationFromMinInt -
func (wc *Config) GetDurationFromMinInt(key string, init int) time.Duration {
	return time.Duration(wc.GetInt(key, init)) * time.Minute
}

// GetDurationFromHourInt -
func (wc *Config) GetDurationFromHourInt(key string, init int) time.Duration {
	return time.Duration(wc.GetInt(key, init)) * time.Hour
}

// GetBool -
func (wc *Config) GetBool(key string, init bool) bool {
	return utils.Some(utils.LeftV(wc.Bool(key)), init).(bool)
}

// GetStrList -
func (wc *Config) GetStrList(key string, init []string) []string {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToStrArray(value)
}

// GetListInt -
func (wc *Config) GetListInt(key string, init []int) []int {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToIntArray(value)
}

