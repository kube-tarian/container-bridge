CLIENT_APP_NAME := container-bridge-client
AGENT_APP_NAME := container-bridge-agent
BUILD := 0.0.1

start-docker-compose:
	docker compose up -d --no-recreate

stop-docker-compose:
	docker compose down -v

build:
	go build -o build/client client/main.go
	go build -o build/agent agent/main.go

docker-build:
	docker build -f Dockerfile.client -t ${CLIENT_APP_NAME}:${BUILD} .
	docker build -f Dockerfile.agent -t ${AGENT_APP_NAME}:${BUILD} .
