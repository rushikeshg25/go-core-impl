VERSION := 1.0

.PHONY: all
all: build run clean

build:
	go build -o mini-git main.go

clean:
	rm -f mini-git

run:
	./mini-git

version:
	@echo "mini-git v$(VERSION)"