package openskindb_parsers

import (
	"fmt"
	"os"

	"github.com/baldurstod/vdf"
	models "github.com/openskindb/openskindb-csitems/models"
)

// merge keys at the root level into its own vdf.KeyValue
func MergeKeysAtRootLevel(root *vdf.KeyValue) {
	// key_count := len(kv.GetChilds())
	if root == nil {
		panic("root KeyValue is nil")
	}

	if len(root.GetChilds()) == 0 {
		panic("root KeyValue has no child keys")
	}

	cached := make(map[string][]*vdf.KeyValue)

	for _, item := range root.GetChilds() {
		sectionName := item.Key

		if cached[sectionName] == nil {
			cached[sectionName] = []*vdf.KeyValue{}
		}

		// append the children of the item to the new root value
		cached[sectionName] = append(cached[sectionName], item.GetChilds()...)
	}

	// create a new KeyValue for each section
	newRoot := &vdf.KeyValue{
		Key:    "items_game",
		Value:  make([]*vdf.KeyValue, 0),
	}

	for sectionName, items := range cached {
		if len(items) == 0 {
			continue
		}
		
		// create a new KeyValue with the section name
		newKV := &vdf.KeyValue{
			Key:   sectionName,
			Value: items,
		}
		// add the new KeyValue to the root
		newRoot.Value = append(newRoot.Value.([]*vdf.KeyValue), newKV)
	}

	// replace the root value with the new root value
	root.Value = newRoot.Value

	// root.Value = root
}

func LoadItemsGame(path string) *models.ItemsGame {
	fileData, err := os.ReadFile(path)
	
	if err != nil {
		panic(err)
	}

	vdf := vdf.VDF{}
	parsed := vdf.Parse(fileData)

	kv, _ := parsed.Get("items_game")

	if kv == nil {
		panic("items_game.txt does not contain 'items_game' section")
	}
	
	// kv.RemoveDuplicates()
	MergeKeysAtRootLevel(kv)

	json, _ := kv.MarshalJSON()

	err = os.WriteFile("exported/items_game.json", json, 0644)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
	}

	return &models.ItemsGame{
		KeyValue: kv,
	}
}