package conf

type Env struct {
	NodeName     string
	NodeId       int64
	HttpEndpoint string
	DatabaseURL  string
}

func (f *Env) Load() {
	f.NodeName = "a01"
	f.NodeId = 1
	f.HttpEndpoint = ":8080"
	f.DatabaseURL = "postgres://postgres:password@192.168.1.7:5432/tarantula_user"
}
