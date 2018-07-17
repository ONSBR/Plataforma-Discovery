package api

import (
	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/echo"
)

func eventGatekeeper(c echo.Context) error {
	rep, err := models.GetReprocessingBySystemIDWithStatus(c.Param("systemId"), models.Running)
	if err != nil {
		return err
	}
	if len(rep) == 0 {
		return c.NoContent(200)
	}
	return c.NoContent(403)
}
