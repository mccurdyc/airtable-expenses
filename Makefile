default: build

build:
	go fmt ./...
	go build -o bin/migrate cmd/migrate/main.go

install: build
	cp bin/migrate /usr/local/bin

run: build
	./bin/migrate

test:
	go test ./...

.PHONY: build install run test
