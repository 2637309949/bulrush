// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
	"testing"
)

func Test_typeExists(t *testing.T) {
	type args struct {
		arr    []interface{}
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		struct {
			name string
			args args
			want bool
		}{
			name: "typeExists",
			args: args{
				arr:    []interface{}{"1", 1, true},
				target: "test",
			},
			want: true,
		},
		struct {
			name string
			args args
			want bool
		}{
			name: "typeExists",
			args: args{
				arr:    []interface{}{"1", 1, true},
				target: 1.1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := typeExists(tt.args.arr, tt.args.target); got != tt.want {
				t.Errorf("typeExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createSlice(t *testing.T) {
	type args struct {
		target interface{}
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
			name: "createSlice",
			args: args{
				target: 12,
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSlice(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testRR struct{}

func Test_reflectObjectAndCall(t *testing.T) {
	type args struct {
		target interface{}
		params []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		struct {
			name string
			args args
		}{
			name: "createSlice",
			args: args{
				target: &testRR{},
				params: []interface{}{1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reflectObjectAndCall(tt.args.target, tt.args.params)
		})
	}
}
