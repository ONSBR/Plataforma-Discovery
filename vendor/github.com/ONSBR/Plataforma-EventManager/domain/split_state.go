package domain

//SplitState is a struct to control splited events execution state
type SplitState struct {
	Tag    string        `json:"tag"`
	Events []*EventState `json:"events"`
}

//Pending is a pending status
const Pending = "pending"

//Error is a error status when events ends with error or exception event
const Error = "error"

//Success is a succes status when a event is executed with no errros with done event
const Success = "success"

//EventState represents an execution state for some event
type EventState struct {
	Name     string `json:"name"`
	EventOut string `json:"eventOut"`
	Version  string `json:"version"`
	Tag      string `json:"tag"`
	Status   string `json:"status"`
	Branch   string `json:"branch"`
	Scope    string `json:"scope"`
}

//AddEventState to event state list
func (split *SplitState) AddEventState(evt *EventState) {
	split.Events = append(split.Events, evt)
}

//IsComplete returns true if all event state has status equal success
func (split *SplitState) IsComplete() bool {
	complete := true
	for _, evt := range split.Events {
		complete = complete && evt.Status == Success
	}
	return complete
}

//NewSplitState creates a new splitstate instance
func NewSplitState() *SplitState {
	split := new(SplitState)
	split.Events = make([]*EventState, 0, 0)
	return split
}
