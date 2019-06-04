/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf := LoadConfig("config_test2.yaml")
	t.Log(conf.Mongo.ReadTimeout)
}
