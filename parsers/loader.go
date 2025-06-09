package openskindb_parsers

import (
	"os"

	"github.com/baldurstod/vdf"
	models "github.com/openskindb/openskindb-csitems/models"
)


func LoadItemsGame(path string) *models.ItemsGame {
	fileData, err := os.ReadFile(path)
	
	if err != nil {
		panic(err)
	}

	vdf := vdf.VDF{}
	parsed := vdf.Parse(fileData)

	kv, _ := parsed.Get("items_game")

	if kv == nil {
		panic("items_game.txt does not contain items_game section")
	}

	return &models.ItemsGame{
		KeyValue: kv,
	}
}