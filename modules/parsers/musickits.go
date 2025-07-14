package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseMusicKits(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.MusicKit {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing music kits...")

	music_definitions, err := ig.Get("music_definitions")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get music_definitions")
		return nil
	}

	var musicKits []models.MusicKit

	for _, musicKit := range music_definitions.GetChilds() {
		definition_index, _ := strconv.Atoi(musicKit.Key)
		name, _ := musicKit.GetString("name")
		loc_name, _ := musicKit.GetString("loc_name")
		image_inventory, _ := musicKit.GetString("image_inventory")
		// pedestal_display_model, _ := musicKit.GetString("pedestal_display_model")

		musicKits = append(musicKits, models.MusicKit{
			DefinitionIndex: definition_index,
			Name:            name,
			ImageInventory:  image_inventory,
			MarketHashName:  modules.GenerateMarketHashName(t, loc_name, "music_kit"),
			// ItemName:        loc_name,
			// Model:           pedestal_display_model,
		})
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' music-kits in %s", len(musicKits), duration)

	return musicKits
}
