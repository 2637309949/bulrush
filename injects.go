// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"github.com/2637309949/bulrush-utils/maps"
	"github.com/kataras/go-events"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/robfig/cron.v2"
)

// Injects defined some entitys that can be inject to middle
// , Injects would panic if repetition
// , Injects can be go base tyle or struct or ptr or interface{}
type Injects []interface{}

// Append defined array concat
func (src *Injects) Append(target *Injects) *Injects {
	injects := append(*src, *target...)
	return &injects
}

// Has defined inject type is existed or not
func (src *Injects) Has(item interface{}) bool {
	return typeExists(*src, item)
}

func builtInInjects(bul *rush) Injects {
	emmiter := events.New()
	status := maps.NewSafeMap()
	validate := validator.New()
	schedule := cron.New()
	reverseInject := &ReverseInject{
		injects: bul.injects,
		inspect: func(ret ...interface{}) {
			bul.Inject(ret...)
		},
	}
	return Injects{
		emmiter,
		status,
		validate,
		schedule,
		reverseInject,
	}
}
