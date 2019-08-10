// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

type (
	executor struct {
		pluginContexts *[]PluginContext
		injects        *Injects
		inspect        func(...interface{})
		index          int
	}
)

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (exec *executor) next() {
	for exec.index < len(*exec.pluginContexts) {
		// roback if error panic in plugin
		if err := CatchError(func() {
			pv := (*exec.pluginContexts)[exec.index]
			debugPrint("next plugin:%v", pv.Plugin.Type())
			pv.inputsFrom(*exec.injects)
			pv.runPre()
			exec.inspect(pv.runPlugin()...)
			pv.runPost()
			exec.index++
		}); err != nil {
			// next plugin
			exec.index++
		}
	}
}

// execute defined all plugin execute in orderly
// ,inspect defined cb for runPlugin
func (exec *executor) execute(inspect func(...interface{})) {
	exec.inspect = inspect
	exec.next()
}
