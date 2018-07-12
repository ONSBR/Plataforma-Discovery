package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ONSBR/Plataforma-Discovery/db"

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
}

type AnalyticsResult struct {
	Units []ReprocessingUnit
}

type ReprocessingUnit struct {
	InstanceID string `json:"instanceId"`
	Branch     string `json:"branch"`
}

type InstanceSummary struct {
	SystemID        string           `json:"systemId"`
	ProcessInstance string           `json:"process"`
	Branch          string           `json:"branch"`
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
	Branch string
	RID    string
}

//NewEntitiesAnalytics creates a new EntitiesAnalytics object
func NewEntitiesAnalytics() *EntitiesAnalytics {
	analytics := new(EntitiesAnalytics)
	analytics.entitiesSet = util.NewStringSet()
	analytics.queryTreeSet = util.NewStringTreeSet()
	analytics.groupedEntitiesByType = make(map[string]EntitiesList)
	analytics.queryMap = make(map[string][]string)
	return analytics
}

//AddEntity map all types to analyze
func (analytic *EntitiesAnalytics) AddEntity(entity string) {
	analytic.entitiesSet.Add(entity)
}

func (analytic *EntitiesAnalytics) MapEntityToQuery(entitiesSummary []*EntitySummary) {
	for _, en := range entitiesSummary {
		for k, v := range en.Parameters {
			query := ""
			switch t := v.(type) {
			case string:
				query = strings.Replace(en.Query, ":"+k, fmt.Sprintf("'%s'", t), -1)
			case float64:
				query = strings.Replace(en.Query, ":"+k, strconv.FormatFloat(t, 'E', -1, 64), -1)
			case int64:
				query = strings.Replace(en.Query, ":"+k, strconv.FormatInt(t, 10), -1)
			case bool:
				query = strings.Replace(en.Query, ":"+k, strconv.FormatBool(t), -1)
			}
			_, ok := analytic.queryMap[en.EntityName]
			if !ok {
				analytic.queryMap[en.EntityName] = make([]string, 0)
			}
			if query != "" {
				query = fmt.Sprintf(`
					(select rid, branch
						from %s
						where rid=$1 and %s)
						union
					(select rid, branch
					from %s
					where from_id=$2 and %s)
					`, en.EntityName, query, en.EntityName, query)
				analytic.queryMap[en.EntityName] = append(analytic.queryMap[en.EntityName], query)
			}
		}
	}
}

func (analytic *EntitiesAnalytics) SearchOnPostgres(systemID string, obj map[string]interface{}) *util.StringSet {
	t, _ := helpers.ExtractFieldFromEntity(obj, "type")
	rid, _ := helpers.ExtractFieldFromEntity(obj, "rid")
	queries := analytic.queryMap[t]
	set := util.NewStringSet()
	for _, query := range queries {
		db.Query(systemID, func(scan db.Scan) {
			var row PostgresRowData
			scan(&row.RID, &row.Branch)
			set.Add(row.Branch)
		}, query, rid, rid)
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
	for _, obj := range entities {
		set := analytics.SearchOnPostgres(systemID, obj)
		if set.Len() > 0 {
			//registros impactados incluindo em branches impactadas
			r := AnalyticsResult{Units: []ReprocessingUnit{}}
			for _, branch := range set.List() {
				r.Units = append(r.Units, ReprocessingUnit{InstanceID: instanceID, Branch: branch})
			}
			channel <- &r
			return
		}
	}
	channel <- &AnalyticsResult{Units: []ReprocessingUnit{}}
}
