package sdk

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/PMoneda/http"
)

func getProcessMemoryUrl() string {
	return fmt.Sprintf("%s://%s:%s/",
		infra.GetEnv("PROCESS_MEMORY_SCHEME", "http"),
		infra.GetEnv("PROCESS_MEMORY_HOST", "localhost"),
		infra.GetEnv("PROCESS_MEMORY_PORT", "9091"))
}

//SaveDocument to process memory
func SaveDocument(collection string, data interface{}) (err error) {
	url := fmt.Sprintf("%s%s?app_origin=event_manager", getProcessMemoryUrl(), collection)
	_, err = http.Post(url, data)
	return
}

//GetDocument from process memory
func GetDocument(collection string, query map[string]string) (string, error) {
	queryString := ""
	for k, v := range query {
		queryString += fmt.Sprintf("%s=%s&", k, v)
	}
	url := fmt.Sprintf("%s%s?%sapp_origin=event_manager", getProcessMemoryUrl(), collection, queryString)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	return string(response.Body), nil
}

//ReplaceDocument update full document based on query and collection
func ReplaceDocument(collection string, query map[string]string, document interface{}) error {
	queryString := ""
	for k, v := range query {
		queryString += fmt.Sprintf("%s=%s&", k, v)
	}
	url := fmt.Sprintf("%s%s?%sapp_origin=event_manager", getProcessMemoryUrl(), collection, queryString)

	_, err := http.Put(url, document)
	return err
}
