package axtools

import (
	"strings"
)
import strUtil "github.com/agrison/go-commons-lang/stringUtils"

func NameFromUnderlinedToCamel(name string) string {
	s := strings.Split(strings.ReplaceAll(name, "-", "_"), "_")
	res := ""
	for _, v := range s {
		res += strUtil.Capitalize(v)
	}
	return res
}

func NameFromCamelToUnderline(name string) string {
	res := ""
	for _, c := range name {
		if strUtil.IsAllUpperCase(string(c)) && res != "" {
			res += "_"
		}
		res += strUtil.Uncapitalize(string(c))
	}
	return res
}

//func UInt64ArrayToString(array []uint64) string {
//	b := make([]string, len(array))
//	for i, v := range array {
//		b[i] = strconv.FormatUint(v, 10)
//	}
//	return strings.Join(b, ",")
//}
//
//func StringToUInt64Array(str string) []uint64 {
//	items := strings.Split(str, ",")
//	b := make([]uint64, len(items))
//	for i, v := range items {
//		number, err := strconv.ParseUint(v, 10, 64)
//		if err == nil {
//			b[i] = number
//		}
//	}
//	return b
//}
