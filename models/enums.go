package openskindb_models

type CollectibleType int

const (
	CollectibleTypeUnknown            CollectibleType = -1
	CollectibleTypeServiceMedal       CollectibleType = 0
	CollectibleTypeMapContributor     CollectibleType = 1
	CollectibleTypeMapPin             CollectibleType = 2
	CollectibleTypeOperation          CollectibleType = 3
	CollectibleTypePickEm             CollectibleType = 4
	CollectibleTypeOldPickEm          CollectibleType = 5
	CollectibleTypeFantasyTrophy      CollectibleType = 6
	CollectibleTypeTournamentFinalist CollectibleType = 7
)