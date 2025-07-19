package main

import (
	"context"
	"fmt"
	"maps"
	"os"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"
	"go-csitems-parser/modules/parsers"

	"github.com/rs/zerolog"
)

type ItemSchema struct {
	Collections    []models.Collection                   `json:"collections"`
	Rarities       []models.SchemaRarity                 `json:"rarities"`
	Stickers       map[int]string                        `json:"stickers"`
	Keychains      map[int]string                        `json:"keychains"`
	Collectibles   map[int]models.SchemaGenericeMap      `json:"collectibles"`
	Containers     map[int]string                        `json:"containers"`
	Agents         map[int]models.SchemaGenericeMap      `json:"agents"`
	CustomStickers map[string]models.SchemaCustomSticker `json:"custom_stickers"`
	MusicKits      map[int]models.SchemaGenericeMap      `json:"music_kits"`
	Weapons        map[int]modules.SchemaWeaponSkinMap   `json:"weapons"`

	HighlightReels []models.HighlightReel `json:"highlight_reels"`
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
	keychains := parsers.ParseKeychains(ctx, itemsGame, translator)
	weapons := parsers.ParseWeapons(ctx, itemsGame, translator)
	gloves := parsers.ParseGloves(ctx, itemsGame, translator)
	knives := parsers.ParseKnives(ctx, itemsGame, translator)
	highlight_reels := parsers.ParseHighlightReels(ctx, itemsGame, translator)

	sticker_capsules := parsers.ParseStickerCapsules(ctx, itemsGame, translator)

	misc_capsules := parsers.ParseSelfOpeningCrates(ctx, itemsGame, translator)
	sticker_kits := parsers.ParseStickerKits(ctx, itemsGame, translator)
	custom_stickers := parsers.ParseCustomStickers(ctx, itemsGame, sticker_kits, translator)

	// Paint kits are tricky
	item_sets := parsers.ParseItemSets(ctx, itemsGame, souvenir_packages, weapon_cases, translator)
	paint_kits := parsers.ParsePaintKits(ctx, itemsGame, translator)

	// ExportToJsonFile(paint_kits, "paint_kits")
	// ExportToJsonFile(weapons, "weapons")
	// ExportToJsonFile(sticker_capsules, "sticker_capsules")
	// ExportToJsonFile(custom_stickers, "custom_stickers")
	// ExportToJsonFile(sticker_kits, "sticker_kits")
	// ExportToJsonFile(misc_capsules, "misc_capsules")

	// Special parsing for collections
	collections := parsers.ParseCollections(ctx, itemsGame, souvenir_packages, weapon_cases, translator)
	// Now we need to map whether or not an item has a souvenir variant or not, same for stattrak

	duration := time.Since(start)
	logger.Debug().Msgf("[go-items] Parsed all items in %s", duration)

	// Some knife stuff
	knife_skin_map := modules.LoadKnifeSkinsMap("./files/knife_skins.json")
	knife_skins := modules.GetKnifePaintKits(&knives, &paint_kits, knife_skin_map)
	weapon_skins := modules.GetWeaponPaintKits(&weapons, &paint_kits, &item_sets)
	glove_skins := modules.GetGlovePaintKits(&gloves, &paint_kits)

	// Im new to Go, so idk
	paint_kits_list := make(map[int]modules.SchemaWeaponSkinMap, 0)
	maps.Copy(paint_kits_list, weapon_skins)
	maps.Copy(paint_kits_list, knife_skins)
	maps.Copy(paint_kits_list, glove_skins)

	// Final schema
	itemSchema := ItemSchema{
		Collections:  collections,
		Rarities:     modules.MapRarities(&rarities),
		Stickers:     modules.MapStickerKits(&sticker_kits),
		Keychains:    modules.MapKeychains(&keychains),
		Collectibles: modules.MapCollectibles(&collectibles),
		Containers: modules.MapContainers(
			&weapon_cases,
			&souvenir_packages,
			&sticker_capsules,
			&misc_capsules,
		),
		Agents:         modules.MapAgents(&player_agents),
		HighlightReels: highlight_reels,
		CustomStickers: modules.MapCustomStickers(&custom_stickers),
		MusicKits:      modules.MapMusicKits(&musicKits),
		Weapons:        paint_kits_list,
	}

	// Export the final schema to a JSON file
	ExportToJsonFile(itemSchema, "item_schema")

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
