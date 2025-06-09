package openskindb_models

import "github.com/baldurstod/vdf"

type ItemsGame struct {
	*vdf.KeyValue
}

type MusicKit struct {
	DefinitionIndex int
	Name						string 
	ItemName				string
	ImageInventory	string
	DisplayModel		string
}