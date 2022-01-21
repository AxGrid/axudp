package axtools

import (
	"reflect"
	"strconv"
)

func getIntFromObjectMap(m map[string]interface{}, key string, def int) int {
	i, ok := m[key]
	if !ok {
		return def
	}
	switch i.(type) {
	case string:
		v, err := strconv.Atoi(i.(string))
		if err != nil {
			return def
		}
		return v
	case int:
		return i.(int)
	default:
		return def
	}
}

func getIntFromStringMap(m map[string]string, key string, def int) int {
	i, ok := m[key]
	if !ok {
		return def
	}
	v, err := strconv.Atoi(i)
	if err != nil {
		return def
	}
	return v
}

func getStringFromObjectMap(m map[string]interface{}, key string, def string) string {
	i, ok := m[key]
	if !ok {
		return def
	}
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	default:
		return def
	}
}

func getStringFromStringMap(m map[string]string, key string, def string) string {
	i, ok := m[key]
	if !ok {
		return def
	}
	return i
}

func GetIntFromMap(m interface{}, key string, def int) int {
	if reflect.ValueOf(m).Kind() != reflect.Map {
		return def
	}
	if reflect.TypeOf(m).Elem().Kind() == reflect.String {
		return getIntFromStringMap(m.(map[string]string), key, def)
	}
	return getIntFromObjectMap(m.(map[string]interface{}), key, def)
}

func GetStringFromMap(m interface{}, key string, def string) string {
	if reflect.ValueOf(m).Kind() != reflect.Map {
		return def
	}
	if reflect.TypeOf(m).Elem().Kind() == reflect.String {
		return getStringFromStringMap(m.(map[string]string), key, def)
	}
	return getStringFromObjectMap(m.(map[string]interface{}), key, def)
}
