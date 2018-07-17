package sdk

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/PMoneda/http"
	log "github.com/sirupsen/logrus"
)

//EventCanProceed to be processed by plataform
func EventCanProceed(event *domain.Event) bool {
	scheme := env.Get("MAESTRO_SCHEME", "http")
	host := env.Get("MAESTRO_HOST", "localhost")
	port := env.Get("MAESTRO_PORT", "6971")
	url := fmt.Sprintf("%s://%s:%s/v1.0.0/gateway/%s/proceed", scheme, host, port, event.SystemID)
	log.Debug(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return false
	}
	return resp.Status == 200
}

//InitPersistHandler init maestro persist handling for system
func InitPersistHandler(event *domain.Event) bool {
	scheme := env.Get("MAESTRO_SCHEME", "http")
	host := env.Get("MAESTRO_HOST", "localhost")
	port := env.Get("MAESTRO_PORT", "6971")
	url := fmt.Sprintf("%s://%s:%s/v1.0.0/handler/persist?queue=%s", scheme, host, port, fmt.Sprintf("persist.%s.queue", event.SystemID))
	log.Debug(url)
	resp, err := http.Post(url, nil)
	if err != nil {
		log.Error(err)
		return false
	}
	return resp.Status == 200
}
