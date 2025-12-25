package riot

import "encoding/json"

type Event struct {
	ID   int64           `json:"EventID"`
	Name event           `json:"EventName"`
	Raw  json.RawMessage `json:"-"`
}

type EventList struct {
	Items []Event `json:"Events"`
}

type RawEventList struct {
	Items []json.RawMessage `json:"Events"`
}

type EventContainer struct {
	List EventList
	Raw  RawEventList
}

func NewEventContainer(data []byte) (*EventContainer, error) {
	var list EventList
	var rawList RawEventList

	if err := json.Unmarshal(data, &rawList); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return &EventContainer{
		List: list,
		Raw:  rawList,
	}, nil
}

func (e EventContainer) GetLast() (Event, bool) {
	size := len(e.List.Items)

	if size <= 0 {
		return Event{}, false
	}

	return e.List.Items[len(e.List.Items)-1], true
}

func (e *EventContainer) FilterActiveEvents() {
	events := make([]Event, 0, len(e.List.Items))
	raw := make([]json.RawMessage, 0, len(e.List.Items))

	for i, evt := range e.List.Items {
		if _, shouldWatch := EventsWatch[evt.Name]; shouldWatch {
			events = append(events, evt)
			raw = append(raw, e.Raw.Items[i])
		}
	}

	e.List.Items = events
	e.Raw.Items = raw

}

type Player struct {
	ChampionName string `json:"championName"`
	RiotId       string `json:"riotId"`
	RiotGameName string `json:"riotIdGameName"`
	Team         string `json:"team"`
	Position     string `json:"position"`
}
