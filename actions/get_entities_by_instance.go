package actions

import (
	"encoding/json"
	"fmt"

	"github.com/PMoneda/http"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Discovery/models"
)

//GetEntitiesByInstance returns all entities that need to saved on domain to complete a process instance
func GetEntitiesByInstance(systemID, processInstance string) (models.EntitiesList, error) {
	list := make(models.EntitiesList, 0)
	domainHost, err := getDomainHost(systemID)
	log.Info(domainHost)
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
		log.Error(errJ)
		return nil, errJ
	}
	log.Info("entities from domain: ", len(list))
	return list, nil
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
		//FIXME remover estas linha
		obj["host"] = "localhost"
		obj["port"] = float64(8087)
		return fmt.Sprintf("http://%s:%d", obj["host"], uint(obj["port"].(float64))), nil
	}
	return "", fmt.Errorf("no app found for %s id", systemID)
}
