package actions

import (
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/sdk/apicore"
	"github.com/ONSBR/Plataforma-Discovery/util"
)

//GetAvailableBranches returns just open branches on platform
//this is important to filter reprocessing instances by just open branches
func GetAvailableBranches(systemID string) (*util.StringSet, error) {
	result := make([]map[string]interface{}, 0)
	err := apicore.Query(apicore.Filter{
		Entity: "branch",
		Map:    "core",
		Name:   "bySystemIdAndStatus",
		Params: []apicore.Param{
			apicore.Param{
				Key:   "systemId",
				Value: systemID,
			},
			apicore.Param{
				Key:   "status",
				Value: "open",
			},
		},
	}, &result)
	if err != nil {
		return nil, err
	}
	set := util.NewStringSet()
	set.Add("master")
	for _, branch := range result {
		name, ok := branch["name"]
		if ok {
			switch t := name.(type) {
			case string:
				set.Add(t)
			default:
				return nil, fmt.Errorf("invalid type for branch name: expected string")
			}

		} else {
			return nil, fmt.Errorf("invalid contract expected name attribute")
		}
	}
	return set, nil
}
