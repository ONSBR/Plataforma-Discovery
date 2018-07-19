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
	for _, entity := range entities {
		typ, err := helpers.ExtractFieldFromEntity(entity, "type")
		if err != nil {
			return nil, err
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
	if len(instancesAfter) == 0 {
		log.Info("No instances found after ", util.ToISOString(util.TimeFromMilliTimestamp(timestamp)))
		return &models.AnalyticsResult{Units: []models.ReprocessingUnit{}}, nil
	}
	log.Debug("Qtd instances after ", util.ToISOString(util.TimeFromMilliTimestamp(timestamp)), " = ", len(instancesAfter))
	instancesStr := make([]string, len(instancesAfter))
	for i, ins := range instancesAfter {
		instancesStr[i] = ins.ID
		log.Debug("Instance ID=", ins.ID, " EventName=", ins.OriginEventName, " StartExecution=", ins.StartExecution, " Status=", ins.Status, " Scope=", ins.Scope)
	}
	summaries, err := GetSummaryBySystem(systemID, strings.Join(analytics.ListEntitiesTypes(), ","), strings.Join(instancesStr, ","), event.Tag)
	log.Info(len(summaries))
	if err != nil {
		return nil, err
	}
	return dispatchWorker(originInstanceID, entities, summaries), nil
}

func dispatchWorker(originInstanceID string, entities models.EntitiesList, summaries []*models.InstanceSummary) *models.AnalyticsResult {
	result := make(chan *models.AnalyticsResult)
	stack := 0
	groupedTag := util.NewStringSet()
	for _, summary := range summaries {
		if summary.ProcessInstance == originInstanceID {
			//skip same instance summary
			continue
		}
		if groupedTag.Exist(summary.Tag) {
			continue
		}
		groupedTag.Add(summary.Tag)
		go models.RunAnalyticsForInstance(summary.SystemID, summary.ProcessInstance, entities, result, summary.Entities)
		stack++
	}

	if stack == 0 {
		close(result)
		return &models.AnalyticsResult{Units: []models.ReprocessingUnit{}}
	}
	finalList := new(models.AnalyticsResult)
	finalList.Units = make([]models.ReprocessingUnit, 0)
	finalSet := util.NewStringSet()
	for r := range result {
		for _, un := range r.Units {
			key := fmt.Sprintf("%s:%s", un.Branch, un.InstanceID)
			if !finalSet.Exist(key) {
				finalList.Units = append(finalList.Units, un)
			}
		}
		stack--
		if stack == 0 {
			break
		}
	}
	close(result)
	return finalList
}
