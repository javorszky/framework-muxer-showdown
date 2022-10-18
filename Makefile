SHELL=/bin/bash
DOCKER_BUILDKIT=1

.PHONY: test build lint lintfix mocks start

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

start:
	go run main.go

hammer:
	hey -z 30s -m GET -cpus 4 -c 200 -H "Authorization: icandowhatiwant" http://localhost:9000/performance

hammer-2:
	hey -z 30s -m GET -cpus 4 -c 200 http://localhost:9000/smol-perf

yeet:
	mkdir -p perftests
	$(MAKE) hammer > perftests/run1.log
	$(MAKE) hammer-2 > perftests/smol-run1.log
	$(MAKE) hammer > perftests/run2.log
	$(MAKE) hammer-2 > perftests/smol-run2.log
	$(MAKE) hammer > perftests/run3.log
	$(MAKE) hammer-2 > perftests/smol-run3.log
	$(MAKE) hammer > perftests/run4.log
	$(MAKE) hammer-2 > perftests/smol-run4.log

reqs:
	@grep -rw "Requests/sec" perftests/*
