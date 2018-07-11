package api

import (
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
	c.JSON(200, result.Units)
	return nil
}
