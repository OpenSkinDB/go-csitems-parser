package models

type StickerCapsule struct {
	DefinitionIndex int     `json:"definition_index"`
	Name            string  `json:"name"`
	MarketHashName  string  `json:"market_hash_name"`
	ItemDescription string  `json:"item_description"`
	ImageInventory  string  `json:"image_inventory"`
	ItemSetId       *string `json:"item_set_id"`
}

func MapStickerCapsulesToSchema(stickerCapsules []StickerCapsule) map[int]string {
	stickerCapsuleMap := make(map[int]string)
	for _, capsule := range stickerCapsules {
		stickerCapsuleMap[capsule.DefinitionIndex] = capsule.GetMarketHashName()
	}
	return stickerCapsuleMap
}

// Getters for StickerCapsule
func (s StickerCapsule) GetDefinitionIndex() int {
	return s.DefinitionIndex
}

func (s StickerCapsule) GetName() string {
	return s.Name
}

func (s StickerCapsule) GetMarketHashName() string {
	return s.MarketHashName
}

func (s StickerCapsule) GetItemDescription() string {
	return s.ItemDescription
}

func (s StickerCapsule) GetImageInventory() string {
	return s.ImageInventory
}

func (s StickerCapsule) GetItemSetId() *string {
	return s.ItemSetId
}
