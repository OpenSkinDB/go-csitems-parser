package openskindb_models

import "github.com/baldurstod/vdf"

type ItemsGame struct {
	*vdf.KeyValue
}

type MusicKit struct {
	DefinitionIndex 		int 		`json:"definition_index"`
	Name								string  `json:"name"`
	Prefab 							string 	`json:"prefab"`
	ItemName						string 	`json:"item_name"`
	ImageInventory			string 	`json:"image_inventory"`
	DisplayModel				string  `json:"display_model"`
}

type WeaponCaseItemSet struct {
	Tag								string  	`json:"tag"`
	TagText 					string 		`json:"tag_text"`
	TagGroup					string 		`json:"tag_group"`
	TagGroupText			string 		`json:"tag_group_text"`
}

type WeaponCaseKey struct {
	DefinitionIndex 		int 		`json:"definition_index"`
	Name            		string 	`json:"name"`
	ItemName        		string 	`json:"item_name"`
	ItemDescription 		string 	`json:"item_description"`
	FirstSaleDate 			string 	`json:"first_sale_date"`
	Prefab 							string 	`json:"prefab"`
	ImageInventory			string 	`json:"image_inventory"`
}

type WeaponCase struct {
	DefinitionIndex 		int 									`json:"definition_index"`
	Name            		string 								`json:"name"`
	ItemName        		string 								`json:"item_name"`
	ItemDescription 		string	 							`json:"item_description"`
	Prefab 							string	 							`json:"prefab"`
	ImageInventory			string	 							`json:"image_inventory"`
	ModelPlayer					string	 							`json:"model_player"`
	FirstSaleDate 			string	 							`json:"first_sale_date"`
	ItemSet 					 	*WeaponCaseItemSet 		`json:"item_set"`
	Key 								*WeaponCaseKey 				`json:"key"`
}

type Collectible struct {
	DefinitionIndex 		int 							`json:"definition_index"`
	Name            		string 						`json:"name"`
	Prefab 							string 						`json:"prefab"`
	ItemName        		string 						`json:"item_name"`
	ItemDescription 		string						`json:"item_description"`
	ImageInventory  		string						`json:"image_inventory"`
	DisplayModel    		string						`json:"display_model"`
	TournamentEventId 	int 							`json:"tournament_event_id"`
	Type 					 			CollectibleType 	`json:"type"`
}