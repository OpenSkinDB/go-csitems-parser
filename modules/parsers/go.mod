module github.com/openskindb/openskindb-csitems/modules/parsers

go 1.24.4

replace github.com/openskindb/openskindb-csitems/models => ../../models
replace github.com/openskindb/openskindb-csitems/modules => ../

require (
	github.com/baldurstod/vdf v0.0.8
	github.com/openskindb/openskindb-csitems/models v0.0.0-00010101000000-000000000000
	github.com/openskindb/openskindb-csitems/modules v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
