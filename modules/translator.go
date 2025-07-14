package modules

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/baldurstod/vdf"
	"github.com/rs/zerolog"
)

type Translator struct {
	Language string
	Tokens   *map[string]string
}

type TranslatorFactory struct {
	Translators map[string]*Translator
}

func (tf *TranslatorFactory) GetTranslator(language string) *Translator {
	if t, ok := tf.Translators[language]; ok {
		return t
	}
	return nil
}

func LoadAllTranslations(ctx context.Context, folderPath string) *TranslatorFactory {
	logger := zerolog.Ctx(ctx)

	files, err := os.ReadDir(folderPath)
	logger.Info().Msgf("Found '%d' language files", len(files))

	if err != nil {
		logger.Error().Msgf("Error reading directory %s: %v", folderPath, err)
		return nil
	}

	start := time.Now()

	// Create a map to store the loaded translations, so we can
	// access them by language code later.
	lang_map := make(map[string]*Translator)
	vdf := vdf.VDF{}

	// Loop through all files in the folder
	for _, dir_entry := range files {
		logger.Debug().Msgf("Processing file %s", dir_entry.Name())

		if dir_entry.IsDir() {
			logger.Info().Msgf("Skipping directory %s", dir_entry.Name())
			continue
		}

		file_name := dir_entry.Name()
		if !strings.HasPrefix(file_name, "csgo_") || !strings.HasSuffix(file_name, ".txt") {
			logger.Info().Msgf("Skipping file %s", file_name)
			continue
		}

		// now we only have files, not folders, so we can load them
		path := fmt.Sprintf("%s/%s", folderPath, file_name)
		fileData, err := os.ReadFile(path)

		data := RemoveBOMFromFile(fileData)

		if err != nil {
			logger.Error().Msgf("Error reading file %s: %v", path, err)
			continue
		}

		keyvalues := vdf.Parse(data)

		if keyvalues.Value == nil {
			logger.Error().Msgf("Error parsing file %s", path)
			continue
		}

		t, lang_name := LoadLanguage(&keyvalues)

		if t == nil {
			logger.Error().Msgf("Error loading language from file %s", path)
			continue
		}

		lang_map[lang_name] = t
		logger.Info().Msgf("Loaded '%d' tokens for language '%s'", len(*t.Tokens), lang_name)
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' language files in %s", len(files), duration)

	return &TranslatorFactory{
		Translators: lang_map,
	}
}

func LoadLanguage(vdf *vdf.KeyValue) (*Translator, string) {
	if vdf == nil {
		panic("vdf is nil")
	}

	kv, _ := vdf.Get("lang")
	lang_name, _ := kv.GetString("Language")

	tokens, _ := kv.Get("Tokens")

	if kv == nil || tokens == nil {
		panic("translation file does not contain 'lang.Tokens' section")
	}

	token_map, err := tokens.ToStringMap()

	if err != nil {
		panic(fmt.Sprintf("Error parsing tokens: %v", err))
	}

	translator := &Translator{
		Language: lang_name,
		Tokens:   token_map,
	}

	return translator, lang_name
}

func (t *Translator) GetValueByKey(key string) (string, error) {
	token_key := strings.Replace(key, "#", "", -1)
	if t == nil {
		fmt.Println("Translator is nil")
		return "", errors.New("key not found")
	}

	return (*t.Tokens)[token_key], nil
}
