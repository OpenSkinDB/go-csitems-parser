package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseKeychains(ctx context.Context, ig *models.ItemsGame) []models.Keychain {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing keychains...")

	keychain_definitions, err := ig.Get("keychain_definitions")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get keychain_definitions, is items_game.txt valid?")
		return nil
	}

	var agents []models.Keychain
	for _, mk := range keychain_definitions.GetChilds() {
		definition_index, _ := strconv.Atoi(mk.Key)
		name, _ := mk.GetString("name")
		loc_name, _ := mk.GetString("loc_name")
		loc_description, _ := mk.GetString("loc_description")
		image_inventory, _ := mk.GetString("image_inventory")
		item_rarity, _ := mk.GetString("item_rarity")
		// item_quality, _ := mk.GetString("item_quality")
		pedestal_display_model, _ := mk.GetString("pedestal_display_model")

		// this is stupid..
		tags, _ := mk.Get("tags")

		loot_list_id := ""
		if tags != nil {
			keychain_capsule, _ := tags.Get("KeychainCapsule")

			if keychain_capsule != nil {
				loot_list_id, _ = keychain_capsule.GetString("tag_value")
			}
		}

		keychain_capsule := modules.GetSubKey(mk, "tags.KeychainCapsule")

		if keychain_capsule != nil {
			loot_list_id, _ = keychain_capsule.GetString("tag_value")
			logger.Debug().Msgf("Found KeychainCapsule tag with loot_list_id: %s", loot_list_id)
		}

		// Create a new Keychain instance
		current := models.Keychain{
			DefinitionIndex: definition_index,
			Name:            name,
			LocName:         loc_name,
			LocDescription:  loc_description,
			Rarity:          item_rarity,
			ImageInventory:  image_inventory,
			Model:           pedestal_display_model,
			LootListId:      loot_list_id,
		}

		// Append the current keychains to the slice
		agents = append(agents, current)
	}

	// Save keychains to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' keychains in %s", len(agents), duration)

	return agents
}
