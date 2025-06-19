package item

type Category struct {
	Scope        string     `json:"Scope"`
	Name         string     `json:"Name"`
	Rechargeable bool       `json:"Rechargeable"`
	Version      string     `json:"Version"`
	Description  string     `json:"Description"`
	Properties   []Property `json:"Properties"`
}

type Property struct {
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	Reference    string `json:"Reference"`
	Nullable     bool   `json:"Nullable"`
	Downloadable bool   `json:"Downloadable"`
}

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
	SaveCategory(c Category) error
	LoadCategory(cname string) (Category, error)

	Save(c Configuration) error
	LoadWithName(cname string, limit int) ([]Configuration, error)
	LoadWithId(cid int32) (Configuration, error)

	DeleteWithName(cname string) error
	DeleteWithId(cid int32) error
}
