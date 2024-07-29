BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

generate:
	protoc --proto_path=api --go_out=paths=source_relative:./api --go-grpc_out=paths=source_relative:./api api/*.proto

#Генерация UML диаграммы проекта
#предварительно требуется установка goplantuml
#go install github.com/jfeliu007/goplantuml/cmd/goplantuml@latest
generate-uml-diagram:
	 goplantuml -recursive  -show-connection-labels ./ > docs/diagram.puml

#генерация swagger документации
swagger:
	swag init -g ./internal/server/internalhttp/server.go -o ./api

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	GOROOT=$(shell go env GOROOT) golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
