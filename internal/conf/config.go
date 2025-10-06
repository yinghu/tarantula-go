package conf

type Config struct {
	Name     string `json:"-"`
	Used     bool   `json:"Used"`
	Sequence int    `json:"Sequence"`
	Sql
}
