package actions

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-Maestro/sdk/appdomain"

	"github.com/PMoneda/carrot"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Maestro/broker"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

var cache map[string]bool
var once sync.Once
var declareMut sync.Mutex

//DispatchReprocessing dispatches an approved reprocessing to system queue
func DispatchReprocessing(reprocessing models.Reprocessing) {
	once.Do(func() {
		cache = make(map[string]bool)
	})
	if _, ok := cache[reprocessing.SystemID]; !ok {
		if err := MountReprocessingInfraBySystem(reprocessing.SystemID); err != nil {
			log.Error(fmt.Sprintf("cannot mount reprocessing infra to system %s: ", reprocessing.SystemID), err)
			return
		}
	}
	publisher := broker.GetPublisher()
	msg, _ := broker.GetMessageFrom(reprocessing)
	publisher.Publish("reprocessing", reprocessing.SystemID, msg)

	go StartReprocessing(reprocessing.SystemID)
}

//StartReprocessing picks first reprocessing in queue
func StartReprocessing(systemID string) {

	log.Debug("starting reprocessing")
	//Start reprocessing process
	context, proceed, reprocessing := pickReprocessing(systemID)
	if !proceed {
		log.Debug("reprocessing queue is empty or cannot get reprocessing")
		return
	}
	defer context.Nack(true)

	errFnc := func(err error) {
		if err := models.SetStatusReprocessing(reprocessing.ID, reprocessing.Status, ""); err != nil {
			log.Error(fmt.Sprintf("cannot set status aborted:persist-domain-failure on reprocessing %s: ", reprocessing.ID), err)
		}
		if err := context.RedirectTo("reprocessing", fmt.Sprintf("error_%s", reprocessing.ID)); err != nil {
			log.Error(fmt.Sprintf("cannot redirect reprocessing %s to error queue: ", reprocessing.ID), err)
		}
	}

	log.Debug("set reprocessing to running")
	reprocessing.Running()
	if err := models.SaveReprocessing(reprocessing); err != nil {
		log.Error("cannot update reprocessing on process memory: ", err)
		return
	}

	log.Debug("splitting reprocessing events")
	if err := SplitReprocessingIntoEvents(reprocessing); err != nil {
		reprocessing.AbortedSplitEventsFailure()
		log.Error(fmt.Sprintf("cannot split event for reprocessing %s: ", reprocessing.ID), err)
		errFnc(err)
		return
	}
	log.Debug("save pending commit to domain")
	if err := appdomain.PersistEntitiesByInstance(reprocessing.SystemID, reprocessing.PendingEvent.InstanceID); err != nil {
		log.Error("cannot persist pending event on domain: ", err)
		reprocessing.AbortedPersistDomainFailure()
		errFnc(err)
		return
	}

	go ReprocessEvent(systemID)

}

func pickReprocessing(id string) (*carrot.MessageContext, bool, *models.Reprocessing) {
	picker := broker.GetPicker()
	queue := fmt.Sprintf(models.ReprocessingQueue, id)
	context, isReprocessing, err := picker.Pick(queue)
	if err != nil {
		log.Error(err)
		return nil, false, nil
	}
	if isReprocessing {
		reprocessing := &models.Reprocessing{}
		err := json.Unmarshal(context.Message.Data, reprocessing)
		if err != nil {
			log.Error(fmt.Sprintf("cannot unmarshall reprocessing"), err)
			return context, false, nil
		}
		return context, true, reprocessing
	}
	return nil, false, nil
}
