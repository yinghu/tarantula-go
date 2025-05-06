module gameclustering.com/main

go 1.24.2

replace gameclustering.com/internal/auth => ../../../internal/auth

replace gameclustering.com/internal/persistence => ../../../internal/persistence

require gameclustering.com/internal/auth v0.0.0-00010101000000-000000000000

require (
	gameclustering.com/internal/persistence v0.0.0-00010101000000-000000000000 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)
