// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

// Option defined implement of option
type (
	Option func(*rush) *rush
)

// option defined implement of option
func (o Option) apply(r *rush) *rush { return o(r) }

// Empty defined Option of rush
func Empty() Option {
	return Option(func(r *rush) *rush {
		r.injects = new(Injects)
		r.prePlugins = new(Plugins)
		r.plugins = new(Plugins)
		r.postPlugins = new(Plugins)
		return r
	})
}
