package openskindb_parsers

import (
	"time"

	models "github.com/openskindb/openskindb-csitems/models"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ParseMusicKits(ig *models.ItemsGame) []models.MusicKit {
	// "music_definitions"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano;

	startedAt := time.Now()
	log.Info().Msg("Parsing music kits...")

	music_definitions, err := ig.Get("music_definitions")

	if err != nil {
		log.Error().Err(err).Msg("Failed to get music_definitions")
		return nil
	}

	var musicKits []models.MusicKit

	for _, musicKit := range music_definitions.GetChilds() {
		definitionIndex, _ := musicKit.GetInt("definition_index")
		name, _ := musicKit.GetString("name")
		loc_name, _ := musicKit.GetString("loc_name")
		image_inventory, _ := musicKit.GetString("image_inventory")
		pedestal_display_model, _ := musicKit.GetString("pedestal_display_model")

		if name == "" {
			log.Warn().Msg("Music Kit name is empty")
			continue
		}

		musicKits = append(musicKits, models.MusicKit{
			DefinitionIndex: definitionIndex,
			Name:            name,
			ItemName:        loc_name,
			ImageInventory:  image_inventory,
			DisplayModel:    pedestal_display_model,
		})
	}

	// Save music kits to the database
	endedAt := time.Now()
	log.Info().Msgf("Parsed %d music kits in %s", len(musicKits), endedAt.Sub(startedAt).String())

	return musicKits
}