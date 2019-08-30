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
	}
)

// next defined foreach
func (exec *executor) next(in func(s Scope) error, index ...int) (errs []error) {
	var i int
	if len(index) > 0 {
		i = index[0]
	}
	if len(*exec.scopes) >= (i + 1) {
		errs = append(errs, in((*exec.scopes)[i]))
		i++
		exec.next(in, i)
	}
	return
}

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (exec *executor) traverse() []error {
	return exec.next(func(pv Scope) error {
		return CatchError(func() {
			debugPrint("next plugin:%v", reflect.TypeOf(pv.Value))
			pv.inFrom(exec.injects)
			pv.methodCall(pv.indirectFunc(preHookName), *exec.injects)
			exec.inspect(pv.reflectCall(pv.indirectPlugin(), pv.Inputs)...)
			pv.methodCall(pv.indirectFunc(postHookName), *exec.injects)
		})
	})
}

// execute defined all plugin execute in orderly
// ,inspect defined cb for runPlugin
func (exec *executor) execute(inspect func(...interface{})) {
	exec.inspect = inspect
	exec.traverse()
}
