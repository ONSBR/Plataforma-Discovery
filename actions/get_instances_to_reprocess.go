package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/ONSBR/Plataforma-EventManager/domain"

	"github.com/ONSBR/Plataforma-Maestro/sdk/processmemory"

	"github.com/ONSBR/Plataforma-Discovery/helpers"
	"github.com/ONSBR/Plataforma-Discovery/models"
	"github.com/ONSBR/Plataforma-Discovery/util"
	"github.com/labstack/gommon/log"
)

//GetInstancesToReprocess returns all instances to reprocess based on systemID and instanceID
//instanceID is the processInstance that is requiring reprocessing
func GetInstancesToReprocess(systemID, instanceID string) (*models.AnalyticsResult, error) {
	entities, err := GetEntitiesByInstance(systemID, instanceID)
	if err != nil {
		return nil, err
	}
	event, err := processmemory.GetEventByInstance(instanceID)
	if err != nil {
		return nil, err
	}
	list, err := run(systemID, instanceID, entities, event)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func run(systemID, originInstanceID string, entities models.EntitiesList, event *domain.Event) (*models.AnalyticsResult, error) {
	analytics := models.NewEntitiesAnalytics()
	timestamp := util.Timestamp(time.Now())
	branches := util.NewStringSet()
	for _, entity := range entities {
		typ, err := helpers.ExtractFieldFromEntity(entity, "type")
		if err != nil {
			return nil, err
		}
		branch := entity["branch"].(string)
		if branch != "master" {
			branches.Add(branch)
		}
		tmp, err := helpers.ExtractModifiedTimestamp(entity)
		if err == nil && tmp <= timestamp {
			timestamp = tmp
		}
		analytics.AddEntity(typ)
	}
	log.Debug("Minimum timestamp: ", util.ToISOString(util.TimeFromMilliTimestamp(timestamp)))
	instancesAfter, err := helpers.GetFinalizedInstancesAfter(systemID, originInstanceID, util.TimeFromMilliTimestamp(timestamp))
	if err != nil {
		return nil, err
	}
	instancesStr := make([]string, len(instancesAfter))
	for i, ins := range instancesAfter {
		instancesStr[i] = ins.ID
	}
	summaries, err := GetSummaryBySystem(systemID, strings.Join(analytics.ListEntitiesTypes(), ","), strings.Join(instancesStr, ","))
	log.Info(len(summaries))
	if err != nil {
		return nil, err
	}
	return dispatchWorker(originInstanceID, entities, summaries, branches), nil
}

func dispatchWorker(originInstanceID string, entities models.EntitiesList, summaries []*models.InstanceSummary, branches *util.StringSet) *models.AnalyticsResult {
	result := make(chan *models.AnalyticsResult)
	stack := 0
	for _, summary := range summaries {
		if summary.ProcessInstance == originInstanceID {
			//skip same instance summary
			continue
		}
		go models.RunAnalyticsForInstance(summary.SystemID, summary.ProcessInstance, entities, result, summary.Entities)
		stack++
	}

	if stack == 0 {
		close(result)
		return &models.AnalyticsResult{Units: []models.ReprocessingUnit{}}
	}
	set := util.NewStringSet()
	for r := range result {
		for i := 0; i < len(r.Units); i++ {
			key := fmt.Sprintf("%s:%s", r.Units[i].InstanceID, r.Units[i].Branch)
			if !set.Exist(key) {
				if r.Units[i].Branch == "master" {
					//Aplica a execução da mesma instancia para os branches que vieram no dataset
					for _, branch := range branches.List() {
						set.Add(fmt.Sprintf("%s:%s", r.Units[i].InstanceID, branch), models.ReprocessingUnit{
							Branch:     branch,
							InstanceID: r.Units[i].InstanceID,
						})
					}
				}
				set.Add(key, r.Units[i])
			}
		}
		stack--
		if stack == 0 {
			break
		}
	}
	close(result)
	toReprocess := models.AnalyticsResult{Units: make([]models.ReprocessingUnit, set.Len())}
	for i, key := range set.List() {
		toReprocess.Units[i] = set.Get(key).(models.ReprocessingUnit)
	}
	return &toReprocess
}
