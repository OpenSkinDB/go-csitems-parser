package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"
	"go-csitems-parser/modules/parsers"

	"github.com/jedib0t/go-pretty/list"
	"github.com/rs/zerolog"
)

type ItemSchema struct {
	Collections      []models.Collection      `json:"collections"`
	Rarities         []models.Rarity          `json:"rarities"`
	Stickers         []models.StickerKit      `json:"stickers"`
	Keychains        []models.Keychain        `json:"keychains"`
	Collectibles     []models.Collectible     `json:"collectibles"`
	Containers       []models.WeaponCase      `json:"containers"`
	SouvenirPackages []models.SouvenirPackage `json:"souvenir_packages"`
	Agents           []models.PlayerAgent     `json:"agents"`
	MusicKits        []models.MusicKit        `json:"music_kits"`
	HighlightReels   []models.HighlightReel   `json:"highlight_reels"`
	PaintKits        []modules.WeaponSkinMap  `json:"paint_kits"`
	Tournaments      map[string]string        `json:"tournaments"`
	Locations        map[string]string        `json:"locations"`
}

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	// Set the global logger to use the console writer
	itemsGame := modules.LoadItemsGame("./files/items_game.txt")

	if itemsGame == nil {
		logger.Error().Msg("Failed to load items_game.txt, please check the file path and format.")
		panic("items_game.txt is nil, exiting...")
	} else {
		logger.Info().Msgf("Successfully loaded items_game.txt")

		l := list.NewWriter()
		l.SetStyle(list.StyleConnectedRounded)

		sorted := itemsGame.GetChilds()

		// Sort based on the number of their children
		sort.Slice(sorted, func(i, j int) bool {
			return len(sorted[i].GetChilds()) > len(sorted[j].GetChilds())
		})

		for _, item := range sorted {
			fmtKey := GetFormattedItemName(item.Key, len(item.GetChilds()), 35)
			l.AppendItem(fmtKey)
		}
		// fmt.Printf("%s\n", l.Render())
	}

	// Attach the Logger to the context.Context
	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	factory := modules.LoadAllTranslations(ctx, "./files/translations")

	if factory == nil {
		logger.Error().Msg("Failed to load translations")
		return
	}

	translator := factory.GetTranslator("English")
	start := time.Now()

	player_agents := parsers.ParseAgents(ctx, itemsGame, translator)
	souvenir_packages := parsers.ParseSouvenirPackages(ctx, itemsGame, translator)
	musicKits := parsers.ParseMusicKits(ctx, itemsGame, translator)
	collectibles := parsers.ParseCollectibles(ctx, itemsGame, translator)
	weapon_cases := parsers.ParseWeaponCases(ctx, itemsGame, translator)
	rarities := parsers.ParseRarities(ctx, itemsGame, translator)
	paint_kits := parsers.ParsePaintKits(ctx, itemsGame, translator)
	item_sets := parsers.ParseItemSets(ctx, itemsGame, translator)
	sticker_kits := parsers.ParseStickerKits(ctx, itemsGame, translator)
	keychains := parsers.ParseKeychains(ctx, itemsGame, translator)
	loot_lists := parsers.ParseClientLootLists(ctx, itemsGame, translator)
	weapons := parsers.ParseWeapons(ctx, itemsGame, translator)
	gloves := parsers.ParseGloves(ctx, itemsGame, translator)
	knives := parsers.ParseKnives(ctx, itemsGame, translator)
	highlight_reels := parsers.ParseHighlightReels(ctx, itemsGame, translator)
	tournaments := parsers.ParseTournaments(ctx, translator)

	// Special parsing for collections
	collections := parsers.ParseCollections(ctx, itemsGame, souvenir_packages, weapon_cases, translator)

	duration := time.Since(start)
	logger.Debug().Msgf("[go-items] Parsed all items in %s", duration)

	// Export all parsed data to JSON files
	// debugging
	ExportToJsonFile(musicKits, "music_kits")
	ExportToJsonFile(collectibles, "collectibles")
	ExportToJsonFile(weapon_cases, "weapon_cases")
	ExportToJsonFile(player_agents, "player_agents")
	ExportToJsonFile(rarities, "rarities")
	ExportToJsonFile(paint_kits, "paint_kits")
	ExportToJsonFile(item_sets, "item_sets")
	ExportToJsonFile(sticker_kits, "sticker_kits")
	ExportToJsonFile(keychains, "keychains")
	ExportToJsonFile(loot_lists, "client_loot_lists")
	ExportToJsonFile(weapons, "weapons")
	ExportToJsonFile(gloves, "gloves")
	ExportToJsonFile(knives, "knives")
	ExportToJsonFile(souvenir_packages, "souvenir_packages")
	ExportToJsonFile(highlight_reels, "highlight_reels")
	ExportToJsonFile(tournaments, "tournaments")
	ExportToJsonFile(collections, "collections")

	var itemsGameCdn = modules.LoadItemsGameCdn("./files/items_game_cdn.txt")
	ExportToJsonFile(itemsGameCdn, "items_game_cdn")

	// Some knife stuff
	knife_skin_map := modules.LoadKnifeSkinsMap("./files/knife_skins.json")
	knife_skins := modules.GetKnifePaintKits(&knives, &paint_kits, knife_skin_map)
	ExportToJsonFile(knife_skins, "knives_with_paint_kits")

	// Weapon map
	weapon_skins := modules.GetWeaponPaintKits(&weapons, &paint_kits, &item_sets)
	ExportToJsonFile(weapon_skins, "weapons_with_paint_kits")
	// knife_skins := modules.GetKnifePaintKits(&knives, &paint_kits)

	// Some glove stuff
	glove_skins := modules.GetGlovePaintKits(&gloves, &paint_kits)
	ExportToJsonFile(glove_skins, "gloves_with_paint_kits")

	paint_kits_list := make([]modules.WeaponSkinMap, 0)
	paint_kits_list = append(paint_kits_list, weapon_skins...)
	paint_kits_list = append(paint_kits_list, knife_skins...)
	paint_kits_list = append(paint_kits_list, glove_skins...)

	// Final schema
	itemSchema := ItemSchema{
		Collections:      collections,
		Rarities:         rarities,
		Stickers:         sticker_kits,
		Keychains:        keychains,
		Collectibles:     collectibles,
		Containers:       weapon_cases,
		SouvenirPackages: souvenir_packages,
		Agents:           player_agents,
		MusicKits:        musicKits,
		HighlightReels:   highlight_reels,
		PaintKits:        paint_kits_list,
		Tournaments:      tournaments.Tournaments,
		Locations:        tournaments.Locations,
	}

	// Export the final schema to a JSON file
	ExportToJsonFile(itemSchema, "item_schema")

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
