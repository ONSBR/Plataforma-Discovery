package domain

//Operation is a event binding config
type Operation struct {
	Event         string `json:"event_in"`
	EventOut      string `json:"event_out"`
	Version       string `json:"version"`
	SystemID      string `json:"systemId"`
	ProcessID     string `json:"processId"`
	Reprocessable bool   `json:"reprocessable"`
}

//OperationInstance hanlde each instance of an operation executed on platform
type OperationInstance struct {
	ProcessID         string `json:"processId,omitempty"`
	SystemID          string `json:"systemId,omitempty"`
	ProcessInstanceID string `json:"processInstanceId,omitempty"`
	OperationID       string `json:"operationId,omitempty"`
	Status            string `json:"status,omitempty"`
	MustCommit        bool   `json:"mustCommit,omitempty"`
	Image             string `json:"image,omitempty"`
	Version           string `json:"version,omitempty"`
	EventName         string `json:"eventName,omitempty"`
}
