package discovery

import (
	"fmt"

	"github.com/PMoneda/http"

	"github.com/ONSBR/Plataforma-Deployer/env"

	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
)

//GetReprocessingInstances from discovery service
func GetReprocessingInstances(systemID, instanceID string, entites appdomain.EntitiesList) ([]models.ReprocessingUnit, error) {
	//TODO get data from Discovey
	scheme := env.Get("DISCOVERY_SCHEME", "http")
	host := env.Get("DISCOVERY_HOST", "localhost")
	port := env.Get("DISCOVERY_PORT", "8090")

	url := fmt.Sprintf("%s://%s:%s/v1.0.0/discovery/entities?systemID=%s&instanceID=%s", scheme, host, port, systemID, instanceID)
	units := make([]models.ReprocessingUnit, 0)
	err := http.GetJSON(url, &units)
	if err != nil {
		return nil, err
	}
	return units, nil
}
