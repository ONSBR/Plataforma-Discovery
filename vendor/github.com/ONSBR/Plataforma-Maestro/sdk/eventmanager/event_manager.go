package eventmanager

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/PMoneda/http"

	"github.com/ONSBR/Plataforma-EventManager/domain"
)

//Push event to event manager
func Push(event *domain.Event) error {
	scheme := env.Get("EVENT_MANAGER_SCHEME", "http")
	host := env.Get("EVENT_MANAGER_HOST", "localhost")
	port := env.Get("EVENT_MANAGER_PORT", "8081")
	url := fmt.Sprintf("%s://%s:%s/sendevent", scheme, host, port)
	resp, err := http.Put(url, event)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		r := make(map[string]string)
		err := json.Unmarshal(resp.Body, &r)
		if err == nil {
			msg, ok := r["message"]
			if ok {
				return fmt.Errorf(msg)
			}
			return fmt.Errorf(string(resp.Body))
		}
	}
	return nil
}
