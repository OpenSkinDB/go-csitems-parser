package modules

import (
	"context"
	"fmt"
	"go-csitems-parser/models"
	"strconv"
	"strings"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

func GetLogger() *zerolog.Logger {
	logger := zerolog.Ctx(context.Background())

	return logger
}

func GetStringMapKeySlice(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetStringMapValueSlice(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetTournamentEventId(item *vdf.KeyValue) (int, error) {
	attributes, err := item.Get("attributes")
	if err != nil {
		return -1, err
	}

	tournament, err := attributes.Get("tournament event id")
	if err != nil {
		return -1, err
	}

	tournament_event_id, err := tournament.GetInt("value")
	if err != nil {
		return -1, err
	}

	return tournament_event_id, nil
}

func GetContainerItemSet(item *vdf.KeyValue, t *Translator, key string) *string {
	tags, err := item.Get("tags")

	if err != nil {
		return nil
	}

	var container_item_set_key string = "ItemSet"
	if key != "" {
		container_item_set_key = key
	}

	item_set, err := tags.Get(container_item_set_key)
	if err != nil {
		return nil
	}

	tag, _ := item_set.GetString("tag_value")
	// tagText, _ := item_set.GetString("tag_text")

	// translated, _ := t.GetValueByKey(tagText)

	return &tag
}

func GetSupplyCrateSeries(item *vdf.KeyValue, ig *models.ItemsGame) *string {
	attributes, err := item.Get("attributes")

	if err != nil {
		return nil
	}

	set_supply_crate_series, err := attributes.Get("set supply crate series")
	if err != nil {
		return nil
	}

	series_id, err := set_supply_crate_series.GetString("value")
	if err != nil {
		return nil
	}

	revolving_loot_lists, _ := ig.Get("revolving_loot_lists")

	for _, list := range revolving_loot_lists.GetChilds() {
		if list.Key == series_id {
			value, _ := list.ToString()

			return &value
		}
	}

	return nil
}

type ItemWear struct {
	Name    string  `json:"name"`
	MinWear float64 `json:"min_wear"`
	MaxWear float64 `json:"max_wear"`
}

var ItemWears = map[string]ItemWear{
	"Factory New": {
		Name:    "Factory New",
		MinWear: 0.00,
		MaxWear: 0.07,
	},
	"Minimal Wear": {
		Name:    "Minimal Wear",
		MinWear: 0.07,
		MaxWear: 0.15,
	},
	"Field-Tested": {
		Name:    "Field-Tested",
		MinWear: 0.15,
		MaxWear: 0.38,
	},
	"Well-Worn": {
		Name:    "Well-Worn",
		MinWear: 0.38,
		MaxWear: 0.45,
	},
	"Battle-Scarred": {
		Name:    "Battle-Scarred",
		MinWear: 0.45,
		MaxWear: 1.00,
	},
}

var hashNamePrefixes = map[string]string{
	"sticker_kit": "Sticker | ",
	"music_kit":   "Music Kit | ",
	"keychain":    "Charm | ",
	// "highlight_reel": "Souvenir Charm | Austin 2025 Highlight | ",
}

var sticker_effect_names = map[string]string{
	"glitter":    " (Glitter)",
	"holo":       " (Holo)",
	"foil":       " (Foil)",
	"gold":       " (Gold)",
	"lenticular": " (Lenticular)",
	"normal":     " (Normal)",
}

func GenerateCustomStickerMarketHashName_Team(t *Translator, team_id int, effect *string) string {
	lang_key := fmt.Sprintf("CSGO_TeamID_%d", team_id)
	team_name, _ := t.GetValueByKey(lang_key)

	if effect == nil {
		return fmt.Sprintf("Sticker | %s", team_name) // Fallback if effect or type is nil
	}

	effect_name := sticker_effect_names[*effect]
	return fmt.Sprintf("Sticker | %s%s", team_name, effect_name)
}

func GenerateCustomStickerMarketHashName_Player(t *Translator, player *models.TournamentData, effect *string) string {
	if effect == nil {
		return fmt.Sprintf("Sticker | %s", player.Name) // Fallback if effect is nil
	}

	effect_name := sticker_effect_names[*effect]
	return fmt.Sprintf("Sticker | %s%s", player.Name, effect_name)
}

func GenerateCustomStickerMarketHashName_Event(t *Translator, event_id int, effect *string) string {
	lang_key := fmt.Sprintf("CSGO_Tournament_Event_Location_%d", event_id)

	translated, _ := t.GetValueByKey(lang_key)

	if effect == nil {
		return fmt.Sprintf("Sticker | %s", translated) // Fallback if effect or type is nil
	}

	effect_name := sticker_effect_names[*effect]
	return fmt.Sprintf("Sticker | %s%s", translated, effect_name)
}

func GenerateHighlightReelMarketHashName(t *Translator, name string, event int) string {
	value, err := t.GetValueByKey("HighlightReel_" + name)

	if err != nil {
		fmt.Printf("Error translating name '%s': %v\n", name, err)
		value = name // Fallback to original name if translation fails
	}

	// split name by "_", first part is the tournament id for the keychain capsule
	tournament_id := ""
	if len(name) > 0 {
		parts := strings.Split(name, "_")

		if len(parts) > 0 {
			tournament_id = parts[0]
		}
	}

	// Get the capsule name "keychain_kc_%s"
	capsule, err := t.GetValueByKey("keychain_kc_" + tournament_id)

	if err != nil {
		fmt.Printf("Error translating capsule name '%s': %v\n", tournament_id, err)
		capsule = "Unknown Capsule" // Fallback to a default value
	}

	return fmt.Sprintf("Souvenir Charm | %s | %s", capsule, value)
}

func GetTournamentData(t *Translator, id int) *models.TournamentData {
	lang_key := fmt.Sprintf("CSGO_Tournament_Event_NameShort_%d", id)

	name, _ := t.GetValueByKey(lang_key)

	if id == 0 || name == "" {
		return nil
	}

	return &models.TournamentData{
		Id:   id,
		Name: name,
	}
}

func GetTournamentStageData(t *Translator, id int) *models.TournamentData {
	lang_key := fmt.Sprintf("CSGO_Tournament_Event_Stage_%d", id)

	name, _ := t.GetValueByKey(lang_key)

	if id == 0 || name == "" {
		return nil
	}

	return &models.TournamentData{
		Id:   id,
		Name: name,
	}
}

func GetPlayerByAccountId(ig *models.ItemsGame, account_id int) *models.TournamentData {
	pro_players, _ := ig.Get("pro_players")

	for _, player := range pro_players.GetChilds() {
		current_aid, _ := strconv.Atoi(player.Key)

		if current_aid != account_id {
			continue
		}
		name, _ := player.GetString("name")

		return &models.TournamentData{
			Id:   account_id,
			Name: name,
		}
	}
	return nil
}

func GetTournamentTeamData(t *Translator, id int) *models.TournamentData {
	lang_key := fmt.Sprintf("CSGO_TeamID_%d", id)

	name, _ := t.GetValueByKey(lang_key)

	if id == 0 || name == "" {
		return nil
	}

	return &models.TournamentData{
		Id:   id,
		Name: name,
	}
}

var dopplerPhaseMap = map[string]string{
	// Doppler phases
	"am_ruby_marbleized":         "Ruby",
	"am_ruby_marbleized_b":       "Ruby",
	"am_sapphire_marbleized":     "Sapphire",
	"am_sapphire_marbleized_b":   "Sapphire",
	"am_blackpearl_marbleized":   "Black Pearl",
	"am_blackpearl_marbleized_b": "Black Pearl",

	// Phase 1-4, with and without "b" suffix???
	"am_doppler_phase1":   "Phase 1",
	"am_doppler_phase2":   "Phase 2",
	"am_doppler_phase3":   "Phase 3",
	"am_doppler_phase4":   "Phase 4",
	"am_doppler_phase1_b": "Phase 1",
	"am_doppler_phase2_b": "Phase 2",
	"am_doppler_phase3_b": "Phase 3",
	"am_doppler_phase4_b": "Phase 4",

	// Gamma Doppler phases
	"am_emerald_marbleized":   "Emerald",
	"am_emerald_marbleized_b": "Emerald",
	"am_gamma_doppler_phase1": "Phase 1",
	"am_gamma_doppler_phase2": "Phase 2",
	"am_gamma_doppler_phase3": "Phase 3",
	"am_gamma_doppler_phase4": "Phase 4",

	// Gamma Doppler Glock phases
	"am_emerald_marbleized_glock":   "Emerald",
	"am_gamma_doppler_phase1_glock": "Phase 1",
	"am_gamma_doppler_phase2_glock": "Phase 2",
	"am_gamma_doppler_phase3_glock": "Phase 3",
	"am_gamma_doppler_phase4_glock": "Phase 4",
}

func GenerateMarketHashName(t *Translator, name string, extra *string, item_type string) string {
	value, err := t.GetValueByKey(name)

	if err != nil {
		fmt.Printf("Error translating name '%s': %v\n", name, err)
		value = name // Fallback to original name if translation fails
	}

	// Special case for the vanilla paint kit
	if name == "#PaintKit_Default_Tag" {
		value = "Vanilla"
	}

	// If the item type is a doppler, we need to add the phase to the name
	if extra != nil && *extra != "" {
		if phase, ok := dopplerPhaseMap[*extra]; ok {
			value = fmt.Sprintf("%s (%s)", value, phase)
		}
	}

	if item_type == "knife" || item_type == "glove" {
		value = fmt.Sprintf("★ %s", value)
	}

	if prefix, ok := hashNamePrefixes[item_type]; ok {
		value = prefix + value
	}

	return value
}

// Performance might be bad, but it is a one-time operation once every few months so its fine
func AddPaintKitMappings(
	item_sets *[]models.ItemSet,
	paint_kits *[]models.PaintKit,
) []models.PaintKit {
	// create a clone of paint_kits to avoid modifying the original slice
	paint_kits_clone := make([]models.PaintKit, len(*paint_kits))
	copy(paint_kits_clone, *paint_kits)

	for _, set := range *item_sets {

		// If the item set has a crate, we assume it can be StatTrak
		items_can_be_stattrak := set.HasCrate

		// If the item set has a souvenir, we assume it can be Souvenir
		items_can_be_souvenir := set.HasSouvenir

		for _, item := range set.Items {
			for _, paint_kit := range paint_kits_clone {
				if paint_kit.Name == item.PaintKitName {

					paint_kit.StatTrak = items_can_be_stattrak
					paint_kit.Souvenir = items_can_be_souvenir
				}
			}
		}
	}

	return paint_kits_clone
}
