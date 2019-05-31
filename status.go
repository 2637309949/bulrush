/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush status]
 */

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
