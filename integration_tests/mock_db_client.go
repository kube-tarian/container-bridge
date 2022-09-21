package integrationtests

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/kube-tarian/container-bridge/client/pkg/clickhouse"
	"github.com/kube-tarian/container-bridge/client/pkg/config"
)

type MockDBClient struct {
	events []map[string]interface{}
	mutex  sync.Mutex
}

func NewMockDBClient(cfg *config.Config) (clickhouse.DBInterface, error) {
	return &MockDBClient{
		events: []map[string]interface{}{},
		mutex:  sync.Mutex{},
	}, nil
}

func (m *MockDBClient) InsertEvent(event string) {
	v := map[string]interface{}{}
	json.Unmarshal([]byte(event), &v)
	m.events = append(m.events, v)
}

func (m *MockDBClient) FetchEvents() []map[string]interface{} {
	return m.events
}

func (m *MockDBClient) Close() {}

func CheckExists(event string, events []map[string]interface{}) bool {
	for _, v := range events {
		actualEvent, _ := json.Marshal(v)
		log.Printf("%v", string(actualEvent))
		if string(actualEvent) == event {
			return true
		}
	}
	return false
}
