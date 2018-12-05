package bulrush

import (
	"fmt"
	"strings"
	"reflect"
	"github.com/thoas/go-funk"
	"github.com/gin-gonic/gin"
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

// invokeObject -
func invokeObject(target interface{}, injectParams []interface{}) {
	getType  := reflect.TypeOf(target)
	getValue := reflect.ValueOf(target)

	if getValue.Kind() != reflect.Ptr {
		panic("target must be a ptr")
	}
	for i := 0; i < getType.NumMethod(); i++ {
		valid	   := true
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
			r := funk.Find(injectParams, func(x interface{}) bool {
				return ptype == reflect.TypeOf(x)
			})
			if r != nil {
				inputs = append(inputs, reflect.ValueOf(r))
			} else {
				valid = false
				break
			}
		}
		if method.IsValid() && valid {
			method.Call(inputs)
		} else {
			panic(fmt.Errorf("Invalid method: %s in inject", methodName))
		}
	}
}

// invokeMethod -
func invokeMethod(target interface{}, injectParams []interface{}) interface {} {
	valid 	   := true
	getType    := reflect.TypeOf(target)
	methodName := getType.Name()
	getValue   := reflect.ValueOf(target)
	inputs 	   := make([]reflect.Value, 0)
	numIn	   := getType.NumIn()
	for index := 0; index < numIn; index ++ {
		ptype := getType.In(index)
		r := funk.Find(injectParams, func(x interface{}) bool {
			return ptype == reflect.TypeOf(x)
		})
		if r != nil {
			inputs = append(inputs, reflect.ValueOf(r))
		} else {
			valid = false
			break
		}
	}
	if getValue.IsValid() && valid {
		rs := getValue.Call(inputs)
		if rs[0].IsValid() {
			return rs[0].Interface()
		}
	} else {
		panic(fmt.Errorf("Invalid method: %s in inject", methodName))
	}
	return nil
}

// injectInvoke -
func injectInvoke(injects []interface{}, bulrush *Bulrush) {
	injectParams := []interface{}{
		bulrush.engine,
		bulrush.router,
		bulrush.mongo,
		bulrush.config,
		bulrush.redis,
		map[string]interface{} {
			"Engine": bulrush.engine,
			"Router": bulrush.router,
			"Mongo":  bulrush.mongo,
			"Config": bulrush.config,
			"Redis":  bulrush.redis,
		},
	}
	for _, target := range injects {
		invokeObject(target, injectParams)
	}
}

// inspectInvoke -
func inspectInvoke(target interface{}, bulrush *Bulrush) interface {}{
	injectParams := []interface{}{
		bulrush.engine,
		bulrush.router,
		map[string]interface{}{
			"DebugPrintRouteFunc": gin.DebugPrintRouteFunc,
			"SetMode": 			   gin.SetMode,
		},
	}
	return invokeMethod(target, injectParams)
}

// createSlice -
// return slice
func createSlice(target interface{}) interface{} {
	tagetType 	:= reflect.TypeOf(target)
	if tagetType.Kind() == reflect.Ptr {
		tagetType = tagetType.Elem()
	}
	targetSlice := reflect.MakeSlice(reflect.SliceOf(tagetType), 0, 0).Interface()
	return targetSlice
}

// createObject -
// return ptr
func createObject(target interface{}) interface{} {
	tagetType 	 := reflect.TypeOf(target)
	if tagetType.Kind() == reflect.Ptr {
		tagetType = tagetType.Elem()
	}
	targetObject := reflect.New(tagetType).Interface()
	return targetObject
}
