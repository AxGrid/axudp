package axtools

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadInterfaceAsArray Create from array from map
func ReadInterfaceAsArray(i interface{}) ([]interface{}, error) {
	if reflect.ValueOf(i).Kind() == reflect.Map {
		return []interface{}{i}, nil
	}
	if reflect.ValueOf(i).Kind() == reflect.Array {
		return i.([]interface{}), nil
	}
	if reflect.ValueOf(i).Kind() == reflect.Slice {
		return i.([]interface{}), nil
	}
	return nil, errors.New(fmt.Sprintf("undefined type %s", reflect.ValueOf(i).Kind()))
}

func ReadInterfaceAsArrayInt(t interface{}) ([]int, error) {
	var res []int
	array, err := ReadInterfaceAsArray(t)
	if err != nil {
		return res, err
	}
	res = make([]int, len(array))
	for i, e := range array {
		res[i] = e.(int)
	}
	return res, nil
}

func ReadInterfaceAsArrayString(t interface{}) ([]string, error) {
	var res []string
	array, err := ReadInterfaceAsArray(t)
	if err != nil {
		return res, err
	}
	res = make([]string, len(array))
	for i, e := range array {
		res[i] = e.(string)
	}
	return res, nil
}
