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

type ItemSchemaPaintKits struct {
	Weapons *[]modules.WeaponSkinMap `json:"weapons"`
	Knives  *[]modules.KnifeSkinMap  `json:"knives"`
	Gloves  *[]modules.GloveSkinMap  `json:"gloves"`
}

type ItemSchema struct {
	Collections  []models.ItemSet      `json:"collections"`
	Rarities     []models.Rarity       `json:"rarities"`
	Stickers     []*models.StickerKit  `json:"stickers"`
	Keychains    []models.Keychain     `json:"keychains"`
	Collectibles []*models.Collectible `json:"collectibles"`
	Containers   []models.WeaponCase   `json:"containers"`
	Agents       []models.PlayerAgent  `json:"agents"`
	MusicKits    []models.MusicKit     `json:"music_kits"`
	PaintKits    ItemSchemaPaintKits   `json:"paint_kits"`
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

	translator := modules.LoadAllTranslations(ctx, "./files/translations")

	if translator == nil {
		logger.Error().Msg("Failed to load translations")
		return
	}

	start := time.Now()

	musicKits := parsers.ParseMusicKits(ctx, itemsGame)
	collectibles := parsers.ParseCollectibles(ctx, itemsGame)
	weapon_cases := parsers.ParseWeaponCases(ctx, itemsGame)
	player_agents := parsers.ParseAgents(ctx, itemsGame)
	rarities := parsers.ParseRarities(ctx, itemsGame)
	paint_kits := parsers.ParsePaintKits(ctx, itemsGame)
	item_sets := parsers.ParseItemSets(ctx, itemsGame)
	sticker_kits := parsers.ParseStickerKits(ctx, itemsGame)
	keychains := parsers.ParseKeychains(ctx, itemsGame)
	loot_lists := parsers.ParseClientLootLists(ctx, itemsGame)
	weapons := parsers.ParseWeapons(ctx, itemsGame)
	gloves := parsers.ParseGloves(ctx, itemsGame)
	knives := parsers.ParseKnives(ctx, itemsGame)
	souvenir_packages := parsers.ParseSouvenirPackages(ctx, itemsGame)

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

	// Final schema
	itemSchema := ItemSchema{
		Collections:  item_sets,
		Rarities:     rarities,
		Stickers:     sticker_kits,
		Keychains:    keychains,
		Collectibles: collectibles,
		Containers:   weapon_cases,
		Agents:       player_agents,
		MusicKits:    musicKits,
		PaintKits: ItemSchemaPaintKits{
			Weapons: &weapon_skins,
			Knives:  &knife_skins,
			Gloves:  &glove_skins,
		},
	}

	// Export the final schema to a JSON file
	ExportToJsonFile(itemSchema, "item_schema")

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
