package integrationtests

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestEvent(t *testing.T) {
	data := setup()

	stop := startAagentAndClient()

	event := `{"events": [{"hash": "abcd"}]}`
	// Post docker event to agent
	resp, err := callHTTPRequest(http.MethodPost, "http://localhost:8090", "localregistry/event/docker", []byte(event))
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	log.Println("Sleeping now")
	time.Sleep(5 * time.Second)

	// Verify the event persisted in clickhouse database
	events := data.dbClient.FetchEvents()
	expectedEvent := `{"event":{"hash":"abcd"},"repoName":"docker registry"}`
	assert.True(t, CheckExists(expectedEvent, events))

	log.Println("Starting teardown")
	tearDown(data)
	stop <- true
}
