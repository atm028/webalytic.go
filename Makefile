# note: call scripts from /scripts
collector:
	go build -o ./build/collector ./cmd/collector/server.go

docker_collector:
	docker build -f ./cmd/collector/Dockerfile -t collector .


docker_handler:
	docker build -f ./cmd/handler/Dockerfile -t handler .

handler:
	go build -o ./build/handler ./cmd/handler/server.go

docker: docker_collector docker_handler

all: collector handler

.PHONY: unit-test
unit-test:
	@docker build . --target unit-test


