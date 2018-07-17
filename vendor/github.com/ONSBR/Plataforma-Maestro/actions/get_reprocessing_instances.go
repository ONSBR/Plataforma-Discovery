package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/sdk/discovery"
)

//GetReprocessingInstances from discovery service based on perist event
func GetReprocessingInstances(event *domain.Event) ([]models.ReprocessingUnit, error) {
	if event.InstanceID == "" {
		return nil, fmt.Errorf("event %s should have instance id", event.Name)
	}
	if event.SystemID == "" {
		return nil, fmt.Errorf("event %s should have system id", event.Name)
	}
	entities, err := appdomain.GetEntitiesByProcessInstance(event.SystemID, event.InstanceID)
	if err != nil {
		return nil, err
	}
	return discovery.GetReprocessingInstances(event.SystemID, event.InstanceID, entities)
}
