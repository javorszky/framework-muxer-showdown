SHELL=/bin/bash
DOCKER_BUILDKIT=1

.PHONY: test lint lintfix mocks

test:
	go test ./...

build:
	docker build -f ops/Dockerfile -t suborbital/muxer-util:latest .

lint:
	docker compose up linter

lintfix:
	docker compose up lintfixer

mocks:
	docker compose up mocks
