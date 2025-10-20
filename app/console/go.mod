module gameclustering.com/tarantula

go 1.24.2

replace gameclustering.com/cmd => ./cmd

replace gameclustering.com/cmd/player => ./cmd/player

replace gameclustering.com/cmd/admin => ./cmd/admin

replace gameclustering.com/cmd/admin/node => ./cmd/admin/node

replace gameclustering.com/internal/bootstrap => ../../internal/bootstrap

replace gameclustering.com/internal/persistence => ../../internal/persistence

replace gameclustering.com/internal/util => ../../internal/util

replace gameclustering.com/internal/cluster => ../../internal/cluster

replace gameclustering.com/internal/conf => ../../internal/conf

replace gameclustering.com/internal/event => ../../internal/event

replace gameclustering.com/internal/core => ../../internal/core

replace gameclustering.com/internal/item => ../../internal/item

replace gameclustering.com/internal/metrics => ../../internal/metrics

require gameclustering.com/cmd v0.0.0-00010101000000-000000000000

require (
	gameclustering.com/cmd/admin v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/cmd/admin/node v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/cmd/player v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/bootstrap v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/cluster v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/conf v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/core v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/event v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/item v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/metrics v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/persistence v0.0.0-00010101000000-000000000000 // indirect
	gameclustering.com/internal/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/0xc0d/encoding v0.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/dgraph-io/badger/v4 v4.7.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/flatbuffers v25.2.10+incompatible // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/spf13/cobra v1.10.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.etcd.io/etcd/api/v3 v3.6.5 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.6.5 // indirect
	go.etcd.io/etcd/client/v3 v3.6.5 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/grpc v1.71.1 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)
