package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/list"
	"github.com/rs/zerolog"
	modules "github.com/zwolof/open-csitems-parser/modules"
	parsers "github.com/zwolof/open-csitems-parser/modules/parsers"
	// models "github.com/zwolof/open-csitems-parser/models"
)

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
		fmt.Printf("%s\n", l.Render())
	}

	// Attach the Logger to the context.Context
	ctx := context.Background()
	ctx = logger.WithContext(ctx)
		
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

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
};