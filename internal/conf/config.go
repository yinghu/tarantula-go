package conf

type Config struct {
	Used     bool `json:"Used"`
	Sequence int  `json:"Sequence"`
	Sql
}
