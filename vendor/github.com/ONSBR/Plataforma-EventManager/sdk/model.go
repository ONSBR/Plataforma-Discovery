package sdk

type Metadata struct {
	Type        string `json:"type"`
	InstanceID  string `json:"instance_id"`
	Origin      string `json:"origin"`
	ChangeTrack string `json:"changeTrack"`
}

type Model struct {
	Metadata Metadata `json:"_metadata"`
	ID       string   `json:"id"`
}
