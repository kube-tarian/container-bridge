package integrationtests

import "testing"

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tearDown()
}

func TestEventAtAgent(t *testing.T) {
	// Start nats consumer
	// Post docker event to agent
	// Verify the event in nats consumer
}

func TestEventAtClient(t *testing.T) {
	// Post docker event to agent
	// Verify the event persisted in clickhouse database
}
