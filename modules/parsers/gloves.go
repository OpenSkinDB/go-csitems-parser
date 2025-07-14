package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseGloves(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.GloveItem {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil
	}

	var gloves []models.GloveItem

	for _, w := range items.GetChilds() {
		prefab, _ := w.GetString("prefab")

		if prefab != "hands_paintable" {
			// Skip non-glove items
			continue
		}

		definition_index, _ := strconv.Atoi(w.Key)

		item_name, _ := w.GetString("item_name")
		name, _ := w.GetString("name")
		item_description, _ := w.GetString("item_description")

		current := models.GloveItem{
			DefinitionIndex: definition_index,
			ItemName:        item_name,
			Name:            name,
			ItemDescription: item_description,
			Prefab:          prefab,
		}

		gloves = append(gloves, current)
	}

	// Save gloves to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' gloves in %s", len(gloves), duration)

	return gloves
}
