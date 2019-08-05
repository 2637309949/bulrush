// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"reflect"
)

type (
	executor struct {
		pluginValues *[]PluginContext
		injects      *Injects
		inspect      func(...interface{})
		index        int
	}
)

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (exec *executor) next() {
	for exec.index < len(*exec.pluginValues) {
		pv := (*exec.pluginValues)[exec.index]
		debugPrint("exec plugin:%v", reflect.TypeOf(pv.Plugin.Interface()))
		pv.inputsFrom(*exec.injects)
		pv.runPre()
		exec.inspect(pv.runPlugin()...)
		pv.runPost()
		exec.index++
	}
}

func (exec *executor) execute(inspect func(...interface{})) {
	exec.inspect = inspect
	exec.next()
}
