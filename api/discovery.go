package api

import (
	"encoding/json"

	"github.com/ONSBR/Plataforma-Discovery/actions"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

func discovery(c echo.Context) error {
	systemID := c.QueryParam("systemID")
	instanceID := c.QueryParam("instanceID")
	result, err := actions.GetInstancesToReprocess(systemID, instanceID)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Total de instancias para serem reprocessadas: ", len(result.Units))
	d, _ := json.Marshal(result.Units)
	log.Info(string(d))
	c.JSON(200, result.Units)
	return nil
}
