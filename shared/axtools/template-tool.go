package axtools

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"text/template"
)

func Render(templateByte []byte, data interface{}) ([]byte, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"intRange":          intRange,
		"intArray":          intArrayAsString,
		"arrayAsString":     arrayAsString,
		"arrayEnumAsString": arrayEnumAsString,
		"toCamel":           NameFromUnderlinedToCamel,
		"toUnderline":       NameFromCamelToUnderline,
	}).Parse(string(templateByte))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create template")
		return nil, err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to render template")
		return nil, err
	}
	return tpl.Bytes(), nil
}

func intRange(start, end int) []int {
	var result []int
	for i := start; i < end; i++ {
		result = append(result, i)
	}
	return result
}

func intArrayAsString(data []int) string {
	res := "[]int{ "
	for i, v := range data {
		res += strconv.Itoa(v)
		if i < len(data)-1 {
			res += ","
		}
	}
	return res + " }"
}

func arrayAsString(data interface{}, t string) string {
	res := fmt.Sprintf("[]%s{", t)
	switch dt := data.(type) {
	case []int:
		res += IntArrayToString(dt)
		break
	case []int32:
		res += Int32ArrayToString(dt)
		break
	case []int64:
		res += Int64ArrayToString(dt)
		break
	case []uint64:
		res += Uint64ArrayToString(dt)
		break
	case []uint32:
		res += Uint32ArrayToString(dt)
		break
	}
	return res + "}"
}

func arrayEnumAsString(data []interface{}, t string, prefix string) string {
	res := fmt.Sprintf("[]%s{", t)
	for i, v := range data {
		res += fmt.Sprintf("%s%v", prefix, v)
		if i < len(data)-1 {
			res += ","
		}
	}
	return res + "}"
}
