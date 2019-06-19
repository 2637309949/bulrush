// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"github.com/bamzi/jobrunner"
	job "github.com/bamzi/jobrunner"
	"github.com/kataras/go-events"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/robfig/cron.v2"
)

// Jobrunner defind wrap Jobrunner to a struct
type Jobrunner struct {
	Start      func(...int)
	Schedule   func(string, cron.Job) error
	StatusPage func() []jobrunner.StatusData
	StatusJSON func() map[string]interface{}
}

func defaultInjects(bul *rush) Injects {
	emmiter := events.New()
	status := statusStorage(emmiter)
	validate := validator.New()
	jobrunner := &Jobrunner{
		Start:      job.Start,
		Schedule:   job.Schedule,
		StatusPage: job.StatusPage,
		StatusJSON: job.StatusJson,
	}
	rs := &ReverseInject{
		injects: bul.injects,
	}
	return Injects{
		emmiter,
		status,
		validate,
		jobrunner,
		rs,
	}
}
