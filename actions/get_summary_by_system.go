package actions

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/PMoneda/http"
)

type InstanceSummary struct {
	SystemID        string           `json:"systemId"`
	ProcessInstance string           `json:"process"`
	Entities        []*EntitySummary `json:"entities"`
}

type EntitySummary struct {
	EntityName  string                 `json:"name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Query       string                 `json:"query"`
	EntitiesIds []string               `json:"data"`
}

//GetSummaryBySystem returns all instances query summary from process memory
func GetSummaryBySystem(systemID, entities string) ([]*InstanceSummary, error) {
	summary := make([]*InstanceSummary, 0)
	scheme := env.Get("PROCESS_MEMORY_SCHEME", "http")
	host := env.Get("PROCESS_MEMORY_HOST", "localhost")
	port := env.Get("PROCESS_MEMORY_PORT", "9091")

	url := fmt.Sprintf("%s://%s:%s/instances/byEntities?systemId=%s&entities=%s", scheme, host, port, systemID, entities)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("status: %d message: %s", resp.Status, string(resp.Body))
	}
	err = json.Unmarshal(resp.Body, &summary)
	if err != nil {
		return nil, err
	}
	return summary, nil
}
