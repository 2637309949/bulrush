// Copyright (c) 2018-2020 Double All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bulrush

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

var (
	// ErrPlugin is used when bul.Use().
	ErrPlugin = &Error{Code: uint64(1 << 20)}
	// ErrInject is used when bul.Inject() fails.
	ErrInject = &Error{Code: uint64(1 << 21)}
	// ErrUnaddressable unaddressable value
	ErrUnaddressable = &Error{Code: uint64(1 << 22)}
	// ErrPrivate indicates a private error.
	ErrPrivate = &Error{Code: uint64(1 << 23)}
	// ErrPublic indicates a public error.
	ErrPublic = &Error{Code: uint64(1 << 24)}
	// ErrAny indicates any other error.
	ErrAny = &Error{Code: uint64(1 << 25)}
	// ErrNu indicates any other error.
	ErrNu = &Error{Code: uint64(1 << 55)}
)

// Error represents a error's specification.
type Error struct {
	Err  error
	Code uint64
	Meta interface{}
}

// SetCode sets the error's code
func (msg *Error) SetCode(code uint64) *Error {
	msg.Code = code
	return msg
}

// SetMeta sets the error's meta data.
func (msg *Error) SetMeta(data interface{}) *Error {
	msg.Meta = data
	return msg
}

// JSON creates a properly formatted JSON
func (msg *Error) JSON() interface{} {
	json := map[string]interface{}{}
	if msg.Meta != nil {
		value := reflect.ValueOf(msg.Meta)
		switch value.Kind() {
		case reflect.Struct:
			return msg.Meta
		case reflect.Map:
			funk.ForEach(value.MapKeys(), func(kv reflect.Value) {
				json[kv.String()] = value.MapIndex(kv).Interface()
			})
		default:
			json["meta"] = msg.Meta
		}
	}
	if _, ok := json["error"]; !ok {
		json["error"] = msg.Error()
	}
	return json
}

// MarshalJSON implements the json.Marshaller interface.
func (msg *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.JSON())
}

// Error implements the error interface.
func (msg *Error) Error() (message string) {
	if msg.Err != nil {
		message = msg.Err.Error()
	}
	return
}

// IsCode judges one error.
func (msg *Error) IsCode(code uint64) bool {
	return (msg.Code & code) > 0
}

// ErrWith defined wrap error
func ErrWith(err *Error, msg string) error {
	return errors.Wrap(err, msg)
}

// ErrOut defined unwrap error
func ErrOut(err error) (bulError *Error, ok bool) {
	bulError, ok = errors.Cause(err).(*Error)
	return
}

// ErrMsgs defined split wrap msg
func ErrMsgs(err error) []string {
	return strings.Split(err.Error(), ":")
}

// ErrCode defined return error code
func ErrCode(err error) (code uint64) {
	if bulErr, ok := ErrOut(err); ok {
		code = bulErr.Code
	}
	code = ErrNu.Code
	return
}
