// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import "github.com/2637309949/bulrush-addition/logger"

// rushLogger just for console log
var rushLogger *logger.Journal

// reloadUrhsLogger for reload logger level after setting pro mode
func reloadRushLogger(mode string) {
	var level = logger.SILLYLevel
	if mode == "release" {
		level = logger.WARNLevel
	}
	rushLogger = logger.CreateLogger(level, nil,
		[]*logger.Transport{
			&logger.Transport{
				Level: level,
			},
		},
	)
}

func init() {
	reloadRushLogger(Conf.Mode)
}
