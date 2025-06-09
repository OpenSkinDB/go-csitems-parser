package openskindb_parsers

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/baldurstod/vdf"
	models "github.com/openskindb/openskindb-csitems/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ParseWeaponCases(ctx context.Context, ig *models.ItemsGame) []models.WeaponCase {
	logger := zerolog.Ctx(ctx);

	start := time.Now()
	logger.Info().Msg("Parsing weapon cases...")

	items, err := ig.Get("items")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to get collectibles from items_game.txt")
		return nil
	}

	var weapon_cases []models.WeaponCase

	// Iterate through all items in the "items" section
	for _, item := range items.GetChilds() {
		prefab, _ := item.GetString("prefab")
		
		if prefab != "weapon_case" {
			continue
		}
		
		definition_index, _ := strconv.Atoi(item.Key)
		item_name, _ := item.GetString("item_name")
		name, _ := item.GetString("name")
		item_description, _ := item.GetString("item_description")
		image_inventory, _ := item.GetString("image_inventory")

		if name == "" {
			log.Warn().Msg("Collectible name is empty")
			continue
		}

		// Get child key called "attributes"
		associated_items, _ := item.Get("associated_items")

		case_key_def_idx := -1 // Default to -1 if not found
		if associated_items != nil {
			value := associated_items.GetChilds()[0].Key;
			
			if value != "" {
				case_key_def_idx, _ = strconv.Atoi(value)
			}
		}
		
		// If case_key_def_idx is still -1, we cannot find the key for this case
		case_key := GetWeaponCaseKeyByDefIndex(ig, case_key_def_idx)
		item_set := GetWeaponCaseItemSet(item)

		// Create the weapon case model
		var current = models.WeaponCase{
			DefinitionIndex: 		definition_index,
			Prefab: 				 		prefab,
			Name:            		name,
			ItemName:        		item_name,
			ItemDescription: 		item_description,
			ImageInventory:  		image_inventory,
			Key: case_key,
			ItemSet: item_set,
		};

		weapon_cases = append(weapon_cases, current)
	}

	// Save music kits to the database
	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' weapon cases in %s", len(weapon_cases), duration)

	return weapon_cases
}

func GetWeaponCaseKeyByDefIndex(ig *models.ItemsGame, definitionIndex int) *models.WeaponCaseKey {
	items, err := ig.Get("items")

	if err != nil {
		log.Error().Err(err).Msg("gg Failed to get items from items_game.txt")
		return nil
	}

	var current models.WeaponCaseKey
	for _, item := range items.GetChilds() {
		def_idx, _ := strconv.Atoi(item.Key)
		if def_idx != definitionIndex {
			continue
		}

		prefab, _ := item.GetString("prefab")

		if !strings.Contains(prefab, "weapon_case_key") {
			continue
		}

		name, _ := item.GetString("name")
		item_name, _ := item.GetString("item_name")
		item_description, _ := item.GetString("item_description")
		first_sale_date, _ := item.GetString("first_sale_date")
		image_inventory, _ := item.GetString("image_inventory")

		current = models.WeaponCaseKey{
			DefinitionIndex:  def_idx,
			Prefab:            prefab,
			Name:              name,
			ItemName:          item_name,
			ItemDescription:   item_description,
			FirstSaleDate:     first_sale_date,
			ImageInventory:    image_inventory,
		}

		break // We found the item, no need to continue
	}

	if current.Prefab == "" {
		log.Error().Msgf("No weapon case key found for definition index %d", definitionIndex)
		return nil
	}

	return &current
}

func GetWeaponCaseItemSet(item *vdf.KeyValue) *models.WeaponCaseItemSet {
	tags, err := item.Get("tags")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get tags for weapon case item")
		return nil
	}

	item_set, err := tags.Get("ItemSet")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get ItemSet for weapon case item")
		return nil
	}

	tag, _ := item_set.GetString("tag_value")
	tagText, _ := item_set.GetString("tag_text")
	tagGroup, _ := item_set.GetString("tag_group")
	tagGroupText, _ := item_set.GetString("tag_group_text")

	return &models.WeaponCaseItemSet{
		Tag:           tag,
		TagText:       tagText,
		TagGroup:      tagGroup,
		TagGroupText:  tagGroupText,
	}
}