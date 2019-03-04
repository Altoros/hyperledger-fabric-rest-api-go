.PHONY: all dev clear run test unit_test
.PHONY: build basic_up basic_down basic_restart basic_clear basic_newman_test basic_docker_up basic_docker_down basic_e2e_test
.PHONY: byfn_up byfn_test byfn_clear byfn
.PHONY: docker_build docker_up docker_restart

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

##### ALL TESTS
test: unit_test basic_e2e_test byfn_e2e_test

##### UNIT TESTS
unit_test:
	@go test ./pkg... ./cmd...

##### CLEAR ALL
clear: byfn_clear basic_clear basic_docker_down

##### BASIC NETWORK testing
basic_up:
	@echo "Starting basic network ..."
	@./scripts/basic-up.sh
	@echo "Basic network up"

basic_down:
	@echo "Stoping basic network ..."
	@./scripts/basic-down.sh
	@echo "Basic network down"

basic_clear: basic_down
	@echo "Clean up basic network ..."
	@./scripts/basic-clear.sh
	@echo "Clean up done"

basic_restart: basic_clear basic_up

basic_docker_up:
	@./scripts/basic-docker-up.sh

basic_docker_down:
	@docker stop frag
	@docker rm frag

basic_newman_test:
	@./scripts/basic-newman-test.sh

basic_e2e_test: docker_build basic_up basic_docker_up basic_newman_test basic_docker_down basic_clear

##### BYFN testing
byfn_up:
	@./scripts/byfn-up.sh

byfn_test: build
	@./scripts/byfn-test.sh

byfn_clear:
	@./scripts/byfn-clear.sh

byfn_docker_up:
	@./scripts/byfn-docker-up.sh

byfn_docker_down:
	@docker stop frag
	@docker rm frag

byfn_newman_test:
	@./scripts/byfn-newman-test.sh

byfn: byfn_up byfn_test byfn_clear

byfn_e2e_test: docker_build byfn_up byfn_docker_up byfn_newman_test byfn_docker_down byfn_clear

##### Docker
docker_build:
	@docker build -t frag:dev .

docker_restart: docker_clear docker_up

