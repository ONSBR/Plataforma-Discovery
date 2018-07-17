package models

import (
	"github.com/ONSBR/Plataforma-Discovery/helpers"
	"github.com/ONSBR/Plataforma-Discovery/util"
)

type ProcessDependency struct {
}

type Process struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Version  string        `json:"version"`
	SystemID string        `json:"systemId"`
	Sources  []*DataSource `json:"sources"`
}

type DataSource struct {
	Entity    string     `json:"entity"`
	Filter    string     `json:"filter"`
	Processes []*Process `json:"processes"`
}

func (ds *DataSource) AddProcess(dep *helpers.DependencyDomain) {
	if ds.Processes == nil {
		ds.Processes = make([]*Process, 1)
		ds.Processes[0] = &Process{
			ID:       dep.ProcessID,
			Name:     dep.Name,
			SystemID: dep.SystemID,
			Version:  dep.Version,
		}
		return
	}
	ds.Processes = append(ds.Processes, &Process{
		ID:       dep.ProcessID,
		Name:     dep.Name,
		SystemID: dep.SystemID,
		Version:  dep.Version,
	})
}

func NewDataSourceChain(dependencyDomain []*helpers.DependencyDomain) []*DataSource {
	set := util.NewStringSet()
	for _, dependency := range dependencyDomain {
		if !set.Exist(dependency.DataSourceKey()) {
			source := &DataSource{
				Entity: dependency.Entity,
				Filter: dependency.Filter,
			}
			source.AddProcess(dependency)
			set.Add(dependency.DataSourceKey(), source)
		} else {
			source := set.Get(dependency.DataSourceKey()).(*DataSource)
			source.AddProcess(dependency)
			set.Add(dependency.DataSourceKey(), source)
		}
	}
	dataSources := make([]*DataSource, set.Len())
	for i, k := range set.List() {
		dataSources[i] = set.Get(k).(*DataSource)
	}
	return dataSources
}
