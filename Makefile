APP_NAME=backend-challenge-092025
BINARY=server

.PHONY: build run test clean docker-build docker-up docker-down

build:
	go build -o $(BINARY) ./cmd/server/main.go

run: build
	./$(BINARY)

test:
	go test ./... -v

clean:
	rm -f $(BINARY)

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down
