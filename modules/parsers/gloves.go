package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseGloves(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.BaseWeapon {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil
	}

	var gloves []models.BaseWeapon

	for _, w := range items.GetChilds() {
		prefab, _ := w.GetString("prefab")

		if prefab != "hands_paintable" {
			// Skip non-glove items
			continue
		}

		definition_index, _ := strconv.Atoi(w.Key)

		classname, _ := w.GetString("name")
		item_name, _ := w.GetString("item_name")
		// item_description, _ := w.GetString("item_description")

		current := models.BaseWeapon{
			DefinitionIndex: definition_index,
			ClassName:       classname,
			Name:            modules.GenerateMarketHashName(t, item_name, nil, "glove"),
			// ItemName:        item_name,
			// ItemDescription: item_description,
			// Prefab:          prefab,
		}

		gloves = append(gloves, current)
	}

	// Save gloves to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' gloves in %s", len(gloves), duration)

	return gloves
}
