// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	addition "github.com/2637309949/bulrush-addition"
	"github.com/2637309949/bulrush-addition/logger"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// rushLogger just for console log
var rushLogger = addition.RushLogger

// SetLogger defined set logger for bulrush
// you can append transport to addition.RushLogger
// It is not recommended to create new one Journal
// I recommendation to customize different output from addition.RushLogger, but not SetLogger
func SetLogger(logger *logger.Journal) {
	rushLogger = logger
}

// GetLogger defined get logger for bulrush
// you can append transport to addition.RushLogger
// It is not recommended to create new one Journal
// I recommendation to customize different output from addition.RushLogger, but not SetLogger
func GetLogger() *logger.Journal {
	return rushLogger
}
