// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"github.com/gin-gonic/gin"
)

const (
	// DebugMode indicates bul mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates bul mode is release.
	ReleaseMode = "release"
	// TestMode indicates bul mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

var bulMode = debugCode
var modeName = DebugMode

// SetMode sets gin mode according to input string.
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		bulMode = debugCode
	case ReleaseMode:
		bulMode = releaseCode
	case TestMode:
		bulMode = testCode
	default:
		panic("bul mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	modeName = value
	gin.SetMode(value)
}

// Mode returns currently gin mode.
func Mode() string {
	return modeName
}
