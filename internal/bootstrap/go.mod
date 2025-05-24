module gameclustering.com/internal/bootstrap

go 1.24.2

replace gameclustering.com/internal/persistence => ../persistence

replace gameclustering.com/internal/util => ../util

replace gameclustering.com/internal/cluster => ../cluster

replace gameclustering.com/internal/conf => ../conf

replace gameclustering.com/internal/event => ../event

replace gameclustering.com/internal/core => ../core

replace gameclustering.com/internal/metrics => ../metrics

require (
	gameclustering.com/internal/cluster v0.0.0-00010101000000-000000000000
	gameclustering.com/internal/conf v0.0.0-00010101000000-000000000000
	gameclustering.com/internal/event v0.0.0-00010101000000-000000000000
)

require (
	gameclustering.com/internal/core v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	go.etcd.io/etcd/api/v3 v3.5.21 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.21 // indirect
	go.etcd.io/etcd/client/v3 v3.5.21 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
