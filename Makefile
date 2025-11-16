container_runtime := $(shell which podman || which docker)

$(info using ${container_runtime})

up:
	${container_runtime} compose --env-file ./env/.env.example up --build -d 

up-dev:
	${container_runtime} compose --profile dev up --build -d

down:
	${container_runtime} compose down

clean:
	${container_runtime} compose down -v

lint:
	golangci-lint run -v ./...

tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.4.0

build:
	go build -o bin/reviewer ./cmd/reviewer/

run:
	go run ./cmd/reviewer/

deps:
	go mod download
	go mod verify

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet: 
	go vet ./...

.PHONY: up up-dev down clean lint tools build run deps tidy fmt vet