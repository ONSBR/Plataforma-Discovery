package domain

import (
	"github.com/google/uuid"
)

type CeleryMessage struct {
	Task    string  `json:"task"`
	ID      string  `json:"id"`
	Args    []Event `json:"args"`
	Retries uint    `json:"retries"`
	UTC     bool    `json:"utc"`
}

func getCeleryMessage(event *Event) *CeleryMessage {
	c := new(CeleryMessage)
	c.UTC = true
	c.Args = make([]Event, 1, 1)
	c.Args[0] = *event
	uuid, _ := uuid.NewUUID()
	c.ID = uuid.String()
	c.Task = "tasks.process"
	return c
}
