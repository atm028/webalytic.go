FROM golang:1.16.4

RUN mkdir /webalytic
RUN mkdir /mnt/log
COPY ./ /webalytic
WORKDIR /webalytic

RUN go mod download
RUN go mod verify
RUN rm -rf ./build/*
RUN go build -o ./build/handler ./cmd/handler/server.go
RUN cp -rf ./build/handler /opt/handler
WORKDIR /opt
RUN rm -rf /webalytic
RUN chmod +x ./handler

ENV COLLECTOR_REDIS_PORT=6379

EXPOSE 8091

ENTRYPOINT ["./handler"]