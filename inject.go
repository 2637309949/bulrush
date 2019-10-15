// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"fmt"
	"reflect"

	"github.com/2637309949/bulrush-utils/maps"
	"github.com/2637309949/bulrush-utils/sync"
	"github.com/thoas/go-funk"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/robfig/cron.v2"
)

// Injects defined some entitys that can be inject to middle
// , Injects would panic if repetition
// , Injects can be go base tyle or struct or ptr or interface{}
type (
	Injects []interface{}
	// InjectOption defined inject option
	InjectOption interface {
		apply(r *rush) *rush
		check(r *rush) interface{}
	}
)

// InjectsValidOption defined Option of valid
func InjectsValidOption(injects ...interface{}) InjectOption {
	return Option(func(r *rush) interface{} {
		funk.ForEach(injects, func(item interface{}) {
			assert1(!r.injects.Has(item), ErrWith(ErrInject, fmt.Sprintf("inject %v has existed", reflect.TypeOf(item))))
		})
		return injects
	})
}

// InjectsOption defined Option of Injects
func InjectsOption(injects ...interface{}) InjectOption {
	return Option(func(r *rush) interface{} {
		r.lock.Acquire("injects", func(async sync.Async) {
			funk.ForEach(injects, func(item interface{}) {
				r.injects.Put(item)
			})
		})
		return r
	})
}

func newInjects(items ...interface{}) *Injects {
	inject := make(Injects, 0)
	inject = append(inject, items...)
	return &inject
}

// Wire defined wire ele from type
func (src *Injects) Wire(target interface{}) (err error) {
	// tv := (*interface{})(unsafe.Pointer(targetValue.Pointer()))
	// va := reflect.ValueOf(&a).Elem()
	// va.Set(reflect.New(va.Type().Elem()))
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Ptr && !tv.IsNil() {
		err = ErrWith(ErrUnaddressable, fmt.Sprintf("type %v should be pointer", reflect.TypeOf(target)))
		return
	}
	if v := src.Acquire(tv.Elem().Type()); v != nil {
		tv = tv.Elem()
		if tv.Type() == reflect.TypeOf(v) && tv.CanSet() {
			tv.Set(reflect.ValueOf(v))
			return
		}
	}
	err = ErrWith(ErrUnaddressable, fmt.Sprintf("type %v not found in ct", reflect.TypeOf(target)))
	return
}

// Acquire defined acquire inject ele from type
func (src *Injects) Acquire(ty reflect.Type) interface{} {
	ele := typeMatcher(ty, *src)
	if ele == nil {
		ele = duckMatcher(ty, *src)
	}
	if ele != nil {
		ele = ele.(reflect.Value).Interface()
	}
	return ele
}

// Append defined array concat
func (src *Injects) Append(target *Injects) *Injects {
	injects := append(*src, *target...)
	return &injects
}

// Put defined array Put
func (src *Injects) Put(target interface{}) *Injects {
	*src = append(*src, target)
	return src
}

// Has defined inject type is existed or not
func (src *Injects) Has(item interface{}) bool {
	return typeExists(*src, item)
}

// Size defined inject Size
func (src *Injects) Size() int {
	return len(*src)
}

// Get defined index of inject
func (src *Injects) Get(pos int) interface{} {
	return (*src)[pos]
}

// Swap swaps the two values at the specified positions.
func (src *Injects) Swap(i, j int) {
	(*src)[i], (*src)[j] = (*src)[j], (*src)[i]
}

func builtInInjects(bul *rush) Injects {
	lifecycle := bul.lifecycle
	emmiter := bul.EventEmmiter
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
		lifecycle,
		status,
		validate,
		schedule,
		reverseInject,
	}
}
