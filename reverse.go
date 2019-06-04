/**
 * @author [Double]
 * @email [2637309949@qq.com.com]
 * @create date 2019-01-12 22:46:31
 * @modify date 2019-01-12 22:46:31
 * @desc [bulrush reverse]
 */

package bulrush

import (
	"fmt"
	"reflect"
)

// ReverseInject Inject
type ReverseInject struct {
	injects *Injects
}

// Register function for Reverse Injects
// If the function you're injecting is a black box,
// then you can try this
func (r *ReverseInject) Register(rFunc interface{}) interface{} {
	kind := reflect.TypeOf(rFunc).Kind()
	if kind != reflect.Func {
		panic(fmt.Errorf("rFunc should to be func type"))
	}
	return reflectMethodAndCall(rFunc, *r.injects)
}
