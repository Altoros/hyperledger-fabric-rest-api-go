.PHONY: all dev basic_clear build basic_up basic_down run basic_restart byfn_up byfn_test byfn_down byfn clear

help:
	@echo "Fabric REST API GoLang"
	@echo "up: start testnet"
	@echo "down: stop testnet"
	@echo "clean: clean testnet"
	@echo "restart: clean and start over testnet"

all: basic_clear build basic_up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@go get
	@go build -o=./build/rest-api main.go
	@echo "Build done"

##### RUN
run: build
	@echo "Start app ..."
	@./build/rest-api

##### BASIC NETWORK testing
basic_up:
	@echo "Start environment ..."
	@cd test/fixtures && docker-compose up --force-recreate -d
	@echo "Environment up"

basic_down:
	@echo "Stop environment ..."
	@cd test/fixtures && docker-compose down
	@echo "Environment down"

basic_restart: basic_clear basic_up

basic_clear: basic_down
	@echo "Clean up ..."
	@docker rm -f -v `docker ps -a --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "heroes-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"

##### CLEAR ALL
clear: basic_clear byfn_down

##### BYFN testing
byfn_up:
	@./scripts/byfn-up.sh

byfn_test: build
	@./scripts/byfn-test.sh

byfn_down:
	@./scripts/byfn-down.sh

byfn: byfn_up byfn_test byfn_down
