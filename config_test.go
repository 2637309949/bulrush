/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import (
	"path"
	"testing"
	"time"
)

func TestConfig_GetString(t *testing.T) {
	type fields struct {
		Path string
	}
	type args struct {
		key  string
		init string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		struct {
			name   string
			fields fields
			args   args
			want   string
		}{
			name: "GetString",
			fields: fields{
				Path: path.Join(".", "./config_test.yaml"),
			},
			args: args{
				key:  "mode",
				init: "debug",
			},
			want: "debug",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewCfg(tt.fields.Path)
			if got := cfg.GetString(tt.args.key, tt.args.init); got != tt.want {
				t.Errorf("Config.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetDurationFromHourInt(t *testing.T) {
	type fields struct {
		Path string
	}
	type args struct {
		key  string
		init int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
	}{
		struct {
			name   string
			fields fields
			args   args
			want   time.Duration
		}{
			name: "GetDurationFromHourInt",
			fields: fields{
				Path: path.Join(".", "./config_test.yaml"),
			},
			args: args{
				key:  "mongo.opts.timeout",
				init: 10,
			},
			want: time.Duration(60) * time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewCfg(tt.fields.Path)
			if got := cfg.GetDurationFromHourInt(tt.args.key, tt.args.init); got != tt.want {
				t.Errorf("Config.GetDurationFromHourInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
