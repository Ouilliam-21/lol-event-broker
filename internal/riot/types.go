package riot

type status string

type event string

type Events struct {
	Events []struct {
		ID   int64 `json:"EventID"`
		Name event `json:"EventName"`
	} `json:"Events"`
}

type Player struct {
	ChampionName string `json:"championName"`
	RiotId       string `json:"riotId"`
	RiotGameName string `json:"riotIdGameName"`
	Team         string `json:"team"`
	Position     string `json:"position"`
}
