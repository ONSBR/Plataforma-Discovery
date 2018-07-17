package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

//MountReprocessingInfraBySystem mounts rabbitmq infra to handle reprocessing by system
func MountReprocessingInfraBySystem(systemID string) error {
	defer declareMut.Unlock()
	declareMut.Lock()
	queue := fmt.Sprintf(models.ReprocessingQueue, systemID)
	routingKey := fmt.Sprintf("#.%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing", queue, routingKey); err != nil {
		return err
	}
	queue = fmt.Sprintf(models.ReprocessingEventsQueue, systemID)
	routingKey = fmt.Sprintf("#.%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing-events", queue, routingKey); err != nil {
		return err
	}

	queue = fmt.Sprintf(models.ReprocessingEventsControlQueue, systemID)
	routingKey = fmt.Sprintf("#.control_%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing-events", queue, routingKey); err != nil {
		return err
	}

	queue = fmt.Sprintf(models.ReprocessingErrorQueue, systemID)
	routingKey = fmt.Sprintf("#.error_%s.#", systemID)
	if err := broker.DeclareQueue("reprocessing", queue, routingKey); err != nil {
		return err
	}
	return nil
}
