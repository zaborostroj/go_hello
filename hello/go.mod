module example.com/hello

go 1.25.3

replace example.com/greetings => ../greetings

replace example.com/repository => ../repository

require example.com/greetings v0.0.0-00010101000000-000000000000

require (
	example.com/repository v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx-gofrs-uuid v0.0.0-20230224015001-1d428863c2e2
)

require (
	github.com/gofrs/uuid/v5 v5.4.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)
