package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"

	"github.com/rs/zerolog"
)

func ParseAgents(ctx context.Context, ig *models.ItemsGame) []models.PlayerAgent {
	logger := zerolog.Ctx(ctx);

	start := time.Now()
	// logger.Info().Msg("Parsing agents...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items, is items_game.txt valid?")
		return nil
	}

	var agents []models.PlayerAgent
	for _, agent := range items.GetChilds() {
		prefab, _ := agent.GetString("prefab")

		// Skip if prefab is not "customplayertradable", xd
		if prefab != "customplayertradable" {
			continue 
		}

		definition_index, _ := strconv.Atoi(agent.Key)
		name, _ := agent.GetString("name")
		loc_name, _ := agent.GetString("loc_name")
		item_description, _ := agent.GetString("item_description")
		image_inventory, _ := agent.GetString("image_inventory")
		item_rarity, _ := agent.GetString("item_rarity")
		
		// this is stupid..
		used_by_classes, _ := agent.Get("used_by_classes")

		team_str := "all" // Default to "all" if not found
		if used_by_classes != nil {
			value := used_by_classes.GetChilds()[0].Key;
			
			if value != "" {
				team_str = value
			}
		}

		// Create a new MusicKit instance
		current := models.PlayerAgent{
			DefinitionIndex: 	definition_index,
			Name:            	name,
			Prefab: 					prefab,
			ModelPlayer: 			prefab,
			ItemName:        	loc_name,
			ItemDescription: 	item_description,
			ImageInventory: 	image_inventory,
			ItemRarity: 			item_rarity,
			UsedByTeam: 			team_str,
		}

		// Append the current music kit to the slice
		agents = append(agents, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' agents in %s", len(agents), duration)

	return agents
}