package infra

import (
	"encoding/json"
)

//Clone make a deep clone from objects
func Clone(from interface{}, to interface{}) error {
	if data, err := json.Marshal(from); err != nil {
		return NewArgumentException(err.Error())
	} else if err := json.Unmarshal(data, to); err != nil {
		return NewArgumentException(err.Error())
	}
	return nil
}
