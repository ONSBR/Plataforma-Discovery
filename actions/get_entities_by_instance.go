package actions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PMoneda/http"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

//EntitiesList maps entities that domain app will save based on process memory
type EntitiesList []map[string]interface{}

//GetEntitiesByInstance returns all entities that need to saved on domain to complete a process instance
func GetEntitiesByInstance(systemID, processInstance string) (EntitiesList, error) {
	list := make(EntitiesList, 0)
	domainHost, err := getDomainHost(systemID)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/instance/%s/entities", domainHost, processInstance)

	resp, errR := http.Get(url)
	if errR != nil {
		return nil, errR
	}
	errJ := json.Unmarshal(resp.Body, &list)
	if errJ != nil {
		return nil, errJ
	}
	return list, nil
}

//PersistEntitiesByInstance call domain to persist data based on process instance
func PersistEntitiesByInstance(systemID, instanceID string) error {
	domainHost, err := getDomainHost(systemID)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/instance/%s/persist", domainHost, instanceID)
	if resp, err := http.Post(url, nil); err != nil {
		return err
	} else if !strings.Contains(string(resp.Body), "ok") {
		return fmt.Errorf("%s", resp)
	}
	return nil
}

func getDomainHost(systemID string) (string, error) {
	result := make([]map[string]interface{}, 1)
	filter := apicore.Filter{
		Entity: "installedApp",
		Map:    "core",
		Name:   "bySystemIdAndType",
		Params: []apicore.Param{apicore.Param{
			Key:   "systemId",
			Value: systemID,
		}, apicore.Param{
			Key:   "type",
			Value: "domain",
		},
		},
	}
	err := apicore.Query(filter, &result)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	if len(result) > 0 {
		obj := result[0]
		return fmt.Sprintf("http://%s:%d", obj["host"], uint(obj["port"].(float64))), nil
	}
	return "", fmt.Errorf("no app found for %s id", systemID)
}
