package riot

const (
	NotStarted status = "NOT STARTED"
	Running    status = "RUNNING"
)

const (
	MultiKill       event = "MultiKill"
	Ace             event = "Ace"
	GameStart       event = "GameStart"
	MinionsSpawning event = "MinionsSpawning"
	ChampionKill    event = "ChampionKill"
	BaronKill       event = "BaronKill"
	HeraldKill      event = "HeraldKill"
	DragonKill      event = "DragonKill"
	InhibKilled     event = "InhibKilled"
	TurretKilled    event = "TurretKilled"
	FirstBrick      event = "FirstBrick"
)

var EventsWatch = map[event]struct{}{
	MultiKill:    {},
	Ace:          {},
	ChampionKill: {},
	BaronKill:    {},
	HeraldKill:   {},
	DragonKill:   {},
	InhibKilled:  {},
	TurretKilled: {},
	FirstBrick:   {},
}
