package parsers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

var valid_misc_capsule_prefabs = []string{
	"patch_capsule",
	"graffiti_box",
	"weapon_case_selfopening_collection",
}

var prefab_to_tagtype = map[string]string{
	"patch_capsule":                      "PatchCapsule",
	"stockh2021_patch_capsule_prefab":    "PatchCapsule",
	"graffiti_box":                       "SprayCapsule",
	"weapon_case_selfopening_collection": "ItemSet",
	"weapon_case_base":                   "ItemSet",
}

func IsValidMiscSelfOpeningCapsule(prefab string, name string) bool {
	// Dumb check for the P250 X-Ray package..
	if strings.HasPrefix(name, "crate_xray_") {
		return true
	}

	// Dumb check for music kits
	if strings.HasPrefix(name, "crate_musickit_") {
		return true
	}

	for _, valid_prefab := range valid_misc_capsule_prefabs {
		if strings.Contains(prefab, valid_prefab) {
			return true
		}
	}
	return false
}

func ParseSelfOpeningCrates(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.StickerCapsule {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing weapon cases...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil
	}

	var capsules []models.StickerCapsule

	// Iterate through all items in the "items" section
	for _, item := range items.GetChilds() {
		// prefab, _ := item.GetString("prefab")
		prefab, _ := item.GetString("prefab")
		name, _ := item.GetString("name")

		if !IsValidMiscSelfOpeningCapsule(prefab, name) {
			continue // Skip non-self-opening capsules
		}

		item_name, _ := item.GetString("item_name")
		definition_index, _ := strconv.Atoi(item.Key)
		image_inventory, _ := item.GetString("image_inventory")
		item_description, _ := item.GetString("item_description")

		// Get the item set ID from the item tags
		tag_type := prefab_to_tagtype[prefab]

		if tag_type == "" {
			logger.Warn().Msgf("Unknown tag type for item: %s", item_name)
			continue // Skip items with unknown tag types
		}

		item_set := modules.GetContainerItemSet(item, t, tag_type)

		if item_set == nil {
			item_set = modules.GetSupplyCrateSeries(item, ig)

			if item_set == nil {
				fmt.Println("Item set is nil again, skipping item:", item_name)
				continue
			}
		}

		// Create the sticker capsule model
		var current = models.StickerCapsule{
			DefinitionIndex: definition_index,
			Name:            name,
			ImageInventory:  image_inventory,
			ItemDescription: item_description,
			ItemSetId:       item_set,
			MarketHashName:  modules.GenerateMarketHashName(t, item_name, nil, "self_opening_capsule"),
		}

		capsules = append(capsules, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' self-opening capsules in %s", len(capsules), duration)

	return capsules
}
