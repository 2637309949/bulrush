package bulrush

import (
	"github.com/2637309949/bulrush/utils"
	"errors"
	"strings"
	"github.com/olebedev/config"
)

// WellConfig -
type WellConfig struct {
	config.Config
	Path string
}

// LoadFile -
func (wc *WellConfig) LoadFile(path string) (*WellConfig, error) {
	var (
		jsonSuffix = ".json"
		yamlSuffix = ".yaml"
		ErrUNSupported = errors.New("unsupported file type")
		readFile func(filename string) (*config.Config, error)
	)
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

// getBool -
func (wc *WellConfig) getBool(key string, init bool) bool {
	return utils.Some(utils.LeftV(wc.Bool(key)), init).(bool)
}

// getListStr
func (wc *WellConfig) getStrList(key string, init []string) []string {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToStrArray(value)
}
