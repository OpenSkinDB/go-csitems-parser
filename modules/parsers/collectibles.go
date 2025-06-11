package openskindb_parsers

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	models "github.com/openskindb/openskindb-csitems/models"
	modules "github.com/openskindb/openskindb-csitems/modules"

	"github.com/rs/zerolog"
)

func ParseCollectibles(ctx context.Context, ig *models.ItemsGame) []models.Collectible {
	logger := zerolog.Ctx(ctx);

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
		name, _ := item.GetString("name")
		item_description, _ := item.GetString("item_description")
		image_inventory, _ := item.GetString("image_inventory")

		attributes := modules.GetSubKey(item, "attributes")

		// Get subkeys for attributes
		tournament_event_id, _ := attributes.GetInt("value")
		pedestal_display_model, _ := attributes.GetString("pedestal display model")

		// Get the pedestal display model from attributes
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
	logger.Info().Msgf("Parsed '%d' collectibles in %s", len(collectibles), duration.String())

	return collectibles
}

func GetCollectibleType(
	image_inventory string,
	prefab string,
	item_name string,
	tournament_event_id int,
) models.CollectibleType {
	if prefab == "" {
		return models.CollectibleTypeUnknown
	}

	if image_inventory == "" {
		return models.CollectibleTypeUnknown
	}

	if prefab == "premier_season_coin" {
		return models.CollectibleTypePremierSeasonCoin
	}

	if strings.Contains(image_inventory, "service_medal") {
		return models.CollectibleTypeServiceMedal
	}

	// image_inventory looks like "10yearcoin", "5yearcoin", etc.
	reg1 := regexp.MustCompile(`\d+yearcoin`)
	if reg1.MatchString(image_inventory) {
		return models.CollectibleTypeYearsOfService	
	}

	if strings.Contains(item_name, "#CSGO_Collectible_Map") {
		return models.CollectibleTypeMapContributor
	}

	if strings.HasPrefix(item_name, "#CSGO_TournamentJournal") || strings.HasPrefix(item_name, "#CSGO_CollectibleCoin") {
		return models.CollectibleTypePickEm
	}

	if strings.HasPrefix(item_name, "#CSGO_Collectible_Pin") {
		return models.CollectibleTypeMapPin
	}

	// This is a bit odd, idk what Valve was thinking 
	if strings.HasPrefix(item_name, "#CSGO_Collectible_CommunitySeason") {
		return models.CollectibleTypeMapPin
	}

	// Create a regex for season1_coin, // season2_coin, etc.
	reg2 := regexp.MustCompile(`season\d+_coin`)
	if reg2.MatchString(prefab) {
		return models.CollectibleTypeOperation
	}

	if prefab == "majors_trophy" {
		return models.CollectibleTypeTournamentFinalist
	}

	// katowice_2014_finalist

	return models.CollectibleTypeUnknown
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