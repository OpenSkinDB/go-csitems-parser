module github.com/zwolof/go-csitems-parser

go 1.24.4

require github.com/baldurstod/vdf v0.0.8 // indirect

// import from local folder
replace github.com/zwolof/go-csitems-parser/modules/parsers => ./modules/parsers

replace github.com/zwolof/go-csitems-parser/modules => ./modules

replace github.com/zwolof/go-csitems-parser/models => ./models

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jedib0t/go-pretty/v6 v6.6.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/zwolof/go-csitems-parser/models v0.0.0-00010101000000-000000000000 // indirect
	github.com/zwolof/go-csitems-parser/modules v0.0.0-00010101000000-000000000000 // indirect
	github.com/zwolof/go-csitems-parser/modules/parsers v0.0.0-00010101000000-000000000000 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/rs/zerolog v1.34.0
)
