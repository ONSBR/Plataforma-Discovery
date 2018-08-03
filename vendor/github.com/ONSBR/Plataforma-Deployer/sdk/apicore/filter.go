package apicore

//Filter is the object to query on apicore
type Filter struct {
	Map      string
	Entity   string
	Name     string
	Page     int
	PageSize int
	Params   []Param
}

//Param is a query string params
type Param struct {
	Key   string
	Value string
}
