package item


type Configuration struct {
	Name        string `json:"ConfigurationName"`
	Type        string `json:"ConfigurationType"`
	TypeId      string `json:"ConfigurationTypeId"`
	Category    string `json:"ConfigurationCategory"`
	Version     string `json:"ConfigurationVersion"`
	Header      map[string]any     `json:"header"`
	Application map[string][]int64 `json:"application"`
	Reference   []int64            `json:"reference"`
}
