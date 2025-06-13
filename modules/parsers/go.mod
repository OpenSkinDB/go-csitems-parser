module github.com/zwolof/go-csitems-parser/modules/parsers

go 1.24.4

replace github.com/zwolof/go-csitems-parser/models => ../../models
replace github.com/zwolof/go-csitems-parser/modules => ../

require (
	github.com/baldurstod/vdf v0.0.8
	github.com/zwolof/go-csitems-parser/models v0.0.0-00010101000000-000000000000
	github.com/zwolof/go-csitems-parser/modules v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
