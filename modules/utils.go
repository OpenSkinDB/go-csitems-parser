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

// func GetFilteredKeyValues()
