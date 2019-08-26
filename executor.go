// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import "reflect"

type (
	executor struct {
		scopes  *[]Scope
		injects *Injects
		inspect func(...interface{})
		index   int
	}
)

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (exec *executor) traverse() {
	for exec.index < len(*exec.scopes) {
		// roback if error panic in plugin
		if err := CatchError(func() {
			pv := (*exec.scopes)[exec.index]
			debugPrint("next plugin:%v", reflect.TypeOf(pv.Value))
			pv.inFrom(exec.injects)
			pv.methodCall(pv.indirectFunc(preHookName), *exec.injects)
			exec.inspect(pv.reflectCall(pv.indirectPlugin(), pv.Inputs)...)
			pv.methodCall(pv.indirectFunc(postHookName), *exec.injects)
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
	exec.traverse()
}
