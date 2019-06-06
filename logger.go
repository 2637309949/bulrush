/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import "github.com/2637309949/bulrush-addition/logger"

// rushLogger just for console log
var rushLogger *logger.Journal

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
	reloadRushLogger(Mode)
}
