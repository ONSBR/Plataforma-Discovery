package helpers

import (
	"fmt"
	"time"

	"github.com/ONSBR/Plataforma-Discovery/util"
	"github.com/labstack/gommon/log"
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
	} else {
		log.Error("_metada not found on entity")
	}
	return "", fmt.Errorf(fmt.Sprintf("cannot find field %s", field))
}

func ExtractModifiedTimestamp(entity map[string]interface{}) (int64, error) {
	timestamp, _ := ExtractFieldFromEntity(entity, "modified_at")
	return parseStringToTime(timestamp)
}

func parseStringToTime(str string) (int64, error) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, str)
	if err != nil {
		return 0, err
	}
	return util.Timestamp(t), nil
}
