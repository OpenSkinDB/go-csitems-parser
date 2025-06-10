package openskindb_parsers

import (
	"context"
	"time"

	models "github.com/openskindb/openskindb-csitems/models"

	"github.com/rs/zerolog"
)

func ParseRarities(ctx context.Context, ig *models.ItemsGame) []models.Rarity {
	logger := zerolog.Ctx(ctx);

	start := time.Now()
	logger.Info().Msg("Parsing skin rarities...")

	rarities, err := ig.Get("rarities")
	
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get rarities from items_game.txt")
		return nil
	}
	
	colors, _ := ig.Get("colors")
	cmap := make(map[string]models.GenericColor)
	for _, clr := range colors.GetChilds() {
		color_name, _ := clr.GetString("color_name")
		hex_color, _ := clr.GetString("hex_color")

		curr := models.GenericColor{
			Key: clr.Key,
			ColorName: color_name,
			HexColor: hex_color,
		}

		cmap[clr.Key] = curr
	}

	logger.Debug().Msgf("Found '%d' colors", len(cmap))

	var items []models.Rarity
	for _, r := range rarities.GetChilds() {
		loc_key, _ := r.GetString("loc_key")
		loc_key_weapon, _ := r.GetString("loc_key_weapon")
		loc_key_character, _ := r.GetString("loc_key_character")
		drop_sound, _ := r.GetString("drop_sound")
		
		current := models.Rarity{
			Key: r.Key,
			LocKey: loc_key,
			LocKeyWeapon: loc_key_weapon,	
			LocKeyCharacter: loc_key_character,
			DropSound: drop_sound,
		}
		
		// Get color Data
		color_str, _ := r.GetString("color")

		// loop through the color map to find the matching color
		if color_str != "" {
			// If the color exists in the map, use it
			if colorData, exists := cmap[color_str]; exists {
				logger.Debug().Msgf("Found color '%s' for rarity '%s'", colorData.ColorName, current.LocKey)
				
				current.HexColor = colorData.HexColor
				current.ColorName = colorData.ColorName
			}
		}

		items = append(items, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' rarities in %s", len(items), duration)

	return items
}

// func