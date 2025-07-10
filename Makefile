GO_CMD := go
CONFIG_FILE := configs/local.yaml
PROTO_FILE := internal/api/proto/auth/v1/auth.proto

SERVER_APP_NAME := server
CLIENT_APP_NAME := client

SERVER_RUN_CMD := $(GO_CMD) run cmd/$(SERVER_APP_NAME)/main.go -config $(CONFIG_FILE)
CLIENT_RUN_CMD := $(GO_CMD) run cmd/$(CLIENT_APP_NAME)/main.go -config $(CONFIG_FILE)

.PHONY: runserver
runserver:
	$(SERVER_RUN_CMD)

.PHONY: server
server: runserver

.PHONY: runclient
runclient:
	$(CLIENT_RUN_CMD)

.PHONY: client
client: runclient


.PHONY: build
build:
	$(GO_CMD) build -o bin/$(SERVER_APP_NAME) cmd/$(SERVER_APP_NAME)/main.go

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: test
test:
	$(GO_CMD) test ./...

.PHONY: generate
generate:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_FILE)

.PHONY: all
all: clean build
