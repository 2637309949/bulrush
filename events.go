// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"

	utils "github.com/2637309949/bulrush-utils"
	"github.com/kataras/go-events"
)

// EventName defined bulrush events name
func EventName(name string) events.EventName {
	return events.EventName(fmt.Sprintf("BULRUSH::%v", name))
}

// Events defined built-in events
var (
	EventsStarting = EventName(utils.ToString(1<<49 + 0))
	EventsRunning  = EventName(utils.ToString(1<<49 + 1))
	EventsShutdown = EventName(utils.ToString(1<<49 + 2))
)
