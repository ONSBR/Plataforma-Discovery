package domain

type Reprocessing struct {
	ID     string   `json:"id"`
	Done   bool     `json:"done"`
	Events []*Event `json:"events"`
}
