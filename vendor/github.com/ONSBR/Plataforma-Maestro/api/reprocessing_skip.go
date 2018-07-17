package api

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Maestro/models"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

//reprocessingSkip skips current top reprocessing and suppress commit
func reprocessingSkip(c echo.Context) error {
	defer lock.Unlock()
	lock.Lock()
	approver := new(approve)
	if err := c.Bind(approver); err != nil {
		return err
	}
	rep, err := models.GetReprocessing(c.Param("id"))
	if err != nil {
		return err
	}
	if rep.IsPendingApproval() {
		rep.Skipped(approver.User)
		log.Debug(fmt.Sprintf("reprocessing skipped by %s", approver.User))
		return models.SaveReprocessing(rep)
	}
	return fmt.Errorf("You cannot skip a not pending approval reprocessing")
}
