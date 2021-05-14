FROM golang:1.16.4

RUN mkdir /webalytic
COPY ./ /webalytic
WORKDIR /webalytic

RUN go mod download
RUN go mod verify
RUN rm -rf ./build/*
RUN go build -o ./build/collector ./cmd/collector/server.go
RUN cp -rf ./build/collector /opt/collector
WORKDIR /opt
RUN mkdir /opt/logs
RUN rm -rf /webalytic
RUN chmod +x ./collector

ENV COLLECTOR_REDIS_PORT=6379

EXPOSE 8090

ENTRYPOINT ["./collector"]