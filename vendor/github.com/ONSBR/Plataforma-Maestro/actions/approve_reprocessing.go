package actions

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/sdk"
	"github.com/ONSBR/Plataforma-Maestro/etc"
	"github.com/ONSBR/Plataforma-Maestro/models"
)

//ApproveReprocessing approve reprocessing
func ApproveReprocessing(reprocessingID, user string) (*models.Reprocessing, error) {
	j, err := sdk.GetDocument("reprocessing", map[string]string{"id": reprocessingID, "status": "pending_approval"})
	data := make([]models.Reprocessing, 0)
	json.Unmarshal([]byte(j), &data)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("reprocessing not found")
	}
	status := "approved"
	data[0].Status = status
	data[0].HistoryStatus = append(data[0].HistoryStatus, models.ReprocessingStatus{User: user, Status: status, Timestamp: etc.GetStrTimestamp()})
	if err := sdk.ReplaceDocument("reprocessing", map[string]string{"id": reprocessingID, "status": "pending_approval"}, data[0]); err != nil {
		return nil, err
	}
	go DispatchReprocessing(data[0])
	return &data[0], nil
}
