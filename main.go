package main

import (
	"encoding/json"
	"fmt"
	"os"

	parsers "github.com/openskindb/openskindb-csitems/parsers"
	// models "github.com/openskindb/openskindb-csitems/models"
)

func main() {
	itemsGame := parsers.LoadItemsGame("items_game.txt")
	fmt.Println("Items Game Loaded Successfully")

	musicKits := parsers.ParseMusicKits(itemsGame)

	// dump to json
	jsonData, err := json.MarshalIndent(musicKits, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling music kits to JSON:", err)
		return
	}
	fmt.Println(string(jsonData))

	// dump to file
	err = os.WriteFile("music_kits.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing music kits to file:", err)
		return
	}

	// keep alive
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
};