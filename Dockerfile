FROM --platform=${BUILDPLATFORM} golang:alpine3.14 as base
RUN mkdir -p /go/github.com/webalytic.go
RUN mkdir -p /mnt/log
RUN touch /mnt/log/webalytic_test_buildkitsandbox.log
COPY ./ /go/github.com/webalytic.go/

ENV CGO_ENABLED=0
WORKDIR /go/github.com/webalytic.go
COPY go.* .
RUN go mod download
COPY . .

#FROM base as build
#ARG GOS
#ARG GOARCH
#RUN --mount=type=cache,target=/root/.cache/go-build GOS=${GOS} GOARCH=${GOARCH} go build -o ./ ./app/main.go

FROM base as unit-test

ENV REDIS_HOST="192.168.2.26"
ENV REDIS_STREAM=collector-stream
ENV REDIS_HANDLER=log-handlers
ENV REDIS_GROUP=log-handlers
ENV CLICKHOUSE_HOST="192.168.2.26"
ENV CLICKHOUSE_HTTP_PORT=8123
ENV CLICKHOUSE_SERVICE_PORT=9000
ENV CLICKHOUSE_DBNAME=webalytic
ENV CLICKHOUSE_FLUSH_INTERVAL=1000
ENV CLICKHOUSE_FLUSH_LIMIT=500
ENV LOG_PATH=/mnt/log
ENV LOG_LEVEL=DEBUG
ENV CONSUL_ENABLED="true"
ENV CONSUL_ADDRESS="192.168.2.26:8500"

#RUN --mount=type=cache,target=/root/.cache/go-build go test -v --run ^TestHandlePaymentRequest$ ./cmd/collector/...
RUN --mount=type=cache,target=/root/.cache/go-build go test -v ./...
