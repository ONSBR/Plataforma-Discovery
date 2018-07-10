package apicore

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/PMoneda/http"
)

func getURL() string {
	return fmt.Sprintf("%s://%s:%s", env.Get("APICORE_SCHEME", "http"), env.Get("APICORE_HOST", "localhost"), env.Get("APICORE_PORT", "9110"))
}

//Persist data on APICORE
func Persist(entities interface{}) error {
	_, err := http.Post(fmt.Sprintf("%s/core/persist", getURL()), entities)
	if err != nil {
		return err
	}
	return nil
}

//PersistOne single entity to API Core
func PersistOne(entity ...interface{}) error {
	return Persist(entity)
}

//Query data on apicore
func Query(filter Filter, response interface{}) error {
	url := fmt.Sprintf("%s/%s/%s?filter=%s", getURL(), filter.Map, filter.Entity, filter.Name)
	for _, param := range filter.Params {
		url += fmt.Sprintf("&%s=%s", param.Key, param.Value)
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Body, response)
	if err != nil {
		return err
	}
	return nil
}
