// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"sync"

	"github.com/kataras/go-events"
)

// Status store and share status Between different plug-ins
type Status struct {
	events events.EventEmmiter
	m      map[string]interface{}
	l      *sync.RWMutex
}

// Get value
func (s *Status) Get(key string) (interface{}, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	ret, ok := s.m[key]
	return ret, ok
}

// ALL value
func (s *Status) ALL() map[string]interface{} {
	s.l.RLock()
	defer s.l.RUnlock()
	return s.m
}

// Set value
func (s *Status) Set(key string, value interface{}) *Status {
	s.l.Lock()
	defer s.l.Unlock()
	s.m[key] = value
	return s
}

// new statusStorage
func statusStorage(events events.EventEmmiter) *Status {
	return &Status{
		events: events,
		m:      make(map[string]interface{}, 0),
	}
}
