package parsers

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"

	"github.com/baldurstod/vdf"
)

type StickerTypeParams struct {
	TournamentPlayerId int `json:"tournament_player_id"`
	TournamentTeamId   int `json:"tournament_team_id"`
	TournamentEventId  int `json:"tournament_event_id"`
}

var sticker_types = []string{
	"normal",
	"player",
	"team",
	"event",
}

var sticker_effects = []string{
	"normal",
	"foil",
	"holo",
	"glitter",
	"gold",
}

var group_id_to_sub_id = map[int]string{
	1: "E",
	2: "P",
	3: "T",
}

func ParseCustomStickers(ctx context.Context, ig *models.ItemsGame, sticker_kits []models.StickerKit, t *modules.Translator) []models.CustomStickers {
	logger := modules.GetLogger()

	start := time.Now()

	// Store all custom stickers
	var items []models.CustomStickers

	if len(sticker_kits) == 0 {
		logger.Warn().Msg("No sticker kits found, skipping custom stickers parsing")
		return items
	}

	var unique_tournament_ids []int
	var unique_player_ids []int
	var unique_team_ids []int

	// Handle all the event stickers
	for _, sticker_kit := range sticker_kits {

		// Tournament stickers only
		if sticker_kit.Tournament != nil && sticker_kit.Player == nil {
			tournament_id := sticker_kit.Tournament.Id

			if tournament_id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid tournament ID: %s", sticker_kit.Name)
				continue // Skip if the tournament ID is invalid
			}

			// We need to get the sticker type, and based on that, we can process it
			current_type := modules.GetStickerType(0, tournament_id, 0)
			current_effect := modules.GetStickerEffect(sticker_kit.StickerMaterial)

			// Get the count of stickers for this event and type
			count := GetStickerCountByTournamentId(&sticker_kits, tournament_id, current_effect, false)

			if count == 0 {
				continue // Skip if no stickers found for the team and type
			}

			generated_id := GenerateCustomStickerId(tournament_id, group_id_to_sub_id[1], &current_effect, &current_type)

			if CustomStickerExists(items, generated_id) {
				continue // Skip if the custom sticker already exists
			}

			items = append(items, models.CustomStickers{
				GeneratedId: generated_id,
				Group:       2,
				Count:       count,
				Name:        modules.GenerateCustomStickerMarketHashName_Event(t, tournament_id, &current_effect),
			})

			if !slices.Contains(unique_tournament_ids, tournament_id) {
				unique_tournament_ids = append(unique_tournament_ids, tournament_id)
			}
		}

		// Team stickers only
		if sticker_kit.Team != nil && sticker_kit.Player == nil {
			team_id := sticker_kit.Team.Id

			if team_id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid team ID: %s", sticker_kit.Name)
				continue // Skip if the team ID is invalid
			}

			// We need to get the sticker type, and based on that, we can process it
			current_type := modules.GetStickerType(0, 0, team_id)
			current_effect := modules.GetStickerEffect(sticker_kit.StickerMaterial)

			// Get the count of stickers for this event and type
			count := GetStickerCountByTeamId(&sticker_kits, team_id, current_effect, false)

			if count == 0 {
				continue // Skip if no stickers found for the team and type
			}

			generated_id := GenerateCustomStickerId(team_id, group_id_to_sub_id[3], &current_effect, &current_type)

			if CustomStickerExists(items, generated_id) {
				continue // Skip if the custom sticker already exists
			}

			items = append(items, models.CustomStickers{
				GeneratedId: generated_id,
				Group:       3,
				Count:       count,
				Name:        modules.GenerateCustomStickerMarketHashName_Team(t, team_id, &current_effect),
			})

			if !slices.Contains(unique_team_ids, team_id) {
				unique_team_ids = append(unique_team_ids, team_id)
			}
		}

		// Player stickers only
		if sticker_kit.Player != nil {
			player_id := sticker_kit.Player.Id

			if player_id <= 0 {
				logger.Debug().Msgf("Skipping sticker kit with invalid player ID: %s", sticker_kit.Name)
				continue // Skip if the player ID is invalid
			}

			// We need to get the sticker type, and based on that, we can process it
			current_type := modules.GetStickerType(player_id, 0, 0)
			current_effect := modules.GetStickerEffect(sticker_kit.StickerMaterial)

			// Get the count of stickers for this event and type
			count := GetStickerCountByPlayerId(&sticker_kits, player_id, current_effect, false)

			if count == 0 {
				continue // Skip if no stickers found for the player and type
			}

			generated_id := GenerateCustomStickerId(player_id, group_id_to_sub_id[2], &current_effect, &current_type)

			if CustomStickerExists(items, generated_id) {
				continue // Skip if the custom sticker already exists
			}

			items = append(items, models.CustomStickers{
				GeneratedId: generated_id,
				Group:       2,
				Count:       count,
				Name:        modules.GenerateCustomStickerMarketHashName_Player(t, sticker_kit.Player, &current_effect),
			})

			if !slices.Contains(unique_player_ids, player_id) {
				unique_player_ids = append(unique_player_ids, player_id)
			}
		}
	}

	// now we need to get the total per player/team/event
	for _, tournament_id := range unique_tournament_ids {
		curr := GetTotalStickerForSubId(ig, t, sticker_kits, 1, tournament_id)
		items = append(items, curr)
	}

	for _, player_id := range unique_player_ids {
		curr := GetTotalStickerForSubId(ig, t, sticker_kits, 2, player_id)

		items = append(items, curr)
	}

	for _, team_id := range unique_team_ids {
		curr := GetTotalStickerForSubId(ig, t, sticker_kits, 3, team_id)
		items = append(items, curr)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' custom stickers in %s", len(items), duration)

	return items
}

func GetTotalStickerForSubId(ig *models.ItemsGame, t *modules.Translator, sticker_kits []models.StickerKit, group_id int, sub_id int) models.CustomStickers {
	var count int = 0

	// Generate a unique ID for the custom sticker
	generated_id := GenerateCustomStickerId(sub_id, group_id_to_sub_id[group_id], nil, nil)

	for _, kit := range sticker_kits {
		switch group_id {

		// Events
		case 1:
			if kit.Tournament == nil {
				continue // Skip if the sticker kit does not have a tournament
			}

			if kit.Tournament.Id == sub_id {
				count++
			}

		// Players
		case 2:
			if kit.Player == nil {
				continue // Skip if the sticker kit does not have a player
			}

			if kit.Player.Id == sub_id {
				count++
			}

		// Teams
		case 3:
			// Teams are special, we need to ignore ones with player-id
			if kit.Team == nil || kit.Player != nil {
				continue // Skip if the sticker kit does not have a team
			}

			if kit.Team.Id == sub_id {
				count++
			}
		}
	}

	var name string
	switch group_id {
	case 1:
		event := modules.GetTournamentData(t, sub_id)
		name = modules.GenerateCustomStickerMarketHashName_Event(t, event.Id, nil)
	case 2:
		player := modules.GetPlayerByAccountId(ig, sub_id)
		name = modules.GenerateCustomStickerMarketHashName_Player(t, player, nil)
	case 3:
		team := modules.GetTournamentTeamData(t, sub_id)
		name = modules.GenerateCustomStickerMarketHashName_Team(t, team.Id, nil)
	}

	return models.CustomStickers{
		GeneratedId: generated_id,
		Group:       group_id,
		Count:       count,
		Name:        name,
	}
}

func CustomStickerExists(items []models.CustomStickers, generated_id string) bool {
	for _, item := range items {
		if item.GeneratedId == generated_id {
			return true // Found a duplicate
		}
	}
	return false // No duplicates found
}

var sticker_effect_id_suffix = map[string]string{
	"glitter":    "G",
	"holo":       "H",
	"foil":       "F",
	"gold":       "G",
	"lenticular": "L",
	"normal":     "N",
}

var sticker_type_id_suffix = map[string]string{
	"normal": "N",
	"player": "P",
	"team":   "T",
	"event":  "E",
}

func GenerateCustomStickerId(team_id int, _type string, effect *string, sticker_type *string) string {
	// Fallback if effect or type is nil
	if effect == nil || sticker_type == nil {
		return fmt.Sprintf("C%dA%s", team_id, _type)
	}

	// Default format for custom sticker ID
	return fmt.Sprintf("C%d%s%s%s", team_id, sticker_effect_id_suffix[*effect], sticker_type_id_suffix[*sticker_type], _type)
}

func GetStickerKitsBySubId(sticker_kits *vdf.KeyValue, sub_key string, sub_id int) []*vdf.KeyValue {
	var items = make([]*vdf.KeyValue, 0)

	for _, cs := range sticker_kits.GetChilds() {
		subkey_value, _ := cs.GetInt(sub_key)
		sticker_material, _ := cs.GetString("sticker_material")
		name, _ := cs.GetString("name")

		if strings.Contains(name, "_graffiti") || strings.Contains(sticker_material, "patch_") {
			continue // Skip graffitis and patches
		}

		if strings.Contains(sticker_material, "_graffiti") || strings.Contains(sticker_material, "patch_") {
			continue // Skip graffiti and patches
		}

		// Skip if the subkey value does not match
		if subkey_value != sub_id {
			continue
		}

		items = append(items, cs)
	}

	return items
}

func GetStickerKitsByPlayerId(sticker_kits *vdf.KeyValue, player_id int) []*vdf.KeyValue {
	var items = make([]*vdf.KeyValue, 0)

	for _, cs := range sticker_kits.GetChilds() {
		tournament_player_id, _ := cs.GetInt("tournament_player_id")
		sticker_material, _ := cs.GetString("sticker_material")

		if strings.Contains(sticker_material, "_graffiti") || strings.Contains(sticker_material, "patch_") {
			continue // Skip graffiti and patches
		}

		// Skip if the player ID does not match
		if tournament_player_id != player_id {
			continue
		}

		items = append(items, cs)
	}

	return items
}

func GetStickerCountByPlayerId(sticker_kits *[]models.StickerKit, player_id int, sticker_effect string, ignore_effect bool) int {
	var count int

	for _, cs := range *sticker_kits {
		if cs.Player == nil {
			continue // Skip if the sticker kit does not have a player
		}

		// Get the tournament player ID
		if cs.Player.Id != player_id || cs.Effect != sticker_effect && !ignore_effect {
			continue // Skip if the player ID does not match
		}
		count++
	}

	return count
}

func GetStickerCountByTeamId(sticker_kits *[]models.StickerKit, team_id int, sticker_effect string, ignore_effect bool) int {
	var count int

	for _, cs := range *sticker_kits {
		if cs.Team == nil {
			continue // Skip if the sticker kit does not have a team
		}

		// Get the tournament team ID
		if cs.Team.Id != team_id || cs.Effect != sticker_effect && !ignore_effect {
			continue // Skip if the team ID does not match
		}
		count++
	}

	return count
}

func GetStickerCountByTournamentId(sticker_kits *[]models.StickerKit, tournament_id int, sticker_effect string, ignore_effect bool) int {
	var count int

	for _, cs := range *sticker_kits {
		if cs.Tournament == nil {
			continue // Skip if the sticker kit does not have a tournament
		}

		// Get the tournament ID
		if cs.Tournament.Id != tournament_id || cs.Effect != sticker_effect && !ignore_effect {
			continue // Skip if the tournament ID does not match
		}
		count++
	}

	return count
}

func GetCountByParameters(sticker_kits *[]models.StickerKit, subkey string, custom_id int, sticker_effect string) (int, int) {
	var count int
	var total int

	for _, cs := range *sticker_kits {
		switch subkey {
		case "team":
			if cs.Team == nil {
				continue // Skip if the sticker kit does not have a team
			}
			tournament_team_id := cs.Team.Id

			if cs.Player != nil {
				continue // Skip if tournament_player_id is present, we only want team stickers
			}

			if tournament_team_id == custom_id && sticker_effect != "" && cs.Effect == sticker_effect {
				count++
			}
			if tournament_team_id == custom_id {
				total++
			}
		case "player":
			if cs.Player == nil {
				continue // Skip if the sticker kit does not have a player
			}
			tournament_player_id := cs.Player.Id

			if tournament_player_id == custom_id && sticker_effect != "" && cs.Effect == sticker_effect {
				count++
			}
			if tournament_player_id == custom_id {
				total++
			}
		case "event":
			if cs.Tournament == nil {
				continue // Skip if the sticker kit does not have a tournament
			}
			tournament_event_id := cs.Tournament.Id

			if tournament_event_id == custom_id && sticker_effect != "" && cs.Effect == sticker_effect {
				count++
			}

			if tournament_event_id == custom_id {
				total++
			}
		default:
			fmt.Printf("Unknown subkey: %s\n", subkey)
			// Handle unknown subkeys gracefully
			return 0, 0 // Return zero counts for unknown subkeys
		}
	}

	return count, total
}
