package apicore

//FindByID finds entity by ID on apicore
func FindByID(entity, id string, response interface{}) error {
	filter := Filter{
		Entity: entity,
		Map:    "core",
		Name:   "byId",
		Params: []Param{Param{
			Key:   "id",
			Value: id,
		},
		},
	}
	return Query(filter, response)
}

//FindBySystemID finds entity by SystemID on apicore
func FindBySystemID(entity, id string, response interface{}) error {
	filter := Filter{
		Entity: entity,
		Map:    "core",
		Name:   "bySystemId",
		Params: []Param{Param{
			Key:   "systemId",
			Value: id,
		},
		},
	}
	return Query(filter, response)
}

//FindByProcessID finds entity by ProcessID on apicore
func FindByProcessID(entity, id string, response interface{}) error {
	filter := Filter{
		Entity: entity,
		Map:    "core",
		Name:   "byProcessId",
		Params: []Param{Param{
			Key:   "processId",
			Value: id,
		},
		},
	}
	return Query(filter, response)
}
