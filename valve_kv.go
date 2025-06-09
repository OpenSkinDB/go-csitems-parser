package main

import (
	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog/log"
)

func  GetKeyValueSubKeyValueRecursive(kv *vdf.KeyValue, key string) (string, error) {
	attributes, err := kv.Get("attributes")

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get attributes for key '%s'", key)
		return "", err
	}

	if attributes == nil {
		log.Warn().Msgf("Key '%s' has no attributes, skipping", key)
		return "", nil
	}

	for _, attribute := range attributes.GetChilds() {
		if attribute.Key == key {
			value, err := attribute.GetString("value")
			if err != nil {
				log.Error().Err(err).Msgf("Failed to get value for key '%s'", key)
				return "", err
			}
			return value, nil
		}
	}

	log.Warn().Msgf("Key '%s' not found in attributes", key)
	return "", nil
}