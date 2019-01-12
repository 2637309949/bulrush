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

// WellCfg -
type WellCfg struct {
	config.Config
	Path string
}

// NewWc -
func NewWc(path string) *WellCfg{
	wellCfg := &WellCfg{ Path: path }
	return utils.LeftSV(wellCfg.LoadFile(path)).(*WellCfg)
}

// LoadFile -
func (wc *WellCfg) LoadFile(path string) (*WellCfg, error) {
	if strings.HasSuffix(wc.Path, ".json") {
		if cfg, err := config.ParseJsonFile(wc.Path); err == nil {
			return &WellCfg{ *cfg, wc.Path }, nil
		}
	} else if strings.HasSuffix(wc.Path, ".yaml") {
		if cfg, err := config.ParseYamlFile(wc.Path); err == nil {
			return &WellCfg{ *cfg, wc.Path }, nil
		}
	}
	return nil, errors.New("unsupported file type")
}

// getString -
func (wc *WellCfg) getString(key string, init string) string {
	return utils.Some(utils.LeftV(wc.String(key)), init).(string)
}

// getInt -
func (wc *WellCfg) getInt(key string, init int) int {
	return utils.Some(utils.LeftV(wc.Int(key)), init).(int)
}

// getDurationFromSecInt -
func (wc *WellCfg) getDurationFromSecInt(key string, init int) time.Duration {
	return time.Duration(wc.getInt(key, init)) * time.Second
}

// getDurationFromMinInt -
func (wc *WellCfg) getDurationFromMinInt(key string, init int) time.Duration {
	return time.Duration(wc.getInt(key, init)) * time.Minute
}

// getDurationFromHourInt -
func (wc *WellCfg) getDurationFromHourInt(key string, init int) time.Duration {
	return time.Duration(wc.getInt(key, init)) * time.Hour
}

// getBool -
func (wc *WellCfg) getBool(key string, init bool) bool {
	return utils.Some(utils.LeftV(wc.Bool(key)), init).(bool)
}

// getListStr -
func (wc *WellCfg) getStrList(key string, init []string) []string {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToStrArray(value)
}

// getListInt -
func (wc *WellCfg) getListInt(key string, init []int) []int {
	value := utils.LeftV(wc.List(key)).([]interface{})
	if value == nil {
		return init
	}
	return utils.ToIntArray(value)
}

