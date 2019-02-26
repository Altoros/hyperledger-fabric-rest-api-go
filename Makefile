.PHONY: all dev basic_clear build basic_up basic_down run test basic_restart byfn_up byfn_test byfn_clear byfn clear

help:
	@echo "Fabric REST API GoLang"
	@echo ""
	@echo "build: build project"
	@echo "run: build and run project"
	@echo "test: run unit tests"
	@echo ""
	@echo "clear: stop and clear all test networks"
	@echo ""
	@echo "basic_up: start basic test network"
	@echo "basic_down: stop basic test network"
	@echo "basic_clear: clear basic test network"
	@echo "basic_restart: stop, clear and start over basic test network"
	@echo ""
	@echo "byfn: start, run API tests and clear BYFN network"
	@echo "byfn_up: start BYFN network"
	@echo "byfn_test: run API tests on BYFN network"
	@echo "byfn_down: stop and clear BYFN network"
	@echo ""

all: basic_clear build basic_up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@go get -d ./cmd/...
	@go build -o=./build/frag ./cmd/frag/main.go
	@echo "Build done"

##### RUN
run: build
	@echo "Start app ..."
	@./build/frag

##### UNTI TEST
test:
	@go test . ./api... ./notifications...

##### CLEAR ALL
clear: basic_clear byfn_clear

##### BASIC NETWORK testing
basic_up:
	@echo "Starting basic network ..."
	@cd test/basic/fixtures && docker-compose up --force-recreate -d
	@echo "Basic network up"

basic_down:
	@echo "Stoping basic network ..."
	@cd test/basic/fixtures && docker-compose down
	@echo "Basic network down"

basic_restart: basic_clear basic_up

basic_clear: basic_down
	@echo "Clean up basic network ..."
	@docker rm -f -v `docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"

##### BYFN testing
byfn_up:
	@./scripts/byfn-up.sh

byfn_test: build
	@./scripts/byfn-test.sh

byfn_clear:
	@./scripts/byfn-clear.sh

byfn: clear byfn_up byfn_test byfn_clear
