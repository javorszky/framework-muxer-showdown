SHELL=/bin/bash
DOCKER_BUILDKIT=1

.PHONY: test build lint lintfix mocks

test:
	go test ./...

build:
	docker build -f ops/Dockerfile -t suborbital/muxer-util:latest .

lint: build
	docker compose up linter

lintfix: build
	docker compose up lintfixer

mocks: build
	docker compose up mocks
