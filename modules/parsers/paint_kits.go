package openskindb_parsers

import (
	"context"
	"strconv"
	"time"

	models "github.com/openskindb/openskindb-csitems/models"

	"github.com/rs/zerolog"
)

func ParsePaintKits(ctx context.Context, ig *models.ItemsGame) []models.PaintKit {
	logger := zerolog.Ctx(ctx);

	start := time.Now()
	// logger.Info().Msg("Parsing paintkits...")

	paint_kits, err := ig.Get("paint_kits")
	
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get paint_kits from items_game.txt")
		return nil
	}
	
	var raritymap = GetPaintKitRarityStringMap(ig)

	if(raritymap == nil) {
		logger.Error().Msg("Failed to get paint_kits_rarity from items_game.txt")
		return nil
	}
	
	// logger.Debug().Msgf("Found '%d' entires in 'paint_kits_rarity'", len(raritymap))

	var items []models.PaintKit
	for _, r := range paint_kits.GetChilds() {
		name, _ := r.GetString("name")

		// skip if name equal to "default" or "workshop_default"
		if name == "default" || name == "workshop_default" {
			// Log a warning and skip this paint kit
			// logger.Warn().Msg("Skipping paint kit with name 'default'")
			continue
		}

		definition_index, _ := strconv.Atoi(r.Key)
		use_legacy_model, _ := r.GetBool("use_legacy_model")
		description_string, _ := r.GetString("description_string")
		description_tag, _ := r.GetString("description_tag")
		style, _ := r.GetInt("style")
		wear_remap_min, _ := r.GetFloat64("wear_remap_min")
		wear_remap_max, _ := r.GetFloat64("wear_remap_max")
		
		current := models.PaintKit{
			DefinitionIndex: definition_index,
			Name:            name,
			DescriptionString: description_string,
			DescriptionTag: description_tag,
			UseLegacyModel: use_legacy_model,
			Style: style,
			WearRemapMin: wear_remap_min,
			WearRemapMax: wear_remap_max,
		}

		// Get the skin rarity from the paint_kits_rarity map
		val, exists := raritymap[current.Name]
		if !exists {
			logger.Warn().Msgf("No rarity found for paint kit '%s' (definition index: %d)", current.Name, current.DefinitionIndex)
		}
		current.Rarity = val

		items = append(items, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' paintkits in %s", len(items), duration)

	return items
}

func GetPaintKitRarityStringMap(ig *models.ItemsGame) map[string]string {
	paint_kits_rarity, err := ig.Get("paint_kits_rarity")
	logger := zerolog.Ctx(context.Background())
	
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get paint_kits_rarity from items_game.txt")
		return nil
	}

	// Create a map to hold the rarity strings
	rmap, err := paint_kits_rarity.ToStringMap()

	if err != nil {
		logger.Error().Err(err).Msg("Failed to convert paint_kits_rarity to string map")
		return nil
	}

	return *rmap
}