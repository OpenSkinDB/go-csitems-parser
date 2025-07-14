package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseAgents(ctx context.Context, ig *models.ItemsGame, tf *modules.TranslatorFactory) []models.PlayerAgent {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing agents...")

	items, err := ig.Get("items")
	t := tf.GetTranslator("english")

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
		item_name, _ := agent.GetString("item_name")
		item_description, _ := agent.GetString("item_description")
		item_rarity, _ := agent.GetString("item_rarity")

		market_hash_name, err := modules.GenerateMarketHashName(t, item_name, "agent")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to generate market hash name for agent")
			continue
		}

		// Create a new MusicKit instance
		current := models.PlayerAgent{
			DefinitionIndex: definition_index,
			MarketHashName:  market_hash_name,
			Name:            item_name,
			Description:     item_description,
			Rarity:          item_rarity,
		}

		// Append the current music kit to the slice
		agents = append(agents, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' agents in %s", len(agents), duration)

	return agents
}
