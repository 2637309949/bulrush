// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import "github.com/kataras/go-events"

type (
	// Status store and share status Between different plug-ins
	Status struct {
		events events.EventEmmiter
		data   map[string]interface{}
	}
)

// Get value
func (s *Status) Get(key string) (interface{}, bool) {
	ret, ok := s.data[key]
	s.events.Emit("bulrush:status:get", key, ret)
	return ret, ok
}

// ALL value
func (s *Status) ALL() map[string]interface{} {
	s.events.Emit("bulrush:status:all", s.data)
	return s.data
}

// Set value
func (s *Status) Set(key string, value interface{}) *Status {
	s.events.Emit("bulrush:status:set", key, value)
	s.data[key] = value
	return s
}

// new statusStorage
func statusStorage(events events.EventEmmiter) *Status {
	return &Status{
		events: events,
		data:   make(map[string]interface{}, 0),
	}
}
