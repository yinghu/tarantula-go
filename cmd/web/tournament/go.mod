module gameclustering.com/tournament

go 1.24.2

replace gameclustering.com/internal/persistence => ../../../internal/persistence

replace gameclustering.com/internal/util => ../../../internal/util

replace gameclustering.com/internal/cluster => ../../../internal/cluster

replace gameclustering.com/internal/conf => ../../../internal/conf

replace gameclustering.com/internal/event => ../../../internal/event

replace gameclustering.com/internal/core => ../../../internal/core

replace gameclustering.com/internal/metrics => ../../../internal/metrics

replace gameclustering.com/internal/bootstrap => ../../../internal/bootstrap

require (
	gameclustering.com/internal/bootstrap v0.0.0-00010101000000-000000000000
	gameclustering.com/internal/cluster v0.0.0-00010101000000-000000000000
	gameclustering.com/internal/conf v0.0.0-00010101000000-000000000000
)

require (
	gameclustering.com/internal/core v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/event v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/metrics v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/persistence v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/0xc0d/encoding v0.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/dgraph-io/badger/v4 v4.7.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/flatbuffers v25.2.10+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.21 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.21 // indirect
	go.etcd.io/etcd/client/v3 v3.5.21 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
