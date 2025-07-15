package parsers

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"
)

func ParseCollectibles(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.Collectible {
	logger := modules.GetLogger()

	start := time.Now()
	// logger.Info().Msg("Parsing collectibles...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil
	}

	var collectibles []models.Collectible

	// Iterate through all items in the "items" section
	for _, item := range items.GetChilds() {
		item_name, _ := item.GetString("item_name")

		if item_name == "" || !IsItemCollectible(item_name) {
			continue
		}

		definition_index, _ := strconv.Atoi(item.Key)
		prefab, _ := item.GetString("prefab")
		image_inventory, _ := item.GetString("image_inventory")

		tournament_event_id, _ := modules.GetTournamentEventId(item)
		collectible_type := GetCollectibleType(image_inventory, prefab, item_name, tournament_event_id)

		// item_description, _ := item.GetString("item_description")
		// attributes := modules.GetSubKey(item, "attributes")

		// Get subkeys for attributes
		// pedestal_display_model, _ := attributes.GetString("pedestal display model")

		// Get the pedestal display model from attributes
		// Determine the type of collectible

		collectibles = append(collectibles, models.Collectible{
			DefinitionIndex:   definition_index,
			MarketHashName:    modules.GenerateMarketHashName(t, item_name, nil, "collectible"),
			ImageInventory:    image_inventory,
			Type:              collectible_type,
			TournamentEventId: tournament_event_id,

			// Prefab:            prefab,
			// Name:              item_name,
			// Description:       item_description,
			// Model:             pedestal_display_model,
		})
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' collectibles in %s", len(collectibles), duration)

	return collectibles
}

func GetCollectibleType(
	image_inventory string,
	prefab string,
	item_name string,
	tournament_event_id int,
) string {
	if prefab == "" {
		return "unknown"
	}

	if image_inventory == "" {
		return "unknown"
	}

	if prefab == "premier_season_coin" {
		return "premier_season_coin"
	}

	if strings.Contains(image_inventory, "service_medal") {
		return "service_medal"
	}

	// image_inventory looks like "10yearcoin", "5yearcoin", etc.
	reg1 := regexp.MustCompile(`\d+yearcoin`)
	if reg1.MatchString(image_inventory) {
		return "years_of_service"
	}

	if strings.Contains(item_name, "#CSGO_Collectible_Map") {
		return "map_contributor"
	}

	if strings.HasPrefix(item_name, "#CSGO_TournamentJournal") || strings.HasPrefix(item_name, "#CSGO_CollectibleCoin") {
		return "pickem"
	}

	if strings.HasPrefix(item_name, "#CSGO_Collectible_Pin") {
		return "collectible_pin"
	}

	// This is a bit odd, idk what Valve was thinking
	if strings.HasPrefix(item_name, "#CSGO_Collectible_CommunitySeason") {
		return "collectible_pin"
	}

	// Create a regex for season1_coin, // season2_coin, etc.
	reg2 := regexp.MustCompile(`season\d+_coin`)
	if reg2.MatchString(prefab) {
		return "operation_coin"
	}

	if prefab == "majors_trophy" {
		return "tournament_trophy"
	}

	// katowice_2014_finalist
	return "unknown"
}

func IsItemCollectible(item_name string) bool {
	if len(item_name) == 0 {
		return false
	}

	//if it starts with "Collectible" or "Collectible_"
	if strings.HasPrefix(item_name, "#CSGO_Collectible") || strings.HasPrefix(item_name, "#CSGO_TournamentJournal") {
		return true
	}

	return false
}
