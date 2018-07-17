package actions

import (
	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"
)

//ProceedToCommit process commiting to domain events by solution
func ProceedToCommit(event *domain.Event) error {
	if err := appdomain.PersistEntitiesByInstance(event.SystemID, event.InstanceID); err != nil {
		return err
	}
	return nil
}
