package main

import (
	"runtime"
	"reflect"
	"fmt"
)

// User -
type User struct {
	name string
}

func (user *User) getRouter() {
	fmt.Print("getRouter")
}

func main() {
	user := User{ name: "OK", }
	userValue := reflect.ValueOf(user)
	runtime.FuncForPC(userValue.Pointer()).Name()
}