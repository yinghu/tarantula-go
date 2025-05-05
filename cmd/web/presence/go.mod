module gameclustering.com/main

go 1.24.2

replace gameclustering.com/internal/auth => ../../../internal/auth

replace gameclustering.com/internal/util => ../../../internal/util

require gameclustering.com/internal/auth v0.0.0-00010101000000-000000000000

require (
	gameclustering.com/internal/util v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/crypto v0.37.0 // indirect
)
