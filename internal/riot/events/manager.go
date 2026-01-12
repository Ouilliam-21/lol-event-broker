package events

import "encoding/json"

type RawEvents struct {
	Events []json.RawMessage `json:"Events"`
}

type EventManager struct {
	events         []IBaseEvent
	eventsToWatch  map[EventName]struct{}
	playersToWatch map[string]struct{}
	eventFactory   *EventFactory
}

func NewEventManager(watchedPlayers []string, watchedEvents []string) *EventManager {

	eventsToWatch := make(map[EventName]struct{})
	for _, event := range watchedEvents {
		eventsToWatch[EventName(event)] = struct{}{}
	}

	playersToWatch := make(map[string]struct{})
	for _, player := range watchedPlayers {
		playersToWatch[player] = struct{}{}
	}

	return &EventManager{
		events:         make([]IBaseEvent, 0),
		playersToWatch: playersToWatch,
		eventsToWatch:  eventsToWatch,
		eventFactory:   NewEventFactory(),
	}
}

func (m *EventManager) ProcessEvent(rawEvent []byte) error {
	m.ClearEvents()

	var rawEvents RawEvents

	if err := json.Unmarshal(rawEvent, &rawEvents); err != nil {
		return err
	}

	for _, eventData := range rawEvents.Events {
		event, err := m.eventFactory.CreateEvent(eventData)
		if err != nil {
			return err
		}
		m.events = append(m.events, event)
	}

	return nil

}

func (m *EventManager) ClearEvents() {
	m.events = make([]IBaseEvent, 0)
}

func (m *EventManager) GetLast() IBaseEvent {
	if len(m.events) <= 0 {
		return nil
	}

	return m.events[len(m.events)-1]
}

func (m *EventManager) FilterEvents() []IBaseEvent {
	lastEventId := m.GetLast().GetEventID() //By default remove last id
	events := make([]IBaseEvent, 0)
	for _, event := range m.events {

		if event.GetEventID() == lastEventId {
			continue
		}

		if !m.isWatchedEvent(event.GetEventName()) {
			continue
		}

		involved := false

		switch event.GetEventName() {
		case ChampionKill:
			involved = m.isWatchedPlayer(event.GetInvolvedPlayer())
		case BaronKill,
			HeraldKill,
			DragonKill,
			InhibKilled,
			TurretKilled,
			Ace:
			involved = !m.isWatchedPlayer(event.GetInvolvedPlayer())
		}

		if involved {
			events = append(events, event)
		}

	}
	return events
}

func (m *EventManager) IsEmpty() bool {
	return len(m.events) <= 0
}

func (m *EventManager) isWatchedPlayer(playerName string) bool {
	_, shouldWatch := m.playersToWatch[playerName]
	return shouldWatch
}

func (m *EventManager) isWatchedEvent(eventName EventName) bool {
	_, shouldWatch := m.eventsToWatch[eventName]
	return shouldWatch
}
