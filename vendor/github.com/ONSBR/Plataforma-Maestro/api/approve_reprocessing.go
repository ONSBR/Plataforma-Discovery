package api

import (
	"sync"

	"github.com/ONSBR/Plataforma-Maestro/actions"
	"github.com/labstack/echo"
)

var lock sync.Mutex

//approveReprocessing approve reprocessing
func approveReprocessing(c echo.Context) error {
	defer lock.Unlock()
	lock.Lock()

	approver := new(approve)
	if err := c.Bind(approver); err != nil {
		return err
	}
	data, err := actions.ApproveReprocessing(c.Param("id"), approver.User)
	if err != nil {
		return err
	}
	return c.JSON(202, data)
}
