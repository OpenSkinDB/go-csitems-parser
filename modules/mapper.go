package modules

import (
	"go-csitems-parser/models"
	"strings"
)

type WeaponToPaintKitMap struct {
	DefinitionIndex int    `json:"definition_index"`
	Name            string `json:"name"`
	StickerCount    int    `json:"sticker_count"`
	Type            string `json:"type"`
}

type WeaponSkinMap struct {
	BaseItem  models.BaseWeapon `json:"base_item"`
	PaintKits []models.PaintKit `json:"paint_kits"`
}

func GetWeaponPaintKits(
	weapons *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	item_sets *[]models.ItemSet,
) []WeaponSkinMap {
	weapon_skin_map := make([]WeaponSkinMap, 0)

	for _, weapon := range *weapons {
		// Create a new glove skin map entry
		current := WeaponSkinMap{
			BaseItem:  weapon,
			PaintKits: make([]models.PaintKit, 0),
		}

		item_set_paint_kits := GetItemSetPaintKitsForWeapon(item_sets, weapon.ItemClass)
		for _, paint_kit_name := range item_set_paint_kits {
			paint_kit := GetPaintKitByName(paint_kits, paint_kit_name)

			current.PaintKits = append(current.PaintKits, *paint_kit)
		}

		weapon_skin_map = append(weapon_skin_map, current)
	}

	return weapon_skin_map
}

func GetItemSetPaintKitsForWeapon(
	item_sets *[]models.ItemSet,
	weapon_name string,
) []string {
	paint_kits := make([]string, 0)

	for _, item_set := range *item_sets {
		if item_set.Type != models.ItemSetTypePaintKits {
			continue
		}

		for _, item := range item_set.Items {
			if item.WeaponClass == weapon_name {
				paint_kits = append(paint_kits, item.PaintKitName)
			}
		}
	}

	return paint_kits
}

type KnifeSkinMap struct {
	BaseItem  models.KnifeItem  `json:"base_item"`
	PaintKits []models.PaintKit `json:"paint_kits"`
}

func GetKnifePaintKits(
	knives *[]models.KnifeItem,
	paint_kits *[]models.PaintKit,
	knife_map map[string][]string,
) []KnifeSkinMap {
	knife_skin_map := make([]KnifeSkinMap, 0)

	for _, knife := range *knives {
		// Create a new glove skin map entry
		current := KnifeSkinMap{
			BaseItem:  knife,
			PaintKits: make([]models.PaintKit, 0),
		}

		knife_map_value, ok := knife_map[knife.Name]
		if !ok {
			continue
		}

		for _, pk_name := range knife_map_value {
			for _, paint_kit := range *paint_kits {
				if paint_kit.Name != pk_name {
					continue
				}

				// Add the paint kit to the current glove skin map
				current.PaintKits = append(current.PaintKits, paint_kit)
			}
		}

		knife_skin_map = append(knife_skin_map, current)
	}

	return knife_skin_map
}

type GloveSkinMap struct {
	BaseItem  models.GloveItem  `json:"base_item"`
	PaintKits []models.PaintKit `json:"paint_kits"`
}

// Because the schema is weird..
var glovePrefixMap = map[string]string{
	"studded_brokenfang_gloves": "operation10_",
	"studded_hydra_gloves":      "bloodhound_hydra_",
	"leather_handwraps":         "handwrap_",
	"studded_bloodhound_gloves": "bloodhound_",
}

func GetGlovePaintKits(gloves *[]models.GloveItem, paint_kits *[]models.PaintKit) []GloveSkinMap {
	glove_skin_map := make([]GloveSkinMap, 0)

	for _, glove := range *gloves {
		// Create a new glove skin map entry
		current := GloveSkinMap{
			BaseItem:  glove,
			PaintKits: make([]models.PaintKit, 0),
		}

		for _, paint_kit := range *paint_kits {
			// We need to remove "_gloves" from the item name
			var newGlovePrefix string

			value, ok := glovePrefixMap[glove.Name]
			if !ok {
				newGlovePrefix = strings.Replace(glove.Name, "_gloves", "", -1)
			} else {
				newGlovePrefix = value
			}

			// if the paintkit name starts with the glove name, we can assume it's a paint kit for that glove
			if !strings.HasPrefix(paint_kit.Name, newGlovePrefix) {
				continue
			}

			// Add the paint kit to the current glove skin map
			current.PaintKits = append(current.PaintKits, paint_kit)
		}

		glove_skin_map = append(glove_skin_map, current)
	}

	return glove_skin_map
}

func GetPaintKitByName(paint_kits *[]models.PaintKit, name string) *models.PaintKit {
	for _, paint_kit := range *paint_kits {
		if paint_kit.Name == name {
			return &paint_kit
		}
	}
	return nil
}

func GetWeaponByClass(weapons *[]models.BaseWeapon, weapon_class string) *models.BaseWeapon {
	for _, wpn := range *weapons {
		if wpn.ItemClass == weapon_class {
			return &wpn
		}
	}
	return nil
}
