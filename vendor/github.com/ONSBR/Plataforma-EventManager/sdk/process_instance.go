package sdk

type ProcessInstance struct {
	Model
	Status string `json:"status"`
}

func NewProcessInstance() *ProcessInstance {
	procInst := new(ProcessInstance)
	procInst.Metadata.Type = "processInstance"
	return procInst
}
