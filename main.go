package main

import (
	"fmt"

	parsers "github.com/openskindb/openskindb-csitems/parsers"
	"github.com/rs/zerolog/log"
	// models "github.com/openskindb/openskindb-csitems/models"
)

func main() {
	itemsGame := parsers.LoadItemsGame("./files/items_game.txt")
	log.Info().Msg("Loaded items_game.txt successfully")

	musicKits := parsers.ParseMusicKits(itemsGame)
	collectibles := parsers.ParseCollectibles(itemsGame)

	ExportToJsonFile(musicKits, "music_kits")
	ExportToJsonFile(collectibles, "collectibles")

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
};