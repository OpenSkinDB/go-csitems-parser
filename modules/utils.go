package modules

import (
	"context"
	"go-csitems-parser/models"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

func GetLogger() *zerolog.Logger {
	logger := zerolog.Ctx(context.Background())

	return logger
}

func GetStringMapKeySlice(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetStringMapValueSlice(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetTournamentEventId(item *vdf.KeyValue) (int, error) {
	attributes, err := item.Get("attributes")
	if err != nil {
		return -1, err
	}

	tournament, err := attributes.Get("tournament event id")
	if err != nil {
		return -1, err
	}

	tournament_event_id, err := tournament.GetInt("value")
	if err != nil {
		return -1, err
	}

	return tournament_event_id, nil
}

func GetContainerItemSet(item *vdf.KeyValue) *models.WeaponCaseItemSet {
	tags, err := item.Get("tags")

	if err != nil {
		return nil
	}

	item_set, err := tags.Get("ItemSet")
	if err != nil {
		return nil
	}

	tag, _ := item_set.GetString("tag_value")
	tagText, _ := item_set.GetString("tag_text")

	return &models.WeaponCaseItemSet{
		Id:   tag,
		Name: tagText,
	}
}

type ItemWear struct {
	Name    string  `json:"name"`
	MinWear float64 `json:"min_wear"`
	MaxWear float64 `json:"max_wear"`
}

var ItemWears = map[string]ItemWear{
	"Factory New": {
		Name:    "Factory New",
		MinWear: 0.00,
		MaxWear: 0.07,
	},
	"Minimal Wear": {
		Name:    "Minimal Wear",
		MinWear: 0.07,
		MaxWear: 0.15,
	},
	"Field-Tested": {
		Name:    "Field-Tested",
		MinWear: 0.15,
		MaxWear: 0.38,
	},
	"Well-Worn": {
		Name:    "Well-Worn",
		MinWear: 0.38,
		MaxWear: 0.45,
	},
	"Battle-Scarred": {
		Name:    "Battle-Scarred",
		MinWear: 0.45,
		MaxWear: 1.00,
	},
}

func GenerateMarketHashName(t *Translator, item_name string, item_type string) (string, error) {

	switch item_type {
	case "agent":
		value, err := t.GetValueByKey(item_name)

		if err != nil {
			return "", err
		}

		return value, nil
	}

	return "", nil
}

// func GetFilteredKeyValues()
