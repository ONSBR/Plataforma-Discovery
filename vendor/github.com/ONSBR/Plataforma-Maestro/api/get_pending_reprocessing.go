package api

import (
	"encoding/json"

	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/labstack/echo"
)

func getPendingReprocessing(c echo.Context) error {
	j, err := sdk.GetDocument("reprocessing", map[string]string{"systemId": c.Param("systemId"), "status": "pending_approval"})
	if err != nil {
		return err
	}
	data := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(j), &data)
	if err != nil {
		return err
	}
	for _, item := range data {
		delete(item, "_id")
	}
	return c.JSON(200, data)
}
