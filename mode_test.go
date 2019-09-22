// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import "testing"

func TestSetMode(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		struct {
			name string
			args args
		}{
			name: "TestSetMode",
			args: struct {
				value string
			}{
				value: DebugMode,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if SetMode(tt.args.value); modeName != tt.args.value {
				t.Errorf("SetMode() = %v, want %v", modeName, tt.args.value)
			}
		})
	}
}
