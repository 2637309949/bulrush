package bulrush

import (
	"strings"
	"reflect"
)

// InjectGroup -
type InjectGroup struct {
	InjectMongo  func(interface{})
	InjectRouter func(interface{})
	InjectConfig func(interface{})
	InjectEngine func(interface{})
	InjectRedis  func(interface{})
	Inject 	     func(interface{})
}

// invoke -
func invoke(target interface{}, bulrush *Bulrush) {
	getType  := reflect.TypeOf(target)
	getValue := reflect.ValueOf(target)

	if getValue.Kind() != reflect.Ptr {
		panic("target must be a ptr")
	}
	for i := 0; i < getType.NumMethod(); i++ {
		inputs 	   := make([]reflect.Value, 0)
		methodType := getType.Method(i)
		methodName := methodType.Name
		method 	   := getValue.Method(i)
		numIn	   := methodType.Type.NumIn()
		if !strings.HasPrefix(methodName, "Inject") {
			continue
		}
		for index := 1; index < numIn; index ++ {
			ptype := methodType.Type.In(index)
			switch {
				case ptype == reflect.TypeOf(map[string]interface{}{}):
					inputs = append(inputs, reflect.ValueOf(map[string]interface{} {
						"Engine": bulrush.engine,
						"Router": bulrush.router,
						"Mongo":  bulrush.mongo,
						"Config": bulrush.config,
						"Redis":  bulrush.redis,
					}))
				case ptype == reflect.TypeOf(bulrush.engine):
					inputs = append(inputs, reflect.ValueOf(bulrush.engine))
				case ptype == reflect.TypeOf(bulrush.mongo):
					inputs = append(inputs, reflect.ValueOf(bulrush.mongo))
				case ptype == reflect.TypeOf(bulrush.router):
					inputs = append(inputs, reflect.ValueOf(bulrush.router))
				case ptype == reflect.TypeOf(bulrush.config):
					inputs = append(inputs, reflect.ValueOf(bulrush.config))
				case ptype == reflect.TypeOf(bulrush.redis):
					inputs = append(inputs, reflect.ValueOf(bulrush.redis))
				default:
			}
		}
		if method.IsValid() {
			method.Call(inputs)
		}
	}
}

// injectInvoke -
func injectInvoke(injects []interface{}, bulrush *Bulrush) {
	for _, target := range injects {
		invoke(target, bulrush)
	}
}