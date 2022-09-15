CLIENT_APP_NAME := container-bridge-client
AGENT_APP_NAME := container-bridge-agent
BUILD := 0.1.1

OPEN_API_CODEGEN := github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

${OPEN_API_CODEGEN}:
	$(eval TOOL=$(@:%=%))
	@echo Installing ${TOOL}...
	go install $(@:%=%)

tools: ${OPEN_API_CODEGEN}

OPEN_API_DIR = ./api

oapi-gen: tools oapi-gen-agent oapi-gen-client

oapi-gen-agent:
	$(eval APP_NAME=agent)
	@echo Generating server for ${APP_NAME}
	@mkdir -p ${APP_NAME}/${OPEN_API_DIR}
	${GOBIN}/oapi-codegen -config ./${APP_NAME}/cfg.yaml ./${APP_NAME}/openapi.yaml

oapi-gen-client:
	$(eval APP_NAME=client)
	@echo Generating server for ${APP_NAME}
	@mkdir -p ${APP_NAME}/${OPEN_API_DIR}
	${GOBIN}/oapi-codegen -config ./${APP_NAME}/cfg.yaml ./${APP_NAME}/openapi.yaml

start-docker-compose:
	docker compose up -d --no-recreate

stop-docker-compose:
	docker compose down -v

build:
	go mod vendor
	go build -o build/client client/main.go
	go build -o build/agent agent/main.go

clean:
	rm -rf build
docker-build:
	docker build -f Dockerfile.client -t ${CLIENT_APP_NAME}:${BUILD} .
	docker build -f Dockerfile.agent -t ${AGENT_APP_NAME}:${BUILD} .
