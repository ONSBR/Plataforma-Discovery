package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

//TODO maybe will be better put this events on apicore
var SystemEvents = []string{
	"system.reprocessing.error",
	"system.executor.enable.debug",
	"system.executor.disable.debug",
	"system.process.persist.error",
	"system.events.reprocessing.request",
	"system.events.reproduction.request",
	"system.deploy.finished",
}

//Events ia a lista of pointers to Event that implements sort interface
type Events []*Event

func (s Events) Len() int {
	return len(s)
}
func (s Events) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Events) Less(i, j int) bool {
	a := s[i]
	b := s[j]
	ta, err := a.GetTimestamp()
	if err != nil {
		return false
	}
	tb, err := b.GetTimestamp()
	if err != nil {
		return false
	}

	return tb.After(ta)
}

//Event define a basic platform event contract
type Event struct {
	Timestamp      string                 `json:"timestamp"`
	Branch         string                 `json:"branch"`
	SystemID       string                 `json:"systemId,omitempty"`
	Name           string                 `json:"name,omitempty"`
	Version        string                 `json:"version,omitempty"`
	Image          string                 `json:"image,omitempty"`
	IdempotencyKey string                 `json:"idempotencyKey"`
	Tag            string                 `json:"tag"`
	AppOrigin      string                 `json:"appOrigin,omitempty"`
	Owner          string                 `json:"owner,omitempty"`
	InstanceID     string                 `json:"instanceId,omitempty"`
	Scope          string                 `json:"scope,omitempty"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
	Reproduction   map[string]interface{} `json:"reproduction,omitempty"`
	Reprocessing   *ReprocessingInfo      `json:"reprocessing,omitempty"`
	Bindings       []*Operation           `json:"-"`
}

//ReprocessingInfo store all reprocessing information on event
type ReprocessingInfo struct {
	ID         string `json:"id,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
	SystemID   string `json:"system_id,omitempty"`
	Image      string `json:"image,omitempty"`
}

//NewEvent creates a new Event Instance
func NewEvent() *Event {
	event := new(Event)
	event.Bindings = make([]*Operation, 0, 0)
	event.Payload = make(map[string]interface{})
	return event
}

//IsEndingEvent checkis event's name ends with .done, .error or .exception
func (e *Event) IsEndingEvent() bool {
	return strings.HasSuffix(e.Name, ".done") || strings.HasSuffix(e.Name, ".error") || strings.HasSuffix(e.Name, ".exception")
}

//ToCeleryMessage transform event to a celery compatible message
func (e *Event) ToCeleryMessage() *CeleryMessage {
	return getCeleryMessage(e)
}

//GetTimestamp returns a parsed timestamp event value
func (e *Event) GetTimestamp() (time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04:05.999", e.Timestamp)
	return t, err
}

//ApplyDefaultFields apply default fields for branch, scope and tag
func (e *Event) ApplyDefaultFields() {
	if e.Branch == "" {
		e.Branch = "master"
	}
	if e.Scope == "" {
		e.Scope = "execution"
	}
	if e.Tag == "" {
		u, _ := uuid.NewUUID()
		e.Tag = u.String()
	}
}

//IsSystemEvent returns true if this event is a internal platform event
func (e *Event) IsSystemEvent() bool {
	for _, sysEvt := range SystemEvents {
		if sysEvt == e.Name {
			return true
		}
	}
	return false
}

//IsReprocessing returns true if event's scope is reprocesing
func (e *Event) IsReprocessing() bool {
	return e.Scope == "reprocessing"
}

//IsExecution returns true if event's scope is execution
func (e *Event) IsExecution() bool {
	return e.Scope == "" || e.Scope == "execution"
}

//IsReproduction returns true if event's scope is reproduction
func (e *Event) IsReproduction() bool {
	return e.Scope == "reproduction"
}

//ToEventState converts an event to state event
func (e *Event) ToEventState() *EventState {
	return &EventState{
		Branch: e.Branch,
		Name:   e.Name,
		Scope:  e.Scope,
		Status: Pending,
		Tag:    e.Tag,
	}
}
