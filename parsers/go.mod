module github.com/openskindb/openskindb-csitems/parsers

go 1.24.4

require github.com/rs/zerolog v1.34.0

// import from local folder
replace github.com/openskindb/openskindb-csitems/models => ../models

require (
	github.com/baldurstod/vdf v0.0.8 // indirect
	github.com/golang-collections/collections v0.0.0-20130729185459-604e922904d3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/openskindb/openskindb-csitems/models v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.12.0 // indirect
)
