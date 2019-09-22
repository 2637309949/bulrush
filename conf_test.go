// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		struct {
			name string
			args args
			want *Config
		}{
			name: "load cfg",
			args: struct {
				path string
			}{
				path: "conf_test_cfg.yaml",
			},
			want: &Config{Name: "bulrush-test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoadConfig(tt.args.path); !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
