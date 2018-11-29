package utils

import (
	"math/rand"
	"fmt"
)

// Some get or a default value
func Some(target interface{}, initValue interface{}) interface{}{
	if target != nil && target != "" && target != 0 {
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

// LeftV -
func LeftV(left interface{}, right interface{}) interface{} {
	return left
}

// RightV -
func RightV(left interface{}, right interface{}) interface{} {
	return right
}

// LeftOkV -
func LeftOkV(left interface{}, right... bool) interface{} {
	var (
		l interface{}
		r = true
	)
	if len(right) == 0 && (l == "" || l == nil || l == 0){
		r = false
	} else {
		r = right[0]
	}
	if r {
		return left
	}
	return nil
}

// LeftSV -
func LeftSV(left interface{}, right error) interface{} {
	if right != nil {
		panic(right)
	}
	return left
}
