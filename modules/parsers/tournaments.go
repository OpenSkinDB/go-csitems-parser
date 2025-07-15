package parsers

import (
	"context"
	"regexp"
	"time"

	"go-csitems-parser/modules"

	"github.com/rs/zerolog"
)

type TournamentSchema struct {
	Locations   map[string]string `json:"locations"`
	Tournaments map[string]string `json:"tournaments"`
}

func ParseTournaments(ctx context.Context, t *modules.Translator) *TournamentSchema {
	logger := zerolog.Ctx(ctx)

	start := time.Now()

	// CSGO_Tournament_Event_Location_(\d+)$
	var locations = make(map[string]string, 0)

	// CSGO_Tournament_Event_Location_(\d+)$
	regex := `^csgo_tournament_event_location_(\d+)$`
	re := regexp.MustCompile(regex)
	for key := range *t.Tokens {
		if !re.MatchString(key) {
			continue
		}
		// Extract the location ID from the key
		matches := re.FindStringSubmatch(key)
		if len(matches) < 2 {
			logger.Warn().Msgf("Key '%s' does not match expected format, skipping", key)
			continue
		}
		locationID := matches[1]
		locations[locationID] = (*t.Tokens)[key]
	}

	// CSGO_TeamID_<teamID>
	var tournaments = make(map[string]string, 0)

	// t.Tokens is a *map[string]string, get the keys
	// CSGO_TeamID_<teamID>
	regex = `^csgo_teamid_(\d+)$`

	re = regexp.MustCompile(regex)
	for key := range *t.Tokens {
		if !re.MatchString(key) {
			continue
		}

		// Extract the team ID from the key
		matches := re.FindStringSubmatch(key)
		if len(matches) < 2 {
			logger.Warn().Msgf("Key '%s' does not match expected format, skipping", key)
			continue
		}
		teamID := matches[1]

		tournaments[teamID] = (*t.Tokens)[key]
	}

	duration := time.Since(start)
	logger.Info().Msgf("Parsed '%d' tournaments and '%d' locations in %s", len(tournaments), len(locations), duration)

	return &TournamentSchema{
		Locations:   locations,
		Tournaments: tournaments,
	}
}
