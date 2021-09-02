.PHONY: build

build:
	go build -o  photouploader ./cmd


.PHONY: run

run:
	go run ./cmd/main.go

.DEFAULT_GOAL := run