.PHONY: all build run swag test clean

all: build

build:
	go build -o bin/app ./cmd/distributed-task-scheduler

run:
	go run ./cmd/distributed-task-scheduler/main.go

swag:
	swag init -g cmd/distributed-task-scheduler/main.go -o cmd/distributed-task-scheduler/docs

test:
	go test ./... -v

docker-up:
	docker-compose up --build

clean:
	go clean
	rm -rf bin


