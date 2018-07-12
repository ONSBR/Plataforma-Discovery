package helpers

import (
	"time"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Discovery/util"
)

type ProcessInstance struct {
	ID string `json:"id"`
}

func GetFinalizedInstancesAfter(systemID, instanceID string, startedAfter time.Time) ([]*ProcessInstance, error) {
	// /executedAfter
	instances := make([]*ProcessInstance, 0)
	err := apicore.Query(apicore.Filter{
		Entity: "processInstance",
		Map:    "core",
		Name:   "executedAfter",
		Params: []apicore.Param{
			apicore.Param{
				Key:   "systemId",
				Value: systemID,
			},
			apicore.Param{
				Key:   "instanceId",
				Value: instanceID,
			},
			apicore.Param{
				Key:   "date",
				Value: util.ToISOString(startedAfter),
			},
		},
	}, &instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}
