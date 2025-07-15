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

		item_set_paint_kits := GetItemSetPaintKitsForWeapon(item_sets, weapon.ClassName)
		for _, paint_kit_name := range item_set_paint_kits {
			paint_kit := GetPaintKitByName(paint_kits, paint_kit_name)

			if paint_kit == nil {
				continue
			}

			data := GetPaintKitWeaponCombinationData(item_sets, weapon.ClassName, paint_kit.Name)

			if data != nil {
				paint_kit.StatTrak = data.CanBeStatTrak
				paint_kit.Souvenir = data.CanBeSouvenir
			}

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

func GetKnifePaintKits(
	knives *[]models.BaseWeapon,
	paint_kits *[]models.PaintKit,
	knife_map map[string][]string,
) []WeaponSkinMap {
	knife_skin_map := make([]WeaponSkinMap, 0)

	for _, knife := range *knives {
		// Create a new glove skin map entry
		current := WeaponSkinMap{
			BaseItem:  knife,
			PaintKits: make([]models.PaintKit, 0),
		}

		knife_map_value, ok := knife_map[knife.ClassName]
		if !ok {
			continue
		}

		for _, pk_name := range knife_map_value {
			for _, paint_kit := range *paint_kits {
				if paint_kit.Name != pk_name {
					continue
				}

				paint_kit.Souvenir = false // Knives can NOT be Souvenir
				paint_kit.StatTrak = true  // Knives can always be StatTrak

				// Add the paint kit to the current glove skin map
				current.PaintKits = append(current.PaintKits, paint_kit)
			}
		}

		knife_skin_map = append(knife_skin_map, current)
	}

	return knife_skin_map
}

type GloveSkinMap struct {
	BaseItem  models.BaseWeapon `json:"base_item"`
	PaintKits []models.PaintKit `json:"paint_kits"`
}

// Because the schema is weird..
var glovePrefixMap = map[string]string{
	"studded_brokenfang_gloves": "operation10_",
	"studded_hydra_gloves":      "bloodhound_hydra_",
	"leather_handwraps":         "handwrap_",
	"studded_bloodhound_gloves": "bloodhound_",
}

func GetGlovePaintKits(gloves *[]models.BaseWeapon, paint_kits *[]models.PaintKit) []WeaponSkinMap {
	glove_skin_map := make([]WeaponSkinMap, 0)

	for _, glove := range *gloves {
		// Create a new glove skin map entry
		current := WeaponSkinMap{
			BaseItem:  glove,
			PaintKits: make([]models.PaintKit, 0),
		}

		for _, paint_kit := range *paint_kits {
			// We need to remove "_gloves" from the item name
			var newGlovePrefix string

			value, ok := glovePrefixMap[glove.ClassName]
			if !ok {
				newGlovePrefix = strings.Replace(glove.ClassName, "_gloves", "", -1)
			} else {
				newGlovePrefix = value
			}

			// if the paintkit name starts with the glove name, we can assume it's a paint kit for that glove
			if !strings.HasPrefix(paint_kit.Name, newGlovePrefix) {
				continue
			}

			paint_kit.StatTrak = false
			paint_kit.Souvenir = false

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
		if wpn.ClassName == weapon_class {
			return &wpn
		}
	}
	return nil
}

func GetPaintKitWeaponCombinationData(item_sets *[]models.ItemSet, cn string, pk string) *models.PaintKitWeaponCombinationData {
	for _, item_set := range *item_sets {
		for _, item := range item_set.Items {
			if item.WeaponClass == cn && item.PaintKitName == pk {
				return &models.PaintKitWeaponCombinationData{
					ItemSetId:     item_set.Key,
					CanBeStatTrak: item_set.HasCrate,
					CanBeSouvenir: item_set.HasSouvenir,
				}
			}
		}
	}
	return nil
}
