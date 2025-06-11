package openskindb_parsers

import (
	"context"
	"strconv"
	"time"

	models "github.com/openskindb/openskindb-csitems/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ParseClientLootLists(ctx context.Context, ig *models.ItemsGame) []models.Collectible {
	logger := zerolog.Ctx(ctx);

	start := time.Now()
	// logger.Info().Msg("Parsing client loot lists...")

	client_loot_lists, err := ig.Get("client_loot_lists")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get client_loot_lists from items_game.txt")
		return nil
	}

	var collectibles []models.Collectible

	// Iterate through all items in the "items" section
	for _, item := range client_loot_lists.GetChilds() {
		item_name, _ := item.GetString("item_name")

		if item_name == "" || !IsItemCollectible(item_name) {
			continue
		}

		definition_index, _ := strconv.Atoi(item.Key)
		prefab, _ := item.GetString("prefab")
		name, _ := item.GetString("name")
		item_description, _ := item.GetString("item_description")
		image_inventory, _ := item.GetString("image_inventory")

		if name == "" {
			log.Warn().Msg("Collectible name is empty")
			continue
		}

		// Get child key called "attributes"
		attributes, _ := item.Get("attributes")

		if attributes == nil {
			continue
		}

		// Get the tournament event id from attributes
		tournament_event_data, _ := attributes.Get("tournament event id");
		tournament_event_id := -1 // Default to -1 if not found

		if(tournament_event_data != nil) {
			tournament_id, err := tournament_event_data.GetInt("value")

			if err == nil {
				// logger.Warn().Msgf("Found tournament event id '%d' for item '%s'", tournament_id, item_name)
				tournament_event_id = tournament_id
			}
		}

		// Get the pedestal display model from attributes
		pedestal_display_model, _ := attributes.GetString("pedestal display model")

		// Determine the type of collectible
		collectible_type := GetCollectibleType(image_inventory, prefab, item_name, tournament_event_id)

		collectibles = append(collectibles, models.Collectible{
			DefinitionIndex: 		definition_index,
			Prefab: 				 		prefab,
			Name:            		name,
			ItemName:        		item_name,
			ItemDescription: 		item_description,
			ImageInventory:  		image_inventory,
			DisplayModel:    		pedestal_display_model,
			Type: 							collectible_type,
			TournamentEventId: 	tournament_event_id,
		})
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' client loot lists in %s", len(collectibles), duration)

	return collectibles
}