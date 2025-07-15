package parsers

import (
	"context"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

func ParseCollections(
	ctx context.Context,
	ig *models.ItemsGame,
	sv []models.SouvenirPackage,
	cs []models.WeaponCase,
	t *modules.Translator,
) []models.Collection {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	item_sets, err := ig.Get("item_sets")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get item_sets from items_game.txt")
		return nil
	}

	var collections []models.Collection
	for _, s := range item_sets.GetChilds() {
		name, _ := s.GetString("name")

		if strings.Contains(name, "_characters") {
			// Skip if the name contains "_characters"
			continue
		}

		current := models.Collection{
			Key:  s.Key,
			Name: modules.GenerateMarketHashName(t, name, "collection"),
		}

		// We're done here, add the current item set to the list
		collections = append(collections, current)
	}

	// Check if a souvenir package exists with the same itemset.id
	for _, c := range cs {
		for i, col := range collections {
			if c.ItemSetId == nil {
				continue
			}
			if c.ItemSetId == &col.Key {
				collections[i].HasCrate = true
				break
			}
		}
	}

	// Same for weapon cases
	for _, sv_pkg := range sv {
		for i, col := range collections {

			if sv_pkg.ItemSetId == nil {
				continue
			}

			if sv_pkg.ItemSetId == &col.Key {
				collections[i].HasSouvenir = true
				break
			}
		}
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' collections in %s", len(collections), duration)

	return collections
}
