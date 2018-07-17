package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/ONSBR/Plataforma-EventManager/domain"
	"github.com/ONSBR/Plataforma-EventManager/infra"
	"github.com/PMoneda/http"
)

func getUrl() string {
	return fmt.Sprintf("%s://%s:%s/core",
		infra.GetEnv("APICORE_SCHEME", "http"),
		infra.GetEnv("APICORE_HOST", "localhost"),
		infra.GetEnv("APICORE_PORT", "9110"))
}

//GetSolutionIDByEventName on apicore
func GetSolutionIDByEventName(eventName string) (string, error) {
	url := fmt.Sprintf("%s/operation?filter=byEvent&event=%s", getUrl(), eventName)

	arr := make([]*domain.Operation, 0, 0)

	if err := http.GetJSON(url, &arr); err != nil {
		return "", err
	}
	if len(arr) == 0 {
		return "", infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", eventName))
	}
	return arr[0].SystemID, nil
}

//EventBindings any operation that have inbound or outbound event binding with same name
func EventBindings(eventName string) ([]*domain.Operation, error) {
	url := fmt.Sprintf("%s/operation?filter=bindingEvent&event=%s", getUrl(), eventName)
	arr := make([]*domain.Operation, 0, 0)
	err := http.GetJSON(url, &arr)
	return arr, err
}

//EventHasBindings verify if exist any operation that have inbound or outbound event binding with same name
func EventHasBindings(eventName string) (bool, error) {
	arr, err := EventBindings(eventName)
	if len(arr) == 0 {
		return false, infra.NewSubscriberNotFoundException(fmt.Sprintf("Event %s has no subscribers", eventName))
	}
	if err != nil {
		return false, infra.NewComponentException(err.Error())
	}
	return true, nil
}

//UpdateProcessInstance updates process instance status on apicore
func UpdateProcessInstance(instanceID, status string) error {

	procInstance := NewProcessInstance()
	procInstance.ID = instanceID
	procInstance.Status = status
	procInstance.Metadata.ChangeTrack = "update"

	url := fmt.Sprintf("%s/persist", getUrl())
	arr := make([]*ProcessInstance, 1, 1)
	arr[0] = procInstance
	if _, err := http.Post(url, arr); err != nil {
		return err
	}
	return nil
}

//GetOpenBranchesBySystem returns all branches open on apicore
func GetOpenBranchesBySystem(systemID string) ([]Branch, error) {
	url := fmt.Sprintf("%s/branch?filter=bySystemIdAndStatus&systemId=%s&status=open", getUrl(), systemID)
	if response, err := http.Get(url); err != nil {
		return nil, infra.NewArgumentException(err.Error())
	} else {
		branches := make([]Branch, 0, 0)
		if err := json.Unmarshal(response.Body, &branches); err != nil {
			return nil, infra.NewComponentException(err.Error())
		}
		return branches, err
	}

}
