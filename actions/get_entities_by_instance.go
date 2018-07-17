package actions

import (
	"encoding/json"
	"fmt"

	"github.com/PMoneda/http"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Discovery/helpers"
	"github.com/ONSBR/Plataforma-Discovery/models"
)

//GetEntitiesByInstance returns all entities that need to saved on domain to complete a process instance
func GetEntitiesByInstance(systemID, processInstance string) (models.EntitiesList, error) {
	list := make(models.EntitiesList, 0)
	domainHost, err := helpers.GetDomainHost(systemID)
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
	return list, nil
}
