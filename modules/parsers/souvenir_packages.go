package parsers

import (
	"context"
	"strconv"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseSouvenirPackages(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.SouvenirPackage {
	logger := zerolog.Ctx(ctx)

	start := time.Now()
	// logger.Info().Msg("Parsing music kits...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get items")
		return nil
	}

	var souvenir_packages []models.SouvenirPackage
	for _, c := range items.GetChilds() {
		prefab, _ := c.GetString("prefab")

		if prefab != "weapon_case_souvenirpkg" {
			continue
		}

		definition_index, _ := strconv.Atoi(c.Key)
		item_name, _ := c.GetString("item_name")
		name, _ := c.GetString("name")
		image_inventory, _ := c.GetString("image_inventory")

		item_set := modules.GetContainerItemSet(c)
		tournament_event_id, _ := modules.GetTournamentEventId(c)

		current := models.SouvenirPackage{
			DefinitionIndex:   definition_index,
			ItemName:          item_name,
			Name:              name,
			ImageInventory:    image_inventory,
			ItemSet:           item_set,
			TournamentEventId: tournament_event_id,
			MarketHashName:    modules.GenerateMarketHashName(t, item_name, "souvenir_package"),
		}

		souvenir_packages = append(souvenir_packages, current)
	}

	// Save knives to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' souvenir packages in %s", len(souvenir_packages), duration)

	return souvenir_packages
}
