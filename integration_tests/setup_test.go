package integrationtests

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	agentapp "github.com/kube-tarian/container-bridge/agent/pkg/application"
	clientapp "github.com/kube-tarian/container-bridge/client/pkg/application"
	"github.com/kube-tarian/container-bridge/client/pkg/clickhouse"
	"github.com/kube-tarian/container-bridge/client/pkg/config"
)

type TestContextData struct {
	clientConf *config.Config
	dbClient   clickhouse.DBInterface
}

func setupENV() {
	os.Setenv("NATS_TOKEN", "UfmrJOYwYCCsgQvxvcfJ3BdI6c8WBbnD")
	os.Setenv("NATS_ADDRESS", "nats://localhost:4222")
	os.Setenv("STREAM_NAME", "CONTAINERMETRICS")
	os.Setenv("DB_ADDRESS", "localhost:9000")
}

func setup() *TestContextData {
	setupENV()
	cfg := &config.Config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Could not parse env Config: %v", err)
	}

	// Create a db client
	// dbClient, err := NewMockDBClient(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	dbClient, err := clickhouse.NewDBClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &TestContextData{
		clientConf: cfg,
		dbClient:   dbClient,
	}
}

func tearDown(t *TestContextData) {
	// Close a db client
	t.dbClient.Close()
}

func startAagentAndClient() chan bool {
	stopCh := make(chan bool)

	// Start agent and client
	go startAgent(stopCh)
	time.Sleep(2 * time.Second)
	go startClient(stopCh)

	// Wait till Agent and Client healthy
	isAgentHealthy := false
	isClientHealthy := false
	for {
		select {
		// wait till 1min, after that exit 1
		case <-time.After(1 * time.Minute):
			log.Fatalf("Agent/Client not healthy")
		case <-time.After(2 * time.Second):
			// Check Agent health
			isAgentHealthy = getHealth(http.MethodGet, "http://localhost:8090", "status", "agent")
			// Check Client health
			isClientHealthy = getHealth(http.MethodGet, "http://localhost:8091", "status", "client")
		}
		if isAgentHealthy && isClientHealthy {
			break
		}
	}
	return stopCh
}

func getHealth(method, url, path, serviceName string) bool {
	resp, err := callHTTPRequest(method, url, path, nil)
	if err != nil {
		log.Printf("%v health check call failed: %v", serviceName, err)
		return false
	}

	return checkResponse(resp, http.StatusOK)
}

func checkResponse(resp *http.Response, statusCode int) bool {
	return resp.StatusCode == statusCode
}

func startAgent(stop chan bool) {
	os.Setenv("PORT", "8090")
	app := agentapp.New()
	go app.Start()

	<-stop
}

func startClient(stop chan bool) {
	os.Setenv("PORT", "8091")
	app := clientapp.New()
	go app.Start()
	<-stop
}

func callHTTPRequest(method, url, path string, body []byte) (*http.Response, error) {
	finalURL := fmt.Sprintf("%s/%s", url, path)
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, finalURL, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, finalURL, nil)
	}
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return client.Do(req)
}
