package openskindb_models

type CollectibleType int

const (
	CollectibleTypeUnknown            CollectibleType = -1
	CollectibleTypeServiceMedal       CollectibleType = 0
	CollectibleTypeMapContributor     CollectibleType = 10
	CollectibleTypeMapPin             CollectibleType = 20
	CollectibleTypeOperation          CollectibleType = 30
	CollectibleTypePickEm             CollectibleType = 40
	CollectibleTypeOldPickEm          CollectibleType = 50
	CollectibleTypeFantasyTrophy      CollectibleType = 60
	CollectibleTypeTournamentFinalist CollectibleType = 70
	CollectibleTypePremierSeasonCoin  CollectibleType = 80
	CollectibleTypeYearsOfService     CollectibleType = 90
)