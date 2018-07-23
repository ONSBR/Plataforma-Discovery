package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func ParseQuery(query string, parameters map[string]interface{}) string {
	for prop, v := range parameters {
		switch value := v.(type) {
		case string:
			query = compileParams(query, value, prop)
		case float64:
			query = strings.Replace(query, ":"+prop, strconv.FormatFloat(value, 'E', -1, 64), -1)
		case int64:
			query = strings.Replace(query, ":"+prop, strconv.FormatInt(value, 10), -1)
		case bool:
			query = strings.Replace(query, ":"+prop, strconv.FormatBool(value), -1)
		}
	}
	return removeUnsedParams(query)
}

func compileParams(query string, value, prop string) string {
	if strings.Contains(query, ":"+prop) {
		query = strings.Replace(query, ":"+prop, fmt.Sprintf("'%s'", value), -1)
	} else if strings.Contains(query, "$"+prop) {
		convertType := value[len(value)-1] != '!'
		if !convertType {
			value = value[0 : len(value)-1]
		}
		parts := strings.Split(value, ";")

		if len(parts) > 0 && IsNumeric(parts[0]) && convertType {
			joined := strings.Join(parts, ",")
			query = strings.Replace(query, "$"+prop, joined, -1)
		} else {
			for i := 0; i < len(parts); i++ {
				parts[i] = fmt.Sprintf("'%s'", parts[i])
			}
			query = strings.Replace(query, "$"+prop, strings.Join(parts, ","), -1)
		}
	}

	return query
}

var onceCompile sync.Once
var regex *regexp.Regexp
var regexOptional *regexp.Regexp

func removeUnsedParams(query string) string {
	//TODO
	onceCompile.Do(func() {
		regex, _ = regexp.Compile(`:\w*`)
		regexOptional, _ = regexp.Compile(`\[(.*?)\]`)
	})
	optionalParameters := regexOptional.FindAll([]byte(query), -1)
	for _, optional := range optionalParameters {
		value := string(optional)
		notPopulated := regex.FindAll(optional, -1)
		if len(notPopulated) > 0 {
			query = strings.Replace(query, value, "", 1)
		}
	}
	query = strings.Replace(query, "[", "", -1)
	query = strings.Replace(query, "]", "", -1)
	removeTrailing := strings.Trim(query, " ")
	return removeTrailing
}
