// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"github.com/kataras/go-events"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/robfig/cron.v2"
)

// Schedule defined Cron job
type Schedule struct {
	Cron *cron.Cron
}

// ScheduleJob add job and run
func (s *Schedule) ScheduleJob(spec string, job func()) *cron.Cron {
	c := cron.New()
	c.AddFunc(spec, job)
	c.Start()
	return c
}

func defaultInjects(bul *rush) Injects {
	emmiter := events.New()
	status := statusStorage(emmiter)
	validate := validator.New()
	schedule := &Schedule{}
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
