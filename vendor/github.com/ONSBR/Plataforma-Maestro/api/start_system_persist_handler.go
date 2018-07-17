package api

import (
	"github.com/ONSBR/Plataforma-Maestro/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func startSystemPersistHandler(c echo.Context) error {
	log.Info("starting persist handling on queue: ", c.QueryParam("queue"))
	err := handlers.SubscribeToReceiveEventsBySystem(c.QueryParam("queue"))
	if err != nil {
		return err
	}
	c.NoContent(200)
	return nil
}
