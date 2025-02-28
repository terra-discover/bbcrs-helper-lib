package sqlib

import (
	"errors"
	"reflect"
)

/*
CheckBeforeCreateInBatches

check input of interface{} before add it to gorm.CreateInBatches()
*/
func CheckBeforeCreateInBatches(input interface{}) (err error) {
	isNil := IsNilFixed(input)
	if isNil {
		err = errors.New("input is nil")
		return
	}

	// check length
	switch reflect.TypeOf(input).Kind() {
	case reflect.Map, reflect.Array, reflect.Slice:
		isValid := reflect.ValueOf(input).Len() > 0
		if !isValid {
			err = errors.New("input length is zero")
			return
		}
	}

	return
}

/*
IsNilFixed

check Nil interface the right way

Source: https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1
*/
func IsNilFixed(input interface{}) (result bool) {
	if input == nil {
		result = true
		return
	}

	switch reflect.TypeOf(input).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		result = reflect.ValueOf(input).IsNil()
		return
	}

	return
}
