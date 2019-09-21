// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import "reflect"

type (
	engine struct {
		scopes  *[]Scope
		injects *Injects
		inspect func(...interface{})
	}
)

// next defined foreach
func (e *engine) next(cb func(s Scope) error, index ...int) (errs []error) {
	i := 0
	if len(index) > 0 {
		i = index[0]
	}
	if len(*e.scopes) >= (i + 1) {
		errs = append(errs, cb((*e.scopes)[i]))
		i++
		e.next(cb, i)
	}
	return
}

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (e *engine) traverse() (errs []error) {
	errs = e.next(func(pv Scope) (err error) {
		err = CatchError(func() {
			debugPrint("next plugin:%v", reflect.TypeOf(pv.Value))
			pv.arguments(e.injects)
			pv.methodCall(pv.indirectFunc(preHookName), *e.injects)
			e.inspect(pv.reflectCall(pv.indirectPlugin(), pv.Inputs)...)
			pv.methodCall(pv.indirectFunc(postHookName), *e.injects)
		})
		return
	})
	return
}

// execute defined all plugin execute in orderly
// ,inspect defined cb for runPlugin
func (e *engine) exec(inspect func(...interface{})) {
	e.inspect = inspect
	e.traverse()
}
