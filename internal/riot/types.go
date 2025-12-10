package riot

type status string

type event string

type Event struct {
	ID   int64 `json:"EventID"`
	Name event `json:"EventName"`
}

type Events struct {
	Events []Event `json:"Events"`
}

func (e Events) GetLast() Event {
	return e.Events[len(e.Events)-1]
}

func (e Events) FilterActiveEvents() []Event {
	res := make([]Event, 0)
	for _, evt := range e.Events {
		if _, shouldWatch := EventsWatch[evt.Name]; shouldWatch {
			res = append(res, evt)
		}
	}

	return res
}

type Player struct {
	ChampionName string `json:"championName"`
	RiotId       string `json:"riotId"`
	RiotGameName string `json:"riotIdGameName"`
	Team         string `json:"team"`
	Position     string `json:"position"`
}
