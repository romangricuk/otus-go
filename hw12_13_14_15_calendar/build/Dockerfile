# Собираем в гошке
FROM golang:1.22 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статические бинарники Go (без зависимостей на Си API),
# иначе они не будут работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o /opt/calendar/calendar-app cmd/calendar/*
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o /opt/calendar/calendar-sender cmd/calendar_sender/*
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -o /opt/calendar/calendar-scheduler cmd/calendar_scheduler/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

# Установка curl для healthcheck
RUN apk --no-cache add curl

ENV BIN_DIR "/opt/calendar/"
ENV CONFIG_FILE /etc/calendar/config.toml

# Копируем все собранные бинарники
COPY --from=build /opt/calendar/calendar-app ${BIN_DIR}
COPY --from=build /opt/calendar/calendar-sender ${BIN_DIR}
COPY --from=build /opt/calendar/calendar-scheduler ${BIN_DIR}