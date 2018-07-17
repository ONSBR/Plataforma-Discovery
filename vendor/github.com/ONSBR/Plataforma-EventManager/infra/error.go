package infra

import (
	"fmt"
)

var mapErrorTagHTTPStatus = map[string]int{RunningReprocessing: 200, PlatformLocked: 401, PlatformComponentError: 500, SystemError: 500, InvalidArguments: 400, PersistEventQueueEmpty: 201, SubscriberNotFound: 404}

//var mapErrorTagHTTPStatus = map[string]int{RunningReprocessing: 200, PlatformLocked: 200, PlatformComponentError: 200, SystemError: 200, InvalidArguments: 200, PersistEventQueueEmpty: 200, SubscriberNotFound: 200}

//PlatformComponentError is a error tag to map error on platform components
const PlatformComponentError = "platform_component_error"

//PlatformLocked is a error tag to map reprocessing running on platform
const PlatformLocked = "platform_locked"

//SystemError is a error tag to map general errors
const SystemError = "system_error"

//InvalidArguments is a error tag to map bad request params
const InvalidArguments = "invalid_arguments"

const PersistEventQueueEmpty = "empty_queue"

const RunningReprocessing = "running_reprocessing"

const SubscriberNotFound = "subscriber_not_found"

//Exception is a basic struct for errors
type Exception struct {
	Tag     string `json:"error_code"`
	Message string `json:"message"`
}

//NewException creates a new Exception Object
func NewException(tag, message string) *Exception {
	ex := Exception{Tag: tag, Message: message}
	//log.Error(fmt.Sprintf("%s %s", tag, message))
	return &ex
}

func NewSubscriberNotFoundException(message string) *Exception {
	return NewException(SubscriberNotFound, message)
}

func NewPlatformLockedException(message string) *Exception {
	return NewException(PlatformLocked, message)
}

func NewEmptyQueueException(message string) *Exception {
	return NewException(PersistEventQueueEmpty, message)
}

func NewRunningReprocessingException(message string) *Exception {
	return NewException(RunningReprocessing, message)
}

func NewSystemException(message string) *Exception {
	return NewException(SystemError, message)
}

func NewArgumentException(message string) *Exception {
	return NewException(InvalidArguments, message)
}

func NewComponentException(message string) *Exception {
	return NewException(PlatformComponentError, message)
}

func (e *Exception) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

//HTTPStatus returns equivalent http status and tag error
func (e *Exception) HTTPStatus() int {
	st, ok := mapErrorTagHTTPStatus[e.Tag]
	if ok {
		return st
	}
	return 500
}

//BindHTTPStatusToErrorTag binding custom http status to error tag
func BindHTTPStatusToErrorTag(tag string, status int) {
	mapErrorTagHTTPStatus[tag] = status
}
