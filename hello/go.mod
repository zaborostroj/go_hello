module example.com/hello

go 1.25.3

replace example.com/greetings => ../greetings

replace example.com/repository => ../repository

require example.com/greetings v0.0.0-00010101000000-000000000000

require example.com/repository v0.0.0-00010101000000-000000000000

require (
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/text v0.30.0 // indirect
)
