package openskindb_modules

import (
	"strings"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog/log"
)

// GetSubKey recursively retrieves a deep subkey from a VDF KeyValue structure
func GetSubKey(vdfRoot *vdf.KeyValue, path string) *vdf.KeyValue {
	if vdfRoot == nil || path == "" {
		return nil
	}

	parts := strings.Split(path, ".")
	return getSubKeyRecursive(vdfRoot, parts)
}

// getSubKeyRecursive is a helper that walks through the KeyValue tree
func getSubKeyRecursive(current *vdf.KeyValue, pathParts []string) *vdf.KeyValue {
	if current == nil || len(pathParts) == 0 {
		return current
	}

	for _, child := range current.GetChilds() {
		if child.Key == pathParts[0] {
			return getSubKeyRecursive(child, pathParts[1:])
		}
	}

	// Key not found
	return nil
}

func GetKeyValueSubKeyValueRecursive(kv *vdf.KeyValue, key string) (string, error) {
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