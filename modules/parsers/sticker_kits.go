package openskindb_parsers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	models "github.com/zwolof/go-csitems-parser/models"
)

type StickerTypeParams struct {
	TournamentPlayerId 	int 	`json:"tournament_player_id"`
	TournamentTeamId 		int 	`json:"tournament_team_id"`
	TournamentEventId 	int 	`json:"tournament_event_id"`
}

func ParseStickerKits(ctx context.Context, ig *models.ItemsGame) []models.StickerKit {
	logger := zerolog.Ctx(ctx);

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
		description_string, _ := item.GetString("description_string")
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
			DefinitionIndex: 		definition_index,
			Name:            		name,
			ItemName:        		item_name,
			DescriptionString: 	description_string,
			StickerMaterial:    sticker_material,
			Rarity:         		item_rarity,
			Effect: 			 			sticker_effect,
			Type:             	sticker_type,
			TournamentEventId: 	tournament_event_id,
			TournamentTeamId: 	tournament_team_id,
		})
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' sticker kits in %s", len(items), duration)

	return items
}

func GetStickerType(player int, event int, team int) models.StickerType {
	if player > 0 {
		return models.StickerTypeAutograph
	}

	if team > 0 {
		return models.StickerTypeTeam
	}

	if event > 0 {
		return models.StickerTypeEvent
	}

	return models.StickerTypeUnknown
}

func GetStickerEffect(sticker_material string) models.StickerEffect {
	if strings.HasSuffix(sticker_material, "_glitter") {
		return models.StickerEffectGlitter
	}

	if strings.HasSuffix(sticker_material, "_holo") {
		return models.StickerEffectHolo
	}

	if strings.HasSuffix(sticker_material, "_foil") {
		return models.StickerEffectFoil
	}

	if strings.HasSuffix(sticker_material, "_gold") {
		return models.StickerEffectGold
	}

	return models.StickerEffectUnknown
}