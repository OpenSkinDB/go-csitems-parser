package main

import (
	"context"
	"fmt"
	"os"
	"time"

	parsers "github.com/openskindb/openskindb-csitems/parsers"
	"github.com/rs/zerolog"
	// models "github.com/openskindb/openskindb-csitems/models"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().		
		Logger()

	// Set the global logger to use the console writer
	itemsGame := parsers.LoadItemsGame("./files/items_game.txt")

	logger.Info().Msgf("Loaded items_game.txt with %d keys", len(itemsGame.KeyValue.GetChilds()))

	// Attach the Logger to the context.Context
	ctx := context.Background()
	ctx = logger.WithContext(ctx)
		
	musicKits := parsers.ParseMusicKits(ctx, itemsGame)
	collectibles := parsers.ParseCollectibles(ctx, itemsGame)

	ExportToJsonFile(musicKits, "music_kits")
	ExportToJsonFile(collectibles, "collectibles")

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
};