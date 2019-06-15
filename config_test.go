// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf := LoadConfig("config_test2.yaml")
	t.Log(conf.Version)
}
