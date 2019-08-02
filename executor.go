// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"

	"github.com/thoas/go-funk"
)

type (
	executor struct {
		pluginValues *[]PluginValue
		injects      *Injects
	}
)

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (exec *executor) execute(inspect func(...interface{})) {
	funk.ForEach(*exec.pluginValues, func(pv PluginValue) {
		debugPrint("Exec plugin:%v", reflect.TypeOf(pv.Plugin.Interface()))
		pv.inputsFrom(*exec.injects)
		pv.runPre()
		inspect(pv.runPlugin()...)
		pv.runPost()
	})
}
