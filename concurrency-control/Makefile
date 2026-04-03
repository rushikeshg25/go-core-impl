.PHONY: build clean run test docker-up docker-down

build:
	mkdir -p bin
	go build -o bin/server .

clean:
	rm -rf bin

run: build
	./bin/server

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
