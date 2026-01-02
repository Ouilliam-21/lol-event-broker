package events

import (
	"encoding/json"
	"fmt"
)

type EventFactory struct{}

func NewEventFactory() *EventFactory {
	return &EventFactory{}
}

func (f *EventFactory) CreateEvent(rawEvent json.RawMessage) (IBaseEvent, error) {
	var base BaseEvent
	if err := json.Unmarshal(rawEvent, &base); err != nil {
		return nil, fmt.Errorf("failed to parse base event: %w", err)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal(rawEvent, &dat); err != nil {
		panic(err)
	}
	fmt.Println("Data", dat)

	switch EventName(base.EventName) {
	case ChampionKill:
		var evt ChampionKillEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse ChampionKill event: %w", err)
		}
		return &evt, nil

	case BaronKill:
		var evt BaronKillEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse BaronKill event: %w", err)
		}
		return &evt, nil

	case DragonKill:
		var evt DragonKillEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse DragonKill event: %w", err)
		}
		return &evt, nil

	case TurretKilled:
		var evt TurretKilledEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse TurretKilled event: %w", err)
		}
		return &evt, nil

	case InhibKilled:
		var evt InhibKilledEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse InhibKilled event: %w", err)
		}
		return &evt, nil

	case MultiKill:
		var evt MultiKillEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse MultiKill event: %w", err)
		}
		return &evt, nil

	case Ace:
		var evt AceEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse Ace event: %w", err)
		}
		return &evt, nil

	case FirstBrick:
		var evt FirstBrickEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse FirstBrick event: %w", err)
		}
		return &evt, nil
	case HeraldKill:
		var evt HeraldKillEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse HeraldKill event: %w", err)
		}
		return &evt, nil

	case GameStart:
		var evt GameStartEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse GameStart event: %w", err)
		}
		return &evt, nil

	case MinionsSpawning:
		var evt MinionsSpawningEvent
		if err := json.Unmarshal(rawEvent, &evt); err != nil {
			return nil, fmt.Errorf("failed to parse MinionsSpawning event: %w", err)
		}
		return &evt, nil

	default:
		return &base, nil
	}
}
