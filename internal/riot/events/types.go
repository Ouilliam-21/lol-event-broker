package events

import "encoding/json"

type IBaseEvent interface {
	GetEventID() int64
	GetEventName() EventName
	GetEventTime() float64
	GetInvolvedPlayers() []string
	GetInvolvedPlayer() string
	ToJson() (json.RawMessage, error)
}

type BaseEvent struct {
	EventID   int64     `json:"EventID"`
	EventName EventName `json:"EventName"`
	EventTime float64   `json:"EventTime"`
}

func (b *BaseEvent) GetEventID() int64 {
	return b.EventID
}

func (b *BaseEvent) GetEventName() EventName {
	return b.EventName
}

func (b *BaseEvent) GetEventTime() float64 {
	return b.EventTime
}

func (b *BaseEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(b)
}

func (b *BaseEvent) GetInvolvedPlayers() []string {
	return []string{}
}

func (b *BaseEvent) GetInvolvedPlayer() string {
	return ""
}

type GameStartEvent struct {
	BaseEvent
}

type MinionsSpawningEvent struct {
	BaseEvent
}

type HeraldKillEvent struct {
	BaseEvent
	Stolen     bool     `json:"Stolen"`
	KillerName string   `json:"KillerName"`
	Assisters  []string `json:"Assisters"`
}

func (h *HeraldKillEvent) GetInvolvedPlayers() []string {
	players := []string{h.KillerName}
	players = append(players, h.Assisters...)
	return players
}

func (h *HeraldKillEvent) GetInvolvedPlayer() string {
	return h.KillerName
}

func (h *HeraldKillEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(h)
}

type ChampionKillEvent struct {
	BaseEvent
	KillerName string   `json:"KillerName"`
	VictimName string   `json:"VictimName"`
	Assisters  []string `json:"Assisters"`
}

func (c *ChampionKillEvent) GetInvolvedPlayers() []string {
	players := []string{c.KillerName, c.VictimName}
	players = append(players, c.Assisters...)
	return players
}

func (c *ChampionKillEvent) GetInvolvedPlayer() string {
	return c.VictimName
}

func (c *ChampionKillEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(c)
}

type BaronKillEvent struct {
	BaseEvent
	KillerName string   `json:"KillerName"`
	Assisters  []string `json:"Assisters"`
	Stolen     bool     `json:"Stolen"`
}

func (b *BaronKillEvent) GetInvolvedPlayers() []string {
	players := []string{b.KillerName}
	players = append(players, b.Assisters...)
	return players
}

func (b *BaronKillEvent) GetInvolvedPlayer() string {
	return b.KillerName
}

func (b *BaronKillEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(b)
}

type DragonKillEvent struct {
	BaseEvent
	DragonType string   `json:"DragonType"`
	KillerName string   `json:"KillerName"`
	Assisters  []string `json:"Assisters"`
	Stolen     bool     `json:"Stolen"`
}

func (d *DragonKillEvent) GetInvolvedPlayer() string {
	return d.KillerName
}

func (d *DragonKillEvent) GetInvolvedPlayers() []string {
	players := []string{d.KillerName}
	players = append(players, d.Assisters...)
	return players
}

func (d *DragonKillEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(d)
}

type TurretKilledEvent struct {
	BaseEvent
	TurretKilled string   `json:"TurretKilled"`
	KillerName   string   `json:"KillerName"`
	Assisters    []string `json:"Assisters"`
}

func (t *TurretKilledEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(t)
}

func (t *TurretKilledEvent) GetInvolvedPlayer() string {
	return t.KillerName
}

func (t *TurretKilledEvent) GetInvolvedPlayers() []string {
	players := []string{t.KillerName}
	players = append(players, t.Assisters...)
	return players
}

type InhibKilledEvent struct {
	BaseEvent
	InhibKilled string   `json:"InhibKilled"`
	KillerName  string   `json:"KillerName"`
	Assisters   []string `json:"Assisters"`
}

func (i *InhibKilledEvent) GetInvolvedPlayer() string {
	return i.KillerName
}

func (i *InhibKilledEvent) GetInvolvedPlayers() []string {
	players := []string{i.KillerName}
	players = append(players, i.Assisters...)
	return players
}

func (i *InhibKilledEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(i)
}

type MultiKillEvent struct {
	BaseEvent
	KillerName string `json:"KillerName"`
	KillStreak int    `json:"KillStreak"`
}

func (m *MultiKillEvent) GetInvolvedPlayer() string {
	return m.KillerName
}

func (m *MultiKillEvent) GetInvolvedPlayers() []string {
	return []string{m.KillerName}
}

func (m *MultiKillEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(m)
}

type AceEvent struct {
	BaseEvent
	Acer      string `json:"Acer"`
	AcingTeam string `json:"AcingTeam"`
}

func (a *AceEvent) GetInvolvedPlayer() string {
	return a.Acer
}

func (a *AceEvent) GetInvolvedPlayers() []string {
	return []string{a.Acer}
}

func (a *AceEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(a)
}

type FirstBrickEvent struct {
	BaseEvent
	KillerName string `json:"KillerName"`
}

func (f *FirstBrickEvent) GetInvolvedPlayer() string {
	return f.KillerName
}

func (f *FirstBrickEvent) GetInvolvedPlayers() []string {
	return []string{f.KillerName}
}

func (f *FirstBrickEvent) ToJson() (json.RawMessage, error) {
	return json.Marshal(f)
}
