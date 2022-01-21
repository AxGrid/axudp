package axtools

import (
	"strconv"
	"strings"
)

func RemoveDuplicate(data []interface{}) []interface{} {
	keys := make(map[interface{}]bool)
	var list []interface{}
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func RemoveIntDuplicate(data []int) []int {
	keys := make(map[interface{}]bool)
	var list []int
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func RemoveStringDuplicate(data []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range data {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func GetFromArrayOrSingle(item interface{}) []interface{} {
	if item == nil {
		return []interface{}{}
	}
	switch item.(type) {
	case int:
		return []interface{}{item}
	default:
		res, err := ReadInterfaceAsArray(item)
		if err != nil {
			panic(err)
		}
		return res
	}
}

func GetStringFromArrayOrSingle(item interface{}) []string {
	if item == nil {
		return []string{}
	}
	switch item.(type) {
	case string:
		return []string{item.(string)}
	default:
		res, err := ReadInterfaceAsArrayString(item)
		if err != nil {
			panic(err)
		}
		return res
	}
}

func GetIntFromArrayOrSingle(item interface{}) []int {
	if item == nil {
		return []int{}
	}
	switch item.(type) {
	case int:
		return []int{item.(int)}
	default:
		res, err := ReadInterfaceAsArrayInt(item)
		if err != nil {
			panic(err)
		}
		return res
	}
}

func IntArrayAsString(data []int) string {
	res := "[]int{ "
	for i, v := range data {
		res += strconv.Itoa(v)
		if i < len(data)-1 {
			res += ","
		}
	}
	return res + " }"
}

func IntArrayAsStringEmptyNil(data []int) string {
	if data == nil || len(data) == 0 {
		return "nil"
	}
	res := "[]int{ "
	for i, v := range data {
		res += strconv.Itoa(v)
		if i < len(data)-1 {
			res += ","
		}
	}
	return res + " }"
}

func StringArrayAsString(data []string) string {
	res := "[]string{ "
	for i, v := range data {
		res += "\"" + v + "\""
		if i < len(data)-1 {
			res += ","
		}
	}
	return res + " }"
}

func StringArrayAsStringEmptyNil(data []string) string {
	if data == nil || len(data) == 0 {
		return "nil"
	}
	res := "[]string{ "
	for i, v := range data {
		res += "\"" + v + "\""
		if i < len(data)-1 {
			res += ","
		}
	}
	return res + " }"
}

func Uint64ArrayContain(data []uint64, e uint64) bool {
	for _, a := range data {
		if a == e {
			return true
		}
	}
	return false
}

func IntArrayToString(data []int) string {
	l := len(data)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range data {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, ",")
}

func Int32ArrayToString(data []int32) string {
	l := len(data)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range data {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, ",")
}

func Uint32ArrayToString(data []uint32) string {
	l := len(data)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range data {
		b[i] = strconv.Itoa(int(v))
	}
	return strings.Join(b, ",")
}

func Int64ArrayToString(data []int64) string {
	l := len(data)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range data {
		b[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(b, ",")
}

func Uint64ArrayToString(data []uint64) string {
	l := len(data)
	if l == 0 {
		return ""
	}
	b := make([]string, l)
	for i, v := range data {
		b[i] = strconv.FormatUint(v, 10)
	}
	return strings.Join(b, ",")
}

func StringToIntArray(text string, d ...string) []int {
	if text == "" {
		return []int{}
	}
	var data []string
	if len(d) > 0 {
		data = strings.Split(text, d[0])
	} else {
		data = strings.Split(text, ",")
	}
	res := make([]int, len(data))
	for i, sv := range data {
		val, err := strconv.Atoi(sv)
		if err != nil {
			res[i] = 0
		} else {
			res[i] = val
		}
	}
	return res
}

func StringToUint64Array(text string, d ...string) []uint64 {
	if text == "" {
		return []uint64{}
	}
	var data []string
	if len(d) > 0 {
		data = strings.Split(text, d[0])
	} else {
		data = strings.Split(text, ",")
	}
	res := make([]uint64, len(data))
	for i, sv := range data {
		val, err := strconv.ParseUint(sv, 10, 64)
		if err != nil {
			res[i] = 0
		} else {
			res[i] = val
		}
	}
	return res
}

func StringToInt64Array(text string, d ...string) []int64 {
	if text == "" {
		return []int64{}
	}
	var data []string
	if len(d) > 0 {
		data = strings.Split(text, d[0])
	} else {
		data = strings.Split(text, ",")
	}
	res := make([]int64, len(data))
	for i, sv := range data {
		val, err := strconv.ParseInt(sv, 10, 64)
		if err != nil {
			res[i] = 0
		} else {
			res[i] = val
		}
	}
	return res
}
