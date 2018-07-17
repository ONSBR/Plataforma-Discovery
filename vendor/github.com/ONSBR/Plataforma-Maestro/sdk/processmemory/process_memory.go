package processmemory

import (
	"fmt"

	"github.com/PMoneda/http"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//GetEventByInstance returns event from process memory
func GetEventByInstance(instanceID string) (*domain.Event, error) {
	evts := make([]*domain.Event, 0)
	scheme := env.Get("PROCESS_MEMORY_SCHEME", "http")
	host := env.Get("PROCESS_MEMORY_HOST", "localhost")
	port := env.Get("PROCESS_MEMORY_PORT", "9091")
	url := fmt.Sprintf("%s://%s:%s/%s/event", scheme, host, port, instanceID)
	err := http.GetJSON(url, &evts)
	if err != nil {
		return nil, err
	}
	if len(evts) > 0 {
		return evts[0], nil
	}
	return nil, fmt.Errorf("event not found for instance %s", instanceID)
}
