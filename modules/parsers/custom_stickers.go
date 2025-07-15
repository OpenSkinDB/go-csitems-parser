package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseCustomStickers(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.CustomStickers {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing keychains...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items, is items_game.txt valid?")
		return nil
	}

	var custom_stickers []models.CustomStickers
	for _, mk := range items.GetChilds() {
		prefab, _ := mk.GetString("prefab")

		if prefab != "self_opening_purchase" {
			continue
		}

		definition_index, _ := strconv.Atoi(mk.Key)
		// name, _ := mk.GetString("name")
		loc_name, _ := mk.GetString("loc_name")

		// Create a new Keychain instance
		current := models.CustomStickers{
			Group: definition_index,
			Name:  loc_name,
			Count: 0, // Default count, can be updated later if needed
		}

		// Append the current keychains to the slice
		custom_stickers = append(custom_stickers, current)
	}

	// Save keychains to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' custom stickers in %s", len(custom_stickers), duration)

	return custom_stickers
}
