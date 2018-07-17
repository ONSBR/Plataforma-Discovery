package handlers

import (
	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/PMoneda/carrot"
)

//SubscribeToReceiveEventsBySystem subscribe on broker to starting receiving persist events isolated by system id
func SubscribeToReceiveEventsBySystem(queue string) error {
	subs := broker.GetSubscriber()
	err := subs.Subscribe(carrot.SubscribeWorker{
		Queue:   queue,
		Scale:   1,
		Handler: PersistHandler,
	})
	if err != nil {
		return err
	}
	return nil
}
