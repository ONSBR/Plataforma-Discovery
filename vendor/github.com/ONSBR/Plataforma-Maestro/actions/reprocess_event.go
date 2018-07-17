package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/ONSBR/Plataforma-Maestro/sdk/eventmanager"
	"github.com/labstack/gommon/log"
)

//ReprocessEvent takes a event to reprocess by system and emit event to event manager
func ReprocessEvent(systemID string) {
	log.Debug(fmt.Sprintf("Reprocessing event from %s", systemID))
	picker := broker.GetPicker()
	context, ok, err := picker.Pick(fmt.Sprintf(models.ReprocessingEventsQueue, systemID))
	if err != nil {
		log.Error("cannot pick event to reprocess ", err)
		return
	}
	if !ok {
		err := FinishReprocessing(systemID)
		if err != nil {
			log.Error(err)
		}
		return
	}
	if context == nil {
		log.Error("Context is nil cannot proceed")
		return
	}
	defer context.Ack()
	log.Debug("getting event to reprocess: ", len(context.Message.Data), " bytes")

	event, err := models.GetEventFromContext(context)
	if err != nil {
		log.Error(err)
		return
	}
	err = eventmanager.Push(event)
	if err != nil {
		log.Error("cannot push event to event manager: ", err)
	}
	log.Debug(fmt.Sprintf("event %s publish to event manager", event.Name))
}
