package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/gommon/log"
)

//SplitReprocessingIntoEvents takes a reprocessing e publish all events to reprocessing queue
func SplitReprocessingIntoEvents(reprocessing *models.Reprocessing) error {
	for _, event := range reprocessing.Events {
		originalInstance := event.InstanceID
		event.SystemID = reprocessing.SystemID
		event.InstanceID = ""
		event.Tag = ""
		event.Scope = "reprocessing"
		event.Reprocessing = new(domain.ReprocessingInfo)
		event.Reprocessing.ID = reprocessing.ID
		event.Reprocessing.InstanceID = originalInstance
		event.Reprocessing.Image = event.Image
		event.Reprocessing.SystemID = reprocessing.SystemID
	}

	return PublishReprocessingEvents(reprocessing.Events)
}

//PublishReprocessingEvents publish reprocessing events to reprocessing queue
func PublishReprocessingEvents(events []*domain.Event) error {
	log.Debug("publishing events to reprocessing events")
	publisher := broker.GetPublisher()
	for _, event := range events {
		msg, _ := broker.GetMessageFrom(event)
		if err := publisher.Publish("reprocessing-events", fmt.Sprintf("%s.control_%s", event.SystemID, event.SystemID), msg); err != nil {
			log.Error(fmt.Sprintf("failure to publish event to reprocessing-events exchange "), err)
			return err
		}
	}
	return nil
}
