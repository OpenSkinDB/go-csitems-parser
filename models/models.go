package openskindb_models

import "github.com/baldurstod/vdf"

type ItemsGame struct {
	*vdf.KeyValue
}

type MusicKit struct {
	DefinitionIndex int
	Name							string 
	ItemName					string
	ImageInventory		string
	DisplayModel			string
}

type Collectible struct {
	DefinitionIndex 	int
	Name            	string
	ItemName        	string
	ItemDescription 	string
	ImageInventory  	string
	DisplayModel    	string
	Prefab 						string
	TournamentEventId int
	Type 					 		CollectibleType
}