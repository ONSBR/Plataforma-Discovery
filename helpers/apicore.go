package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Discovery/util"
)

type ProcessInstance struct {
	ID        string `json:"id"`
	ProcessId string `json:"processId"`
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

type DependencyDomain struct {
	ID        string `json:"id"`
	SystemID  string `json:"systemId"`
	ProcessID string `json:"processId"`
	Entity    string `json:"entity"`
	Version   string `json:"version"`
	Name      string `json:"name"`
	Filter    string `json:"filter"`
}

func (dep *DependencyDomain) ProcessKey() string {
	return fmt.Sprintf("%s/%s/%s", dep.SystemID, dep.ProcessID, dep.Version)
}

func (dep *DependencyDomain) DataSourceKey() string {
	return fmt.Sprintf("%s/%s/%s", dep.SystemID, dep.Entity, dep.Filter)
}

func GetProcessesWithDependsOn(systemID string, entities []string) ([]*DependencyDomain, error) {
	instances := make([]*DependencyDomain, 0)
	err := apicore.Query(apicore.Filter{
		Entity: "dependencyDomain",
		Map:    "core",
		Name:   "byEntites",
		Params: []apicore.Param{
			apicore.Param{
				Key:   "entities",
				Value: strings.Join(entities, ";"),
			},
			apicore.Param{
				Key:   "systemId",
				Value: systemID,
			},
		},
	}, &instances)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func GetProcessInstanceById(instanceID string) (*ProcessInstance, error) {
	list := make([]*ProcessInstance, 0)
	err := apicore.FindByID("processInstance", instanceID, &list)
	if err != nil {
		return nil, err
	}
	return list[0], nil
}
