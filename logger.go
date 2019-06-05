/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush cfg struct]
 */

package bulrush

import "github.com/2637309949/bulrush-addition/logger"

var rushLogger = logger.CreateLogger(logger.INFOLevel, nil,
	[]*logger.Transport{
		&logger.Transport{
			Level: logger.INFOLevel,
		},
	},
)
