.PHONY: build clean run

build:
	mkdir -p bin
	go build -o bin/server .

clean:
	rm -rf bin

run: build
	./bin/server
