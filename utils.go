package bulrush

import (
	"fmt"
	"math/rand"
)

// Some get or a default value
func Some(t interface{}, i interface{}) interface{} {
	if t != nil && t != "" && t != 0 {
		return t
	}
	return i
}

// Find element in arrs with some match condition
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

// ToIntArray -
func ToIntArray(t []interface{}) []int {
	s := make([]int, len(t))
	for i, v := range t {
		s[i] = v.(int)
	}
	return s
}

// RandString gen string
func RandString(n int) string {
	const seeds = "abcdefghijklmnopqrstuvwxyz1234567890"
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
func LeftOkV(left interface{}, right ...bool) interface{} {
	var r = true
	if len(right) != 0 {
		r = right[0]
	} else if left == "" || left == nil || left == 0 {
		r = false
	}
	if r {
		return left
	}
	return nil
}

// LeftSV left value or panic
func LeftSV(left interface{}, right error) interface{} {
	if right != nil {
		panic(right)
	}
	return left
}

func fixedPortPrefix(port string) string {
	if prefix := port[:1]; prefix != ":" {
		port = fmt.Sprintf(":%s", port)
	}
	return port
}
