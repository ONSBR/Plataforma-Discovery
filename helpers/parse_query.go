package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseQuery(query string, parameters map[string]interface{}) string {
	for k, v := range parameters {
		switch t := v.(type) {
		case string:
			query = strings.Replace(query, ":"+k, fmt.Sprintf("'%s'", t), -1)
		case float64:
			query = strings.Replace(query, ":"+k, strconv.FormatFloat(t, 'E', -1, 64), -1)
		case int64:
			query = strings.Replace(query, ":"+k, strconv.FormatInt(t, 10), -1)
		case bool:
			query = strings.Replace(query, ":"+k, strconv.FormatBool(t), -1)
		}
	}
	return query
}
