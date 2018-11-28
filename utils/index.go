package utils

import (
	"math/rand"
	"fmt"
)

// GetOrElse -
func GetOrElse(target map[string]interface{}, key string, initValue interface{}) interface{} {
	if value, ok := target[key]; ok {
		return value
	}
	return initValue
}

// Some get or a default value
func Some(target interface{}, initValue interface{}) interface{}{
	if target != nil && target != "" {
		return target
	}
	return initValue
}

// Find -
func Find(arrs []interface{}, matcher func(interface{}) bool) interface{} {
	var target interface{}
	for _, item := range arrs {
		match := matcher(item)
		if match {
			target = item
			break
		}
	}
	return target
}


// ToStrArray -
func ToStrArray(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}
	return s
}

// RandString -
func RandString(n int) string {
	const seeds = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = seeds[rand.Intn(len(seeds))]
	}
	return string(bytes)
}
