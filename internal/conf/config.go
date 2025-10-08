package conf

type Config struct {
	Used         bool   `json:"Used"`
	Sequence     int    `json:"Sequence"`
	HttpEndpoint string `json:"HttpEndpoint"`
	Sql
	EventEndpoint
}
