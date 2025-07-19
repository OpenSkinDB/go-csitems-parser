package parsers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"
)

func ParseStickerKits(ctx context.Context, ig *models.ItemsGame, t *modules.Translator) []models.StickerKit {
	logger := modules.GetLogger()

	start := time.Now()
	// logger.Info().Msg("Parsing sticker kits...")

	sticker_kits, err := ig.Get("sticker_kits")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get sticker kits from items_game.txt")
		return nil
	}

	var items []models.StickerKit

	// Iterate through all items in the "items" section
	for _, item := range sticker_kits.GetChilds() {
		definition_index, _ := strconv.Atoi(item.Key)
		item_name, _ := item.GetString("item_name")
		name, _ := item.GetString("name")

		if definition_index <= 0 {
			logger.Debug().Msgf("Skipping invalid sticker kit definition index: %d", definition_index)
			continue // Skip invalid definition indices
		}

		if strings.Contains(name, "patch_") {
			continue // Skip non-sticker kit items
		}

		if strings.Contains(name, "_graffiti") {
			continue // Skip non-sticker kit items
		}

		// description_string, _ := item.GetString("description_string")
		sticker_material, _ := item.GetString("sticker_material")

		item_rarity, _ := item.GetString("item_rarity")
		tournament_event_id, _ := item.GetInt("tournament_event_id")
		tournament_team_id, _ := item.GetInt("tournament_team_id")
		tournament_player_id, _ := item.GetInt("tournament_player_id")

		// Get sticker effect, dunno how accurate this is, but it works for now
		sticker_effect := GetStickerEffect(sticker_material)
		sticker_type := GetStickerType(
			tournament_player_id,
			tournament_team_id,
			tournament_event_id,
		)

		items = append(items, models.StickerKit{
			DefinitionIndex: definition_index,
			Name:            name,
			MarketHashName:  modules.GenerateMarketHashName(t, item_name, nil, "sticker_kit"),
			StickerMaterial: sticker_material,
			Rarity:          item_rarity,
			Effect:          sticker_effect,
			Type:            sticker_type,
			Tournament:      modules.GetTournamentData(t, tournament_event_id),
			Team:            modules.GetTournamentTeamData(t, tournament_team_id),
			Player:          modules.GetPlayerByAccountId(ig, tournament_player_id),
		})
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' sticker kits in %s", len(items), duration)

	return items
}

func GetStickerType(player int, event int, team int) string {
	if player > 0 {
		return "autograph"
	}

	if team > 0 {
		return "team"
	}

	if event > 0 {
		return "event"
	}

	return "normal"
}

func GetStickerEffect(sticker_material string) string {
	if strings.HasSuffix(sticker_material, "_glitter") {
		return "glitter"
	}

	if strings.HasSuffix(sticker_material, "_holo") {
		return "holo"
	}

	if strings.HasSuffix(sticker_material, "_foil") {
		return "foil"
	}

	if strings.HasSuffix(sticker_material, "_gold") {
		return "gold"
	}

	if strings.HasSuffix(sticker_material, "_lenticular") {
		return "lenticular"
	}

	return "normal"
}
