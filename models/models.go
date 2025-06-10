package openskindb_models

import "github.com/baldurstod/vdf"

type ItemsGame struct {
	*vdf.KeyValue
}

type GenericColor struct {
	Key 		 		string `json:"key"`
	ColorName 	string `json:"color_name"`
	HexColor  	string `json:"hex_color"`
}

type PaintKit struct {
	DefinitionIndex 		int			`json:"definition_index"`
	Name 		 						string 	`json:"name"`
	UseLegacyModel 			bool 		`json:"use_legacy_model"`
	DescriptionString 	string 	`json:"description_string"`
	DescriptionTag 			string 	`json:"description_tag"`
	Style 							int 		`json:"style"`
	WearRemapMin 				float64 `json:"wear_remap_min"`
	WearRemapMax 				float64 `json:"wear_remap_max"`
	Rarity 							string	`json:"rarity"`
}

type ItemSetItem struct {
	PaintKitName 			string 	`json:"paintkit_name"`
	WeaponClass 	string 	`json:"weapon_class"`
}

type ItemSet struct {
	Key 								string 				`json:"key"`
	Name 								string 				`json:"name"`
	SetDescription 			string 				`json:"set_description"`
	IsCollection 				bool 					`json:"is_collection"`
	Items 							[]ItemSetItem `json:"items"`
}

type Rarity struct {
	Key									string  `json:"key"`
	LocKey 							string 	`json:"loc_key"`
	LocKeyWeapon				string 	`json:"loc_key_weapon"`
	LocKeyCharacter			string 	`json:"loc_key_character"`
	HexColor						string  `json:"hex_color"`
	ColorName 					string  `json:"color_name"`
	DropSound						string  `json:"drop_sound"`
}

type MusicKit struct {
	DefinitionIndex 		int 		`json:"definition_index"`
	Name								string  `json:"name"`
	Prefab 							string 	`json:"prefab"`
	ItemName						string 	`json:"item_name"`
	ImageInventory			string 	`json:"image_inventory"`
	DisplayModel				string  `json:"display_model"`
}

type PlayerAgent struct {
	DefinitionIndex 		int 		`json:"definition_index"`
	Name								string  `json:"name"`
	Prefab 							string 	`json:"prefab"`
	ModelPlayer					string	`json:"model_player"`	
	ItemName						string 	`json:"item_name"`
	ItemDescription 		string 	`json:"item_description"`
	ImageInventory			string 	`json:"image_inventory"`
	ItemRarity					string 	`json:"item_rarity"`
	UsedByTeam					string 	`json:"used_by_team"`
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