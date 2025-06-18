package item

type Configuration struct {
	Id          int32              `json:"ItemId"`
	Name        string             `json:"ConfigurationName"`
	Type        string             `json:"ConfigurationType"`
	TypeId      string             `json:"ConfigurationTypeId"`
	Category    string             `json:"ConfigurationCategory"`
	Version     string             `json:"ConfigurationVersion"`
	Header      map[string]any     `json:"header"`
	Application map[string][]int64 `json:"application"`
	Reference   []int64            `json:"reference"`
}

type ItemService interface {
	Save(c Configuration) error
	LoadWithName(cname string, limit int) ([]Configuration, error)
	LoadWithId(cid int32) (Configuration, error)

	DeleteWithName(cname string) error
	DeleteWithId(cid int32) error
}
