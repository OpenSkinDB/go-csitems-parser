package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go-csitems-parser/models"
	"go-csitems-parser/modules"
	"go-csitems-parser/modules/parsers"

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
	sticker_kits := parsers.ParseStickerKits(ctx, itemsGame, translator)
	keychains := parsers.ParseKeychains(ctx, itemsGame, translator)
	weapons := parsers.ParseWeapons(ctx, itemsGame, translator)
	gloves := parsers.ParseGloves(ctx, itemsGame, translator)
	knives := parsers.ParseKnives(ctx, itemsGame, translator)
	highlight_reels := parsers.ParseHighlightReels(ctx, itemsGame, translator)
	tournaments := parsers.ParseTournaments(ctx, translator)

	// Paint kits are tricky
	item_sets := parsers.ParseItemSets(ctx, itemsGame, souvenir_packages, weapon_cases, translator)
	paint_kits := parsers.ParsePaintKits(ctx, itemsGame, translator)

	ExportToJsonFile(paint_kits, "paint_kits")
	ExportToJsonFile(weapons, "weapons")

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
