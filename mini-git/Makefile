VERSION := 1.0

.PHONY: all
all: build run clean

build:
	go build -o ./bin/mini-git cmd/main.go

clean:
	rm -f ./bin/mini-git

run:
	./bin/mini-git

version:
	@echo "mini-git v$(VERSION)"