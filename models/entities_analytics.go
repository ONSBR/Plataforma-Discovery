package models

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Discovery/db"
	"github.com/labstack/gommon/log"

	"github.com/ONSBR/Plataforma-Discovery/helpers"
	"github.com/ONSBR/Plataforma-Discovery/util"
)

//EntitiesList maps entities that domain app will save based on process memory
type EntitiesList []map[string]interface{}

//EntitiesAnalytics manages all information to respond wich instances should be reprocess
type EntitiesAnalytics struct {
	entitiesSet *util.StringSet

	queryTreeSet *util.StringTreeSet

	groupedEntitiesByType map[string]EntitiesList

	queryMap map[string][]string

	queryIDS map[string]*util.StringSet
}

type AnalyticsResult struct {
	Units []ReprocessingUnit
}

type ReprocessingUnit struct {
	InstanceID string `json:"instanceId"`
	Branch     string `json:"branch"`
	Forking    bool   `json:"forking"`
}

type InstanceSummary struct {
	SystemID        string           `json:"systemId"`
	ProcessInstance string           `json:"process"`
	Tag             string           `json:"tag"`
	Version         string           `json:"version"`
	ProcessApp      string           `json:"processAppId"`
	Reprocessable   bool             `json:"reprocessable"`
	Branch          string           `json:"branch"`
	IdempotencyKey  string           `json:"idempotencyKey"`
	Entities        []*EntitySummary `json:"entities"`
}

type EntitySummary struct {
	EntityName  string                 `json:"name"`
	Parameters  map[string]interface{} `json:"parameters"`
	Query       string                 `json:"query"`
	EntitiesIds []EntityID             `json:"data"`
}

type EntityID struct {
	ID  string `json:"id"`
	RID string `json:"rid"`
}

type PostgresRowData struct {
	Branch         string
	RID            string
	MetaInstanceID string
}

//NewEntitiesAnalytics creates a new EntitiesAnalytics object
func NewEntitiesAnalytics() *EntitiesAnalytics {
	analytics := new(EntitiesAnalytics)
	analytics.entitiesSet = util.NewStringSet()
	analytics.queryTreeSet = util.NewStringTreeSet()
	analytics.groupedEntitiesByType = make(map[string]EntitiesList)
	analytics.queryMap = make(map[string][]string)
	analytics.queryIDS = make(map[string]*util.StringSet)
	return analytics
}

//AddEntity map all types to analyze
func (analytic *EntitiesAnalytics) AddEntity(entity string) {
	analytic.entitiesSet.Add(entity)
}

func (analytic *EntitiesAnalytics) MapEntityToQuery(entitiesSummary []*EntitySummary) {
	for _, en := range entitiesSummary {
		query := helpers.ParseQuery(en.Query, en.Parameters)
		_, ok := analytic.queryMap[en.EntityName]
		if !ok {
			analytic.queryMap[en.EntityName] = make([]string, 0)
		}
		if query != "" {
			query = fmt.Sprintf(`(select rid, branch, meta_instance_id from %s where rid=$1 and branch in ($2,'master') and %s) union (select rid, branch, meta_instance_id from %s where from_id=$3 and %s)`, en.EntityName, query, en.EntityName, query)
			analytic.queryMap[en.EntityName] = append(analytic.queryMap[en.EntityName], query)
			analytic.queryIDS[query] = util.NewStringSet()
			for _, id := range en.EntitiesIds {
				set := analytic.queryIDS[query]
				set.Add(id.RID)
			}

		}

	}
}

func (analytic *EntitiesAnalytics) SearchOnPostgres(systemID string, obj map[string]interface{}) *util.StringSet {
	t, _ := helpers.ExtractFieldFromEntity(obj, "type")
	rid, _ := helpers.ExtractFieldFromEntity(obj, "rid")
	set := util.NewStringSet()
	branchAttr, ok := obj["branch"]
	branch := ""
	if !ok {
		log.Error("not found branch field on entity ", obj)
		return set
	}
	switch t := branchAttr.(type) {
	case string:
		branch = t
	default:
		log.Error("branch attribute is not string: ", t)
		return set
	}
	queries := analytic.queryMap[t]

	if rid == "" {
		return set
	}
	for _, query := range queries {
		found := false
		db.Query(systemID, func(scan db.Scan) {
			var row PostgresRowData
			scan(&row.RID, &row.Branch, &row.MetaInstanceID)
			found = true
			if row.Branch != "master" {
				set.Add(row.Branch)
			} else {
				set.Add(branch)
			}
		}, query, rid, branch, rid)
		if !found {
			ridSet := analytic.queryIDS[query]
			if ridSet.Exist(rid) {
				set.Add(branch)
			}
		}
		log.Debug("query= ", query, " rid=", rid, " branch=", branch)
	}
	return set
}

//ListEntitiesTypes list all entity type on persistence
func (analytic *EntitiesAnalytics) ListEntitiesTypes() []string {
	return analytic.entitiesSet.List()
}

func RunAnalyticsForInstance(systemID, instanceID string, entities EntitiesList, channel chan *AnalyticsResult, entitiesSummary []*EntitySummary) {
	analytics := NewEntitiesAnalytics()
	analytics.MapEntityToQuery(entitiesSummary)
	finalSet := util.NewStringSet()
	r := AnalyticsResult{Units: []ReprocessingUnit{}}
	isForkInstance, err := isFork(instanceID)
	if err != nil {
		log.Error(err)
		channel <- &r
		return
	}
	for _, obj := range entities {
		set := analytics.SearchOnPostgres(systemID, obj)
		if set.Len() > 0 {
			if isForkInstance {
				log.Info(instanceID, " é uma instancia fork")
				finalSet.Add("master")
				r.Units = append(r.Units, ReprocessingUnit{InstanceID: instanceID, Branch: "master", Forking: true})
				break
			}
			//registros impactados incluindo em branches impactadas
			for _, branch := range set.List() {
				if !finalSet.Exist(branch) {
					r.Units = append(r.Units, ReprocessingUnit{InstanceID: instanceID, Branch: branch})
					finalSet.Add(branch)
				}
			}
		}
	}
	if finalSet.Len() > 0 {
		log.Info("total de reprocessamentos da instancia: ", finalSet.Len())
		channel <- &r
		return
	}
	channel <- &AnalyticsResult{Units: []ReprocessingUnit{}}
}

func isFork(instanceID string) (bool, error) {
	list := make([]map[string]interface{}, 0)
	err := apicore.FindByID("processInstance", instanceID, &list)
	if err != nil {
		return false, err
	}
	if len(list) == 0 {
		return false, fmt.Errorf("instance not found")
	}
	isFork, ok := list[0]["isFork"]
	if ok && isFork != nil && isFork.(bool) {
		return true, nil
	}
	return false, nil
}
