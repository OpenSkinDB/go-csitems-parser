package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"go-csitems-parser/models"

	"github.com/baldurstod/vdf"
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
		Key:   "items_game",
		Value: make([]*vdf.KeyValue, 0),
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
}

func LoadKnifeSkinsMap(path string) map[string][]string {
	fileData, err := os.ReadFile(path)

	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %v", path, err))
	}

	if len(fileData) == 0 {
		panic(fmt.Sprintf("File %s is empty", path))
	}

	result := make(map[string][]string)
	json.Unmarshal(fileData, &result)

	return result
}

type ItemsGameCdn struct {
	ItemName string `json:"item_name"`
	URL      string `json:"url"`
}

func LoadItemsGameCdn(path string) map[string]string {
	fileData, err := os.ReadFile(path)

	if err != nil {
		panic(fmt.Sprintf("Error reading file %s: %v", path, err))
	}

	// Each newline looks like: weapon=url\n
	// make a regex to split this into a map[string]string

	if len(fileData) == 0 {
		panic(fmt.Sprintf("File %s is empty", path))
	}

	// reg := `/(\w+)=([^\s]+)/g`
	reg := `(\w+)=([^\s]+)`
	re := regexp.MustCompile(reg)
	matches := re.FindAllStringSubmatch(string(fileData), -1)

	if matches == nil {
		panic(fmt.Sprintf("No matches found in file %s", path))
	}

	items := make(map[string]string, 0)

	for _, match := range matches {
		if len(match) < 3 {
			continue // skip if we don't have enough matches
		}
		itemName := match[1]
		url := match[2]
		if itemName == "" || url == "" {
			continue // skip if item name or url is empty
		}
		items[itemName] = url
	}

	return items
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
