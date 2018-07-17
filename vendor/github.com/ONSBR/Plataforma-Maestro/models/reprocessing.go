package models

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/PMoneda/carrot"
)

const Running string = "running"
const Finished string = "finished"
const PendingApproval string = "pending_approval"
const Skipped string = "skipped"
const AbortedSplitEventsFail string = "aborted:split-events-failure"
const AbortedPersistDomainFail string = "aborted:persist-domain-failure"
const ReprocessingQueue = "reprocessing.%s.queue"
const ReprocessingEventsQueue = "reprocessing.%s.events.queue"
const ReprocessingEventsControlQueue = "reprocessing.%s.events.control.queue"
const ReprocessingErrorQueue = "reprocessing.%s.error.queue"

type ReprocessingUnit struct {
	Branch     string `json:"branch"`
	InstanceID string `json:"instanceId"`
}

//Reprocessing handle data from discovery service
type Reprocessing struct {
	SystemID      string               `json:"systemId,omitempty"`
	PendingEvent  *domain.Event        `json:"pendingEvent,omitempty"`
	Origin        *domain.Event        `json:"origin,omitempty"`
	ID            string               `json:"id,omitempty"`
	Events        []*domain.Event      `json:"events,omitempty"`
	Status        string               `json:"status"`
	HistoryStatus []ReprocessingStatus `json:"history"`
}

//ReprocessingStatus stores user actions over reprocessing
type ReprocessingStatus struct {
	User      string `json:"user,omitempty"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func (rep *Reprocessing) PendingApproval() {
	rep.SetStatus("", PendingApproval)
}

func (rep *Reprocessing) IsPendingApproval() bool {
	return rep.Status == PendingApproval
}

func (rep *Reprocessing) IsRunning() bool {
	return rep.Status == Running
}

func (rep *Reprocessing) Finish() {
	rep.SetStatus("", Finished)
}

func (rep *Reprocessing) Skipped(owner string) {
	rep.SetStatus(owner, Skipped)
}

func (rep *Reprocessing) Running() {
	rep.SetStatus("", Running)
}

func (rep *Reprocessing) Append(events []*domain.Event) {
	for _, event := range events {
		rep.Events = append(rep.Events, event)
	}
}

func (rep *Reprocessing) AbortedSplitEventsFailure() {
	rep.SetStatus("", AbortedSplitEventsFail)
}

func (rep *Reprocessing) AbortedPersistDomainFailure() {
	rep.SetStatus("", AbortedPersistDomainFail)
}

//SetStatus on reprocessing
func (rep *Reprocessing) SetStatus(owner, status string) {
	rep.Status = status
	st := NewReprocessingStatus(rep.Status)
	st.User = owner
	if rep.HistoryStatus == nil {
		rep.HistoryStatus = make([]ReprocessingStatus, 1)
		rep.HistoryStatus[0] = st
	} else {
		rep.HistoryStatus = append(rep.HistoryStatus, st)
	}
}

func NewReprocessing(pendingEvent *domain.Event) *Reprocessing {
	return &Reprocessing{
		PendingEvent: pendingEvent,
		SystemID:     pendingEvent.SystemID,
		ID:           etc.GetUUID(),
	}
}

//NewReprocessingStatus creates a new status object to reprocessing
func NewReprocessingStatus(status string) ReprocessingStatus {
	return ReprocessingStatus{
		Status:    status,
		Timestamp: etc.GetStrTimestamp(),
	}
}

var statusMut sync.Mutex

//SetStatusReprocessing set status of reprocessing on process memory
func SetStatusReprocessing(reprocessingID string, status, owner string) error {
	rep, err := GetReprocessing(reprocessingID)
	if err != nil {
		return err
	}
	rep.SetStatus(owner, status)
	return SaveReprocessing(rep)
}

//SaveReprocessing saves reprocessing on process memory
func SaveReprocessing(reprocessing *Reprocessing) error {
	defer statusMut.Unlock()
	statusMut.Lock()
	return sdk.ReplaceDocument("reprocessing", map[string]string{"id": reprocessing.ID}, reprocessing)
}

//GetReprocessing return reprocessing from process memory
func GetReprocessing(reprocessingID string) (*Reprocessing, error) {
	defer statusMut.Unlock()
	statusMut.Lock()
	sjson, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessingID})
	if err != nil {
		return nil, err
	}
	rep := make([]*Reprocessing, 1)
	err = json.Unmarshal([]byte(sjson), &rep)
	if err != nil {
		return nil, err
	}
	if len(rep) == 0 {
		return nil, fmt.Errorf(fmt.Sprintf("no reprocessing found with id %s", reprocessingID))
	}
	return rep[0], nil
}

//GetReprocessingBySystemIDWithStatus return reprocessing with systemId and status from process memory
func GetReprocessingBySystemIDWithStatus(systemID, status string) ([]*Reprocessing, error) {
	defer statusMut.Unlock()
	statusMut.Lock()
	sjson, err := sdk.GetDocument("reprocessing", map[string]string{"systemId": systemID, "status": status})
	if err != nil {
		return nil, err
	}
	rep := make([]*Reprocessing, 0)
	err = json.Unmarshal([]byte(sjson), &rep)
	return rep, err
}

//GetStatusOfReprocessing return status of reprocessing from process memory
func GetStatusOfReprocessing(reprocessingID string) (string, error) {
	rep, err := GetReprocessing(reprocessingID)
	if err != nil {
		return "", err
	}
	return rep.Status, nil
}

//GetEventFromCeleryMessage returns an event from celery message
func GetEventFromCeleryMessage(context *carrot.MessageContext) (*domain.Event, error) {
	celeryMessage := new(domain.CeleryMessage)
	err := json.Unmarshal(context.Message.Data, celeryMessage)
	if err != nil {
		return nil, err
	}
	eventParsed := celeryMessage.Args[0]
	return &eventParsed, nil
}

//GetEventFromContext returns an event
func GetEventFromContext(context *carrot.MessageContext) (*domain.Event, error) {
	event := new(domain.Event)
	err := json.Unmarshal(context.Message.Data, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

//GetReprocessingFromContext returns an reprocessing
func GetReprocessingFromContext(context *carrot.MessageContext) (*Reprocessing, error) {
	reprocessing := new(Reprocessing)
	err := json.Unmarshal(context.Message.Data, reprocessing)
	if err != nil {
		return nil, err
	}
	return reprocessing, nil
}
