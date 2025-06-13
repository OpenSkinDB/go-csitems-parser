package openskindb_parsers

import (
	"context"
	"strings"
	"time"

	models "github.com/zwolof/go-csitems-parser/models"

	"github.com/rs/zerolog"
)

func ParseWeapons(ctx context.Context, ig *models.ItemsGame) []models.BaseWeapon {
	logger := zerolog.Ctx(ctx);

	start := time.Now()
	// logger.Info().Msg("Parsing music kits...")

	prefabs, err := ig.Get("prefabs")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get prefabs")
		return nil
	}

	var weapons []models.BaseWeapon

	for _, w := range prefabs.GetChilds() {

		if !strings.HasPrefix(w.Key, "weapon_") || !strings.HasSuffix(w.Key, "_prefab") {
			// Skip non-weapon prefabs
			continue
		}

		item_name, _ := w.GetString("item_name")
		item_description, _ := w.GetString("item_description")
		item_class, _ := w.GetString("item_class")
		slot, _ := w.GetString("prefab")
		image_inventory, _ := w.GetString("image_inventory")
		
		current := models.BaseWeapon{
			ItemName: 			 	item_name,
			ItemDescription: 	item_description,
			ItemClass:        item_class,
			Slot:            	slot,
			ImageInventory: 	image_inventory,
		}
		
		attributes, _ := w.Get("attributes")
		
		if attributes != nil {
			priceAttr, _ := attributes.GetInt("in game price")
			current.WeaponPrice = priceAttr
		}

		used_by_classes, _ := w.Get("used_by_classes")

		teams := []string{}
		if used_by_classes != nil {
			for _, team := range used_by_classes.GetChilds() {
				teams = append(teams, team.Key)
			}
		}

		current.Teams = teams

		weapons = append(weapons, current)
	}

	// Save weapons to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' weapons in %s", len(weapons), duration)

	return weapons
}