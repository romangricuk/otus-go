DOCKER_IMG="calendar:develop"
DOCKER_COMPOSE := docker compose --env-file ./deployments/.env -f ./deployments/docker-compose.yaml

# Путь к директории с исходниками
SRC_DIR := ./cmd

# Путь к выходной директории для бинарников
BIN_DIR := ./bin

# Бинарные файлы
BIN_CALENDAR := $(BIN_DIR)/calendar
BIN_SCHEDULER := $(BIN_DIR)/calendar_scheduler
BIN_SENDER := $(BIN_DIR)/calendar_sender

# Команды для сборки каждого бинарника
$(BIN_CALENDAR): $(SRC_DIR)/calendar/main.go
	@echo "Building calendar..."
	@go build -o $(BIN_CALENDAR) $(SRC_DIR)/calendar/main.go

$(BIN_SCHEDULER): $(SRC_DIR)/calendar_scheduler/main.go
	@echo "Building calendar_scheduler..."
	@go build -o $(BIN_SCHEDULER) $(SRC_DIR)/calendar_scheduler/main.go

$(BIN_SENDER): $(SRC_DIR)/calendar_sender/main.go
	@echo "Building calendar_sender..."
	@go build -o $(BIN_SENDER) $(SRC_DIR)/calendar_sender/main.go

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# Обобщенная цель для сборки всех бинарников
build: $(BIN_CALENDAR) $(BIN_SCHEDULER) $(BIN_SENDER)
	@echo "All binaries built successfully!"

run: build
	$(BIN) -config ./configs/config.yaml

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

# запуск приложения в docker
up:
	$(DOCKER_COMPOSE) up -d

# остановка запущенного приложения. Так же удаляет все артефакты: containers, networks, volumes, images
down:
	$(DOCKER_COMPOSE) down -v --rmi all
	$(DOCKER_COMPOSE) --profile tests down -v --rmi all


test-integration:
	$(DOCKER_COMPOSE) --profile tests up --abort-on-container-exit --exit-code-from integration_tests integration_tests
	#$(DOCKER_COMPOSE) --profile tests down
