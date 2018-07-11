package helpers

import (
	"fmt"
)

func ExtractFieldFromEntity(entity map[string]interface{}, field string) (string, error) {
	_, ok := entity["_metadata"]
	if ok {
		switch t := entity["_metadata"].(type) {
		case map[string]interface{}:
			typ, ok := t[field]
			if ok {
				switch e := typ.(type) {
				case string:
					return e, nil
				}
			}
		}
	}
	return "", fmt.Errorf("cannot find entity type")
}
