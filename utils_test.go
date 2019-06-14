// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
	"testing"
)

func TestSome(t *testing.T) {
	type args struct {
		target    interface{}
		initValue interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		struct {
			name string
			args args
			want interface{}
		}{
			name: "Some",
			args: args{
				target:    "test",
				initValue: "_test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Some(tt.args.target, tt.args.initValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Some() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		arrs    []interface{}
		matcher func(interface{}) bool
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		struct {
			name string
			args args
			want interface{}
		}{
			name: "Find",
			args: args{
				arrs: []interface{}{1, 2, 3},
				matcher: func(ele interface{}) bool {
					return ele.(int) == 2
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.args.arrs, tt.args.matcher); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
