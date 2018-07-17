package actions

import (
	"fmt"
	"os"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"
	"github.com/labstack/gommon/log"
)

//SubmitReprocessingToApprove block persistence on domain until reprocessing will be approve
func SubmitReprocessingToApprove(persistEvent *domain.Event, instances []models.ReprocessingUnit) (err error) {
	reprocessing := models.NewReprocessing(persistEvent)
	reprocessing.PendingApproval()
	origin, err := processmemory.GetEventByInstance(persistEvent.InstanceID)
	if err != nil {
		log.Error(err)
		return
	}
	reprocessing.Origin = origin

	events, err := getEventsFromInstances(instances)
	if err != nil {
		log.Error(err)
		return
	}
	reprocessing.Events = events
	err = sdk.SaveDocument("reprocessing", reprocessing)
	if err != nil {
		log.Error(err)
		return
	}
	err = sdk.UpdateProcessInstance(persistEvent.InstanceID, "reprocessing_pending_approval")
	if err != nil {
		log.Error(err)
		return
	}
	if os.Getenv(fmt.Sprintf("AUTO_REPROCESSING_%s", persistEvent.SystemID)) != "" {
		go ApproveReprocessing(reprocessing.ID, "platform")
	}
	return
}

func getEventsFromInstances(instances []models.ReprocessingUnit) ([]*domain.Event, error) {
	events := make([]*domain.Event, 0)
	for _, instance := range instances {
		evt, err := processmemory.GetEventByInstance(instance.InstanceID)
		if err != nil {
			return nil, err
		}
		evt.Branch = instance.Branch
		events = append(events, evt)
	}
	return events, nil
}
