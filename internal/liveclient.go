package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"crypto/tls"
)

type Status string

const (
	NotStarted Status = "NOT STARTED"
	Running    Status = "RUNNING"
)

type EventName string

const (
	MultiKill EventName = "MultiKill"
	Ace       EventName = "Ace"
	GameStart EventName = "GameStart"
	MinionsSpawning EventName = "MinionsSpawning"
)

type Test struct {
	Name EventName `json:"EventName"`
}

type Events struct {
	Events []Test `json:"Events"`
}

type LiveClient struct {
	endpoint   string
	gameStatus Status
	queue      chan<- []byte
}

func NewLiveClient(endpoint string, queue chan<- []byte) *LiveClient {
	return &LiveClient{
		endpoint:   fmt.Sprintf("%s/liveclientdata/eventdata", endpoint),
		gameStatus: NotStarted,
		queue:      queue,
	}
}

func httpGetRequestWithTimeout(url string, timeout int64) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: time.Second * time.Duration(timeout), Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request at %s: %w", url, err)
	}
	defer resp.Body.Close()


	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("url returned status %d: %s", resp.StatusCode, string(raw))
	}

	return raw, nil
}

func (lc LiveClient) poolGameEvent() {

	log.Println("Game started")

	for lc.gameStatus == Running {
		time.Sleep(time.Second * 5)

		//Add last event id to avoid fetch wholes event from beginng
		raw, err := httpGetRequestWithTimeout(lc.endpoint, 5)

		log.Println("raw: ", string(raw))

		if err != nil {
			log.Println("err: ", err)
			continue
		}

		lc.queue <- raw
	}

	lc.gameStatus = NotStarted
}

func (lc LiveClient) Process() {

	var evts Events

	for {
		time.Sleep(time.Second * 5)

		raw, err := httpGetRequestWithTimeout(lc.endpoint, 5)

		log.Println("raw: ", string(raw))

		if err != nil {
			log.Println("err: ", err)
			continue
		}

		if err := json.Unmarshal(raw, &evts); err != nil {
			log.Println("err: ", err)
			continue
		}

		https://127.0.0.1:2999/liveclientdata/eventdata?eventID=1

		for _,evt := range evts.Events{
			if evt.Name == GameStart {
				lc.gameStatus = Running
				lc.poolGameEvent()
			}

			break
		}





	}
}
