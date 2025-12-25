package events

type EventName string

const (
	MultiKill       EventName = "MultiKill"
	Ace             EventName = "Ace"
	GameStart       EventName = "GameStart"
	MinionsSpawning EventName = "MinionsSpawning"
	ChampionKill    EventName = "ChampionKill"
	BaronKill       EventName = "BaronKill"
	HeraldKill      EventName = "HeraldKill"
	DragonKill      EventName = "DragonKill"
	InhibKilled     EventName = "InhibKilled"
	TurretKilled    EventName = "TurretKilled"
	FirstBrick      EventName = "FirstBrick"
	FirstBlood      EventName = "FirstBlood"
)
