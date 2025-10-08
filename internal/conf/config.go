package conf

type Config struct {
	Used         bool   `json:"Used"`
	Sequence     int    `json:"Sequence"`
	HttpEndpoint string `json:"HttpEndpoint"`
	TcpEndpoint  string `json:"TcpEndpoint"`
	SqlEndpoint string `json:"SqlEndpoint"`
}
