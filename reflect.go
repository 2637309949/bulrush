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
func invoke(target interface{}, injects map[string]interface{}) {
	engine, _ := injects["Engine"]
	mongo, _  := injects["Mongo"]
	router, _ := injects["Router"]
	config, _ := injects["Config"]
	redis, _  := injects["Redis"]
	getType   := reflect.TypeOf(target)
	getValue := reflect.ValueOf(target)
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
					inputs = append(inputs, reflect.ValueOf(injects))
				case ptype == reflect.TypeOf(engine):
					inputs = append(inputs, reflect.ValueOf(engine))
				case ptype == reflect.TypeOf(mongo):
					inputs = append(inputs, reflect.ValueOf(mongo))
				case ptype == reflect.TypeOf(router):
					inputs = append(inputs, reflect.ValueOf(router))
				case ptype == reflect.TypeOf(config):
					inputs = append(inputs, reflect.ValueOf(config))
				case ptype == reflect.TypeOf(redis):
					inputs = append(inputs, reflect.ValueOf(redis))
				default:
			}
		}
		if method.IsValid() {
			method.Call(inputs)
		}
	}
}
