package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"
)

//GetEventsFromInstances returns all events from process memory by process instance
func GetEventsFromInstances(instances []models.ReprocessingUnit) ([]*domain.Event, error) {
	list := make([]*domain.Event, 0)
	for _, unit := range instances {
		evt, err := processmemory.GetEventByInstance(unit.InstanceID)
		if err != nil {
			return nil, err
		}
		evt.Branch = unit.Branch
		list = append(list, evt)
	}
	return list, nil
}
