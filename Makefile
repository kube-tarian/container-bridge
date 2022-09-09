CLIENT_APP_NAME := container-bridge-client
AGENT_APP_NAME := container-bridge-agent
BUILD := 0.0.1

TOOLS_DIR		:= .tools/
GOLANGCI_LINT	:= ${TOOLS_DIR}github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.0
OAPI_CODEGEN	:= ${TOOLS_DIR}github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
GOIMPORTS		:= ${TOOLS_DIR}mvdan.cc/gofumpt/gofumports@v0.1.1

${OAPI_CODEGEN} ${GOLANGCI_LINT} ${GOIMPORTS}:
	$(eval TOOL=$(@:%=%))
	@echo Installing ${TOOL}...
	go install $(TOOL:${TOOLS_DIR}%=%)
	@mkdir -p $(dir ${TOOL})
	@cp ${GOBIN}/$(firstword $(subst @, ,$(notdir ${TOOL}))) ${TOOL}

tools: ${OAPI_CODEGEN} ${GOLANGCI_LINT} ${GOIMPORTS}

OAPI_DIR = ./api

oapi-gen: tools oapi-gen-agent oapi-gen-client

oapi-gen-agent:
	$(eval APP_NAME=agent)
	@echo Generating server for ${APP_NAME}
	@mkdir -p ${APP_NAME}/${OAPI_DIR}
	${OAPI_CODEGEN} -config ./${APP_NAME}/cfg.yaml ./${APP_NAME}/openapi.yaml

oapi-gen-client:
	$(eval APP_NAME=client)
	@echo Generating server for ${APP_NAME}
	@mkdir -p ${APP_NAME}/${OAPI_DIR}
	${OAPI_CODEGEN} -config ./${APP_NAME}/cfg.yaml ./${APP_NAME}/openapi.yaml

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
