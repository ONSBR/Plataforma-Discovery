package helpers

import (
	"fmt"
	"sync"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
)

type AppInfo struct {
	Host     string `json:"host"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Port     int64  `json:"port"`
	SystemID string `json:"systemId"`
	Type     string `json:"type"`
}

var cache map[string]*AppInfo
var once sync.Once

func GetDomainAppInfo(systemID string) (*AppInfo, error) {
	once.Do(func() {
		cache = make(map[string]*AppInfo)
	})
	info, ok := cache[systemID]
	if ok {
		return info, nil
	}
	result := make([]*AppInfo, 0)
	filter := apicore.Filter{
		Entity: "installedApp",
		Map:    "core",
		Name:   "bySystemIdAndType",
		Params: []apicore.Param{apicore.Param{
			Key:   "systemId",
			Value: systemID,
		}, apicore.Param{
			Key:   "type",
			Value: "domain",
		},
		},
	}
	err := apicore.Query(filter, &result)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	if len(result) > 0 {
		cache[systemID] = result[0]
		return result[0], nil
	}
	return nil, fmt.Errorf("no app found for %s id", systemID)
}

func GetDomainHost(systemID string) (string, error) {
	info, err := GetDomainAppInfo(systemID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d", info.Host, info.Port), nil
}

func GetDBName(systemID string) (string, error) {
	info, err := GetDomainAppInfo(systemID)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}
