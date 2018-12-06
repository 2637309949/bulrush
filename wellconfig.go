package bulrush

import (
	"time"
	"github.com/2637309949/bulrush/utils"
	"errors"
	"strings"
	"github.com/olebedev/config"
)

const (
	jsonSuffix = ".json"
	yamlSuffix = ".yaml"
)

var (
	// ErrUNSupported -
	ErrUNSupported = errors.New("unsupported file type")
)

// WellConfig -
type WellConfig struct {
	config.Config
	Path string
}

// LoadFile -
func (wc *WellConfig) LoadFile(path string) (*WellConfig, error) {
	var readFile func(filename string) (*config.Config, error) 
	if strings.HasSuffix(wc.Path, jsonSuffix) {
		readFile = config.ParseJsonFile
	} else if strings.HasSuffix(wc.Path, yamlSuffix) {
		readFile = config.ParseYamlFile
	} else {
		return nil, ErrUNSupported
	}
	cfg, err := readFile(wc.Path)
	if err != nil {
		return nil, err
	}
	return &WellConfig{ *cfg, wc.Path }, nil
}

// getString -
func (wc *WellConfig) getString(key string, init string) string {
	return utils.Some(utils.LeftV(wc.String(key)), init).(string)
}

// getInt -
func (wc *WellConfig) getInt(key string, init int) int {
	return utils.Some(utils.LeftV(wc.Int(key)), init).(int)
}

// getDurationFromSecInt -
func (wc *WellConfig) getDurationFromSecInt(key string, init int) time.Duration {
	return time.Duration(wc.getInt(key, init)) * time.Second
}

// getDurationFromMinInt -
func (wc *WellConfig) getDurationFromMinInt(key string, init int) time.Duration {
	return time.Duration(wc.getInt(key, init)) * time.Minute
}

// getDurationFromHourInt -
func (wc *WellConfig) getDurationFromHourInt(key string, init int) time.Duration {
	return time.Duration(wc.getInt(key, init)) * time.Hour
}

// getBool -
func (wc *WellConfig) getBool(key string, init bool) bool {
	return utils.Some(utils.LeftV(wc.Bool(key)), init).(bool)
}

// getListStr -
func (wc *WellConfig) getStrList(key string, init []string) []string {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToStrArray(value)
}

// getListInt -
func (wc *WellConfig) getListInt(key string, init []int) []int {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToIntArray(value)
}

