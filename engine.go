// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

type (
	engine struct {
		scopes  *[]Scope
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
		pv := (*e.scopes)[i]
		debugPrint("next plugin:%v", pv.Type())
		errs = append(errs, cb(pv))
		i++
		e.next(cb, i)
	}
	return
}

// execute defined run app plugin in order
//, if Pre or Post Hook defined in struct, then
//, Pre > Plugin > Post
func (e *engine) traverse() (errs []error) {
	return e.next(func(pv Scope) (err error) {
		return CatchError(func() {
			pv.Pre()
			e.inspect(pv.Plugin()...)
			pv.Post()
		})
	})
}

// execute defined all plugin execute in orderly
// ,inspect defined cb for runPlugin
func (e *engine) exec(inspect func(...interface{})) {
	e.inspect = inspect
	e.traverse()
}
