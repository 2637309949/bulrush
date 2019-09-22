// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"testing"
)

func Test_fixedPortPrefix(t *testing.T) {
	type args struct {
		port string
		plus []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		struct {
			name string
			args args
			want string
		}{
			name: "Test_fixedPortPrefix",
			args: struct {
				port string
				plus []int
			}{
				port: "8080",
				plus: []int{1},
			},
			want: ":8081",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixedPortPrefix(tt.args.port, tt.args.plus...); got != tt.want {
				t.Errorf("fixedPortPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

type test1 struct {
}

type test2 struct {
}

func (t *test2) Plugin() {

}

type test3 struct {
}

func (t test3) Plugin() {

}

func Test_isPlugin(t *testing.T) {

	type args struct {
		src interface{}
	}
	tests := []struct {
		name   string
		args   args
		wantIs bool
	}{
		struct {
			name   string
			args   args
			wantIs bool
		}{
			name: "Test_isPlugin",
			args: struct {
				src interface{}
			}{
				src: func() {},
			},
			wantIs: true,
		},
		struct {
			name   string
			args   args
			wantIs bool
		}{
			name: "Test_isPlugin",
			args: struct {
				src interface{}
			}{
				src: "test",
			},
			wantIs: false,
		},
		struct {
			name   string
			args   args
			wantIs bool
		}{
			name: "Test_isPlugin",
			args: struct {
				src interface{}
			}{
				src: &test1{},
			},
			wantIs: false,
		},
		struct {
			name   string
			args   args
			wantIs bool
		}{
			name: "Test_isPlugin",
			args: struct {
				src interface{}
			}{
				src: &test2{},
			},
			wantIs: true,
		},
		struct {
			name   string
			args   args
			wantIs bool
		}{
			name: "Test_isPlugin",
			args: struct {
				src interface{}
			}{
				src: &test3{},
			},
			wantIs: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIs := isPlugin(tt.args.src); gotIs != tt.wantIs {
				t.Errorf("isPlugin() = %v, want %v", gotIs, tt.wantIs)
			}
		})
	}
}
