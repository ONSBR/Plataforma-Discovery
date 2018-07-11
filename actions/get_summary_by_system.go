package actions

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Discovery/models"
	"github.com/PMoneda/http"
	"github.com/labstack/gommon/log"
)

//GetSummaryBySystem returns all instances query summary from process memory
func GetSummaryBySystem(systemID, entities string) ([]*models.InstanceSummary, error) {
	summary := make([]*models.InstanceSummary, 0)
	scheme := env.Get("PROCESS_MEMORY_SCHEME", "http")
	host := env.Get("PROCESS_MEMORY_HOST", "localhost")
	port := env.Get("PROCESS_MEMORY_PORT", "9091")

	url := fmt.Sprintf("%s://%s:%s/instances/byEntities?branch=master&systemId=%s&entities=%s", scheme, host, port, systemID, entities)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("status: %d message: %s", resp.Status, string(resp.Body))
	}
	err = json.Unmarshal(resp.Body, &summary)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("qtd os summary: ", len(summary))
	return summary, nil
}
