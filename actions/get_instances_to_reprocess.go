package actions

import (
	"strings"
	"time"

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

	list, err := run(systemID, instanceID, entities)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func run(systemID, originInstanceID string, entities models.EntitiesList) (*models.AnalyticsResult, error) {
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
	log.Info(instancesAfter)
	summaries, err := GetSummaryBySystem(systemID, strings.Join(analytics.ListEntitiesTypes(), ","))
	if err != nil {
		return nil, err
	}
	return dispatchWorker(originInstanceID, entities, summaries), nil
}

func dispatchWorker(originInstanceID string, entities models.EntitiesList, summaries []*models.InstanceSummary) *models.AnalyticsResult {
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
	toReprocess := models.AnalyticsResult{Units: []models.ReprocessingUnit{}}
	if stack == 0 {
		close(result)
		return &toReprocess
	}
	for r := range result {
		toReprocess.Units = append(toReprocess.Units, r.Units...)
		stack--
		if stack == 0 {
			break
		}
	}
	close(result)
	return &toReprocess
}
