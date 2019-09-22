// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"testing"
)

func TestErrCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		wantCode uint64
	}{
		struct {
			name     string
			args     args
			wantCode uint64
		}{
			name: "TestErrCode",
			args: struct {
				err error
			}{
				err: ErrPlugin,
			},
			wantCode: ErrPlugin.Code,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCode := ErrCode(tt.args.err); gotCode != tt.wantCode {
				t.Errorf("ErrCode() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}
