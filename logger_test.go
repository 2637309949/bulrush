// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
	"testing"

	addition "github.com/2637309949/bulrush-addition"
	"github.com/2637309949/bulrush-addition/logger"
)

func TestSetLogger(t *testing.T) {
	type args struct {
		logger *logger.Journal
	}
	tests := []struct {
		name string
		args args
	}{
		struct {
			name string
			args args
		}{
			name: "TestSetLogger",
			args: struct {
				logger *logger.Journal
			}{
				logger: addition.RushLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if SetLogger(tt.args.logger); !reflect.DeepEqual(rushLogger, tt.args.logger) {
				t.Errorf("SetMode() = %v, want %v", modeName, tt.args.logger)
			}
		})
	}
}
