APP_NAME := server
CONFIG_FILE := configs/local.yaml
GO_CMD := go
RUN_CMD := $(GO_CMD) run cmd/$(APP_NAME)/main.go -config $(CONFIG_FILE)
PROTO_FILE := internal/api/proto/auth/v1/auth.proto

.PHONY: run
run:
	$(RUN_CMD)

.PHONY: build
build:
	$(GO_CMD) build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

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
