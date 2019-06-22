// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"github.com/kataras/go-events"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/robfig/cron.v2"
)

func defaultInjects(bul *rush) Injects {
	emmiter := events.New()
	status := statusStorage(emmiter)
	validate := validator.New()
	schedule := cron.New()
	reverseInject := &ReverseInject{
		injects: bul.injects,
	}
	return Injects{
		emmiter,
		status,
		validate,
		schedule,
		reverseInject,
	}
}