package modules

import (
	"go-csitems-parser/models"
	"slices"
	"strings"
)

func GetProPlayersAndTeams(ig *models.ItemsGame) ([]int, []int) {
	sticker_kits, err := ig.Get("sticker_kits")

	var player_ids []int
	var team_ids []int

	if err != nil {
		return player_ids, team_ids // Return empty slices if there's an error
	}

	for _, player := range sticker_kits.GetChilds() {
		tournament_player_id, _ := player.GetInt("tournament_player_id")
		tournament_team_id, _ := player.GetInt("tournament_team_id")

		if tournament_player_id > 0 {
			if slices.Contains(player_ids, tournament_player_id) {
				continue // Skip if player ID already exists
			}
			player_ids = append(player_ids, tournament_player_id)
		}

		if tournament_team_id > 0 {
			if slices.Contains(team_ids, tournament_team_id) {
				continue // Skip if team ID already exists
			}
			team_ids = append(team_ids, tournament_team_id)
		}
	}

	return player_ids, team_ids
}

func GetStickerType(player_id int, event_id int, team_id int) string {
	if player_id > 0 {
		return "player"
	}

	if team_id > 0 {
		return "team"
	}

	if event_id > 0 {
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
