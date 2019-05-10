.PHONY: build run unit_test test docker_build clear
.PHONY: basic_up basic_down basic_clear basic_restart basic_docker_up basic_docker_down basic_newman_test basic_e2e_test
.PHONY: byfn_up byfn_test byfn_clear byfn_docker_up byfn_docker_down byfn_newman_test byfn byfn_e2e_test
.PHONY: ci_up ci_clear ca_docker_up ca_docker_down ca_test ca_e2e_test

help:
	@echo "Fabric REST API GoLang"
	@echo ""
	@echo "build: build project"
	@echo "run: build and run project"
	@echo "unit_test: run unit tests"
	@echo "test: run all tests"
	@echo ""
	@echo "docker_build: build Docker container"
	@echo ""
	@echo "clear: stop and clear all test networks"
	@echo ""
	@echo "basic_up: start basic test network"
	@echo "basic_down: stop basic test network"
	@echo "basic_clear: clear basic test network"
	@echo "basic_restart: stop, clear and start over basic test network"
	@echo "basic_docker_up: start Docker container wiht API and inject it in Basic network"
	@echo "basic_docker_down: stop and clear API Docker container"
	@echo "basic_newman_test: run API Newman tests"
	@echo "basic_e2e_test: run end-2-end test with Docker container and API Newman tests on Basic network"
	@echo ""
	@echo "byfn_up: start BYFN network"
	@echo "byfn_test: API shell tests on BYFN network"
	@echo "byfn_clear: stop and clear BYFN network"
	@echo "byfn_docker_up: start Docker container wiht API and inject it in BYFN network"
	@echo "byfn_docker_down: stop and clear API Docker container"
	@echo "byfn_newman_test: run API Newman tests"
	@echo "byfn: start, run API shell tests and clear BYFN network"
	@echo "byfn_e2e_test: run end-2-end test with Docker container and API Newman tests on BYFN network"
	@echo ""
	@echo "ca_up: start CA network"
	@echo "ca_test: API Newman tests on CA network"
	@echo "ca_clear: stop and clear CA network"
	@echo "ca_docker_up: start Docker container wiht API and inject it in CA network"
	@echo "ca_docker_down: stop and clear API Docker container"
	@echo "ca_e2e_test: run end-2-end test with Docker container and API Newman tests on CA network"
	@echo ""

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
test: unit_test basic_e2e_test byfn_e2e_test ca_e2e_test
	@echo "\n\n-=-=-=- All tests passed successfully! -=-=-=-\n\n"

##### UNIT TESTS
unit_test:
	@go test ./pkg... ./cmd...

##### Docker
docker_build:
	@docker build -t frag:dev .

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

##### CA testing
ca_up:
	@./scripts/ca-up.sh

ca_clear:
	@./scripts/ca-clear.sh

ca_docker_up:
	@./scripts/ca-docker-up.sh

ca_docker_down:
	@docker stop frag
	@docker rm frag

ca_test:
	@./scripts/ca-test.sh

ca_e2e_test: docker_build ca_up ca_docker_up ca_test ca_docker_down ca_clear
