package openskindb_models

import "github.com/baldurstod/vdf"

type ItemsGame struct {
	*vdf.KeyValue
}

type MusicKit struct {
	DefinitionIndex int
	Name								string 
	Prefab 							string
	ItemName						string
	ImageInventory			string
	DisplayModel				string
}

type WeaponCaseItemSet struct {
	Tag								string 
	TagText 					string
	TagGroup					string
	TagGroupText			string
}

type WeaponCaseKey struct {
	DefinitionIndex 		int
	Name            		string
	ItemName        		string
	ItemDescription 		string
	FirstSaleDate 			string
	Prefab 							string
	ImageInventory			string
}

type WeaponCase struct {
	DefinitionIndex 		int
	Name            		string
	ItemName        		string
	ItemDescription 		string
	Prefab 							string
	ImageInventory			string
	ModelPlayer					string
	FirstSaleDate 			string
	ItemSet 					 	*WeaponCaseItemSet
	Key 								*WeaponCaseKey
}

type Collectible struct {
	DefinitionIndex 		int
	Name            		string
	Prefab 							string
	ItemName        		string
	ItemDescription 		string
	ImageInventory  		string
	DisplayModel    		string
	TournamentEventId 	int
	Type 					 			CollectibleType
}